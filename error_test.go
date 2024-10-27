package gocerr

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		Name       string
		Code       int
		Message    string
		ErrorField *ErrorField
		Expected   Error
	}{
		{
			Name:    "with error field",
			Code:    400,
			Message: "bad request",
			ErrorField: &ErrorField{
				Field:   "field1",
				Message: "field is required",
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
			Name:       "no error field",
			Code:       500,
			Message:    "internal server error",
			ErrorField: nil,
			Expected: Error{
				Code:        500,
				Message:     "internal server error",
				ErrorFields: []ErrorField{},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actualErr Error = New(testCases[i].Code, testCases[i].Message)
			if testCases[i].ErrorField != nil {
				actualErr = New(testCases[i].Code, testCases[i].Message, *testCases[i].ErrorField)
			}

			if testCases[i].Expected.Code != actualErr.Code {
				t.Errorf("expected code is %d, but got %d", testCases[i].Expected.Code, actualErr.Code)
			}

			if testCases[i].Expected.Message != actualErr.Message {
				t.Errorf("expected message is %s, but got %s", testCases[i].Expected.Message, actualErr.Message)
			}

			if len(testCases[i].Expected.ErrorFields) != len(actualErr.ErrorFields) {
				t.Errorf("expected length of error fields is %d, but got %d", len(testCases[i].Expected.ErrorFields), len(actualErr.ErrorFields))
			}

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

func TestError_Error(t *testing.T) {
	var (
		expectedMessage string
		actualErr       error
	)

	expectedMessage = "internal server errorr"
	actualErr = New(500, expectedMessage)

	if actualErr.Error() != expectedMessage {
		t.Errorf("expected error string return %s, but got %s", expectedMessage, actualErr.Error())
	}
}

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
			Name:  "error is not custom error",
			Error: errors.New("some error"),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   Error{},
				IsCustomError: false,
			},
		},
		{
			Name:  "error is custom error",
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
			Name:  "error is custom error with error fields",
			Error: New(400, "bad request", NewErrorField("field1", "field is required")),
			Expected: struct {
				CustomError   Error
				IsCustomError bool
			}{
				CustomError:   New(400, "bad request", NewErrorField("field1", "field is required")),
				IsCustomError: true,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualCustomError   Error
				actualIsCustomError bool
			)

			actualCustomError, actualIsCustomError = Parse(testCases[i].Error)

			if testCases[i].Expected.IsCustomError != actualIsCustomError {
				t.Errorf("expected is custom error is %t, but got %t", testCases[i].Expected.IsCustomError, actualIsCustomError)
			}

			if testCases[i].Expected.CustomError.Code != actualCustomError.Code {
				t.Errorf("expected custom error code is %d, but got %d", testCases[i].Expected.CustomError.Code, actualCustomError.Code)
			}

			if testCases[i].Expected.CustomError.Message != actualCustomError.Message {
				t.Errorf("expected custom error message is %s, but got %s", testCases[i].Expected.CustomError.Message, actualCustomError.Message)
			}

			if testCases[i].Error != nil && testCases[i].Expected.IsCustomError && testCases[i].Error.Error() != actualCustomError.Message {
				t.Errorf("expected error message is %s, but got %s", testCases[i].Error.Error(), actualCustomError.Message)
			}

			if len(testCases[i].Expected.CustomError.ErrorFields) != len(actualCustomError.ErrorFields) {
				t.Errorf("expected length of custom error error fields is %d, but got %d", len(testCases[i].Expected.CustomError.ErrorFields), len(actualCustomError.ErrorFields))
			}

			for j := 0; j < len(testCases[i].Expected.CustomError.ErrorFields); j++ {
				if testCases[i].Expected.CustomError.ErrorFields[j].Field != actualCustomError.ErrorFields[j].Field {
					t.Errorf("expected field of error fields sub item custom error is %s, but got %s", testCases[i].Expected.CustomError.ErrorFields[j].Field, actualCustomError.ErrorFields[j].Field)
				}
				if testCases[i].Expected.CustomError.ErrorFields[j].Message != actualCustomError.ErrorFields[j].Message {
					t.Errorf("expected message of error fields sub item custom error is %s, but got %s", testCases[i].Expected.CustomError.ErrorFields[j].Message, actualCustomError.ErrorFields[j].Message)
				}
			}
		})
	}
}

func TestIsErrorCodeEqual(t *testing.T) {
	var testCases []struct {
		Name        string
		Code        int
		Error       error
		Expectation bool
	} = []struct {
		Name        string
		Code        int
		Error       error
		Expectation bool
	}{
		{
			Name:        "error is not custom error",
			Code:        http.StatusInternalServerError,
			Error:       nil,
			Expectation: false,
		},
		{
			Name: "error code is equal",
			Code: http.StatusInternalServerError,
			Error: Error{
				Code: http.StatusInternalServerError,
			},
			Expectation: true,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual bool = IsErrorCodeEqual(testCases[i].Error, testCases[i].Code)

			if testCases[i].Expectation != actual {
				t.Errorf("expectation is %t, got %t", testCases[i].Expectation, actual)
			}
		})
	}
}
