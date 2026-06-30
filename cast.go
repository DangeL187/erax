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

// TODO: rewrite it without recursion
func cast(err error) error {
	if err == nil {
		return nil
	}

	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		children := uw.Unwrap()

		errs := make([]error, len(children))
		for i, child := range children {
			errs[i] = cast(child)
		}

		return &errorType{
			msg:  err.Error(),
			errs: errs,
		}
	}

	if uw, ok := err.(interface{ Unwrap() error }); ok {
		return &errorType{
			msg:   err.Error(),
			cause: cast(uw.Unwrap()),
		}
	}

	return err
}
