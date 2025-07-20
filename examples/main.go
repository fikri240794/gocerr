package main

import (
	"fmt"
	"net/http"

	"github.com/fikri240794/gocerr"
)

func main() {
	// --- 1. New: Create a basic custom error ---
	basicErr := gocerr.New(500, "Internal server error")
	fmt.Println("[New]", basicErr.Error()) // Output: Internal server error

	// --- 2. NewErrorField: Create a field validation error ---
	field := gocerr.NewErrorField("username", "Username is required")
	fmt.Printf("[NewErrorField] Field: %s, Message: %s\n", field.Field, field.Message)

	// --- 3. New: Error with multiple fields ---
	multiFieldErr := gocerr.New(422, "Validation failed",
		gocerr.NewErrorField("email", "Invalid email format"),
		gocerr.NewErrorField("password", "Password too short"),
	)
	fmt.Println("[New with fields]", multiFieldErr.String())

	// --- 4. Parse: Convert error to custom error (success & fail) ---
	if custom, ok := gocerr.Parse(multiFieldErr); ok {
		fmt.Printf("[Parse] Code: %d, Message: %s\n", custom.Code, custom.Message)
	}
	if _, ok := gocerr.Parse(fmt.Errorf("not a custom error")); !ok {
		fmt.Println("[Parse] Not a custom error")
	}

	// --- 5. GetErrorCode: Extract error code ---
	fmt.Println("[GetErrorCode]", gocerr.GetErrorCode(multiFieldErr))            // 422
	fmt.Println("[GetErrorCode]", gocerr.GetErrorCode(fmt.Errorf("not custom"))) // 0

	// --- 6. IsErrorCodeEqual: Check error code ---
	fmt.Println("[IsErrorCodeEqual]", gocerr.IsErrorCodeEqual(multiFieldErr, 422)) // true
	fmt.Println("[IsErrorCodeEqual]", gocerr.IsErrorCodeEqual(multiFieldErr, 400)) // false

	// --- 7. HasErrorFields: Check if error has any field errors ---
	fmt.Println("[HasErrorFields]", gocerr.HasErrorFields(multiFieldErr)) // true
	fmt.Println("[HasErrorFields]", gocerr.HasErrorFields(basicErr))      // false

	// --- 8. GetErrorFields: Get all error fields (defensive copy) ---
	fields := gocerr.GetErrorFields(multiFieldErr)
	for _, f := range fields {
		fmt.Printf("[GetErrorFields] Field: %s, Message: %s\n", f.Field, f.Message)
	}

	// --- 9. HasErrorField: Check for a specific field error ---
	fmt.Println("[HasErrorField] email:", gocerr.HasErrorField(multiFieldErr, "email")) // true
	fmt.Println("[HasErrorField] phone:", gocerr.HasErrorField(multiFieldErr, "phone")) // false

	// --- 10. GetErrorFieldMessage: Get message for a specific field ---
	fmt.Println("[GetErrorFieldMessage] email:", gocerr.GetErrorFieldMessage(multiFieldErr, "email"))
	fmt.Println("[GetErrorFieldMessage] phone:", gocerr.GetErrorFieldMessage(multiFieldErr, "phone"))

	// --- 11. ErrorFieldCount: Count number of field errors ---
	fmt.Println("[ErrorFieldCount]", gocerr.ErrorFieldCount(multiFieldErr)) // 2
	fmt.Println("[ErrorFieldCount]", gocerr.ErrorFieldCount(basicErr))      // 0

	// --- 12. String: Detailed debug output ---
	fmt.Println("[String]", multiFieldErr.String())

	// --- 13. IsEmpty: Check if error is zero value ---
	var empty gocerr.Error
	fmt.Println("[IsEmpty]", empty.IsEmpty())         // true
	fmt.Println("[IsEmpty]", multiFieldErr.IsEmpty()) // false

	// --- 14. Complex: Full workflow with HTTP status, parsing, and all helpers ---
	complexErr := gocerr.New(
		http.StatusBadRequest,
		"Request failed",
		gocerr.NewErrorField("username", "Username is required"),
		gocerr.NewErrorField("email", "Invalid email format"),
	)
	if gocerr.IsErrorCodeEqual(complexErr, http.StatusBadRequest) {
		fmt.Println("[Complex] Bad request error detected!")
	}
	if gocerr.HasErrorFields(complexErr) {
		for _, f := range gocerr.GetErrorFields(complexErr) {
			fmt.Printf("[Complex] Field: %s, Error: %s\n", f.Field, f.Message)
		}
	}
	if gocerr.HasErrorField(complexErr, "email") {
		fmt.Println("[Complex] Email error:", gocerr.GetErrorFieldMessage(complexErr, "email"))
	}
	fmt.Println("[Complex] ErrorFieldCount:", gocerr.ErrorFieldCount(complexErr))
	fmt.Println("[Complex] String:", complexErr.String())
}
