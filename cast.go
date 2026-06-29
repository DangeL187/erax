package erax

func Cast(err error) error {
	if err == nil {
		return nil
	}

	if IsErax(err) {
		return err
	}

	return cast(err)
}

func cast(err error) error {
	if err == nil {
		return nil
	}

	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		children := uw.Unwrap()

		errs := make([]error, 0, len(children))
		for _, child := range children {
			if child != nil {
				errs = append(errs, cast(child))
			}
		}

		return &errorType{
			msg:  err.Error(),
			errs: errs,
		}
	}

	if uw, ok := err.(interface{ Unwrap() error }); ok {
		if child := uw.Unwrap(); child != nil {
			return &errorType{
				msg:   err.Error(),
				cause: cast(child),
			}
		}
	}

	return New(err.Error())
}
