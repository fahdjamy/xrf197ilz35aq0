package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"xrf197ilz35aq0/constants"
	"xrf197ilz35aq0/custom"
	"xrf197ilz35aq0/random"
)

const fingerPrintLength = 31

type User struct {
	masked      bool
	Id          int64
	FirstName   string
	LastName    string
	fingerPrint string
	joined      time.Time
	email       custom.Secret[*custom.SerializableString]
	password    custom.Secret[*custom.SerializableString]
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

func NewUser(firstName string, lastName string, email string, password string) *User {
	now := time.Now()

	id := random.PositiveInt64()
	serializableEmail := custom.SerializableString(email)
	serializablePassword := custom.SerializableString(password)

	var secretEmailPtr custom.Secret[*custom.SerializableString]
	var secretPasswordPtr custom.Secret[*custom.SerializableString]

	secretEmailPtr = *custom.NewSecret(&serializableEmail)
	secretPasswordPtr = *custom.NewSecret(&serializablePassword)

	newUser := &User{
		Id:        id,
		joined:    now,
		masked:    false,
		LastName:  lastName,
		FirstName: firstName,
		email:     secretEmailPtr,
		password:  secretPasswordPtr,
	}
	newUser.createFingerPrint()

	return newUser
}

func (u *User) createFingerPrint() {
	uniqueStr, err := random.TimeBasedString(u.joined.Unix(), fingerPrintLength)
	if err != nil {
		uniqueStr = ""
	}

	// remove any "-" in uniqueStr
	uniqueStr = strings.Join(strings.Split(uniqueStr, constants.DASH), constants.EMPTY)
	uniqueStr = strings.Join(strings.Split(uniqueStr, constants.EQUALS), constants.EMPTY)
	uniqueStr = strings.Join(strings.Split(uniqueStr, constants.UNDERSCORE), constants.EMPTY)

	// split string from the third part
	splitParts := splitAtIndex(uniqueStr, 3)

	lastPart := splitParts[1]
	firstPart := splitParts[0][2:] // remove the first 2 letters of the first part

	uniqueStr = fmt.Sprintf("%s%d%s", firstPart, u.Id, lastPart)
	u.fingerPrint = uniqueStr
}

func splitAtIndex(str string, stringPart int) []string {
	length := len(str)
	splitIndex := length - (length / stringPart)

	firstPart := str[:splitIndex]
	secondPart := str[splitIndex:]

	return []string{firstPart, secondPart}
}
