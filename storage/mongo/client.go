package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, dbUri, dbName string) (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the NewClient
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
