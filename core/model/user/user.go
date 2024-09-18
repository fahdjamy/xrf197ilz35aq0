package user

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"xrf197ilz35aq0/core"
	"xrf197ilz35aq0/internal/constants"
	"xrf197ilz35aq0/internal/random"
)

const fingerPrintLength = 55

type User struct {
	fingerPrint string
	Masked      bool      `json:"masked"`
	Id          int64     `json:"id"`
	FirstName   string    `json:"firstName"`
	Email       string    `json:"email"`
	LastName    string    `json:"lastName"`
	Password    string    `json:"password"`
	Joined      time.Time `json:"joined"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (u *User) IsAnonymous() bool {
	return u.Masked
}

func (u *User) FingerPrint() string {
	// send last part of the fingerprint
	return strings.Split(u.fingerPrint, strconv.Itoa(int(u.Id)))[0]
}

func (u *User) UpdatePassword(password string) {
	u.Password = password
}

func (u *User) UnmarshalJSON(bytes []byte) error {
	type Alias User // Create an alias to avoid infinite recursion
	aux := &struct {
		*Alias
		Joined  string `json:"joined"`
		Updated string `json:"updatedAt"`
	}{
		Alias: (*Alias)(u),
	}
	now := time.Now()
	if err := json.Unmarshal(bytes, &aux); err != nil {
		return core.InternalError{
			Err:     err,
			Time:    now,
			Message: "failed to unmarshal json",
			Source:  "core/model/user#UnmarshalJSON",
		}
	}

	var err error
	// RFC3339 -> "YYYY-MM-DDTHH:mm:ssZ"
	u.Joined, err = time.Parse(time.RFC3339, aux.Joined)
	if err != nil {
		return core.InternalError{
			Err:     err,
			Time:    now,
			Message: "failed to parse Joined",
			Source:  "core/model/user#UnmarshalJSON",
		}
	}

	updatedAt, err := time.Parse(time.RFC3339, aux.Updated)
	if err != nil {
		return core.InternalError{
			Err:     err,
			Time:    now,
			Message: "failed to parse Updated",
			Source:  "core/model/user#UnmarshalJSON",
		}
	}
	u.UpdatedAt = updatedAt
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	// dereference u to get the User value, convert it to an Alias, and store the result in the auxAlias variable.
	auxAlias := (Alias)(*u) // Store the converted Alias in a variable

	return json.Marshal(&struct {
		*Alias
		Joined    string `json:"joined"`
		UpdatedAt string `json:"updatedAt"`
	}{
		// take the address of the auxAlias variable, which is a valid *Alias, and assign it to the Alias field in the anonymous struct
		Alias:     &auxAlias,
		Joined:    u.Joined.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	})
}

func (u *User) String() string {
	format := "Id: %d, FirstName: %s, LastName: %s, Anonymous, %t"
	return fmt.Sprintf(format, u.Id, u.FirstName, u.LastName, u.Masked)
}

func NewUser(firstName string, lastName string, email string, password string) *User {
	now := time.Now()

	newUser := &User{
		Joined:    now,
		UpdatedAt: now,
		Masked:    false,
		Email:     email,
		Password:  password,
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