package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SessionID string             `bson:"session_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	IsAdmin   bool               `bson:"isadmin"`
	Issued    time.Time          `bson:"timestamp"`
}
