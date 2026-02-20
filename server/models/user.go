package models

type User struct {
	ID           string `json:"id" dynamodbav:"id"`
	Name         string `json:"name" dynamodbav:"name"`
	Email        string `json:"email" dynamodbav:"email"`
	PasswordHash string `json:"-" dynamodbav:"passwordHash"`
	CreatedAt    string `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt    string `json:"updatedAt" dynamodbav:"updatedAt"`
}
