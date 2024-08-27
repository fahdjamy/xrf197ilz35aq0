package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"xrf197ilz35aq0/constants"
	xrf "xrf197ilz35aq0/custom"
	"xrf197ilz35aq0/random"
)

const fingerPrintLength = 31

type User struct {
	anonymous   bool
	id          int64
	FirstName   string
	LastName    string
	fingerPrint string
	createdAt   time.Time
	email       Secret[*xrf.SerializableString]
	password    Secret[*xrf.SerializableString]
}

func (u *User) String() string {
	format := "id: %d, FirstName: %s, LastName: %s, Anonymous, %t"
	return fmt.Sprintf(format, u.id, u.FirstName, u.LastName, u.anonymous)
}

func (u *User) IsAnonymous() bool {
	return u.anonymous
}

func (u *User) FingerPrint() string {
	// send last part of the fingerprint
	return strings.Split(u.fingerPrint, strconv.Itoa(int(u.id)))[0]
}

func NewUser(firstName string, lastName string, email string, password string) *User {
	now := time.Now()

	id := random.PositiveInt64()
	serializableEmail := xrf.SerializableString(email)
	serializablePassword := xrf.SerializableString(password)

	var secretEmailPtr Secret[*xrf.SerializableString]
	var secretPasswordPtr Secret[*xrf.SerializableString]

	secretEmailPtr = *NewSecret(&serializableEmail)
	secretPasswordPtr = *NewSecret(&serializablePassword)

	newUser := &User{
		id:        id,
		createdAt: now,
		anonymous: false,
		LastName:  lastName,
		FirstName: firstName,
		email:     secretEmailPtr,
		password:  secretPasswordPtr,
	}
	newUser.createFingerPrint()

	return newUser
}

func (u *User) createFingerPrint() {
	uniqueStr, err := random.TimeBasedString(u.createdAt.Unix(), fingerPrintLength)
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

	uniqueStr = fmt.Sprintf("%s%d%s", firstPart, u.id, lastPart)
	u.fingerPrint = uniqueStr
}

func splitAtIndex(str string, stringPart int) []string {
	length := len(str)
	splitIndex := length - (length / stringPart)

	firstPart := str[:splitIndex]
	secondPart := str[splitIndex:]

	return []string{firstPart, secondPart}
}
