package erax

func IsErax(err error) bool {
	_, ok := asErax(err)
	return ok
}

func asErax(err error) (*errorType, bool) {
	if err == nil {
		return nil, false
	}

	stack := [4]error{err}
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
