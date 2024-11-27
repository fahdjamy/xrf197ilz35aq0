package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/storage"
)

type Database struct {
	log          internal.Logger
	db           *mongo.Database
	client       *mongo.Client
	databaseName string
	context      context.Context
	decodeTo     *internal.Serializable
}

func (d *Database) Save(collection string, obj internal.Serializable, ctx context.Context) (any, error) {
	d.log.Debug("saving new object")
	document, err := d.db.Collection(collection).InsertOne(d.context, obj)
	if err != nil {
		return nil, err
	}
	d.log.Debug(fmt.Sprintf("saved new object :: objectID=%v", document.InsertedID))
	return document, nil
}

func (d *Database) FindById(collection string, id int64, ctx context.Context) (*internal.Serializable, error) {
	filter := bson.M{"Id": id}
	coll := d.db.Collection(collection)
	err := coll.FindOne(d.context, filter).Decode(d.decodeTo)
	if err != nil {
		return nil, err
	}

	return d.decodeTo, nil
}

func NewStore(log internal.Logger, client *mongo.Client, dbName string, ctx context.Context) storage.Store {
	database := client.Database(dbName)
	return &Database{
		log:     log,
		context: ctx,
		client:  client,
		db:      database,
	}
}
