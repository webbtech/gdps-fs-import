package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pulpfree/gsales-fs-export/config"
	"github.com/pulpfree/lambda-utils/pkgerrors"
)

// Response data format
type Response struct {
	Code      int         `json:"code"`      // HTTP status code
	Data      interface{} `json:"data"`      // Data payload
	Message   string      `json:"message"`   // Error or status message
	Status    string      `json:"status"`    // Status code (error|fail|success)
	Timestamp int64       `json:"timestamp"` // Machine-readable UTC timestamp in nanoseconds since EPOCH
}

var (
	cfg      *config.Config
	stdError *pkgerrors.StdError
)

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

	fmt.Printf("req.HTTPMethod: %+v\n", req.HTTPMethod)
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
		return gatewayResponse(Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs, nil), nil
	}

	return gatewayResponse(Response{
		Code:      201,
		Data:      "dummy",
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs, nil), nil

}

// old func
/* func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	// Initialize and process request
	exporter := export.New(reqVars, cfg)
	res, err := exporter.Process()
	if err != nil {
		eRes = setErrorResponse(500, "ProcessExport", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	log.Infof("res in exporter.Process(): %+v\n", res)

	body, err := json.Marshal(&res)
	if err != nil {
		eRes = setErrorResponse(500, "ProcessExport", "Unable to marshal JSON")
		return events.APIGatewayProxyResponse{Body: eRes, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
} */

func main() {
	lambda.Start(HandleRequest)
}

// ======================== Helper Function ================================= //

// old func
/* func setErrorResponse(status int, errType, message string) string {

	err := model.ErrorResponse{
		Status:  status,
		Type:    errType,
		Message: message,
	}
	log.Errorf("Error: status: %d, type: %s, message: %s", err.Status, err.Type, err.Message)
	res, _ := json.Marshal(&err)

	return string(res)
} */

func gatewayResponse(resp Response, hdrs map[string]string, err error) events.APIGatewayProxyResponse {

	if err != nil {
		resp.Code = 500
		resp.Status = "error"
		log.Error(err)
		// send friendly error to client
		if ok := errors.As(err, &stdError); ok {
			resp.Message = stdError.Msg
		} else {
			resp.Message = err.Error()
		}
	}
	body, _ := json.Marshal(&resp)

	return events.APIGatewayProxyResponse{Body: string(body), Headers: hdrs, StatusCode: resp.Code}
}
