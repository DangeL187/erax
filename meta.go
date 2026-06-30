package erax

// MetaField represents a key-value pair for error metadata.
type MetaField struct {
	Key, Value string
}

// F is a convenience function to create a MetaField.
func F(k, v string) MetaField {
	return MetaField{Key: k, Value: v}
}

// WithMeta wraps an error with metadata and a new message.
//
// If the error is nil, returns nil.
// If no fields are provided, returns the original error unchanged.
func WithMeta(err error, message string, fields ...MetaField) error {
	if err == nil {
		return nil
	}

	if len(fields) == 0 {
		return err
	}

	return &errorType{
		cause: err,
		msg:   message,
		meta:  fields,
	}
}

// AddMeta is useful when metadata is added incrementally.
//
// If the error is already an erax error, only the key-value pair is added.
// The message argument is ignored in that case, so it should be left empty.
//
// If you already know all fields, prefer WithMeta instead,
// as it performs fewer allocations.
func AddMeta(err error, message string, key, value string) error {
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

// GetMeta searches for a metadata field by key across the entire error chain.
//
// It searches from the most recent error backwards through causes and children.
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

			if e.cause != nil {
				slice = append(slice, e.cause)
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
