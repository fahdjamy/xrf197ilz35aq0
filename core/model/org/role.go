package org

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
	"xrf197ilz35aq0/internal/random"
)

type Role struct {
	Name        string             `json:"name" bson:"name"`
	RoleId      string             `json:"roleId" bson:"roleId"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	Description string             `json:"description" bson:"description"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
}

func CreateRole(name, description string) *Role {
	now := time.Now()
	roleId := createRoleId()
	return &Role{
		Name:        name,
		CreatedAt:   now,
		UpdatedAt:   now,
		RoleId:      roleId,
		Description: description,
	}
}

func createRoleId() string {
	return strconv.FormatInt(random.PositiveInt64(), 10)
}
