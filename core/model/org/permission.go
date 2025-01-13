package org

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
	"xrf197ilz35aq0/internal/random"
)

type Permission struct {
	Name        string             `json:"name" bson:"name"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	Description string             `json:"description" bson:"description"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
	Id          string             `json:"permissionId" bson:"permissionId"`
}

func CreateRole(name, description string) *Permission {
	now := time.Now()
	id := createPermissionId()
	return &Permission{
		CreatedAt:   now,
		UpdatedAt:   now,
		Description: description,
		Id:          id,
		Name:        strings.ToUpper(name),
	}
}

func createPermissionId() string {
	return strconv.FormatInt(random.PositiveInt64(), 10)
}
