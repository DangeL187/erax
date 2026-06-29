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

func (e *errorType) Unwrap() []error {
	if len(e.errs) > 0 {
		return e.errs
	}

	if e.cause != nil {
		return []error{e.cause}
	}

	return nil
}

func (e *errorType) Error() string { return e.msg }

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

func New(message string) error {
	return errors.New(message)
}
