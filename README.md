# About

`erax` is a Go package that enhances error handling with structured metadata and beautiful CLI output.

It provides error chaining, custom metadata, and styled error traces using the [lipgloss](https://github.com/charmbracelet/lipgloss) library.

![image](https://github.com/DangeL187/erax/blob/main/img/demo.png)

# Features

- 🌈 Styled and readable error trace output for CLI
- 🔗 Error chaining with `Unwrap()`
- 🏷️ Attach and retrieve key-value metadata
- 🎨 Configurable colors for trace output

# Usage

```go
func CreateUser() erax.Error {
	err := errors.New("email already in use")
	return erax.New(err, "failed to create user").
		WithMeta("code", 503).
		WithMeta("user_error", "An account with this email already exists.")
}

func Register() erax.Error {
	err := CreateUser()
	return erax.New(err, "failed to register\nbecause of ducks!")
}

...

erax.Trace(err)
```
**Output:**
```
 ▼ [ERROR TRACE]
 ├─ failed to register
 │  because of ducks!
 ├─ failed to create user
 │   ├─ code: 503
 │   ╰─ user_error: An account with this email already exists.
 ╰─ email already in use
```

Check out the full [example](https://github.com/DangeL187/erax/blob/main/examples/main.go)

# API Overview

**Error Creation:**
```go
New(err error, msg string) erax.Error            // Creates a new erax.Error from error
NewFromString(err string, msg string) erax.Error // Creates a new erax.Error from string
```

**Error Methods:**
```go
Msg() string                                  // Retrieves Error's message
Meta(key string) (interface{}, Error)         // Retrieves metadata by key
Metas() map[string]interface{}                // Returns all metadata as a map

WithMeta(key string, value interface{}) Error // Attaches a key-value pair
WithMetas(metas map[string]interface{}) Error // Attaches multiple metadata entries
```

**Error Trace Output:**
```go
Trace(err Error) string // Pretty-prints the full error chain and metadata.

// Customize CLI output colors:
SetErrorColor(color lipgloss.Color)
SetKeyColor(color lipgloss.Color)
SetNormalColor(color lipgloss.Color)
```
