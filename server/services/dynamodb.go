package services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoClient *dynamodb.Client

// InitDynamoDB initializes the DynamoDB client using AWS SDK default configuration.
// It automatically loads credentials and region from environment variables or IAM role.
// When running in Lambda, the AWS_REGION is set automatically.
func InitDynamoDB() error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	return nil
}

// GetDynamoClient returns the initialized DynamoDB client.
func GetDynamoClient() *dynamodb.Client {
	return dynamoClient
}

// GetUsersTableName returns the users table name from environment variable.
func GetUsersTableName() string {
	tableName := os.Getenv("USERS_TABLE_NAME")
	if tableName == "" {
		tableName = "users"
	}
	return tableName
}

// GetJobsTableName returns the jobs table name from environment variable.
func GetJobsTableName() string {
	tableName := os.Getenv("JOBS_TABLE_NAME")
	if tableName == "" {
		tableName = "jobs"
	}
	return tableName
}
