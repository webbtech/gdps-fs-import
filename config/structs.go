package config

// Config struct
type Config struct {
	config
	DefaultsFilePath string
}

// defaults struct
type defaults struct {
	AWSRegion   string  `yaml:"AWSRegion"`
	Dynamo      *Dynamo `yaml:"Dynamo"`
	MongoDBHost string  `yaml:"MongoDBHost"`
	MongoDBName string  `yaml:"MongoDBName"`
	SsmPath     string  `yaml:"SsmPath"`
	Stage       string  `yaml:"Stage"`
}

type config struct {
	AWSRegion         string
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
