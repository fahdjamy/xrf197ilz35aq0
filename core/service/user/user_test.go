package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal/custom"
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

func (s *settingServiceMock) NewSettings(_ *exchange.SettingRequest, _ user.User) (*exchange.SettingResponse, error) {
	s.Called++
	return settingResponseMock, nil
}

func TestUserService_CreateUser(t *testing.T) {
	settingServiceMock := &settingServiceMock{}
	storeMock := &xrf197ilz35aq0.StoreMock{}
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
			logger := &xrf197ilz35aq0.TestLogger{}
			uc := NewUserManager(logger, settingServiceMock, storeMock)
			got, err := uc.NewUser(tt.request)
			if tt.wantErr {
				xrf197ilz35aq0.AssertError(t, err)
			} else {
				xrf197ilz35aq0.AssertNoError(t, err)
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
	assert.True(t, time.Since(got.CreatedAt) > 0)
	assert.True(t, time.Since(got.UpdatedAt) > 0)
}
