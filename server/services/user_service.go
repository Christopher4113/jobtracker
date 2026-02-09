package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/models"
)

func ensureIndexes(ctx context.Context) error {
	col := UsersCollection()
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func UserEmailExists(ctx context.Context, email string) (bool, error) {
	col := UsersCollection()
	count, err := col.CountDocuments(ctx, bson.M{"email": email})
	return count > 0, err
}

func InsertUser(ctx context.Context, u models.User) error {
	_, err := UsersCollection().InsertOne(ctx, u)
	return err
}

func FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := UsersCollection().FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return u, err
}

func FindUserByID(ctx context.Context, idHex string) (models.User, error) {
	var u models.User
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return u, err
	}
	err = UsersCollection().FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	return u, err
}
