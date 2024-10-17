package gocerr

import "testing"

func TestNewErrorField(t *testing.T) {
	field := "field1"
	message := "field is required"

	errField := NewErrorField(field, message)

	if errField.Field != field {
		t.Errorf("expected field is %s, but got %s", field, errField.Field)
	}

	if errField.Message != message {
		t.Errorf("expected message is %s, but got %s", message, errField.Message)
	}
}
