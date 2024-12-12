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
	Called map[string]int
}

func newSettingServiceMock() *settingServiceMock {
	return &settingServiceMock{
		Called: make(map[string]int),
	}
}

func (s *settingServiceMock) GetUserSettings(_ string) (*exchange.SettingResponse, error) {
	method := "getSettingsForUser"
	count, ok := s.Called[method]
	if !ok {
		s.Called[method] = 1
	} else {
		s.Called[method] = count + 1
	}
	return &exchange.SettingResponse{}, nil
}

var settingResponseMock = &exchange.SettingResponse{
	EncryptionKey: *custom.NewSecret[string](string(encryptionTestKey)),
}

func (s *settingServiceMock) NewSettings(_ *exchange.SettingRequest, _ string) (*exchange.SettingResponse, error) {
	method := "newSettings"
	count, ok := s.Called[method]
	if !ok {
		s.Called[method] = 1
	} else {
		s.Called[method] = count + 1
	}
	return settingResponseMock, nil
}

func TestUserServiceCreateUser(t *testing.T) {
	logger := xrf.NewTestLogger()
	userRepo := xrfTest.NewUserRepositoryMock()
	settingServiceMock := newSettingServiceMock()
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
			uc := NewUserService(logger, settingServiceMock, userRepo, context.TODO())
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

func TestGetUserById(t *testing.T) {
	logger := xrf.NewTestLogger()
	userRepo := xrfTest.NewUserRepositoryMock()
	settingServiceMock := newSettingServiceMock()

	tests := []struct {
		name    string
		userId  int64
		want    *exchange.UserResponse
		wantErr bool
	}{
		{name: "get user by id", wantErr: false, userId: 1234567, want: &exchange.UserResponse{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &service{
				log:             logger,
				settingsService: settingServiceMock,
				ctx:             context.TODO(),
				userRepo:        userRepo,
			}
			got, err := uc.GetUserById(tt.userId)
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
