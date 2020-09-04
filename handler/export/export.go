package main

import (
	"encoding/json"
	"time"

	pres "github.com/pulpfree/lambda-go-proxy-response"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/pulpfree/gsales-fs-export/config"
	"github.com/pulpfree/gsales-fs-export/export"
	"github.com/pulpfree/gsales-fs-export/model"
	"github.com/pulpfree/gsales-fs-export/validators"
)

var cfg *config.Config

func init() {
	cfg = &config.Config{}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleRequest function
// NOTE: strange, the error parameter cannot be used or removed... would be good to dig into
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	hdrs["Access-Control-Allow-Origin"] = "*"
	hdrs["Access-Control-Allow-Methods"] = "GET,OPTIONS,POST,PUT"
	hdrs["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{Body: string("null"), Headers: hdrs, StatusCode: 200}, nil
	}

	t := time.Now()

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return pres.ProxyRes(pres.Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs, nil), nil
	}

	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)

	// Validate request params
	reqVars, err := validators.RequestVars(r)
	if err != nil {
		log.Errorf("err in validators.RequestVars: %+v with input of: %+v\n", err, r)
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	// Initialize and process request
	exporter := export.New(reqVars, cfg)
	res, err := exporter.Process()
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}
	log.Infof("res in exporter.Process(): %+v\n", res)

	body, err := json.Marshal(&res)
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	return pres.ProxyRes(pres.Response{
		Code:      201,
		Data:      body,
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs, nil), nil
}

func main() {
	lambda.Start(HandleRequest)
}
