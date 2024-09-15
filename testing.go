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

type TestLogger struct {
	message string
	called  int
}

func (t *TestLogger) Info(message string) {
	t.message = message
	t.called++
}

func (t *TestLogger) Warn(message string) {
	t.message = message
	t.called++
}

func (t *TestLogger) Debug(message string) {
	t.message = message
	t.called++
}

func (t *TestLogger) Error(message string) {
	t.message = message
	t.called++
}

func (t *TestLogger) Fatal(message string) {
	t.message = message
	t.called++
}

func (t *TestLogger) Panic(message string) {
	t.message = message
	t.called++
}

func (t *TestLogger) Message() string {
	return t.message
}

func (t *TestLogger) Called() int {
	return t.called
}
