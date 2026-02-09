package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/models"
)

func ListJobs(ctx context.Context, userID string) ([]models.Job, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})

	cur, err := JobsCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var jobs []models.Job
	if err := cur.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func InsertJob(ctx context.Context, j models.Job) error {
	_, err := JobsCollection().InsertOne(ctx, j)
	return err
}

func UpdateJob(ctx context.Context, jobIDHex string, userID string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(jobIDHex)
	if err != nil {
		return err
	}

	_, err = JobsCollection().UpdateOne(ctx,
		bson.M{"_id": oid, "userId": userID},
		bson.M{"$set": update},
	)
	return err
}

func DeleteJob(ctx context.Context, jobIDHex string, userID string) error {
	oid, err := primitive.ObjectIDFromHex(jobIDHex)
	if err != nil {
		return err
	}

	_, err = JobsCollection().DeleteOne(ctx, bson.M{"_id": oid, "userId": userID})
	return err
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
