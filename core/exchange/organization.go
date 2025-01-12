package exchange

import (
	"time"
)

type OrgResponse struct {
	OrgId        string    `json:"orgId"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	MembersCount int       `json:"membersCount"`
	IsAnonymous  bool      `json:"isAnonymous"`
}

type OrgRequest struct {
	Name        string             `json:"name"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
	Members     []OrgMemberRequest `json:"members"`
	IsAnonymous bool               `json:"isAnonymous"`
}

type OrgMemberRequest struct {
	Owner       bool     `json:"owner"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

type OrgMemberResponse struct {
	Email       string   `json:"email"`
	UserId      string   `json:"userId"`
	Permissions []string `json:"permissions"`
}
