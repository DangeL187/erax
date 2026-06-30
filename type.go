package erax

// IsErax checks whether an error is an erax error type.
func IsErax(err error) bool {
	_, ok := asErax(err)
	return ok
}

// asErax internal helper to safely convert an error to *errorType. Returns the underlying erax error if it is one.
func asErax(err error) (*errorType, bool) {
	e, ok := err.(*errorType)
	return e, ok
}
