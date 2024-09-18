package cmd

import (
	"testing"
)

func TestGenerateRequestId(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "generates request id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRequestId()
			if len(got) < 1 {
				t.Errorf("did not generate request id")
			}
		})
	}
}
