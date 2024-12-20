package org

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Member struct {
	Fingerprint string   `json:"-" bson:"fingerPrint"`
	RoleIds     []string `json:"roleIds" bson:"roleIds"`
}

type Organization struct {
	Id        string             `bson:"orgId" json:"orgId"`
	Name      string             `bson:"name" json:"name"`
	Members   []Member           `bson:"members" json:"members"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	MongoID   primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
}
