package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	yaml "gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	config
	DefaultsFilePath string
}

// defaults struct
type defaults struct {
	AWSRegion       string  `yaml:"AWSRegion"`
	CognitoClientID string  `yaml:"CognitoClientID"`
	CognitoPoolID   string  `yaml:"CognitoPoolID"`
	CognitoRegion   string  `yaml:"CognitoRegion"`
	Dynamo          *Dynamo `yaml:"Dynamo"`
	MongoDBHost     string  `yaml:"MongoDBHost"`
	MongoDBName     string  `yaml:"MongoDBName"`
	MongoDBPassword string  `yaml:"MongoDBPassword"`
	MongoDBUser     string  `yaml:"MongoDBUser"`
	SsmPath         string  `yaml:"SsmPath"`
	Stage           string  `yaml:"Stage"`
}

type config struct {
	AWSRegion         string
	CognitoClientID   string
	CognitoPoolID     string
	CognitoRegion     string
	Dynamo            *Dynamo
	MongoDBConnectURL string
	Stage             StageEnvironment
}

// Dynamo struct
type Dynamo struct {
	APIVersion string `yaml:"APIVersion"`
	Endpoint   string `yaml:"Endpoint"`
	Region     string `yaml:"Region"`
}

// StageEnvironment string
type StageEnvironment string

// DB type constants
const (
	DevEnv   StageEnvironment = "dev"
	StageEnv StageEnvironment = "stage"
	TestEnv  StageEnvironment = "test"
	ProdEnv  StageEnvironment = "prod"
)

const defaultFileName = "defaults.yaml"

var (
	defs = &defaults{}
)

// Load method
func (c *Config) Load() (err error) {

	err = c.setDefaults()
	if err != nil {
		return err
	}
	err = c.setEnvVars()
	if err != nil {
		return err
	}
	err = c.setSSMParams()
	if err != nil {
		return err
	}
	err = c.setFinal()
	if err != nil {
		return err
	}

	c.setDBConnectURL()
	return err
}

// GetStageEnv method
func (c *Config) GetStageEnv() StageEnvironment {
	return c.Stage
}

// GetMongoConnectURL method
func (c *Config) GetMongoConnectURL() string {
	return c.MongoDBConnectURL
}

// this must be called first in c.Load
func (c *Config) setDefaults() (err error) {

	if c.DefaultsFilePath == "" {
		dir, _ := os.Getwd()
		c.DefaultsFilePath = path.Join(dir, defaultFileName)
	}

	file, err := ioutil.ReadFile(c.DefaultsFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), &defs)
	if err != nil {
		return err
	}
	err = c.validateStage()

	return err
}

// validateStage method to validate Stage value
func (c *Config) validateStage() error {

	validEnv := false

	switch defs.Stage {
	case "dev":
		c.Stage = DevEnv
		validEnv = true
	case "stage":
		c.Stage = StageEnv
		validEnv = true
	case "test":
		c.Stage = TestEnv
		validEnv = true
	case "prod":
		c.Stage = ProdEnv
		validEnv = true
	case "production":
		c.Stage = ProdEnv
		validEnv = true
	}

	if !validEnv {
		return errors.New("Invalid Stage type")
	}

	return nil
}

// sets any environment variables that match the default struct fields
func (c *Config) setEnvVars() (err error) {

	vals := reflect.Indirect(reflect.ValueOf(defs))
	for i := 0; i < vals.NumField(); i++ {
		nm := vals.Type().Field(i).Name
		if e := os.Getenv(nm); e != "" {
			vals.Field(i).SetString(e)
		}
		// If field is Stage, validate and return error if required
		if nm == "Stage" {
			err = c.validateStage()
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (c *Config) setSSMParams() (err error) {

	s := []string{"", string(c.GetStageEnv()), defs.SsmPath}
	paramPath := aws.String(strings.Join(s, "/"))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(defs.AWSRegion),
	})
	if err != nil {
		return err
	}

	svc := ssm.New(sess)
	res, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           paramPath,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	paramLen := len(res.Parameters)
	if paramLen == 0 {
		err = fmt.Errorf("Error fetching ssm params, total number found: %d", paramLen)
	}

	// Get struct keys so we can test before attempting to set
	t := reflect.ValueOf(defs).Elem()
	for _, r := range res.Parameters {
		paramName := strings.Split(*r.Name, "/")[3]
		structKey := t.FieldByName(paramName)
		if structKey.IsValid() {
			structKey.Set(reflect.ValueOf(*r.Value))
		}
	}
	return err
}

// Build a url used in mgo.Dial as described in: https://godoc.org/gopkg.in/mgo.v2#Dial
func (c *Config) setDBConnectURL() *Config {

	var db, userPass, authSource string

	if defs.MongoDBName != "" {
		db = "/" + defs.MongoDBName
	}

	if defs.MongoDBUser != "" && defs.MongoDBPassword != "" {
		userPass = defs.MongoDBUser + ":" + defs.MongoDBPassword + "@"
	}

	if userPass != "" {
		authSource = "?authSource=admin"
	}

	c.MongoDBConnectURL = "mongodb://" + userPass + defs.MongoDBHost + db + authSource

	return c
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() (err error) {

	c.AWSRegion = defs.AWSRegion
	c.CognitoClientID = defs.CognitoClientID
	c.CognitoPoolID = defs.CognitoPoolID
	c.CognitoRegion = defs.CognitoRegion
	c.Dynamo = defs.Dynamo
	err = c.validateStage()

	return err
}
