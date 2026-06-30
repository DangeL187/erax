package erax

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type errorType struct {
	cause error
	errs  []error
	meta  []MetaField
	msg   string
}

// Unwrap returns the child errors of this error.
//
// This implements Go's error unwrapping interface.
// Returns either []error containing cause or errors, or nil if there are no children.
func (e *errorType) Unwrap() []error {
	if len(e.errs) > 0 {
		return e.errs
	}

	if e.cause != nil {
		return []error{e.cause}
	}

	return nil
}

// Error returns the error message string.
func (e *errorType) Error() string { return e.msg }

// Format implements fmt.Formatter for custom formatting.
//
// With %+v it prints a detailed error trace with tree visualization.
func (e *errorType) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			var sb strings.Builder
			sb.Grow(512)
			sb.WriteString(message)
			sb.WriteByte('\n')
			formatErrorChain(&sb, e, false, nil)
			_, _ = io.WriteString(s, sb.String())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

// New creates a new standard Go error with the given message.
func New(message string) error {
	return errors.New(message)
}
