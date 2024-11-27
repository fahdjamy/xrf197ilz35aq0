package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"xrf197ilz35aq0/core/exchange"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/custom"
	xrfTest "xrf197ilz35aq0/internal/tests"
)

const (
	strongPassword    = "strongPassword"
	validEmailAddress = "test@xrfaq.com"
)

type settingServiceMock struct {
	Called int
}

var settingResponseMock = &exchange.SettingResponse{
	EncryptionKey: *custom.NewSecret[string](string(encryptionTestKey)),
}

func (s *settingServiceMock) NewSettings(_ *exchange.SettingRequest, _ string) (*exchange.SettingResponse, error) {
	s.Called++
	return settingResponseMock, nil
}

func TestUserService_CreateUser(t *testing.T) {
	logger := xrf.NewTestLogger()
	storeMock := xrf.NewStoreMock()
	userRepo := xrfTest.NewUserRepositoryMock()
	settingServiceMock := &settingServiceMock{}
	tests := []struct {
		name    string
		wantErr bool
		request *exchange.UserRequest
	}{
		{name: "invalid request containingEmail", wantErr: true, request: createUserRequest("wrongMail", strongPassword)},
		{name: "valid user request creates user", wantErr: false, request: createUserRequest(validEmailAddress, strongPassword)},
		{name: "saves user if lastName length is 0", wantErr: false, request: createUserRequest(validEmailAddress, strongPassword)},
		{name: "saves user if firstName length is 0", wantErr: false, request: createUserRequest(validEmailAddress, strongPassword)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserService(logger, settingServiceMock, storeMock, userRepo, context.TODO())
			got, err := uc.CreateUser(tt.request)
			if tt.wantErr {
				xrf.AssertError(t, err)
			} else {
				xrf.AssertNoError(t, err)
				assertUserResponse(t, got)
			}
		})
	}
}

func createUserRequest(email, password string) *exchange.UserRequest {
	secretEmail := custom.NewSecret(email)
	secretPass := custom.NewSecret(password)
	return &exchange.UserRequest{
		LastName:  "lastName",
		FirstName: "firstName",
		Password:  *secretPass,
		Email:     *secretEmail,
	}
}

func assertUserResponse(t *testing.T, got *exchange.UserResponse) {
	t.Helper()
	assert.NotNil(t, got)
	assert.False(t, got.Anonymous)
	assert.Equal(t, got.Anonymous, false)
	assert.True(t, time.Since(got.CreatedAt.Time) > 0)
	assert.True(t, time.Since(got.UpdatedAt.Time) > 0)
}
