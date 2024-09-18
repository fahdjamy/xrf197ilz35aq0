package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRequestId(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "generates request id", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRequestId()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(got) < 1 {
					t.Errorf("did not generate request id")
				}
			}
		})
	}
}
