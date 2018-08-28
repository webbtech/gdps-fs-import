package main

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pulpfree/gales-fuelsale-export/auth"
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

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var err error
	var eRes string

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return events.APIGatewayProxyResponse{Body: "pong", StatusCode: 200}, nil
	}

	// Check for auth header
	if req.Headers["Authorization"] == "" {
		eRes = setErrorResponse(401, "Unauthorized", "Missing Authorization header")
		return events.APIGatewayProxyResponse{Body: eRes, StatusCode: 401}, nil
	}

	// Set auth config
	auth, err := auth.New(&auth.Config{
		ClientID:       cfg.CognitoClientID,
		PoolID:         cfg.CognitoPoolID,
		Region:         cfg.CognitoRegion,
		JwtAccessToken: req.Headers["Authorization"],
	})
	if err != nil {
		eRes = setErrorResponse(500, "Authentication", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, StatusCode: 500}, nil
	}

	// Validate JWT Token
	err = auth.Validate()
	if err != nil {
		eRes = setErrorResponse(401, "Authentication", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, StatusCode: 401}, nil
	}

	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)

	// Validate request params
	reqVars, err := validators.RequestVars(r)
	if err != nil {
		eRes = setErrorResponse(500, "RequestValidation", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, StatusCode: 500}, nil
	}

	// Intitialze and process request
	exporter := export.New(reqVars, cfg)
	res, err := exporter.Process()
	if err != nil {
		eRes = setErrorResponse(500, "ProcessExport", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(&res)
	if err != nil {
		eRes = setErrorResponse(500, "ProcessExport", "Unable to marshal JSON")
		return events.APIGatewayProxyResponse{Body: eRes, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

func main() {
	lambda.Start(handleRequest)
}

// ======================== Helper Function ================================= //

func setErrorResponse(status int, errType, message string) string {

	err := model.ErrorResponse{
		Status:  status,
		Type:    errType,
		Message: message,
	}
	log.Errorf("Error: status: %d, type: %s, message: %s", err.Status, err.Type, err.Message)
	res, _ := json.Marshal(&err)

	return string(res)
}
