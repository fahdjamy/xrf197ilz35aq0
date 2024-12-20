package org

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	RoleId      string             `json:"roleId" bson:"roleId"`
	Description string             `json:"description" bson:"description"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
}
