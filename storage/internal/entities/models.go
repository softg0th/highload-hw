package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	bsonPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectID = bsonPrimitive.ObjectID

type Document struct {
	ID        ObjectID  `bson:"_id,omitempty"`
	UserId    int       `bson:"user_id"`
	Timestamp time.Time `bson:"timestamp"`
	Message   string    `bson:"message"`
}

func NewDocument(userId int64, timestamp time.Time, message string) Document {
	return Document{
		ID:        primitive.NewObjectID(),
		UserId:    int(userId),
		Timestamp: timestamp,
		Message:   message,
	}
}
