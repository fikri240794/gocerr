package gocerr

// ErrorField represents a field-specific validation error.
// It contains the field name and a descriptive error message
// explaining what validation rule was violated.
type ErrorField struct {
	Field   string // Name of the field that failed validation
	Message string // Human-readable validation error message
}

// NewErrorField creates a new ErrorField with the specified field name and message.
// This function provides a convenient way to create field-specific validation errors.
//
// Parameters:
//   - field: The name of the field that failed validation
//   - message: A descriptive error message explaining the validation failure
//
// Returns:
//   - ErrorField: A new ErrorField instance
//
// Example:
//   fieldErr := gocerr.NewErrorField("email", "Invalid email format")
//   fieldErr := gocerr.NewErrorField("age", "Age must be between 18 and 65")
func NewErrorField(field string, message string) ErrorField {
	return ErrorField{
		Field:   field,
		Message: message,
	}
}
