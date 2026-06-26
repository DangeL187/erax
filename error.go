package erax

import (
	"errors"
	"fmt"
	"io"
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
	case 'f':
		_, _ = io.WriteString(s, formatErrorChain(e, true, 0))
	case 'v':
		if s.Flag('+') {
			res := formatDefault(e, 0)
			if res != "" {
				_, _ = io.WriteString(s, res)
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

func formatDefault(e *errorType, nestingLevel int) string {
	msg := e.Error()

	if msg != "" {
		return formatErrorChain(e, false, nestingLevel)
	} else if len(e.errs) > 0 {
		return formatAlienError(e, true)
	}

	return ""
}

func New(message string) error {
	return errors.New(message)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &errorType{
		cause: err,
		msg:   message,
	}
}

func WrapWithErrors(err error, message string, newErrors ...error) error {
	newErrorsLen := len(newErrors)

	if err == nil && newErrorsLen == 0 {
		return nil
	}

	if newErrorsLen == 0 {
		return &errorType{
			cause: err,
			msg:   message,
		}
	}

	if err == nil {
		if newErrorsLen == 1 {
			return &errorType{cause: newErrors[0], msg: message}
		}
		return &errorType{errs: newErrors, msg: message}
	}

	if newErrorsLen == 1 {
		return &errorType{
			cause: err,
			errs:  []error{newErrors[0]},
			msg:   message,
		}
	}

	res := make([]error, newErrorsLen+1)
	copy(res, newErrors)
	res[newErrorsLen] = err

	return &errorType{
		errs: res,
		msg:  message,
	}
}

type MetaField struct {
	Key, Value string
}

func F(k, v string) MetaField {
	return MetaField{Key: k, Value: v}
}

func WithMeta(err error, message string, fields ...MetaField) error {
	if err == nil {
		return nil
	}

	if len(fields) == 0 {
		return err
	}

	e, isErax := asErax(err)
	if isErax {
		oldLen := len(e.meta)
		newLen := oldLen + len(fields)

		if newLen > cap(e.meta) {
			newMeta := make([]MetaField, newLen)
			copy(newMeta, e.meta)
			copy(newMeta[oldLen:], fields)
			e.meta = newMeta
		} else {
			e.meta = append(e.meta, fields...)
		}

		return err
	}

	return &errorType{
		cause: err,
		msg:   message,
		meta:  fields,
	}
}

func AddMeta(err error, message, key, value string) error {
	if err == nil {
		return nil
	}

	e, isErax := asErax(err)
	if isErax {
		e.meta = append(e.meta, MetaField{Key: key, Value: value})
		return err
	}

	return &errorType{
		cause: err,
		msg:   message,
		meta:  []MetaField{{Key: key, Value: value}},
	}
}

func GetMeta(err error, key string) (string, bool) {
	if err == nil {
		return "", false
	}

	stack := [8]error{err}
	slice := stack[:1]

	for len(slice) > 0 {
		current := slice[len(slice)-1]
		slice = slice[:len(slice)-1]

		if current == nil {
			continue
		}

		if e, ok := current.(*errorType); ok {
			for i := len(e.meta) - 1; i >= 0; i-- {
				if e.meta[i].Key == key {
					return e.meta[i].Value, true
				}
			}
			if len(e.errs) > 0 {
				slice = append(slice, e.errs...)
			}
			continue
		}

		if w, ok := current.(interface{ Unwrap() error }); ok {
			if next := w.Unwrap(); next != nil {
				slice = append(slice, next)
			}
		} else if w, ok := current.(interface{ Unwrap() []error }); ok {
			if nextS := w.Unwrap(); len(nextS) > 0 {
				slice = append(slice, nextS...)
			}
		}
	}

	return "", false
}

func GetMetas(err error) []MetaField {
	e, isErax := asErax(err)
	if !isErax {
		return []MetaField{}
	}

	return e.meta
}

func IsErax(err error) bool {
	_, ok := asErax(err)
	return ok
}

func asErax(err error) (*errorType, bool) {
	if err == nil {
		return nil, false
	}

	stack := [2]error{err}
	slice := stack[:1]

	for len(slice) > 0 {
		current := slice[len(slice)-1]
		slice = slice[:len(slice)-1]

		if current == nil {
			continue
		}

		if e, ok := current.(*errorType); ok {
			return e, ok
		}

		if w, ok := current.(interface{ Unwrap() []error }); ok {
			if nextS := w.Unwrap(); len(nextS) > 0 {
				slice = append(slice, nextS...)
			}
			continue
		}

		if w, ok := current.(interface{ Unwrap() error }); ok {
			if next := w.Unwrap(); next != nil {
				slice = append(slice, next)
			}
		}
	}

	return nil, false
}
