package main

import (
	"api-pendencias/utils"
	"context"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/sirupsen/logrus"
)

// Define constants for logging levels and server port.
const (
	InfoLogLevel  = "I"
	ErrorLogLevel = "E"
	PortKey       = "PORT"
	DefaultPort   = "8080"
)

func main() {

	// Set log level from environment or config
	logs := utils.Init()

	dynamoClient, err := driver.NewDynamoClient(context.Background(), logs)
	if err != nil {
		logs.HandleError("F", "Failed to configure AWS", err)
		return
	}

	router := routers.SetupRouter(dynamoClient, logs)

	serverPort := ":8080" // Default port

	if isRunningInLambda() {
		logs.HandleError("E", "Failed to start server", gateway.ListenAndServe(serverPort, router))
	} else {
		logs.HandleError("E", "Failed to start server", http.ListenAndServe(serverPort, router))
	}
}

// isRunningInLambda checks if the application is running within a Lambda environment.
func isRunningInLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// Para compilar o binario do sistema usamos:
//
//	GOARCH=arm64 GOOS=linux  CGO_ENABLED=0 go build -tag lambda.norpc -o bootstrap .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip bootstrap
//
