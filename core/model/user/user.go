package user

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/random"
)

const fingerPrintLength = 55

type Alias User // Create an alias to avoid infinite recursion when marshalling/unMarshalling

type User struct {
	FingerPrint string             `json:"fingerPrint" bson:"fingerPrint"`
	Masked      bool               `json:"masked" bson:"masked"`
	Id          string             `json:"userId" bson:"userId"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	Email       string             `json:"email" bson:"email"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Password    string             `json:"password" bson:"password"`
	Joined      time.Time          `json:"joined" bson:"joined"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	MongoID     primitive.ObjectID `bson:"_id,omitempty" bson:"_id"` // MongoDB's ObjectID (internal)
	// json:"-" signifies that the JSON encoder should ignore this field even though field is exported
}

func (u *User) IsAnonymous() bool {
	return u.Masked
}

func (u *User) UpdatePassword(password string) {
	u.Password = password
}

func (u *User) UnmarshalJSON(bytes []byte) error {
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(bytes, &aux); err != nil {
		return &xrfErr.Internal{
			Err:     err,
			Message: "failed to unmarshal json",
			Source:  "core/model/user#UnmarshalJSON",
		}
	}
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	// dereference u to get the User value, convert it to a UserAlias, and store the result in the auxAlias variable.
	auxAlias := (Alias)(*u) // Store the converted UserAlias in a variable

	return json.Marshal(&struct {
		Id        string    `json:"id"`
		Email     string    `json:"email"`
		Masked    bool      `json:"masked"`
		Joined    time.Time `json:"joined"`
		LastName  string    `json:"lastName"`
		UpdatedAt time.Time `json:"updatedAt"`
		FirstName string    `json:"firstName"`
	}{
		Id:        auxAlias.Id,
		Email:     auxAlias.Email,
		Joined:    auxAlias.Joined,
		Masked:    auxAlias.Masked,
		LastName:  auxAlias.LastName,
		FirstName: auxAlias.FirstName,
		UpdatedAt: auxAlias.UpdatedAt,
	})
}

func (u *User) String() string {
	format := "Id: %s, FirstName: %s, LastName: %s, Anonymous, %t"
	return fmt.Sprintf(format, u.Id, u.FirstName, u.LastName, u.Masked)
}

func NewUser(firstName string, lastName string, email string, password string) *User {
	now := time.Now()

	newUser := &User{
		Masked:    false,
		Email:     email,
		Password:  password,
		LastName:  lastName,
		FirstName: firstName,
		Joined:    now,
		UpdatedAt: now,
	}
	newUser.generateUserId()
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

	u.FingerPrint = fmt.Sprintf("%s%s%s", firstPart, u.Id, lastPart)
}

func (u *User) generateUserId() {
	u.Id = strconv.FormatInt(random.PositiveInt64(), 10)
}

func splitAtIndex(str string, stringPart int) []string {
	length := len(str)
	splitIndex := length - (length / stringPart)

	firstPart := str[:splitIndex]
	secondPart := str[splitIndex:]

	return []string{firstPart, secondPart}
}
