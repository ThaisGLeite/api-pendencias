package utils

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// Logger wraps the logrus Logger
type Logger struct {
	*log.Logger
}

func Init() *Logger {
	logger := log.New()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	return &Logger{logger}
}

// HandleError logs the error based on the provided log level: "I" for Info, "E" for Error, and "W" for Warning.
// If the error is not nil, it logs the error and the associated message at the given log level.
func HandleError(logLevel, msg string, err error) {
	if err != nil {
		// Create an entry for the error
		entry := log.WithFields(log.Fields{
			"error": err,
		})

		// Log the error at the appropriate level
		switch logLevel {
		case "I":
			entry.Info(msg)
		case "E":
			entry.Error(msg)
		case "W":
			entry.Warn(msg)
		default:
			entry.Info(msg)
		}
	}
}

func ConfigAws() (*dynamodb.Client, error) {
	configAws, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{"database/data/credentials.aws"}),
		config.WithSharedConfigFiles([]string{"database/data/config.aws"}),
	)

	return dynamodb.NewFromConfig(configAws), err
}
