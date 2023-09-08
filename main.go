package main

import (
	"api-pendencias/database"
	"api-pendencias/router"
	"api-pendencias/utils"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

// Define constants for logging levels and server port.
const (
	InfoLogLevel  = "I"
	ErrorLogLevel = "E"
	DefaultPort   = ":8080"
)

func main() {
	//Criar o log
	log := utils.Init()
	dynamoCliente := database.NewDynamoClient()
	database := &database.Connection{
		Logs:            log,
		DynamodbCliente: dynamoCliente,
	}
	server := router.StartServer(database)

	if isRunningInLambda() {
		utils.HandleError("E", "Failed to start server", gateway.ListenAndServe(DefaultPort, server))
	} else {
		utils.HandleError("E", "Failed to start server", http.ListenAndServe(DefaultPort, server))
	}
}

// isRunningInLambda checks if the application is running within a Lambda environment.
func isRunningInLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// Para compilar o binario do sistema usamos:
//
//	GOARCH=arm64 GOOS=linux  CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip bootstrap
//
