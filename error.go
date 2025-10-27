package erax

import (
	"errors"
	"fmt"
	"io"
)

type errorType struct {
	err  error
	meta map[string]string
	msg  string
}

func (e *errorType) Unwrap() error { return e.err }

func (e *errorType) Error() string { return e.msg }

func (e *errorType) Format(s fmt.State, verb rune) {
	switch verb {
	case 'f':
		_, _ = io.WriteString(s, formatErrorChain(e, true))
	case 'v':
		if s.Flag('+') {
			unwrapped := e.Unwrap()
			msg := e.Error()

			if msg != "" {
				_, _ = io.WriteString(s, formatErrorChain(e, false))
			} else if unwrapped != nil {
				_, _ = io.WriteString(s, formatAlienError(fmt.Sprintf("%+v", unwrapped), true))
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &errorType{
		err: err,
		msg: message,
	}
}

func WithMeta(err error, key, value string) error {
	if err == nil {
		return nil
	}

	var e *errorType
	if errors.As(err, &e) {
		if e.meta == nil {
			e.meta = make(map[string]string)
		}
		e.meta[key] = value
		return err
	}

	return &errorType{
		err: err,
		msg: "",
		meta: map[string]string{
			key: value,
		},
	}
}

func GetMeta(err error, key string) (string, bool) {
	for err != nil {
		var e *errorType
		if errors.As(err, &e) {
			if e.meta != nil {
				if v, ok := e.meta[key]; ok {
					return v, true
				}
			}
		}
		err = errors.Unwrap(err)
	}
	return "", false
}

func GetMetas(err error) map[string]string {
	var e *errorType
	if !errors.As(err, &e) {
		return map[string]string{}
	}

	if e.meta == nil {
		return map[string]string{}
	}

	out := make(map[string]string, len(e.meta))
	for k, v := range e.meta {
		out[k] = v
	}

	return out
}
