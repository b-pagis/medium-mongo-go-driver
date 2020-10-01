package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User application user.
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
}
