package org

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/random"
)

var externalError *xrfErr.External

type Member struct {
	Fingerprint string   `json:"-" bson:"fingerPrint"`
	Owner       bool     `json:"isOwner" bson:"owner"` // org can have multiple owners
	Permissions []string `json:"permissions" bson:"permissions"`
}

type Organization struct {
	Id          string             `bson:"orgId" json:"orgId"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	DisplayName string             `bson:"displayName" json:"displayName"`
	IsAnonymous bool               `bson:"isAnonymous" json:"isAnonymous"`
	Description string             `bson:"description" json:"description"`
	Members     map[string]Member  `bson:"members" json:"members"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
}

func CreateOrganization(name string, category string, desc string, anonymous bool, members map[string]Member) (*Organization, error) {
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
		CreatedAt:   now,
		UpdatedAt:   now,
		Description: desc,
		DisplayName: name,
		Members:     members,
		Category:    category,
		IsAnonymous: anonymous,
		Name:        strings.ToLower(name),
	}, nil
}

func CreateMember(userFp string, isOwner bool, permissions []string) *Member {
	return &Member{
		Fingerprint: userFp,
		Owner:       isOwner,
		Permissions: permissions,
	}
}

func createOrgId() string {
	return strconv.FormatInt(random.PositiveInt64(), 10)
}
