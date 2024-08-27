package model

import "testing"

func TestUser(t *testing.T) {
	t.Run("Create User", func(t *testing.T) {
		email := "email"
		lastName := "last"
		firstName := "first"
		password := "password"

		user := NewUser(firstName, lastName, email, password)
		if user == nil {
			t.Error("User creation failed")
		}
		if user.FirstName != firstName {
			t.Error("First name does not match")
		}
		if user.LastName != lastName {
			t.Error("Last name does not match")
		}

		if len(user.FingerPrint()) < 1 {
			t.Error("FingerPrint's size is less than 1")
		}

		if user.IsAnonymous() {
			t.Error("User should not be anonymous")
		}
	})
}
