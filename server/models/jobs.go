package models

type JobStatus string

const (
	StatusApplied      JobStatus = "applied"
	StatusInterviewing JobStatus = "interviewing"
	StatusOffer        JobStatus = "offer"
	StatusRejected     JobStatus = "rejected"
)

type Job struct {
	ID              string    `json:"id" dynamodbav:"id"`
	UserID          string    `json:"userId" dynamodbav:"userId"`
	Company         string    `json:"company" dynamodbav:"company"`
	Role            string    `json:"role" dynamodbav:"role"`
	Location        string    `json:"location" dynamodbav:"location"`
	Status          JobStatus `json:"status" dynamodbav:"status"`
	StatusUpdatedAt string    `json:"statusUpdatedAt" dynamodbav:"statusUpdatedAt"`
	Link            string    `json:"link,omitempty" dynamodbav:"link,omitempty"`
	Notes           string    `json:"notes,omitempty" dynamodbav:"notes,omitempty"`
	Source          string    `json:"source,omitempty" dynamodbav:"source,omitempty"`
	CreatedAt       string    `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt       string    `json:"updatedAt" dynamodbav:"updatedAt"`
}
