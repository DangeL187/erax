package erax

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

	return &errorType{
		cause: err,
		msg:   message,
		meta:  fields,
	}
}

func AddMeta(err error, key, value string) error {
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

func GetMetas(err error) []MetaField {
	e, isErax := asErax(err)
	if !isErax {
		return []MetaField{}
	}

	return e.meta
}
