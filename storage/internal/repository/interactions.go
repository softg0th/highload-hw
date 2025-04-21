package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"storage/internal/entities"
)

func (d *DataBase) InsertPostsMongoStream(ctx context.Context, doc entities.Document) error {
	collection := d.Collection
	done := make(chan error, 1)

	go func() {
		defer close(done)
		_, err := collection.InsertOne(ctx, doc)
		if err != nil {
			log.Fatalf("inserting error:%v", err)
		}
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (d *DataBase) GetLastMessages(ctx context.Context) ([]entities.Document, error) {
	collection := d.Collection
	output := make(chan bson.D)
	errChan := make(chan error, 1)

	go func() {
		defer close(output)

		opts := options.Find().
			SetSort(bson.D{{Key: "_id", Value: -1}}).
			SetLimit(100)

		cursor, err := collection.Find(ctx, bson.D{}, opts)
		if err != nil {
			errChan <- err
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var doc bson.D
			if err := cursor.Decode(&doc); err != nil {
				errChan <- err
				return
			}
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case output <- doc:
			}
		}

		if err := cursor.Err(); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	var docs []entities.Document
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
			return docs, nil
		case doc, ok := <-output:
			if !ok {
				continue
			}
			var entity entities.Document
			bsonBytes, _ := bson.Marshal(doc)
			bson.Unmarshal(bsonBytes, &entity)
			docs = append(docs, entity)
		}
	}
}
