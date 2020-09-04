package config

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
	MongoDBName       string
	Stage             StageEnvironment
}

// Dynamo struct
type Dynamo struct {
	APIVersion string `yaml:"APIVersion"`
	Endpoint   string `yaml:"Endpoint"`
	Region     string `yaml:"Region"`
}
