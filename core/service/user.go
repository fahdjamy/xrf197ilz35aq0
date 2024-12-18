package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"net/mail"
	"regexp"
	"runtime"
	"strings"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/custom"
	xrfErr "xrf197ilz35aq0/internal/error"
)

var internalError *xrfErr.Internal

// UserService is a port (Driven side)
type UserService interface {
	GetUserById(userId string) (*exchange.UserResponse, error)
	CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error)
}

type service struct {
	log             internal.Logger
	config          xrf.Security
	settingsService SettingsService
	ctx             context.Context
	userRepo        repository.UserRepository
}

func (uc *service) CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error) {
	internalError = &xrfErr.Internal{}
	internalError.Source = "core/service/user/user#createUser"

	userName := "Unknown name"
	if request.FirstName != "" && request.LastName != "" {
		userName = request.FirstName + " " + request.LastName
	}

	uc.log.Info(fmt.Sprintf("event=creatUser :: action=creatingUser :: username=%s", userName))
	err := uc.validateUser(request)
	if err != nil {
		uc.log.Error(fmt.Sprintf("event=creatUser :: username=%s :: err=%v", userName, err))
		return nil, err
	}

	hashedPassword, err := uc.hashPassword(request.Password.Data())
	if err != nil {
		internalError.Err = err
		internalError.Message = "Something went wrong"
		uc.log.Error(fmt.Sprintf("event=createUser :: action=hashPassword :: err=%v", err))
		return nil, internalError
	}
	newUser := user.NewUser(request.FirstName, request.LastName, request.Email.Data(), hashedPassword)

	// SAVE-USER/DB: ACTION 1 - save user and settings to database
	uc.log.Debug(fmt.Sprintf("event=creatUser :: action=saveUserINDB :: userFP=%s :: userId=%s", newUser.FingerPrint[:7], newUser.Id))
	_, err = uc.userRepo.CreateUser(newUser, uc.ctx)
	if err != nil {
		internalError.Err = err
		internalError.Message = "User creation failed"
		return nil, internalError
	}

	settingRequest := request.Settings
	if settingRequest == nil {
		settingRequest = &exchange.SettingRequest{
			RotateKey:   false,
			RotateAfter: 13,
		}
	}

	// SAVE-USER/DB: ACTION 2 - create user settings
	settings, err := uc.settingsService.NewSettings(settingRequest, newUser.FingerPrint)
	if err != nil {
		return nil, err
	}

	// Return userResponse
	userResponse := toUserResponse(newUser)
	userResponse.Settings = *settings
	return userResponse, nil
}

func (uc *service) GetUserById(userId string) (*exchange.UserResponse, error) {
	userResponse, err := uc.userRepo.GetUserById(userId, uc.ctx)
	uc.log.Debug(fmt.Sprintf("event=getUserById :: action=getUserByIdFromDB :: userId=%s", userId))
	if err != nil {
		uc.log.Error(fmt.Sprintf("event=getUserById :: action=getUserByIdFailure :: err=%v", err))
		return nil, err
	}

	uc.log.Debug(fmt.Sprintf("event=getUserById :: action=fetchedUserByIdFromDB :: userId=%s", userResponse.Id))
	userSettings, err := uc.settingsService.GetUserSettings(userResponse.FingerPrint)

	if err != nil {
		uc.log.Error(fmt.Sprintf("event=getUserById :: action=fetchedUserByIdFromDBFailure :: err=%v", err))
		return nil, err
	}

	response := toUserResponse(userResponse)
	response.Settings = *userSettings
	return response, nil
}

func (uc *service) validateUser(request *exchange.UserRequest) error {
	// is validEmail
	_, err := mail.ParseAddress(request.Email.Data())
	if err != nil {
		return &xrfErr.External{Message: "Invalid email address"}
	}
	lastNameLen := len(request.LastName)
	if lastNameLen != 0 && lastNameLen < 3 {
		return &xrfErr.External{Message: "If last name is specified, it should be at least 3 characters long"}
	}
	firstNameLen := len(request.FirstName)
	if firstNameLen != 0 && firstNameLen < 3 {
		return &xrfErr.External{Message: "If first name is specified, it should be at least 3 characters long"}
	}

	if err := uc.validatePassword(request.Password.Data()); err != nil {
		return err
	}
	return nil
}

func (uc *service) validatePassword(password string) error {
	// Minimum length of 8 characters
	// At least 1 uppercase letter and 1 lowercase letter
	// At least one digit
	externalErr := &xrfErr.External{
		Message: "Password should at least be 8 characters long and must contain a special character",
	}

	if len(password) < 8 {
		return externalErr
	}

	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) || !regexp.MustCompile(`[a-z]`).MatchString(password) {
		externalErr.Message = "password must contain at least one lowercase and an uppercase letter"
		return externalErr
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one digit")
	}

	return nil
}

func (uc *service) hashPassword(password string) (string, error) {
	internalError := &xrfErr.Internal{
		Source: "core/service/settings#hashPassword",
	}
	// Generate a random salt. It's crucial to use a unique salt for each password.
	salt := make([]byte, 16)

	if _, err := rand.Read(salt); err != nil {
		internalError.Err = err
		internalError.Message = "Error generating password salt"
		return "", internalError
	}

	// Use argon2.IDKey to generate the hash. Adjust parameters as needed:
	//   - time:  Number of iterations (higher is slower but more secure).
	//   - memory:  Memory usage in KiB (higher is more resistant to GPU cracking).
	//   - threads: Number of parallel threads (can improve performance).
	//   - keyLen: Length of the generated hash in bytes.
	var argonThreads = uint8(runtime.NumCPU())
	var argonMemory = uc.config.PasswordConfig.Memory
	var argonTime = uint32(uc.config.PasswordConfig.Time)

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, 32)

	// Encode the salt and hash as a single Base64 string for storage.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return b64Salt + "$" + b64Hash, nil
}

func toUserResponse(newUser *user.User) *exchange.UserResponse {
	return &exchange.UserResponse{
		UserId:    newUser.Id,
		LastName:  newUser.LastName,
		FirstName: newUser.FirstName,
		Anonymous: newUser.IsAnonymous(),
		CreatedAt: model.NewTime(newUser.Joined),
		UpdatedAt: model.NewTime(newUser.UpdatedAt),
		Email:     *custom.NewSecret(newUser.Email),
	}
}

func verifyPassword(threads uint8, memory uint32, time uint32, password, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 2 {
		return false, &xrfErr.Internal{Message: "Invalid password format"}
	}

	// Decode from Base64: decode the salt and hash from Base64 back to byte arrays.
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, &xrfErr.Internal{Message: "failed to decode salt", Err: err}
	}

	passHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, &xrfErr.Internal{Message: "failed to decode password", Err: err}
	}

	// Use the same parameters used for hashing:
	testHash := argon2.IDKey([]byte(password), salt, time, memory, threads, 32)

	// Use a constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(testHash, passHash) == 1, nil
}

func NewUserService(
	log internal.Logger,
	userSettings SettingsService,
	userRepo repository.UserRepository,
	ctx context.Context, config xrf.Security) UserService {

	return &service{
		ctx:             ctx,
		log:             log,
		config:          config,
		userRepo:        userRepo,
		settingsService: userSettings,
	}
}
