package tests

import (
	"context"
	"io"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/core/repository"
)

// MockFileDataCopier for os.Open
type MockFileDataCopier struct {
	Content []byte
	err     error
	readPos int // Keep track of the current read position
	Closed  bool
}

// Close since this is a mock, close should do nothing and just return nil error
func (m *MockFileDataCopier) Close() error {
	m.Closed = true
	return nil
}

// MockFileDataCopier's Read needs to fill the p buffer with data from its internal content,
// but starting from the readPos which tracks where we are in the simulated file
func (m *MockFileDataCopier) Read(p []byte) (n int, err error) {
	// In real file, when you call Read,
	// it fills the provided buffer (p) with data from the file (content),
	// starting from the current file position.
	if m.readPos >= len(m.Content) {
		return 0, io.EOF // Simulate end-of-file when all content is read
	}

	// The copy is used to efficiently transfer data from one slice to another.
	// copies a portion of the m.content slice (starting from m.readPos) into the p slice
	// n captures the number of bytes actually copied from m.content to p.
	n = copy(p, m.Content[m.readPos:])

	m.readPos += n
	return n, m.err
}

type userRepositoryMock struct {
	Called map[string]int
}

func (u *userRepositoryMock) FindUsersByEmails(_ []string, _ context.Context) (*[]user.User, error) {
	method := "FindUsersByEmails"
	count, ok := u.Called[method]
	if !ok {
		u.Called[method] = 1
	} else {
		u.Called[method] = count + 1
	}
	return &[]user.User{}, nil
}

func (u *userRepositoryMock) GetUserById(_ string, _ context.Context) (*user.User, error) {
	method := "GetUserById"
	count, ok := u.Called[method]
	if !ok {
		u.Called[method] = 1
	} else {
		u.Called[method] = count + 1
	}

	return &user.User{}, nil
}

func (u *userRepositoryMock) UpdatePassword(_ string, _ string, _ context.Context) (bool, error) {
	method := "UpdatePassword"
	count, ok := u.Called[method]
	if !ok {
		u.Called[method] = 1
	} else {
		u.Called[method] = count + 1
	}

	return true, nil
}

func (u *userRepositoryMock) CreateUser(_ *user.User, _ context.Context) (string, error) {
	method := "CreateUser"
	count, ok := u.Called[method]
	if !ok {
		u.Called[method] = 1
	} else {
		u.Called[method] = count + 1
	}
	return "MockUserId12345", nil
}

func NewUserRepositoryMock() repository.UserRepository {
	return &userRepositoryMock{
		Called: make(map[string]int),
	}
}

type settingsRepositoryMock struct {
	Called map[string]int
}

func (s *settingsRepositoryMock) FetchUserSettings(_ context.Context, _ string) (settings *user.Settings, err error) {
	method := "FetchUserSettings"
	count, ok := s.Called[method]
	if !ok {
		s.Called[method] = 1
	} else {
		s.Called[method] = count + 1
	}
	return &user.Settings{}, nil
}

func (s *settingsRepositoryMock) CreateSettings(settings *user.Settings, _ context.Context) (any, error) {
	method := "CreateSettings"
	count, ok := s.Called[method]
	if !ok {
		s.Called[method] = 1
	} else {
		s.Called[method] = count + 1
	}

	return settings, nil
}

func NewSettingsRepositoryMock() repository.SettingsRepository {
	return &settingsRepositoryMock{
		Called: make(map[string]int),
	}
}
