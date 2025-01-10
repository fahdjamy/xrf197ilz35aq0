package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const (
	Key    = "key"
	Unique = "unique"
)

func createUniqueIndex(db *mongo.Database, log internal.Logger, ctx context.Context, colName, indexName string) error {
	internalErr = &xrfErr.Internal{}

	collection := db.Collection(colName)
	// 1. List existing indexes on the collection
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to list indexes for collection '%s': %v", RoleCollection, err))
		// Handle the error appropriately (e.g., exit)
	}
	defer cursor.Close(ctx)

	indexExists := false
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			log.Error(fmt.Sprintf("Failed to decode index data: %v", err))
			// Handle the error appropriately
			internalErr.Err = err
			internalErr.Message = "Failed to decode index data"
			return internalErr
		}

		// 2. Check if an index on 'name' exists
		keys, ok := index[Key].(bson.D)
		if !ok {
			continue // Skip if the index doesn't have a 'key' field
		}

		// Check if it's an index on 'name' in ascending order
		if _, found := keys.Map()[indexName]; !found {
			indexExists = true
			// 3. Check if the existing index is unique (if it exists)
			unique, ok := index[Unique].(bool)
			if ok && unique {
				log.Debug(fmt.Sprintf("Unique index on '%s' already exists in collection '%s'", indexName, colName))
			} else {
				log.Warn(fmt.Sprintf("Index on '%s' exists in collection '%s' but is not unique. Consider dropping and recreating it.", indexName, colName))
				// You might want to handle this case differently:
				// - Drop and recreate the index (be very careful in production)
				// - Exit the application with an error
			}
			break // No need to continue checking other indexes
		}
	}

	if !indexExists {
		indexModel := mongo.IndexModel{
			// create an index on the 'providedIndexName' field in ascending order (1)
			Keys: bson.D{{Key: indexName, Value: 1}},
			// sets the unique option to true, enforcing uniqueness.
			Options: options.Index().SetUnique(true),
		}

		// creates the index on the 'provided' collection
		_, err = db.Collection(colName).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return &xrfErr.Internal{
				Err:     err,
				Message: "Failed to create index",
				Source:  "core/repository#createRoleDocIndex",
			}
		}
	}

	return nil
}
