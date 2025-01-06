package org

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/random"
)

var externalError *xrfErr.External

type Member struct {
	Fingerprint string   `json:"-" bson:"fingerPrint"`
	Owner       bool     `json:"isOwner" bson:"owner"` // org can have multiple owners
	RoleIds     []string `json:"roleIds" bson:"roleIds"`
}

type Organization struct {
	Id          string             `bson:"orgId" json:"orgId"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Members     []Member           `bson:"members" json:"members"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
}

func CreateOrganization(name string, category string, desc string, members []Member) (*Organization, error) {
	externalError = &xrfErr.External{}
	if members == nil || len(members) == 0 {
		externalError.Message = "At least one member is required"
		return nil, externalError
	}
	if len(name) < 2 || len(name) > 20 {
		externalError.Message = "Name must be between 2 and 20 characters long"
		return nil, externalError
	}
	now := time.Now()
	orgId := createOrgId()

	return &Organization{
		Id:          orgId,
		Name:        name,
		CreatedAt:   now,
		UpdatedAt:   now,
		Description: desc,
		Members:     members,
		Category:    category,
	}, nil
}

func createOrgId() string {
	return strconv.FormatInt(random.PositiveInt64(), 10)
}
