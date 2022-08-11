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

// StageEnvironment string
type StageEnvironment string

// DB type constants
const (
	DevEnv   StageEnvironment = "dev"
	StageEnv StageEnvironment = "stage"
	TestEnv  StageEnvironment = "test"
	ProdEnv  StageEnvironment = "prod"
)

const defaultFileName = "defaults.yml"

var (
	defs = &defaults{}
)

// Load method
func (c *Config) Load() (err error) {

	if err = c.setDefaults(); err != nil {
		return err
	}

	// I want the environment vars to be the final say, but we need them for the SSM Params
	// hence calling it twice
	if err = c.setEnvVars(); err != nil {
		return err
	}

	if err = c.setSSMParams(); err != nil {
		return err
	}

	if err = c.setEnvVars(); err != nil {
		return err
	}

	c.setDBConnectURL()
	c.setFinal()

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
func (c *Config) validateStage() (err error) {

	validEnv := true

	switch defs.Stage {
	case "dev":
	case "development":
		c.Stage = DevEnv
	case "stage":
		c.Stage = StageEnv
	case "test":
		c.Stage = TestEnv
	case "prod":
		c.Stage = ProdEnv
	case "production":
		c.Stage = ProdEnv
	default:
		validEnv = false
	}

	if !validEnv {
		return errors.New(fmt.Sprintf("Invalid StageEnvironment requested: %s", defs.Stage))
	}

	return err
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
		// err = fmt.Errorf("Error fetching ssm params, total number found: %d", paramLen)
		return nil
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

	if c.GetStageEnv() != TestEnv {
		c.setAWSConnectURL()
		return c
	}

	c.MongoDBConnectURL = fmt.Sprintf("mongodb://%s/?readPreference=primary&ssl=false&directConnection=true", defs.MongoDBHost)

	return c
}

func (c *Config) setAWSConnectURL() {
	c.MongoDBConnectURL = fmt.Sprintf("mongodb+srv://%s/%s?authSource=%sexternal&authMechanism=MONGODB-AWS&retryWrites=true&w=majority", defs.MongoDBHost, defs.MongoDBName, "$")
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() (err error) {

	c.AWSRegion = defs.AWSRegion
	c.Dynamo = defs.Dynamo
	c.MongoDBName = defs.MongoDBName

	err = c.validateStage()

	return err
}
