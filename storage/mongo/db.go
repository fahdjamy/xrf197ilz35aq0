package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/storage"
)

type Database struct {
	log          xrf197ilz35aq0.Logger
	db           *mongo.Database
	client       *mongo.Client
	databaseName string
	context      context.Context
	decodeTo     *xrf197ilz35aq0.Serializable
}

func (d *Database) Save(collection string, obj xrf197ilz35aq0.Serializable) (any, error) {
	d.log.Debug("saving new object")
	document, err := d.db.Collection(collection).InsertOne(d.context, obj)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func (d *Database) SetContext(ctx context.Context) {
	d.context = ctx
}

func (d *Database) FindById(collection string, id int64) (*xrf197ilz35aq0.Serializable, error) {
	filter := bson.M{"Id": id}
	coll := d.db.Collection(collection)
	err := coll.FindOne(d.context, filter).Decode(d.decodeTo)
	if err != nil {
		return nil, err
	}

	return d.decodeTo, nil
}

func NewStore(log xrf197ilz35aq0.Logger, client *mongo.Client, dbName string, ctx context.Context) storage.Store {
	database := client.Database(dbName)
	return &Database{
		log:     log,
		context: ctx,
		client:  client,
		db:      database,
	}
}