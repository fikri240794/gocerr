package gocerr

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

// TestNew tests the New function for creating custom errors with various configurations.
// It verifies that the function correctly handles both scenarios with and without error fields.
func TestNew(t *testing.T) {
	testCases := []struct {
		Name        string
		Code        int
		Message     string
		ErrorFields []ErrorField // Changed from single ErrorField pointer to slice for better test coverage
		Expected    Error
	}{
		{
			Name:    "with single error field",
			Code:    400,
			Message: "bad request",
			ErrorFields: []ErrorField{
				{
					Field:   "field1",
					Message: "field is required",
				},
			},
			Expected: Error{
				Code:    400,
				Message: "bad request",
				ErrorFields: []ErrorField{
					{
						Field:   "field1",
						Message: "field is required",
					},
				},
			},
		},
		{
			Name:        "no error fields",
			Code:        500,
			Message:     "internal server error",
			ErrorFields: nil,
			Expected: Error{
				Code:        500,
				Message:     "internal server error",
				ErrorFields: nil, // Updated to match the optimized New function behavior
			},
		},
		{
			Name:    "with multiple error fields",
			Code:    422,
			Message: "validation failed",
			ErrorFields: []ErrorField{
				{
					Field:   "username",
					Message: "username is required",
				},
				{
					Field:   "email",
					Message: "invalid email format",
				},
				{
					Field:   "age",
					Message: "age must be greater than 0",
				},
			},
			Expected: Error{
				Code:    422,
				Message: "validation failed",
				ErrorFields: []ErrorField{
					{
						Field:   "username",
						Message: "username is required",
					},
					{
						Field:   "email",
						Message: "invalid email format",
					},
					{
						Field:   "age",
						Message: "age must be greater than 0",
					},
				},
			},
		},
		{
			Name:        "with empty error fields slice",
			Code:        400,
			Message:     "bad request",
			ErrorFields: []ErrorField{},
			Expected: Error{
				Code:        400,
				Message:     "bad request",
				ErrorFields: []ErrorField{},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actualErr Error
			if testCases[i].ErrorFields == nil {
				actualErr = New(testCases[i].Code, testCases[i].Message)
			} else {
				actualErr = New(testCases[i].Code, testCases[i].Message, testCases[i].ErrorFields...)
			}

			// Verify error code
			if testCases[i].Expected.Code != actualErr.Code {
				t.Errorf("expected code is %d, but got %d", testCases[i].Expected.Code, actualErr.Code)
			}

			// Verify error message
			if testCases[i].Expected.Message != actualErr.Message {
				t.Errorf("expected message is %s, but got %s", testCases[i].Expected.Message, actualErr.Message)
			}

			// Verify error fields length
			expectedFieldsLen := len(testCases[i].Expected.ErrorFields)
			actualFieldsLen := len(actualErr.ErrorFields)
			if expectedFieldsLen != actualFieldsLen {
				t.Errorf("expected length of error fields is %d, but got %d", expectedFieldsLen, actualFieldsLen)
			}

			// Verify each error field
			for j := 0; j < len(testCases[i].Expected.ErrorFields); j++ {
				if testCases[i].Expected.ErrorFields[j].Field != actualErr.ErrorFields[j].Field {
					t.Errorf("expected field of sub item error fields is %s, but got %s", testCases[i].Expected.ErrorFields[j].Field, actualErr.ErrorFields[j].Field)
				}
				if testCases[i].Expected.ErrorFields[j].Message != actualErr.ErrorFields[j].Message {
					t.Errorf("expected message of sub item error fields is %s, but got %s", testCases[i].Expected.ErrorFields[j].Message, actualErr.ErrorFields[j].Message)
				}
			}
		})
	}
}

// TestError_Error tests the Error method which implements the built-in error interface.
// It ensures that the Error method correctly returns the message field.
func TestError_Error(t *testing.T) {
	testCases := []struct {
		Name            string
		Error           Error
		ExpectedMessage string
	}{
		{
			Name: "simple error message",
			Error: Error{
				Code:    500,
				Message: "internal server error",
			},
			ExpectedMessage: "internal server error",
		},
		{
			Name: "error with empty message",
			Error: Error{
				Code:    400,
				Message: "",
			},
			ExpectedMessage: "",
		},
		{
			Name: "error with special characters in message",
			Error: Error{
				Code:    422,
				Message: "validation failed: field 'email' contains invalid characters @#$%",
			},
			ExpectedMessage: "validation failed: field 'email' contains invalid characters @#$%",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actualMessage := testCase.Error.Error()
			if actualMessage != testCase.ExpectedMessage {
				t.Errorf("expected error string return %s, but got %s", testCase.ExpectedMessage, actualMessage)
			}
		})
	}
}

// TestParse tests the Parse function for converting standard Go errors to custom errors.
// It covers all possible scenarios including nil errors, non-custom errors, and custom errors.
func TestParse(t *testing.T) {
	testCases := []struct {
		Name     string
		Error    error
		Expected struct {
			CustomError   Error
			IsCustomError bool
		}
	}{
		{
			Name:  "error is nil",
			Error: nil,
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   Error{},
				IsCustomError: false,
			},
		},
		{
			Name:  "error is standard error (not custom error)",
			Error: errors.New("some standard error"),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   Error{},
				IsCustomError: false,
			},
		},
		{
			Name:  "error is custom error without error fields",
			Error: New(500, "internal server error"),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   New(500, "internal server error"),
				IsCustomError: true,
			},
		},
		{
			Name:  "error is custom error with single error field",
			Error: New(400, "bad request", NewErrorField("field1", "field is required")),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   New(400, "bad request", NewErrorField("field1", "field is required")),
				IsCustomError: true,
			},
		},
		{
			Name: "error is custom error with multiple error fields",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format")),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError: New(422, "validation failed",
					NewErrorField("username", "username is required"),
					NewErrorField("email", "invalid email format")),
				IsCustomError: true,
			},
		},
		{
			Name:  "error is custom error with zero code",
			Error: New(0, "unknown error"),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   New(0, "unknown error"),
				IsCustomError: true,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			actualCustomError, actualIsCustomError := Parse(testCases[i].Error)

			// Verify parsing result
			if testCases[i].Expected.IsCustomError != actualIsCustomError {
				t.Errorf("expected is custom error is %t, but got %t", testCases[i].Expected.IsCustomError, actualIsCustomError)
			}

			// Verify error code
			if testCases[i].Expected.CustomError.Code != actualCustomError.Code {
				t.Errorf("expected custom error code is %d, but got %d", testCases[i].Expected.CustomError.Code, actualCustomError.Code)
			}

			// Verify error message
			if testCases[i].Expected.CustomError.Message != actualCustomError.Message {
				t.Errorf("expected custom error message is %s, but got %s", testCases[i].Expected.CustomError.Message, actualCustomError.Message)
			}

			// Verify that Error() method returns the same message for custom errors
			if testCases[i].Error != nil && testCases[i].Expected.IsCustomError && testCases[i].Error.Error() != actualCustomError.Message {
				t.Errorf("expected error message is %s, but got %s", testCases[i].Error.Error(), actualCustomError.Message)
			}

			// Verify error fields length
			expectedFieldsLen := len(testCases[i].Expected.CustomError.ErrorFields)
			actualFieldsLen := len(actualCustomError.ErrorFields)
			if expectedFieldsLen != actualFieldsLen {
				t.Errorf("expected length of custom error error fields is %d, but got %d", expectedFieldsLen, actualFieldsLen)
			}

			// Verify each error field content
			for j := 0; j < len(testCases[i].Expected.CustomError.ErrorFields); j++ {
				expectedField := testCases[i].Expected.CustomError.ErrorFields[j]
				actualField := actualCustomError.ErrorFields[j]

				if expectedField.Field != actualField.Field {
					t.Errorf("expected field of error fields sub item custom error is %s, but got %s", expectedField.Field, actualField.Field)
				}
				if expectedField.Message != actualField.Message {
					t.Errorf("expected message of error fields sub item custom error is %s, but got %s", expectedField.Message, actualField.Message)
				}
			}
		})
	}
}

// TestGetErrorCode tests the GetErrorCode function for extracting error codes from various error types.
// It ensures the function correctly handles nil errors, standard errors, and custom errors.
func TestGetErrorCode(t *testing.T) {
	testCases := []struct {
		Name        string
		Error       error
		Expectation int
	}{
		{
			Name:        "error is nil",
			Error:       nil,
			Expectation: 0,
		},
		{
			Name:        "error is standard error (not custom error)",
			Error:       errors.New("standard error"),
			Expectation: 0,
		},
		{
			Name: "error is custom error with HTTP status internal server error",
			Error: Error{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
			Expectation: http.StatusInternalServerError,
		},
		{
			Name: "error is custom error with HTTP status bad request",
			Error: Error{
				Code:    http.StatusBadRequest,
				Message: "bad request",
			},
			Expectation: http.StatusBadRequest,
		},
		{
			Name: "error is custom error with zero code",
			Error: Error{
				Code:    0,
				Message: "unknown error",
			},
			Expectation: 0,
		},
		{
			Name: "error is custom error with negative code",
			Error: Error{
				Code:    -1,
				Message: "custom negative error",
			},
			Expectation: -1,
		},
		{
			Name: "error is custom error created with New function",
			Error: New(422, "validation failed",
				NewErrorField("field1", "field is required")),
			Expectation: 422,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			actual := GetErrorCode(testCases[i].Error)

			if testCases[i].Expectation != actual {
				t.Errorf("expectation is %d, got %d", testCases[i].Expectation, actual)
			}
		})
	}
}

// TestIsErrorCodeEqual tests the IsErrorCodeEqual function for comparing error codes.
// It verifies the function correctly identifies matching error codes across different error types.
func TestIsErrorCodeEqual(t *testing.T) {
	testCases := []struct {
		Name        string
		Code        int
		Error       error
		Expectation bool
	}{
		{
			Name:        "error is nil, should return false",
			Code:        http.StatusInternalServerError,
			Error:       nil,
			Expectation: false,
		},
		{
			Name:        "error is standard error (not custom error), should return false",
			Code:        http.StatusBadRequest,
			Error:       errors.New("standard error"),
			Expectation: false,
		},
		{
			Name: "error code matches - internal server error",
			Code: http.StatusInternalServerError,
			Error: Error{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
			Expectation: true,
		},
		{
			Name: "error code does not match",
			Code: http.StatusBadRequest,
			Error: Error{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
			Expectation: false,
		},
		{
			Name: "error code matches - zero code",
			Code: 0,
			Error: Error{
				Code:    0,
				Message: "unknown error",
			},
			Expectation: true,
		},
		{
			Name: "error code matches - negative code",
			Code: -1,
			Error: Error{
				Code:    -1,
				Message: "custom negative error",
			},
			Expectation: true,
		},
		{
			Name: "error code matches - custom error created with New function",
			Code: 422,
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required")),
			Expectation: true,
		},
		{
			Name: "error code does not match - custom error created with New function",
			Code: 400,
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required")),
			Expectation: false,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			actual := IsErrorCodeEqual(testCases[i].Error, testCases[i].Code)

			if testCases[i].Expectation != actual {
				t.Errorf("expectation is %t, got %t", testCases[i].Expectation, actual)
			}
		})
	}
}

// TestHasErrorFields tests the HasErrorFields function for detecting error fields presence.
// It verifies the function correctly identifies when custom errors contain field validation errors.
func TestHasErrorFields(t *testing.T) {
	testCases := []struct {
		Name        string
		Error       error
		Expectation bool
	}{
		{
			Name:        "error is nil",
			Error:       nil,
			Expectation: false,
		},
		{
			Name:        "error is standard error (not custom error)",
			Error:       errors.New("standard error"),
			Expectation: false,
		},
		{
			Name:        "custom error without error fields",
			Error:       New(500, "internal server error"),
			Expectation: false,
		},
		{
			Name:        "custom error with single error field",
			Error:       New(400, "bad request", NewErrorField("field1", "field is required")),
			Expectation: true,
		},
		{
			Name: "custom error with multiple error fields",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format")),
			Expectation: true,
		},
		{
			Name:        "custom error created directly without New function",
			Error:       Error{Code: 400, Message: "bad request", ErrorFields: nil},
			Expectation: false,
		},
		{
			Name: "custom error created directly with error fields",
			Error: Error{
				Code:    400,
				Message: "bad request",
				ErrorFields: []ErrorField{
					{Field: "test", Message: "test error"},
				},
			},
			Expectation: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := HasErrorFields(testCase.Error)

			if testCase.Expectation != actual {
				t.Errorf("expected %t, got %t", testCase.Expectation, actual)
			}
		})
	}
}

// TestGetErrorFields tests the GetErrorFields function for extracting error fields.
// It verifies the function correctly extracts and returns copies of error fields.
func TestGetErrorFields(t *testing.T) {
	testCases := []struct {
		Name     string
		Error    error
		Expected []ErrorField
	}{
		{
			Name:     "error is nil",
			Error:    nil,
			Expected: nil,
		},
		{
			Name:     "error is standard error (not custom error)",
			Error:    errors.New("standard error"),
			Expected: nil,
		},
		{
			Name:     "custom error without error fields",
			Error:    New(500, "internal server error"),
			Expected: nil,
		},
		{
			Name:  "custom error with single error field",
			Error: New(400, "bad request", NewErrorField("field1", "field is required")),
			Expected: []ErrorField{
				{Field: "field1", Message: "field is required"},
			},
		},
		{
			Name: "custom error with multiple error fields",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format"),
				NewErrorField("age", "age must be positive")),
			Expected: []ErrorField{
				{Field: "username", Message: "username is required"},
				{Field: "email", Message: "invalid email format"},
				{Field: "age", Message: "age must be positive"},
			},
		},
		{
			Name: "custom error created directly with empty error fields",
			Error: Error{
				Code:        400,
				Message:     "bad request",
				ErrorFields: []ErrorField{},
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := GetErrorFields(testCase.Error)

			// Check length
			if len(testCase.Expected) != len(actual) {
				t.Errorf("expected length %d, got %d", len(testCase.Expected), len(actual))
				return
			}

			// Check each field
			for i, expectedField := range testCase.Expected {
				if expectedField.Field != actual[i].Field {
					t.Errorf("expected field[%d].Field %s, got %s", i, expectedField.Field, actual[i].Field)
				}
				if expectedField.Message != actual[i].Message {
					t.Errorf("expected field[%d].Message %s, got %s", i, expectedField.Message, actual[i].Message)
				}
			}

			// Test that returned slice is a copy (modification doesn't affect original)
			if len(actual) > 0 {
				originalField := actual[0].Field
				actual[0].Field = "modified"

				// Get fields again and verify they weren't modified
				newFields := GetErrorFields(testCase.Error)
				if len(newFields) > 0 && newFields[0].Field != originalField {
					t.Error("GetErrorFields should return a copy, but original was modified")
				}
			}
		})
	}
}

// TestHasErrorField tests the HasErrorField function for detecting specific error fields.
// It verifies the function correctly identifies when specific field validation errors exist.
func TestHasErrorField(t *testing.T) {
	testCases := []struct {
		Name        string
		Error       error
		FieldName   string
		Expectation bool
	}{
		{
			Name:        "error is nil",
			Error:       nil,
			FieldName:   "username",
			Expectation: false,
		},
		{
			Name:        "error is standard error (not custom error)",
			Error:       errors.New("standard error"),
			FieldName:   "username",
			Expectation: false,
		},
		{
			Name:        "custom error without error fields",
			Error:       New(500, "internal server error"),
			FieldName:   "username",
			Expectation: false,
		},
		{
			Name:        "custom error with error field - field exists",
			Error:       New(400, "bad request", NewErrorField("username", "username is required")),
			FieldName:   "username",
			Expectation: true,
		},
		{
			Name:        "custom error with error field - field does not exist",
			Error:       New(400, "bad request", NewErrorField("username", "username is required")),
			FieldName:   "email",
			Expectation: false,
		},
		{
			Name: "custom error with multiple error fields - existing field",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format"),
				NewErrorField("age", "age must be positive")),
			FieldName:   "email",
			Expectation: true,
		},
		{
			Name: "custom error with multiple error fields - non-existing field",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format")),
			FieldName:   "password",
			Expectation: false,
		},
		{
			Name:        "field name case sensitivity",
			Error:       New(400, "bad request", NewErrorField("Username", "username is required")),
			FieldName:   "username",
			Expectation: false,
		},
		{
			Name:        "empty field name search",
			Error:       New(400, "bad request", NewErrorField("", "field is required")),
			FieldName:   "",
			Expectation: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := HasErrorField(testCase.Error, testCase.FieldName)

			if testCase.Expectation != actual {
				t.Errorf("expected %t, got %t", testCase.Expectation, actual)
			}
		})
	}
}

// TestErrorField_EdgeCases tests edge cases for ErrorField functionality.
// It ensures robust handling of various input scenarios.
func TestErrorField_EdgeCases(t *testing.T) {
	testCases := []struct {
		Name        string
		Field       string
		Message     string
		Description string
	}{
		{
			Name:        "unicode characters in field name",
			Field:       "用户名",
			Message:     "用户名是必需的",
			Description: "should handle unicode characters correctly",
		},
		{
			Name:        "field name with spaces",
			Field:       "first name",
			Message:     "first name is required",
			Description: "should handle field names with spaces",
		},
		{
			Name:        "very long field name",
			Field:       "this_is_a_very_long_field_name_that_might_be_used_in_some_complex_validation_scenarios_with_nested_structures",
			Message:     "field validation failed",
			Description: "should handle very long field names",
		},
		{
			Name:        "special characters in message",
			Field:       "email",
			Message:     "Email must contain @ symbol and valid domain (e.g., user@example.com)",
			Description: "should handle special characters in validation messages",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			errorField := NewErrorField(testCase.Field, testCase.Message)

			if errorField.Field != testCase.Field {
				t.Errorf("expected field %s, got %s", testCase.Field, errorField.Field)
			}

			if errorField.Message != testCase.Message {
				t.Errorf("expected message %s, got %s", testCase.Message, errorField.Message)
			}
		})
	}
}

// TestError_IntegrationScenarios tests real-world integration scenarios.
// These tests simulate how the library would be used in practice.
func TestError_IntegrationScenarios(t *testing.T) {
	t.Run("user registration validation scenario", func(t *testing.T) {
		// Simulate a user registration validation error
		registrationErr := New(422, "User registration failed",
			NewErrorField("username", "Username must be at least 3 characters long"),
			NewErrorField("email", "Email address is already taken"),
			NewErrorField("password", "Password must contain at least one uppercase letter"))

		// Test that we can identify this as a validation error
		if !IsErrorCodeEqual(registrationErr, 422) {
			t.Error("should identify as validation error (422)")
		}

		// Test that we can detect field validation errors
		if !HasErrorFields(registrationErr) {
			t.Error("should have error fields")
		}

		// Test that we can check for specific field errors
		if !HasErrorField(registrationErr, "email") {
			t.Error("should have email field error")
		}

		if HasErrorField(registrationErr, "phone") {
			t.Error("should not have phone field error")
		}

		// Test error message extraction
		if registrationErr.Error() != "User registration failed" {
			t.Errorf("expected error message 'User registration failed', got %s", registrationErr.Error())
		}

		// Test field extraction
		fields := GetErrorFields(registrationErr)
		if len(fields) != 3 {
			t.Errorf("expected 3 error fields, got %d", len(fields))
		}
	})

	t.Run("API error handling scenario", func(t *testing.T) {
		// Simulate different types of API errors
		apiErrors := []error{
			New(404, "Resource not found"),
			New(401, "Unauthorized access"),
			New(500, "Internal server error"),
			errors.New("network timeout"), // Standard error
		}

		expectedCodes := []int{404, 401, 500, 0}

		for i, err := range apiErrors {
			code := GetErrorCode(err)
			if code != expectedCodes[i] {
				t.Errorf("error %d: expected code %d, got %d", i, expectedCodes[i], code)
			}

			// Test parsing
			if customErr, ok := Parse(err); ok {
				if customErr.Code != expectedCodes[i] {
					t.Errorf("parsed error %d: expected code %d, got %d", i, expectedCodes[i], customErr.Code)
				}
			} else if expectedCodes[i] != 0 {
				t.Errorf("error %d: should be parseable as custom error", i)
			}
		}
	})

	t.Run("error chaining and wrapping scenario", func(t *testing.T) {
		// Test that our custom errors work well with error interfaces
		customErr := New(400, "validation failed", NewErrorField("name", "name is required"))

		// Test error interface compliance
		var err error = customErr
		if err.Error() != "validation failed" {
			t.Errorf("error interface: expected 'validation failed', got %s", err.Error())
		}

		// Test that we can still parse it back
		if parsed, ok := Parse(err); ok {
			if parsed.Code != 400 {
				t.Errorf("parsed back error: expected code 400, got %d", parsed.Code)
			}
		} else {
			t.Error("should be able to parse back custom error from error interface")
		}
	})
}

// TestError_String tests the String method for detailed error representation.
// It verifies the method provides comprehensive debugging information.
func TestError_String(t *testing.T) {
	testCases := []struct {
		Name     string
		Error    Error
		Expected string
	}{
		{
			Name: "error without error fields",
			Error: Error{
				Code:    404,
				Message: "not found",
			},
			Expected: `Error{Code: 404, Message: "not found"}`,
		},
		{
			Name: "error with single error field",
			Error: Error{
				Code:    400,
				Message: "validation failed",
				ErrorFields: []ErrorField{
					{Field: "username", Message: "username is required"},
				},
			},
			Expected: `Error{Code: 400, Message: "validation failed", ErrorFields: [{Field: "username", Message: "username is required"}]}`,
		},
		{
			Name: "error with multiple error fields",
			Error: Error{
				Code:    422,
				Message: "multiple validation errors",
				ErrorFields: []ErrorField{
					{Field: "email", Message: "invalid email format"},
					{Field: "age", Message: "age must be positive"},
				},
			},
			Expected: `Error{Code: 422, Message: "multiple validation errors", ErrorFields: [{Field: "email", Message: "invalid email format"}, {Field: "age", Message: "age must be positive"}]}`,
		},
		{
			Name: "error with empty message",
			Error: Error{
				Code:    500,
				Message: "",
			},
			Expected: `Error{Code: 500, Message: ""}`,
		},
		{
			Name: "error with special characters in message",
			Error: Error{
				Code:    400,
				Message: `message with "quotes" and \backslashes`,
			},
			Expected: `Error{Code: 400, Message: "message with \"quotes\" and \\backslashes"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := testCase.Error.String()
			if actual != testCase.Expected {
				t.Errorf("expected:\n%s\ngot:\n%s", testCase.Expected, actual)
			}
		})
	}
}

// TestError_IsEmpty tests the IsEmpty method for detecting zero-value errors.
// It verifies the method correctly identifies empty Error instances.
func TestError_IsEmpty(t *testing.T) {
	testCases := []struct {
		Name     string
		Error    Error
		Expected bool
	}{
		{
			Name:     "completely empty error",
			Error:    Error{},
			Expected: true,
		},
		{
			Name: "error with only code",
			Error: Error{
				Code: 400,
			},
			Expected: false,
		},
		{
			Name: "error with only message",
			Error: Error{
				Message: "some error",
			},
			Expected: false,
		},
		{
			Name: "error with only error fields",
			Error: Error{
				ErrorFields: []ErrorField{
					{Field: "field1", Message: "error"},
				},
			},
			Expected: false,
		},
		{
			Name: "error with code and message",
			Error: Error{
				Code:    404,
				Message: "not found",
			},
			Expected: false,
		},
		{
			Name: "error with all fields populated",
			Error: Error{
				Code:    400,
				Message: "validation failed",
				ErrorFields: []ErrorField{
					{Field: "username", Message: "required"},
				},
			},
			Expected: false,
		},
		{
			Name: "error with zero code but message",
			Error: Error{
				Code:    0,
				Message: "unknown error",
			},
			Expected: false,
		},
		{
			Name: "error with empty slice (not nil)",
			Error: Error{
				Code:        0,
				Message:     "",
				ErrorFields: []ErrorField{},
			},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := testCase.Error.IsEmpty()
			if actual != testCase.Expected {
				t.Errorf("expected %t, got %t", testCase.Expected, actual)
			}
		})
	}
}

// TestStringAndIsEmptyIntegration tests String and IsEmpty methods together.
// This ensures they work correctly with Parse results and other operations.
func TestStringAndIsEmptyIntegration(t *testing.T) {
	t.Run("parse empty error should be empty", func(t *testing.T) {
		// Test that parsing a non-custom error returns empty Error
		standardErr := errors.New("standard error")
		customErr, ok := Parse(standardErr)

		if ok {
			t.Error("should not parse standard error as custom error")
		}

		if !customErr.IsEmpty() {
			t.Error("parsed result should be empty for non-custom error")
		}

		// String representation of empty error should be meaningful
		str := customErr.String()
		expected := `Error{Code: 0, Message: ""}`
		if str != expected {
			t.Errorf("expected empty error string %s, got %s", expected, str)
		}
	})

	t.Run("parse custom error should not be empty", func(t *testing.T) {
		originalErr := New(400, "validation failed", NewErrorField("field1", "required"))
		customErr, ok := Parse(originalErr)

		if !ok {
			t.Error("should parse custom error successfully")
		}

		if customErr.IsEmpty() {
			t.Error("parsed custom error should not be empty")
		}

		// String should contain all information
		str := customErr.String()
		if !strings.Contains(str, "400") {
			t.Error("string should contain error code")
		}
		if !strings.Contains(str, "validation failed") {
			t.Error("string should contain error message")
		}
		if !strings.Contains(str, "field1") {
			t.Error("string should contain error field information")
		}
	})
}

// TestGetErrorFieldMessage tests the GetErrorFieldMessage function for retrieving specific field messages.
// It verifies the function correctly extracts messages for specific fields.
func TestGetErrorFieldMessage(t *testing.T) {
	testCases := []struct {
		Name      string
		Error     error
		FieldName string
		Expected  string
	}{
		{
			Name:      "error is nil",
			Error:     nil,
			FieldName: "username",
			Expected:  "",
		},
		{
			Name:      "error is standard error (not custom error)",
			Error:     errors.New("standard error"),
			FieldName: "username",
			Expected:  "",
		},
		{
			Name:      "custom error without error fields",
			Error:     New(500, "internal server error"),
			FieldName: "username",
			Expected:  "",
		},
		{
			Name:      "custom error with error field - field exists",
			Error:     New(400, "bad request", NewErrorField("username", "username is required")),
			FieldName: "username",
			Expected:  "username is required",
		},
		{
			Name:      "custom error with error field - field does not exist",
			Error:     New(400, "bad request", NewErrorField("username", "username is required")),
			FieldName: "email",
			Expected:  "",
		},
		{
			Name: "custom error with multiple error fields - existing field",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format"),
				NewErrorField("age", "age must be positive")),
			FieldName: "email",
			Expected:  "invalid email format",
		},
		{
			Name: "custom error with multiple error fields - first field",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format")),
			FieldName: "username",
			Expected:  "username is required",
		},
		{
			Name: "custom error with multiple error fields - last field",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format")),
			FieldName: "email",
			Expected:  "invalid email format",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := GetErrorFieldMessage(testCase.Error, testCase.FieldName)

			if testCase.Expected != actual {
				t.Errorf("expected %q, got %q", testCase.Expected, actual)
			}
		})
	}
}

// TestErrorFieldCount tests the ErrorFieldCount function for counting error fields.
// It verifies the function correctly counts the number of field validation errors.
func TestErrorFieldCount(t *testing.T) {
	testCases := []struct {
		Name     string
		Error    error
		Expected int
	}{
		{
			Name:     "error is nil",
			Error:    nil,
			Expected: 0,
		},
		{
			Name:     "error is standard error (not custom error)",
			Error:    errors.New("standard error"),
			Expected: 0,
		},
		{
			Name:     "custom error without error fields",
			Error:    New(500, "internal server error"),
			Expected: 0,
		},
		{
			Name:     "custom error with single error field",
			Error:    New(400, "bad request", NewErrorField("field1", "field is required")),
			Expected: 1,
		},
		{
			Name: "custom error with multiple error fields",
			Error: New(422, "validation failed",
				NewErrorField("username", "username is required"),
				NewErrorField("email", "invalid email format"),
				NewErrorField("age", "age must be positive")),
			Expected: 3,
		},
		{
			Name: "custom error with five error fields",
			Error: New(422, "validation failed",
				NewErrorField("field1", "error1"),
				NewErrorField("field2", "error2"),
				NewErrorField("field3", "error3"),
				NewErrorField("field4", "error4"),
				NewErrorField("field5", "error5")),
			Expected: 5,
		},
		{
			Name: "custom error created directly with empty error fields",
			Error: Error{
				Code:        400,
				Message:     "bad request",
				ErrorFields: []ErrorField{},
			},
			Expected: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := ErrorFieldCount(testCase.Error)

			if testCase.Expected != actual {
				t.Errorf("expected %d, got %d", testCase.Expected, actual)
			}
		})
	}
}

// Benchmark tests to verify performance optimizations

// BenchmarkNew_NoErrorFields benchmarks the New function without error fields.
// This tests the optimization where we avoid slice allocation when no fields are provided.
func BenchmarkNew_NoErrorFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(500, "internal server error")
	}
}

// BenchmarkNew_WithErrorFields benchmarks the New function with error fields.
// This tests the optimized allocation strategy for error fields.
func BenchmarkNew_WithErrorFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(400, "validation failed",
			NewErrorField("username", "username is required"),
			NewErrorField("email", "invalid email format"))
	}
}

// BenchmarkParse_CustomError benchmarks the Parse function with custom errors.
// This tests the type assertion performance optimization.
func BenchmarkParse_CustomError(b *testing.B) {
	err := New(400, "bad request", NewErrorField("field1", "field is required"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(err)
	}
}

// BenchmarkParse_StandardError benchmarks the Parse function with standard errors.
// This tests the early return optimization for non-custom errors.
func BenchmarkParse_StandardError(b *testing.B) {
	err := errors.New("standard error")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(err)
	}
}

// BenchmarkGetErrorCode benchmarks the GetErrorCode function.
// This tests the performance of the optimized error code extraction.
func BenchmarkGetErrorCode(b *testing.B) {
	err := New(404, "not found")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetErrorCode(err)
	}
}

// BenchmarkIsErrorCodeEqual benchmarks the IsErrorCodeEqual function.
// This tests the performance of error code comparison.
func BenchmarkIsErrorCodeEqual(b *testing.B) {
	err := New(400, "bad request")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = IsErrorCodeEqual(err, 400)
	}
}

// BenchmarkHasErrorFields benchmarks the HasErrorFields function.
// This tests the performance of checking for error fields presence.
func BenchmarkHasErrorFields(b *testing.B) {
	err := New(422, "validation failed",
		NewErrorField("field1", "error1"),
		NewErrorField("field2", "error2"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = HasErrorFields(err)
	}
}

// BenchmarkGetErrorFields benchmarks the GetErrorFields function.
// This tests the performance of extracting error fields with copying.
func BenchmarkGetErrorFields(b *testing.B) {
	err := New(422, "validation failed",
		NewErrorField("field1", "error1"),
		NewErrorField("field2", "error2"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetErrorFields(err)
	}
}

// BenchmarkHasErrorField benchmarks the HasErrorField function.
// This tests the performance of searching for specific error fields.
func BenchmarkHasErrorField(b *testing.B) {
	err := New(422, "validation failed",
		NewErrorField("username", "username is required"),
		NewErrorField("email", "invalid email format"),
		NewErrorField("password", "password too weak"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = HasErrorField(err, "email")
	}
}

// BenchmarkError_String benchmarks the String method.
// This tests the performance of detailed string representation.
func BenchmarkError_String(b *testing.B) {
	err := New(422, "validation failed",
		NewErrorField("username", "username is required"),
		NewErrorField("email", "invalid email format"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = err.String()
	}
}

// BenchmarkError_IsEmpty benchmarks the IsEmpty method.
// This tests the performance of empty error detection.
func BenchmarkError_IsEmpty(b *testing.B) {
	err := Error{Code: 400, Message: "test error"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = err.IsEmpty()
	}
}

// BenchmarkGetErrorFieldMessage benchmarks the GetErrorFieldMessage function.
// This tests the performance of retrieving specific field messages.
func BenchmarkGetErrorFieldMessage(b *testing.B) {
	err := New(422, "validation failed",
		NewErrorField("username", "username is required"),
		NewErrorField("email", "invalid email format"),
		NewErrorField("password", "password too weak"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetErrorFieldMessage(err, "email")
	}
}

// BenchmarkErrorFieldCount benchmarks the ErrorFieldCount function.
// This tests the performance of counting error fields.
func BenchmarkErrorFieldCount(b *testing.B) {
	err := New(422, "validation failed",
		NewErrorField("field1", "error1"),
		NewErrorField("field2", "error2"),
		NewErrorField("field3", "error3"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ErrorFieldCount(err)
	}
}
