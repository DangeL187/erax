package erax

import (
	"errors"
)

type Error interface {
	error

	Unwrap() error

	Msg() string
	Meta(key string) (string, error)
	Metas() map[string]string

	WithMeta(key string, value string) Error
	WithMetas(metas map[string]string) Error
}

type ErrorType struct {
	err  error
	meta map[string]string
	msg  string
}

func (e *ErrorType) Unwrap() error {
	return e.err
}

func (e *ErrorType) Error() string {
	return e.err.Error()
}

func (e *ErrorType) Msg() string {
	return e.msg
}

func (e *ErrorType) Meta(key string) (string, error) {
	if e.Metas() != nil && len(e.Metas()) > 0 {
		if val, ok := e.meta[key]; ok {
			return val, nil
		}
	}

	var current Error
	ok := errors.As(e.err, &current)
	if !ok {
		return "", errors.New("key not found in error chain")
	}
	return current.Meta(key)
}

func (e *ErrorType) Metas() map[string]string {
	return e.meta
}

func (e *ErrorType) WithMeta(key, value string) Error {
	if e.meta == nil {
		e.meta = make(map[string]string)
	}

	e.meta[key] = value

	return e
}

func (e *ErrorType) WithMetas(metas map[string]string) Error {
	if e.meta == nil {
		e.meta = make(map[string]string)
	}

	for k, v := range metas {
		e.meta[k] = v
	}

	return e
}

func New(err error, msg string) Error {
	if err == nil {
		return nil
	}

	return &ErrorType{
		err: err,
		msg: msg,
	}
}

func NewFromString(err string, msg string) Error {
	return &ErrorType{
		err: errors.New(err),
		msg: msg,
	}
}
