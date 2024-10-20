# Go Custom Error (gocerr)
Go custom error handler in details that implements the built-in Go `error` interface. It lets you create detailed errors with a code, message, and additional error struct fields validation.

## Installation
```bash
go get github.com/fikri240794/gocerr
```

## Usage
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fikri240794/gocerr"
)

func main() {
	var (
		err           error
		isCustomError bool
		customError   gocerr.Error
	)

	// error from gocerr
	err = gocerr.New(
		http.StatusBadRequest, // not only for http, you can set any error code here
		http.StatusText(http.StatusBadRequest), // set error message
		gocerr.NewErrorField("field1", "field is required"), // additional error struct field validation
		gocerr.NewErrorField("field2", fmt.Sprintf("min value is %d", 50)), // additional error struct field validation
		gocerr.NewErrorField("fieldN", "error message validation"), // additional error struct field validation
	)

	fmt.Println(err.Error()) // print the error message from Error.Message

	// parse error to gocerr
	// if err is gocerr error
	// will return isCustomError true and customError struct with value from err parameter
	// otherwise, will isCustomError false and customError empty struct
	isCustomError, customError = gocerr.Parse(err)

	if isCustomError {
		fmt.Println(customError.Code)    // error code
		fmt.Println(customError.Message) // error message
		for i := 0; i < len(customError.ErrorFields); i++ {
			fmt.Println(customError.ErrorFields[i].Field)   // additional error field name
			fmt.Println(customError.ErrorFields[i].Message) // additional error field message
		}
	}
}
```
