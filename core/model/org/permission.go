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
	RoleId      string             `json:"roleId" bson:"roleId"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	Description string             `json:"description" bson:"description"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
}

func CreateRole(name, description string) *Permission {
	now := time.Now()
	roleId := createPermissionId()
	return &Permission{
		CreatedAt:   now,
		UpdatedAt:   now,
		RoleId:      roleId,
		Description: description,
		Name:        strings.ToUpper(name),
	}
}

func createPermissionId() string {
	return strconv.FormatInt(random.PositiveInt64(), 10)
}
