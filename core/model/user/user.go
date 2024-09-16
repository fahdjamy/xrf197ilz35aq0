package user

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"xrf197ilz35aq0/internal/constants"
	"xrf197ilz35aq0/internal/random"
)

const fingerPrintLength = 55

type User struct {
	masked      bool
	Id          int64
	FirstName   string
	LastName    string
	fingerPrint string
	email       string
	password    string
	Joined      time.Time
	UpdatedAt   time.Time
}

func (u *User) String() string {
	format := "Id: %d, FirstName: %s, LastName: %s, Anonymous, %t"
	return fmt.Sprintf(format, u.Id, u.FirstName, u.LastName, u.masked)
}

func (u *User) IsAnonymous() bool {
	return u.masked
}

func (u *User) FingerPrint() string {
	// send last part of the fingerprint
	return strings.Split(u.fingerPrint, strconv.Itoa(int(u.Id)))[0]
}

func (u *User) UpdatePassword(password string) {
	u.password = password
}

func NewUser(firstName string, lastName string, email string, password string) *User {
	now := time.Now()

	newUser := &User{
		Joined:    now,
		UpdatedAt: now,
		masked:    false,
		email:     email,
		password:  password,
		LastName:  lastName,
		FirstName: firstName,
		Id:        random.PositiveInt64(),
	}
	newUser.createFingerPrint()

	return newUser
}

func (u *User) createFingerPrint() {
	uniqueStr, err := random.TimeBasedString(u.Joined.Unix(), fingerPrintLength)
	if err != nil {
		uniqueStr = ""
	}

	// remove any "-,=,_" in uniqueStr
	uniqueStr = strings.Join(strings.Split(uniqueStr, constants.DASH), constants.EMPTY)
	uniqueStr = strings.Join(strings.Split(uniqueStr, constants.EQUALS), constants.EMPTY)
	uniqueStr = strings.Join(strings.Split(uniqueStr, constants.UNDERSCORE), constants.EMPTY)

	// split string from the third part
	splitParts := splitAtIndex(uniqueStr, 3)

	lastPart := splitParts[1]
	firstPart := splitParts[0][2:] // remove the first 2 letters of the first part

	u.fingerPrint = fmt.Sprintf("%s%d%s", firstPart, u.Id, lastPart)
}

func splitAtIndex(str string, stringPart int) []string {
	length := len(str)
	splitIndex := length - (length / stringPart)

	firstPart := str[:splitIndex]
	secondPart := str[splitIndex:]

	return []string{firstPart, secondPart}
}
