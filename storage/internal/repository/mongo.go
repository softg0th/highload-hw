package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	client *mongo.Client
}

type DataBase struct {
	Collection *mongo.Collection
}

func NewMongoConnection(uri string) (*MongoConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoConnection{client: client}, nil
}

func NewDataBase(conn *MongoConnection, dbName string, collectionName string) *DataBase {
	db := conn.client.Database(dbName)
	collection := db.Collection(collectionName)

	return &DataBase{Collection: collection}
}
