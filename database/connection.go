package database

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// dynamoDBClient is a concrete implementation of the DBClient.
type DynamoDBClient struct {
	Client *dynamodb.Client
}

// NewDynamoClient initializes and returns a new DynamoDB client as a dynamoDBClient.
func NewDynamoClient(ctx context.Context, log *logging.Logger) (*DynamoDBClient, error) {
	configAWS := utils.ConfigAws(ctx, log)
	// Create and return the DynamoDB client
	client := &DynamoDBClient{Client: dynamodb.NewFromConfig(configAWS)}
	return client, nil
}
