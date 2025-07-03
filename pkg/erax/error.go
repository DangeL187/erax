package erax

import (
	"errors"
)

type Error interface {
	error

	Unwrap() error

	Msg() string
	Meta(key string) (interface{}, Error)
	Metas() map[string]interface{}

	WithMeta(key string, value interface{}) Error
	WithMetas(metas map[string]interface{}) Error
}

type MetaProvider interface {
	Meta(key string) (interface{}, Error)
}

type ErrorType struct {
	err  error
	meta map[string]interface{}
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

func (e *ErrorType) Meta(key string) (interface{}, Error) {
	if e.Metas() != nil && len(e.Metas()) > 0 {
		if val, ok := e.meta[key]; ok {
			return val, nil
		}
	}

	var current Error
	ok := errors.As(e.err, &current)
	if !ok {
		return nil, NewFromString("key not found in error chain", "")
	}
	return current.Meta(key)
}

func (e *ErrorType) Metas() map[string]interface{} {
	return e.meta
}

func (e *ErrorType) WithMeta(key string, value interface{}) Error {
	if e.meta == nil {
		e.meta = make(map[string]interface{})
	}

	e.meta[key] = value

	return e
}

func (e *ErrorType) WithMetas(metas map[string]interface{}) Error {
	if e.meta == nil {
		e.meta = make(map[string]interface{})
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
