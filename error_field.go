package gocerr

type ErrorField struct {
	Field   string
	Message string
}

func NewErrorField(field string, message string) ErrorField {
	return ErrorField{
		Field:   field,
		Message: message,
	}
}
