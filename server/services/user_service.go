package services

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"server/models"
)

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// UserEmailExists checks if a user with the given email exists in the users table.
// Since email is the partition key, we use GetItem for efficient lookup.
func UserEmailExists(ctx context.Context, email string) (bool, error) {
	client := GetDynamoClient()
	tableName := GetUsersTableName()

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
		ProjectionExpression: aws.String("email"),
	})
	if err != nil {
		return false, err
	}

	return result.Item != nil, nil
}

// InsertUser inserts a new user into the users table.
// Uses a condition expression to prevent overwriting existing users with the same email.
func InsertUser(ctx context.Context, u models.User) error {
	client := GetDynamoClient()
	tableName := GetUsersTableName()

	item, err := attributevalue.MarshalMap(u)
	if err != nil {
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	})

	return err
}

// FindUserByEmail retrieves a user by email from the users table.
// Email is the partition key, so this is a direct GetItem operation.
func FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	client := GetDynamoClient()
	tableName := GetUsersTableName()

	var u models.User

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return u, err
	}

	if result.Item == nil {
		return u, ErrUserNotFound
	}

	err = attributevalue.UnmarshalMap(result.Item, &u)
	return u, err
}

// FindUserByID retrieves a user by ID using the id-index GSI.
// The GSI has "id" as the partition key.
func FindUserByID(ctx context.Context, id string) (models.User, error) {
	client := GetDynamoClient()
	tableName := GetUsersTableName()

	var u models.User

	keyCond := expression.Key("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return u, err
	}

	result, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String("id-index"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return u, err
	}

	if len(result.Items) == 0 {
		return u, ErrUserNotFound
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &u)
	return u, err
}
