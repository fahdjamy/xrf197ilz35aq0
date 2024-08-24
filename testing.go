package xrf197ilz35aq0

import "testing"

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
