package erax

func IsErax(err error) bool {
	_, ok := asErax(err)
	return ok
}

func asErax(err error) (*errorType, bool) {
	e, ok := err.(*errorType)
	return e, ok
}
