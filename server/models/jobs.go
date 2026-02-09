package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus string

const (
	StatusApplied     JobStatus = "applied"
	StatusInterviewing JobStatus = "interviewing"
	StatusOffer       JobStatus = "offer"
	StatusRejected    JobStatus = "rejected"
)

type Job struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	UserID          string             `bson:"userId" json:"userId"`
	Company         string             `bson:"company" json:"company"`
	Role            string             `bson:"role" json:"role"`
	Location        string             `bson:"location" json:"location"`
	Status          JobStatus           `bson:"status" json:"status"`
	StatusUpdatedAt time.Time          `bson:"statusUpdatedAt" json:"statusUpdatedAt"`

	// optional but useful
	Link   string `bson:"link,omitempty" json:"link,omitempty"`
	Notes  string `bson:"notes,omitempty" json:"notes,omitempty"`
	Source string `bson:"source,omitempty" json:"source,omitempty"` // referral, linkedin, etc

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
