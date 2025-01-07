package exchange

type OrgResponse struct {
	OrgId        string `json:"orgId"`
	Name         string `json:"name"`
	CreatedAt    string `json:"started"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	MembersCount int    `json:"membersCount"`
}

type OrgRequest struct {
	Name        string             `json:"name"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
	Members     []OrgMemberRequest `json:"members"`
}

type OrgMemberRequest struct {
	Owner  bool     `json:"owner"`
	UserId string   `json:"userId"`
	Roles  []string `json:"roles"`
}
