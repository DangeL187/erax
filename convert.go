package erax

import "errors"

func errorToMap(err *errorType) map[string]any {
	m := map[string]any{
		"message": err.msg,
	}

	if len(err.meta) > 0 {
		m["meta"] = err.meta
	}

	hasCause := err.cause != nil
	hasErrs := len(err.errs) > 0

	if !hasCause && !hasErrs {
		return m
	}

	if hasCause && !hasErrs {
		m["cause"] = FormatToJSONMap(err.cause)
	} else {
		totalLen := len(err.errs)
		if hasCause {
			totalLen++
		}

		causeSlice := make([]map[string]any, 0, totalLen)
		if hasCause {
			causeSlice = append(causeSlice, FormatToJSONMap(err.cause))
		}
		for _, ue := range err.errs {
			causeSlice = append(causeSlice, FormatToJSONMap(ue))
		}
		m["cause"] = causeSlice
	}

	return m
}

func mapToError(m map[string]any) error {
	msg, msgOk := m["message"].(string)
	if !msgOk {
		return nil
	}

	meta, metaOk := m["meta"].([]MetaField)
	cause, causeOk := m["cause"]

	if (meta == nil || !metaOk) && !causeOk {
		return errors.New(msg)
	}

	err := &errorType{
		msg: msg,
	}

	if meta != nil && metaOk {
		err.meta = meta
	}

	if causeOk {
		if value, ok := cause.(map[string]any); ok {
			if childErr := mapToError(value); childErr != nil {
				err.cause = childErr
			}
		} else if value, ok := cause.([]map[string]any); ok {
			err.errs = make([]error, len(value))
			for i, c := range value {
				err.errs[i] = mapToError(c)
			}
		}
	}

	return err
}
