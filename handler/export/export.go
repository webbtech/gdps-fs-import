package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pulpfree/gales-fuelsale-export/config"
	"github.com/pulpfree/gales-fuelsale-export/export"
	"github.com/pulpfree/gales-fuelsale-export/model"
	"github.com/pulpfree/gales-fuelsale-export/validators"
)

var cfg *config.Config

const defaultsFilePath = "./defaults.yaml"

func init() {
	cfg = &config.Config{
		DefaultsFilePath: defaultsFilePath,
	}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)

	var err error
	reqVars, err := validators.RequestVars(r)

	exporter := export.New(reqVars, cfg)
	err = exporter.Process()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(&reqVars)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Unable to marshal JSON", StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

func main() {
	lambda.Start(handleRequest)
}
