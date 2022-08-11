package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegSuite struct {
	suite.Suite
	cfg *Config
}

// SetupTest method
func (suite *IntegSuite) SetupTest() {

	suite.cfg = &Config{}

	os.Setenv("Stage", "test")
	suite.cfg.setDefaults()
	suite.cfg.setEnvVars()
}

// TestSetDefaults method
func (suite *IntegSuite) TestSetDefaults() {
	err := suite.cfg.setDefaults()
	tAWSRegion := "ca-central-1"
	suite.NoError(err)
	suite.Equal(tAWSRegion, defs.AWSRegion)
}

// TestSetEnvVars method
func (suite *IntegSuite) TestSetEnvVars() {

	var err error
	err = suite.cfg.setEnvVars()
	suite.NoError(err)

	// Change a var
	os.Setenv("Stage", "noexist")
	err = suite.cfg.setEnvVars()
	suite.EqualError(err, "Invalid StageEnvironment requested: noexist")

	// Reset to valid stage
	os.Setenv("Stage", "test")
	suite.cfg.setEnvVars()
}

// TestValidateStage method
func (suite *IntegSuite) TestValidateStage() {
	err := suite.cfg.validateStage()
	suite.NoError(err)
}

// TestSetSSMParams function
// this test assumes that the CognitoClientID is empty
func (suite *IntegSuite) TestSetSSMParams() {

	// Set Stage to production
	os.Setenv("Stage", "prod")
	suite.cfg.setEnvVars()

	// CognitoClientIDBefore := defs.CognitoClientID
	err := suite.cfg.setSSMParams()
	suite.NoError(err)

	// Reset to test stage
	os.Setenv("Stage", "test")
	suite.cfg.setEnvVars()
}

// TestSetFinal function
func (suite *IntegSuite) TestSetFinal() {

	var se StageEnvironment
	err := suite.cfg.setFinal()

	suite.NoError(err)
	suite.Equal(suite.cfg.AWSRegion, defs.AWSRegion, "Expected Config.AWSRegion (%s) to equal defs.AWSRegion (%s)", suite.cfg.AWSRegion, defs.AWSRegion)
	suite.IsType(se, suite.cfg.Stage)
}

// TestGetMongoConnectURL function
func (suite *IntegSuite) TestGetMongoConnectURL() {
	suite.cfg.setDBConnectURL()
	url := suite.cfg.GetMongoConnectURL()
	fmt.Printf("url: %+v\n", url)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
