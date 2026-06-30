package erax

// Wrap wraps an error with a new message.
//
// If the error is nil, returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &errorType{
		cause: err,
		msg:   message,
	}
}

// WrapCast wraps an error with a new message.
//
// If the error is nil, returns nil.
//
// Implicitly casts errors to erax type.
func WrapCast(err error, message string) error {
	if err == nil {
		return nil
	}

	if !IsErax(err) {
		err = cast(err)
	}

	return &errorType{
		cause: err,
		msg:   message,
	}
}

// WrapWithErrors wraps an error with multiple additional errors.
func WrapWithErrors(err error, message string, newErrors ...error) error {
	newErrorsLen := len(newErrors)

	if newErrorsLen == 0 {
		if err == nil {
			return New(message)
		}

		return &errorType{
			cause: err,
			msg:   message,
		}
	}

	if err == nil {
		return &errorType{
			errs: newErrors,
			msg:  message,
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

// WrapWithErrorsCast wraps an error with multiple additional errors.
//
// Implicitly casts errors to erax type.
func WrapWithErrorsCast(err error, message string, newErrors ...error) error {
	newErrorsLen := len(newErrors)

	if !IsErax(err) {
		err = cast(err)
	}

	if newErrorsLen == 0 {
		if err == nil {
			return New(message)
		}

		return &errorType{
			cause: err,
			msg:   message,
		}
	}

	for i := 0; i < newErrorsLen; i++ {
		if !IsErax(newErrors[i]) {
			newErrors[i] = cast(newErrors[i])
		}
	}

	if err == nil {
		return &errorType{
			errs: newErrors,
			msg:  message,
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
