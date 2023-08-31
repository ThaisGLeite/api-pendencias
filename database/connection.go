package database

import (
	"api-pendencias/utils"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoClient initializes and returns a new DynamoDB client as a dynamoDBClient.
func NewDynamoClient() *dynamodb.Client {
	configAWS := utils.ConfigAws()
	// Create and return the DynamoDB client
	return dynamodb.NewFromConfig(configAWS)
}
