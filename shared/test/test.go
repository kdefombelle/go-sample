package test

import (
	"testing"
	"time"
)

// CheckString checks expected and result string are equals.
// If not equals it sets the testing.T in error and prints an error message containing the field descriptor.
func CheckString(t *testing.T, descriptor string, expected string, result string) {
	if expected != result {
		t.Errorf("%s is wrong: [%s] instead of [%s]", descriptor, result, expected)
	}
}

// CheckBool checks expected and result boolean are equals.
// If not equals it sets the testing.T in error and prints an error message containing the field descriptor.
func CheckBool(t *testing.T, descriptor string, expected bool, result bool) {
	if expected != result {
		t.Errorf("%s is wrong: [%t] instead of [%t]", descriptor, result, expected)
	}
}

// CheckTime checks expected and result time.Time are equals.
// If not equals it sets the testing.T in error and prints an error message containing the field descriptor.
func CheckTime(t *testing.T, descriptor string, expected time.Time, result time.Time) {
	if expected != result {
		t.Errorf("%s is wrong: [%v] instead of [%v]", descriptor, result, expected)
	}
}

// GetDate returns a date in format YYYY-MM-DD.
func GetDate(t *testing.T, date string) time.Time {
	const shortForm = "2006-01-02"
	time, err := time.Parse(shortForm, date)
	if err != nil {
		t.Fatalf("Unexpected error: [%s]", err)
	}
	return time
}

// CheckFloat checks expected and result float64 are equals.
// As float are by nature approximated this method is known fragile.
func CheckFloat(t *testing.T, descriptor string, expected float64, result float64) {
	if expected != result {
		t.Errorf("%s is wrong: [%v] instead of [%v]", descriptor, result, expected)
	}
}

// CheckInt64 checks expected and result int64 are equals.
// If not equals it sets the testing.T in error and prints an error message containing the field descriptor.
func CheckInt64(t *testing.T, descriptor string, expected int64, result int64) {
	if expected != result {
		t.Errorf("%s is wrong: [%d] instead of [%d]", descriptor, result, expected)
	}
}
