package gocerr

type Error struct {
	Code        int
	Message     string
	ErrorFields []ErrorField
}

func New(code int, message string, errorFields ...ErrorField) Error {
	var err Error = Error{
		Code:        code,
		Message:     message,
		ErrorFields: errorFields,
	}

	return err
}

func (e Error) Error() string {
	return e.Message
}

func Parse(err error) (Error, bool) {
	var (
		customError   Error
		isCustomError bool
	)

	if err == nil {
		return Error{}, false
	}

	customError, isCustomError = err.(Error)

	return customError, isCustomError
}

func GetErrorCode(err error) int {
	var (
		customError   Error
		isCustomError bool
	)

	customError, isCustomError = Parse(err)
	if !isCustomError {
		return 0
	}

	return customError.Code
}

func IsErrorCodeEqual(err error, code int) bool {
	return GetErrorCode(err) == code
}
