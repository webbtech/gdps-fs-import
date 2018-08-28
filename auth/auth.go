package auth

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"reflect"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	log "github.com/sirupsen/logrus"
)

// Config struct
type Config struct {
	ClientID       string `json:"clientId"`
	PoolID         string `json:"poolId"`
	Region         string `json:"region"`
	JwtAccessToken string `json:"token"`
}

// Auth struct
type Auth struct {
	cfg           *Config
	WellKnownJWKs *jwk.Set
}

// New function
func New(cfg *Config) (a *Auth, err error) {

	a = &Auth{}

	err = a.validateCfg(cfg)
	if err != nil {
		return a, err
	}

	return a, err
}

// Validate method
func (a *Auth) Validate() (err error) {

	err = a.getWellKnownJWTKs()
	if err != nil {
		return err
	}

	var errStr string

	token, err := jwt.Parse(a.cfg.JwtAccessToken, func(token *jwt.Token) (interface{}, error) {

		// Validate the expected alg
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			errStr = fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			log.Error(errStr)
			return nil, errors.New(errStr)
		}

		// Looking up the key id will return an array of just one key
		keys := a.WellKnownJWKs.LookupKeyID(token.Header["kid"].(string))
		if len(keys) == 0 {
			errStr = "could not find matching `kid` in well known tokens"
			log.Error(errStr)
			return nil, errors.New(errStr)
		}

		// Build the public RSA key
		key, err := keys[0].Materialize()
		if err != nil {
			errStr = fmt.Sprintf("Failed to create public key: %s", err)
			log.Error(errStr)
			return nil, errors.New(errStr)
		}

		return key.(*rsa.PublicKey), nil
	})

	// Populated when you Parse/Verify a token
	// First verify the token itself is a valid format
	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Then check time based claims; exp, iat, nbf
			err = claims.Valid()
			if err != nil {
				errStr = fmt.Sprintf("Invalid claims for id token. Error: %s", err)
				log.Error(errStr)
				return errors.New(errStr)
			}
		}

	} else if ve, ok := err.(*jwt.ValidationError); ok {

		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			errStr = "Invalid malformed token"
			err = errors.New(errStr)

		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			errStr = "Token expired or inactive"
			err = errors.New(errStr)

		} else {
			errStr = fmt.Sprintf("Couldn't handle this token: %s", err)
		}

		if err != nil {
			log.Error(errStr)
			return err
		}

	} else {
		log.Errorf("Couldn't handle this token:", err)
		return err
	}

	return nil
}

func (a *Auth) validateCfg(cfg *Config) (err error) {

	cVals := reflect.ValueOf(cfg).Elem()
	for i := 0; i < cVals.NumField(); i++ {
		nm := cVals.Type().Field(i).Name
		f := cVals.Field(i)
		if reflect.Value(f).String() == "" {
			e := fmt.Sprintf("missing configuration value for key: %s", nm)
			return errors.New(e)
		}
	}
	a.cfg = cfg

	return err
}

// getWellKnownJWTKs gets the well known JSON web token key set for this client's user pool
func (a *Auth) getWellKnownJWTKs() error {
	// https://cognito-idp.<region>.amazonaws.com/<pool_id>/.well-known/jwks.json
	var buffer bytes.Buffer
	buffer.WriteString("https://cognito-idp.")
	buffer.WriteString(a.cfg.Region)
	buffer.WriteString(".amazonaws.com/")
	buffer.WriteString(a.cfg.ClientID)
	buffer.WriteString("/.well-known/jwks.json")
	wkjwksURL := buffer.String()
	buffer.Reset()

	// Use this cool package
	set, err := jwk.Fetch(wkjwksURL)
	if err == nil {
		a.WellKnownJWKs = set
	} else {
		log.Error("There was a problem getting the well known JSON web token key set")
		log.Error(err)
	}
	return err
}
