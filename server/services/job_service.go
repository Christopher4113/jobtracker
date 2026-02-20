package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"server/models"
)

// ListJobs retrieves all jobs for a user from the jobs table.
// Jobs table has userId as partition key and id as sort key.
// Results are returned and can be sorted by frontend (DynamoDB doesn't easily sort by non-key attributes).
func ListJobs(ctx context.Context, userID string) ([]models.Job, error) {
	client := GetDynamoClient()
	tableName := GetJobsTableName()

	keyCond := expression.Key("userId").Equal(expression.Value(userID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	var jobs []models.Job
	err = attributevalue.UnmarshalListOfMaps(result.Items, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// InsertJob inserts a new job into the jobs table.
func InsertJob(ctx context.Context, j models.Job) error {
	client := GetDynamoClient()
	tableName := GetJobsTableName()

	item, err := attributevalue.MarshalMap(j)
	if err != nil {
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	return err
}

// UpdateJob updates a job in the jobs table.
// The update parameter is a map of field names to new values.
// Jobs table has userId as partition key and id as sort key.
func UpdateJob(ctx context.Context, jobID string, userID string, update map[string]interface{}) error {
	client := GetDynamoClient()
	tableName := GetJobsTableName()

	var updateBuilder expression.UpdateBuilder
	first := true
	for key, value := range update {
		if first {
			updateBuilder = expression.Set(expression.Name(key), expression.Value(value))
			first = false
		} else {
			updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(value))
		}
	}

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		return err
	}

	_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID},
			"id":     &types.AttributeValueMemberS{Value: jobID},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	return err
}

// DeleteJob deletes a job from the jobs table.
// Jobs table has userId as partition key and id as sort key.
func DeleteJob(ctx context.Context, jobID string, userID string) error {
	client := GetDynamoClient()
	tableName := GetJobsTableName()

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID},
			"id":     &types.AttributeValueMemberS{Value: jobID},
		},
	})

	return err
}

// NowUTC returns the current UTC time as an RFC3339 formatted string.
// This is used for DynamoDB timestamp fields which are stored as strings.
func NowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}
