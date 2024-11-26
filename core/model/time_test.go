package model

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
	"xrf197ilz35aq0/internal"
)

func TestTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		wantErr   bool
		name      string
		want      []byte
		inputTime Time
	}{
		{
			name:      "marshals to valid json bytes",
			wantErr:   false,
			want:      []byte(`"2009-11-17T20:34:58Z"`),
			inputTime: NewTime(time.Date(2009, 11, 17, 20, 34, 58, 0, time.UTC)),
		},
		{
			name:      "fails to marshall null time",
			want:      []byte(`null`),
			wantErr:   true,
			inputTime: Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			got, err := tt.inputTime.MarshalJSON()
			if tt.wantErr {
				internal.AssertError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				internal.AssertNoError(t, err)
				if !reflect.DeepEqual(got, tt.want) {
					t1.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
