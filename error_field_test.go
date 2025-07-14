package gocerr

import "testing"

// TestNewErrorField tests the NewErrorField function for creating ErrorField instances.
// It verifies that the function correctly creates ErrorField with the provided field name and message.
func TestNewErrorField(t *testing.T) {
	testCases := []struct {
		Name            string
		Field           string
		Message         string
		ExpectedField   string
		ExpectedMessage string
	}{
		{
			Name:            "standard field validation error",
			Field:           "field1",
			Message:         "field is required",
			ExpectedField:   "field1",
			ExpectedMessage: "field is required",
		},
		{
			Name:            "email validation error",
			Field:           "email",
			Message:         "invalid email format",
			ExpectedField:   "email",
			ExpectedMessage: "invalid email format",
		},
		{
			Name:            "empty field name",
			Field:           "",
			Message:         "field is required",
			ExpectedField:   "",
			ExpectedMessage: "field is required",
		},
		{
			Name:            "empty message",
			Field:           "username",
			Message:         "",
			ExpectedField:   "username",
			ExpectedMessage: "",
		},
		{
			Name:            "both field and message empty",
			Field:           "",
			Message:         "",
			ExpectedField:   "",
			ExpectedMessage: "",
		},
		{
			Name:            "field with special characters",
			Field:           "user.profile.name",
			Message:         "nested field validation failed",
			ExpectedField:   "user.profile.name",
			ExpectedMessage: "nested field validation failed",
		},
		{
			Name:            "long validation message",
			Field:           "password",
			Message:         "password must be at least 8 characters long, contain uppercase, lowercase, numbers and special characters",
			ExpectedField:   "password",
			ExpectedMessage: "password must be at least 8 characters long, contain uppercase, lowercase, numbers and special characters",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			errField := NewErrorField(testCase.Field, testCase.Message)

			// Verify field name
			if errField.Field != testCase.ExpectedField {
				t.Errorf("expected field is %s, but got %s", testCase.ExpectedField, errField.Field)
			}

			// Verify message
			if errField.Message != testCase.ExpectedMessage {
				t.Errorf("expected message is %s, but got %s", testCase.ExpectedMessage, errField.Message)
			}
		})
	}
}

// TestNewErrorField_EdgeCases tests edge cases for the NewErrorField function.
// It ensures the function handles various input scenarios correctly.
func TestNewErrorField_EdgeCases(t *testing.T) {
	testCases := []struct {
		Name        string
		Field       string
		Message     string
		Description string
	}{
		{
			Name:        "unicode field name",
			Field:       "пользователь",
			Message:     "пользователь обязателен",
			Description: "should handle unicode characters in both field and message",
		},
		{
			Name:        "field with json path notation",
			Field:       "user.profile.settings.theme",
			Message:     "invalid theme selection",
			Description: "should handle nested field notation",
		},
		{
			Name:        "field with array notation",
			Field:       "items[0].quantity",
			Message:     "quantity must be positive",
			Description: "should handle array field notation",
		},
		{
			Name:        "whitespace handling",
			Field:       "  field_with_spaces  ",
			Message:     "  message with spaces  ",
			Description: "should preserve whitespace as provided",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			errorField := NewErrorField(testCase.Field, testCase.Message)

			if errorField.Field != testCase.Field {
				t.Errorf("expected field '%s', got '%s'", testCase.Field, errorField.Field)
			}

			if errorField.Message != testCase.Message {
				t.Errorf("expected message '%s', got '%s'", testCase.Message, errorField.Message)
			}
		})
	}
}

// BenchmarkNewErrorField benchmarks the NewErrorField function.
// This tests the performance of error field creation.
func BenchmarkNewErrorField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewErrorField("field1", "field is required")
	}
}
