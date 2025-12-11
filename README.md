# About

`erax` is a Go package that enhances error handling with structured metadata and beautiful CLI output.

It provides error chaining, custom metadata, and styled error traces using the [lipgloss](https://github.com/charmbracelet/lipgloss) library.

![image](https://github.com/DangeL187/erax/blob/main/img/demo.png)

# Features

- 🌈 Styled and readable error trace output for CLI
- 🔗 Error chaining with `Unwrap()`
- 🏷️ Attach and retrieve key-value metadata
- 🎨 Configurable colors for trace output
- 🔄 Compatible with standard and third-party errors (e.g., pkg/errors)
- 📄 **JSON compatibility (SOON)**

# Usage

```go
func CreateUser() error {
    err := errors.New("email is already in use")

    err = erax.Wrap(err, "failed to create user")
    err = erax.WithMeta(err, "code", "503")
    err = erax.WithMeta(err, "info", "This is a really\nreally long information.")
    err = erax.WithMeta(err, "user_error", "An account with this email already exists.")
    return err
}

func Register() error {
    err := CreateUser()
    return erax.Wrap(err, "failed to register\nbecause of ducks!")
}

...

// Print trace
fmt.Printf("%f\n", err)

// Print trace without first line (▼ [ERROR TRACE]),
// use it if the final error is not erax
fmt.Printf("%+v\n", err)

errCode, _ := erax.GetMeta(err, "code")
errUserError, _ := erax.GetMeta(err, "user_error")
```
**Output:**
```
 ▼ [ERROR TRACE]
 ├─ failed to register
 │  because of ducks!
 ├─ failed to create user
 │   ├─ code: 503
 │   ├─ info: 
 │   │   This is a really
 │   │   really long information.
 │   ╰─ user_error: An account with this email already exists.
 ╰─ email is already in use
```

Check out the full [example](https://github.com/DangeL187/erax/blob/main/examples/main.go)

# API Overview

**Error Creation:**
```go
Wrap(err error, message string) error  // Wraps an existing error with an additional message
```

**Error Functions:**
```go
WithMeta(err error, key, value string) error   // Adds a key-value pair to the error's metadata
GetMeta(err error, key string) (string, bool)  // Retrieves a value from the error's metadata by key (recursively)
GetMetas(err error) map[string]string          // Returns all metadata from the error as a map
```

**Error Trace Output:**
```go
fmt.Printf("%f\n", err)  // Pretty-prints the full error chain and metadata

// Customize CLI output colors:
SetErrorColor(color lipgloss.Color)
SetKeyColor(color lipgloss.Color)
SetNormalColor(color lipgloss.Color)
SetValueColor(color lipgloss.Color)
```
