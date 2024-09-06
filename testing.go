package xrf197ilz35aq0

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func AssertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	// If err is not nil (meaning an error did occur)
	if err != nil {
		t.Fatalf("did not expect an error but got one, %v", err)
	}
}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("failed to generate random bytes: %v\n", err)

		for index := 0; index < size; index++ {
			b[index] = byte(index)
		}
	}
	return b
}
