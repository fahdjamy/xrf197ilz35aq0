package org

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	type args struct {
		name     string
		category string
		desc     string
		members  []Member
	}
	validOrgName := "xrfOrg"
	roleIds := []string{"admin"}
	validMembers := []Member{{Fingerprint: "", RoleIds: roleIds}}

	tests := []struct {
		name    string
		args    args
		want    *Organization
		wantErr bool
	}{
		{name: "Missing members", args: args{name: validOrgName, category: "", desc: "", members: nil}, wantErr: true},
		{name: "Invalid organization name", args: args{name: "", category: "", desc: "", members: validMembers}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateOrganization(tt.args.name, tt.args.category, tt.args.desc, tt.args.members)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotNil(t, got.CreatedAt)
			}
		})
	}
}
