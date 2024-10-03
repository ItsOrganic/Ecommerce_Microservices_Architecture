package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client     *mongo.Client
	DB         *mongo.Database
	Collection *mongo.Collection
}

var MI MongoInstance

func Connect(uri, dbName, collectionName string) error {
	// Create a new client with the connection options.
	clientOpts := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return err
	}

	// Ping the MongoDB server to ensure the connection is established
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Initialize the database and collection
	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	// Assign the client, database, and collection to the MongoInstance
	MI = MongoInstance{
		Client:     client,
		DB:         db,
		Collection: collection,
	}

	return nil
}
