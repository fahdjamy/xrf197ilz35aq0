package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	xrf "xrf197ilz35aq0/internal"
	xrfTest "xrf197ilz35aq0/internal/tests"
)

var encryptionTestKey = xrf.RandomBytes(32)

func TestNewSettings(t *testing.T) {
	logger := xrf.NewTestLogger()
	settingsRepoMock := xrfTest.NewSettingsRepositoryMock()
	type args struct {
		request   *exchange.SettingRequest
		userModel user.User
	}
	userObj := user.User{Id: fmt.Sprintf("%v", 123433)}
	tests := []struct {
		name           string
		args           args
		want           *exchange.SettingResponse
		wantErr        assert.ErrorAssertionFunc
		assertResponse bool
	}{
		{
			name:           "creates a new setting successfully when RotateKey is false",
			wantErr:        assert.NoError,
			args:           args{userModel: userObj, request: createSettingRequest(string(encryptionTestKey), false, 0)},
			assertResponse: true,
		},
		{
			name:           "creates a new setting successfully when RotateKey is true and key is not provided",
			wantErr:        assert.NoError,
			args:           args{userModel: userObj, request: createSettingRequest(string(encryptionTestKey), true, 3)},
			assertResponse: true,
		},
		{
			name:           "creates a new setting successfully when RotateKey is true and key is not provided",
			wantErr:        assert.NoError,
			args:           args{userModel: userObj, request: createSettingRequest("", true, 3)},
			assertResponse: true,
		},
		{
			name:    "returns an error if RotateKey is true and RotateAfter is 0",
			wantErr: assert.Error,
			args:    args{userModel: userObj, request: createSettingRequest("", true, 0)},
		},
		{
			name:    "returns an error if RotateKey is true and RotateBefore is greater than 12",
			wantErr: assert.Error,
			args:    args{userModel: userObj, request: createSettingRequest("", true, 13)},
		},
		{
			name:    "returns an error if user provided rotation key is less than 31 characters",
			wantErr: assert.Error,
			args:    args{userModel: userObj, request: createSettingRequest("not-very-long", false, 6)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewSettingService(logger, settingsRepoMock, context.TODO())
			got, err := manager.NewSettings(tt.args.request, "userTestVVFingerXXPrintLL")
			if !tt.wantErr(t, err, fmt.Sprintf("NewSettings(%v, %v)", tt.args.request, tt.args.userModel)) {
				return
			}
			if tt.assertResponse {
				assertSettingResponse(t, got)
			}
		})
	}
}

func assertSettingResponse(t *testing.T, response *exchange.SettingResponse) {
	t.Helper()
	assert.NotNil(t, response)
}

func createSettingRequest(key string, rotate bool, rotateAfter int) *exchange.SettingRequest {
	request := &exchange.SettingRequest{
		RotateKey: rotate,
	}
	if key != "" {
		request.EncryptionKey = key
	}
	if rotateAfter != 0 {
		request.RotateAfter = rotateAfter
	}
	return request
}
