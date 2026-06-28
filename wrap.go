package erax

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	if joinErr, ok := err.(interface{ Unwrap() []error }); ok {
		if !IsErax(err) {
			return &errorType{
				errs: joinErr.Unwrap(),
				msg:  message,
			}
		}
	}

	return &errorType{
		cause: err,
		msg:   message,
	}
}

func WrapWithErrors(err error, message string, newErrors ...error) error {
	newLen := len(newErrors)

	isAlienJoin := func(e error) bool {
		if e == nil {
			return false
		}
		if _, ok := e.(interface{ Unwrap() []error }); ok {
			if _, isErax := e.(*errorType); !isErax {
				return true
			}
		}
		return false
	}

	if newLen == 0 {
		if err == nil {
			if message != "" {
				return &errorType{msg: message}
			}
			return nil
		}
		if isAlienJoin(err) {
			return &errorType{errs: err.(interface{ Unwrap() []error }).Unwrap(), msg: message}
		}
		return &errorType{cause: err, msg: message}
	}

	hasJoin := isAlienJoin(err)
	if !hasJoin {
		for i := 0; i < newLen; i++ {
			if isAlienJoin(newErrors[i]) {
				hasJoin = true
				break
			}
		}
	}

	if !hasJoin {
		liveNew := 0
		for i := 0; i < newLen; i++ {
			if newErrors[i] != nil {
				liveNew++
			}
		}

		totalLen := liveNew
		if err != nil {
			totalLen++
		}

		if totalLen == 0 {
			if message != "" {
				return &errorType{msg: message}
			}
			return nil
		}

		if totalLen == 1 {
			var single error
			if err != nil {
				single = err
			} else {
				for i := 0; i < newLen; i++ {
					if newErrors[i] != nil {
						single = newErrors[i]
						break
					}
				}
			}

			return &errorType{cause: single, msg: message}
		}

		res := make([]error, totalLen)
		idx := 0
		for i := 0; i < newLen; i++ {
			if newErrors[i] != nil {
				res[idx] = newErrors[i]
				idx++
			}
		}
		if err != nil {
			res[idx] = err
		}

		return &errorType{errs: res, msg: message}
	}

	totalLen := 0
	countFn := func(e error) int {
		if e == nil {
			return 0
		}
		if isAlienJoin(e) {
			return len(e.(interface{ Unwrap() []error }).Unwrap())
		}
		return 1
	}

	for i := 0; i < newLen; i++ {
		totalLen += countFn(newErrors[i])
	}
	totalLen += countFn(err)

	res := make([]error, totalLen)
	idx := 0

	appendFn := func(e error) {
		if e == nil {
			return
		}
		if isAlienJoin(e) {
			children := e.(interface{ Unwrap() []error }).Unwrap()
			for i := 0; i < len(children); i++ {
				if children[i] != nil {
					res[idx] = children[i]
					idx++
				}
			}
		} else {
			res[idx] = e
			idx++
		}
	}

	for i := 0; i < newLen; i++ {
		appendFn(newErrors[i])
	}
	appendFn(err)

	if idx < totalLen {
		res = res[:idx]
	}

	if idx == 1 {
		return &errorType{cause: res[0], msg: message}
	}

	return &errorType{errs: res, msg: message}
}

// WrapWithErrorsLegacy is deprecated but works well. Use it if you hate everything new
func WrapWithErrorsLegacy(err error, message string, newErrors ...error) error {
	newErrorsLen := len(newErrors)
	isNewErrorsValid := newErrorsLen > 0

	if err == nil && !isNewErrorsValid {
		return nil
	}

	if !isNewErrorsValid {
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

	res := make([]error, newErrorsLen+1)
	copy(res, newErrors)
	res[newErrorsLen] = err

	return &errorType{
		errs: res,
		msg:  message,
	}
}
