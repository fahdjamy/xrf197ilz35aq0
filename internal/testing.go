package internal

import (
	"context"
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
	prefix  string
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

func (t *TestLogger) SetPrefix(prefix string) {
	t.prefix = prefix
}

func (t *TestLogger) Message() string {
	return t.message
}

func (t *TestLogger) Called() int {
	return t.called
}

func NewTestLogger() *TestLogger {
	return &TestLogger{
		prefix:  "",
		message: "",
		called:  0,
	}
}

type StoreMock struct {
	Called     int
	Document   any
	Collection string
	context    context.Context
	calledFunc map[string]int
}

func (s *StoreMock) FindById(collection string, _ int64, ctx context.Context) (*Serializable, error) {
	methodName := "StoreMock.FindById"
	s.Collection = collection
	val, ok := s.calledFunc[methodName]
	if !ok {
		s.calledFunc[methodName] = 0
	}
	s.calledFunc[methodName] += val
	return nil, nil
}

func (s *StoreMock) Save(collection string, obj Serializable, ctx context.Context) (any, error) {
	s.Document = obj
	s.Collection = collection

	methodName := "StoreMock.Save"
	val, ok := s.calledFunc[methodName]
	if !ok {
		s.calledFunc[methodName] = 0
	}
	s.calledFunc[methodName] += val
	return nil, nil
}

func NewStoreMock() *StoreMock {
	return &StoreMock{
		Called:     0,
		Collection: "",
		Document:   nil,
		context:    context.TODO(),
		calledFunc: make(map[string]int),
	}
}
