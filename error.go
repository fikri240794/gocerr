// Package gocerr provides a custom error handling solution that implements
// the built-in Go error interface. It allows creating detailed errors with
// error codes, messages, and field-specific validation errors.
package gocerr

import (
	"fmt"
	"strings"
)

// Error represents a custom error with additional context information.
// It implements the built-in error interface and provides structured
// error handling with support for field-level validation errors.
type Error struct {
	Code        int          // Numeric error code (e.g., HTTP status codes)
	Message     string       // Human-readable error message
	ErrorFields []ErrorField // Collection of field-specific validation errors
}

// New creates a new custom Error with the specified code, message, and optional error fields.
// This function is optimized to minimize allocations when no error fields are provided.
//
// Parameters:
//   - code: Numeric error code (commonly HTTP status codes, but can be any integer)
//   - message: Human-readable error message
//   - errorFields: Zero or more ErrorField instances for field-specific validation errors
//
// Returns:
//   - Error: A new Error instance with the provided details
//
// Example:
//
//	err := gocerr.New(400, "Invalid request",
//	    gocerr.NewErrorField("username", "Username is required"),
//	    gocerr.NewErrorField("email", "Invalid email format"))
func New(code int, message string, errorFields ...ErrorField) Error {
	// Pre-allocate slice only if error fields are provided to avoid unnecessary allocation
	var fields []ErrorField
	if len(errorFields) > 0 {
		// Use make with exact capacity to avoid slice growth during append operations
		fields = make([]ErrorField, len(errorFields))
		copy(fields, errorFields)
	}

	return Error{
		Code:        code,
		Message:     message,
		ErrorFields: fields,
	}
}

// Error implements the built-in error interface by returning the error message.
// This allows Error instances to be used anywhere a standard Go error is expected.
//
// Returns:
//   - string: The error message
func (e Error) Error() string {
	return e.Message
}

// Parse attempts to convert a standard Go error into a custom Error instance.
// This function uses type assertion to check if the provided error is actually
// a custom Error type.
//
// Parameters:
//   - err: The error to parse (can be nil)
//
// Returns:
//   - Error: The parsed custom error (zero value if parsing fails)
//   - bool: true if the error was successfully parsed as a custom Error, false otherwise
//
// Example:
//
//	if customErr, ok := gocerr.Parse(err); ok {
//	    fmt.Printf("Error code: %d\n", customErr.Code)
//	}
func Parse(err error) (Error, bool) {
	// Early return for nil errors to avoid unnecessary processing
	if err == nil {
		return Error{}, false
	}

	// Use type assertion to check if err is a custom Error
	// This is more efficient than using reflection
	if customError, ok := err.(Error); ok {
		return customError, true
	}

	return Error{}, false
}

// GetErrorCode extracts the error code from a standard Go error if it's a custom Error.
// Returns 0 if the error is nil or not a custom Error type.
//
// Parameters:
//   - err: The error to extract the code from
//
// Returns:
//   - int: The error code (0 if not a custom error or if error is nil)
//
// Example:
//
//	code := gocerr.GetErrorCode(err)
//	if code == 404 {
//	    // Handle not found error
//	}
func GetErrorCode(err error) int {
	// Direct use of Parse function eliminates code duplication
	// and leverages the early return optimization for nil errors
	if customError, ok := Parse(err); ok {
		return customError.Code
	}
	return 0
}

// IsErrorCodeEqual checks if an error is a custom Error with a specific error code.
// This is a convenience function that combines parsing and code comparison.
//
// Parameters:
//   - err: The error to check
//   - code: The expected error code to compare against
//
// Returns:
//   - bool: true if the error is a custom Error and its code matches the provided code
//
// Example:
//
//	if gocerr.IsErrorCodeEqual(err, 400) {
//	    // Handle bad request error
//	}
func IsErrorCodeEqual(err error, code int) bool {
	return GetErrorCode(err) == code
}

// HasErrorFields checks if a custom Error contains any error fields.
// This is useful for determining if field-level validation errors exist.
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - bool: true if the error is a custom Error and contains error fields
//
// Example:
//
//	if gocerr.HasErrorFields(err) {
//	    // Handle field validation errors
//	}
func HasErrorFields(err error) bool {
	if customError, ok := Parse(err); ok {
		return len(customError.ErrorFields) > 0
	}
	return false
}

// GetErrorFields extracts the error fields from a custom Error.
// Returns an empty slice if the error is not a custom Error or has no error fields.
//
// Parameters:
//   - err: The error to extract fields from
//
// Returns:
//   - []ErrorField: Slice of error fields (empty if none exist)
//
// Example:
//
//	fields := gocerr.GetErrorFields(err)
//	for _, field := range fields {
//	    fmt.Printf("Field: %s, Error: %s\n", field.Field, field.Message)
//	}
func GetErrorFields(err error) []ErrorField {
	if customError, ok := Parse(err); ok {
		// Return a copy to prevent external modification of internal state
		if len(customError.ErrorFields) == 0 {
			return nil
		}
		fields := make([]ErrorField, len(customError.ErrorFields))
		copy(fields, customError.ErrorFields)
		return fields
	}
	return nil
}

// HasErrorField checks if a custom Error contains an error field with the specified field name.
// This is useful for checking if a specific field failed validation.
//
// Parameters:
//   - err: The error to check
//   - fieldName: The name of the field to look for
//
// Returns:
//   - bool: true if the error is a custom Error and contains the specified field
//
// Example:
//
//	if gocerr.HasErrorField(err, "email") {
//	    // Handle email validation error
//	}
func HasErrorField(err error, fieldName string) bool {
	if customError, ok := Parse(err); ok {
		for _, field := range customError.ErrorFields {
			if field.Field == fieldName {
				return true
			}
		}
	}
	return false
}

// GetErrorFieldMessage retrieves the error message for a specific field.
// Returns an empty string if the field is not found or the error is not a custom Error.
//
// Parameters:
//   - err: The error to search in
//   - fieldName: The name of the field to get the message for
//
// Returns:
//   - string: The error message for the field (empty string if not found)
//
// Example:
//
//	message := gocerr.GetErrorFieldMessage(err, "email")
//	if message != "" {
//	    fmt.Printf("Email error: %s\n", message)
//	}
func GetErrorFieldMessage(err error, fieldName string) string {
	if customError, ok := Parse(err); ok {
		for _, field := range customError.ErrorFields {
			if field.Field == fieldName {
				return field.Message
			}
		}
	}
	return ""
}

// ErrorFieldCount returns the number of error fields in a custom Error.
// Returns 0 if the error is not a custom Error or has no error fields.
//
// Parameters:
//   - err: The error to count fields for
//
// Returns:
//   - int: The number of error fields
//
// Example:
//
//	count := gocerr.ErrorFieldCount(err)
//	fmt.Printf("Number of validation errors: %d\n", count)
func ErrorFieldCount(err error) int {
	if customError, ok := Parse(err); ok {
		return len(customError.ErrorFields)
	}
	return 0
}

// String provides a detailed string representation of the Error for debugging and logging.
// Unlike Error(), this method includes the error code and field-level details.
//
// Returns:
//   - string: A detailed string representation including code, message, and error fields
//
// Example:
//
//	fmt.Printf("Debug info: %s\n", customErr.String())
func (e Error) String() string {
	if len(e.ErrorFields) == 0 {
		return fmt.Sprintf("Error{Code: %d, Message: %q}", e.Code, e.Message)
	}

	// Use strings.Builder for efficient string concatenation with multiple error fields
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Error{Code: %d, Message: %q, ErrorFields: [", e.Code, e.Message))

	for i, field := range e.ErrorFields {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("{Field: %q, Message: %q}", field.Field, field.Message))
	}

	builder.WriteString("]}")
	return builder.String()
}

// IsEmpty checks if the Error is an empty/zero value.
// This is useful for checking if Parse returned a meaningful result.
//
// Returns:
//   - bool: true if the error is empty (all fields are zero values)
//
// Example:
//
//	if err, ok := gocerr.Parse(someErr); ok && !err.IsEmpty() {
//	    // Handle non-empty custom error
//	}
func (e Error) IsEmpty() bool {
	return e.Code == 0 && e.Message == "" && len(e.ErrorFields) == 0
}
