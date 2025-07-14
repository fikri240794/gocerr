# gocerr

**gocerr** is a comprehensive custom error handling library for Go that extends the built-in `error` interface. It provides structured error handling with error codes, detailed messages, and field-level validation errors, making it perfect for APIs, web services, and complex applications.

## ✨ Features

- **🎯 Type Safety**: Full compile-time type checking with Go interfaces
- **📊 Structured Errors**: Support for error codes, messages, and field-specific validation
- **🔍 Rich Introspection**: Comprehensive error analysis and debugging capabilities
- **🛡️ Production Ready**: Extensive test coverage (100%) and battle-tested
- **📝 Well Documented**: Complete Go-standard documentation with examples
- **🔧 Easy Integration**: Drop-in replacement for standard Go errors

## 📦 Installation

```bash
go get github.com/fikri240794/gocerr
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/fikri240794/gocerr"
)

func main() {
    // Create a custom error with validation details
    err := gocerr.New(
        http.StatusBadRequest,
        "User registration failed",
        gocerr.NewErrorField("username", "Username must be at least 3 characters"),
        gocerr.NewErrorField("email", "Email address is already taken"),
        gocerr.NewErrorField("password", "Password must contain uppercase letter"),
    )

    // Standard error interface
    fmt.Println(err.Error()) // "User registration failed"

    // Rich error introspection
    if gocerr.IsErrorCodeEqual(err, 400) {
        fmt.Println("Handling bad request...")
        
        if gocerr.HasErrorFields(err) {
            fields := gocerr.GetErrorFields(err)
            for _, field := range fields {
                fmt.Printf("❌ %s: %s\n", field.Field, field.Message)
            }
        }
    }
}
```

## 📖 Documentation

### Core Types

#### Error
The main error type that implements Go's `error` interface:

```go
type Error struct {
    Code        int          // Numeric error code (e.g., HTTP status codes)
    Message     string       // Human-readable error message  
    ErrorFields []ErrorField // Field-specific validation errors
}
```

#### ErrorField
Represents field-level validation errors:

```go
type ErrorField struct {
    Field   string // Field name that failed validation
    Message string // Validation error message
}
```

### Core Functions

#### Creating Errors

```go
// Simple error
err := gocerr.New(500, "Internal server error")

// Error with field validations
err := gocerr.New(422, "Validation failed",
    gocerr.NewErrorField("email", "Invalid email format"),
    gocerr.NewErrorField("age", "Must be between 18-65"),
)
```

#### Parsing and Checking Errors

```go
// Parse any error to custom error
if customErr, ok := gocerr.Parse(err); ok {
    fmt.Printf("Error code: %d\n", customErr.Code)
}

// Extract error code
code := gocerr.GetErrorCode(err) // Returns 0 for non-custom errors

// Check specific error code
if gocerr.IsErrorCodeEqual(err, 404) {
    // Handle not found
}
```

#### Working with Error Fields

```go
// Check if error has field validations
if gocerr.HasErrorFields(err) {
    // Get all error fields (defensive copy)
    fields := gocerr.GetErrorFields(err)
    
    // Check for specific field
    if gocerr.HasErrorField(err, "email") {
        // Get specific field message
        msg := gocerr.GetErrorFieldMessage(err, "email")
        fmt.Printf("Email error: %s\n", msg)
    }
    
    // Count validation errors
    count := gocerr.ErrorFieldCount(err)
    fmt.Printf("Total validation errors: %d\n", count)
}
```

#### Debugging and Logging

```go
// Standard error message
fmt.Println(err.Error()) // "Validation failed"

// Detailed debug information
fmt.Println(err.String()) // "Error{Code: 422, Message: "Validation failed", ErrorFields: [...]}"

// Check if error is empty/zero value
if !err.IsEmpty() {
    // Handle non-empty error
}
```

## 🎯 Use Cases

### API Error Handling

```go
func CreateUser(req UserRequest) error {
    var validationErrors []gocerr.ErrorField
    
    if req.Username == "" {
        validationErrors = append(validationErrors, 
            gocerr.NewErrorField("username", "Username is required"))
    }
    
    if !isValidEmail(req.Email) {
        validationErrors = append(validationErrors,
            gocerr.NewErrorField("email", "Invalid email format"))
    }
    
    if len(validationErrors) > 0 {
        return gocerr.New(400, "Validation failed", validationErrors...)
    }
    
    // ... create user logic
    return nil
}
```

### HTTP Error Responses

```go
func handleError(w http.ResponseWriter, err error) {
    if customErr, ok := gocerr.Parse(err); ok {
        w.WriteHeader(customErr.Code)
        
        response := ErrorResponse{
            Message: customErr.Message,
            Code:    customErr.Code,
        }
        
        if gocerr.HasErrorFields(err) {
            response.ValidationErrors = gocerr.GetErrorFields(err)
        }
        
        json.NewEncoder(w).Encode(response)
        return
    }
    
    // Handle standard errors
    http.Error(w, "Internal Server Error", 500)
}
```

### Database Operations

```go
func (r *UserRepository) Create(user User) error {
    err := r.db.Create(&user).Error
    if err != nil {
        if isDuplicateError(err) {
            return gocerr.New(409, "User already exists",
                gocerr.NewErrorField("email", "Email address is already registered"))
        }
        return gocerr.New(500, "Database error")
    }
    return nil
}
```

## 🔧 Advanced Usage

### Error Chaining

```go
func processUser(id string) error {
    user, err := getUserFromDB(id)
    if err != nil {
        // Preserve the original error while adding context
        if gocerr.IsErrorCodeEqual(err, 404) {
            return gocerr.New(404, fmt.Sprintf("User with ID %s not found", id))
        }
        return err // Pass through other errors
    }
    
    // ... process user
    return nil
}
```

### Custom Error Codes

```go
const (
    ErrCodeValidation   = 1001
    ErrCodeDuplicate    = 1002
    ErrCodeUnauthorized = 1003
    ErrCodeRateLimit    = 1004
)

err := gocerr.New(ErrCodeRateLimit, "Rate limit exceeded")
if gocerr.IsErrorCodeEqual(err, ErrCodeRateLimit) {
    // Handle rate limiting
}
```

### Integration with Popular Frameworks

#### Gin Framework

```go
func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid JSON"})
        return
    }
    
    if err := userService.CreateUser(req); err != nil {
        if customErr, ok := gocerr.Parse(err); ok {
            c.JSON(customErr.Code, gin.H{
                "message": customErr.Message,
                "fields":  gocerr.GetErrorFields(err),
            })
            return
        }
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(201, gin.H{"message": "User created successfully"})
}
```

## 📋 API Reference

### Functions

| Function | Description |
|----------|-------------|
| `New(code, message, fields...)` | Creates a new custom error |
| `NewErrorField(field, message)` | Creates a new error field |
| `Parse(err)` | Converts standard error to custom error |
| `GetErrorCode(err)` | Extracts error code |
| `IsErrorCodeEqual(err, code)` | Checks if error has specific code |
| `HasErrorFields(err)` | Checks if error has field validations |
| `GetErrorFields(err)` | Gets all error fields (defensive copy) |
| `HasErrorField(err, field)` | Checks for specific error field |
| `GetErrorFieldMessage(err, field)` | Gets message for specific field |
| `ErrorFieldCount(err)` | Counts number of error fields |

### Methods

| Method | Description |
|--------|-------------|
| `Error()` | Returns error message (implements error interface) |
| `String()` | Returns detailed debug representation |
| `IsEmpty()` | Checks if error is zero value |

---