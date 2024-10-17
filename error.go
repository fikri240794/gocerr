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

func Parse(err error) (bool, Error) {
	if err == nil {
		return false, Error{}
	}

	customError, isCustomError := err.(Error)

	return isCustomError, customError
}
