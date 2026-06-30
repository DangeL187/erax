package erax

import "bytes"

// writeErrorJSON writes an error's JSON representation directly to a buffer.
func writeErrorJSON(buf *bytes.Buffer, err error) {
	if err == nil {
		return
	}

	buf.WriteString(`{"message":`)
	writeEscapedString(buf, err.Error())

	if e, ok := err.(*errorType); ok {
		writeEraxJSONFields(buf, e)
	} else if e, isErax := asErax(err); isErax {
		writeEraxJSONFields(buf, e)
	}

	buf.WriteByte('}')
}

// writeEraxJSONFields writes erax-specific JSON fields (metadata and cause) to a buffer.
func writeEraxJSONFields(buf *bytes.Buffer, e *errorType) {
	if len(e.meta) > 0 {
		buf.WriteString(`,"meta":{`)
		for i, field := range e.meta {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeEscapedString(buf, field.Key)
			buf.WriteByte(':')
			writeEscapedString(buf, field.Value)
		}
		buf.WriteByte('}')
	}

	hasCause := e.cause != nil
	hasErrs := len(e.errs) > 0

	if hasCause || hasErrs {
		buf.WriteString(`,"cause":`)

		if hasCause && !hasErrs {
			writeErrorJSON(buf, e.cause)
		} else if !hasCause && len(e.errs) == 1 {
			writeErrorJSON(buf, e.errs[0])
		} else {
			buf.WriteByte('[')
			first := true

			if hasCause {
				writeErrorJSON(buf, e.cause)
				first = false
			}

			for _, ue := range e.errs {
				if !first {
					buf.WriteByte(',')
				}
				writeErrorJSON(buf, ue)
				first = false
			}
			buf.WriteByte(']')
		}
	}
}

// writeEscapedString writes a string to a buffer with JSON escaping applied.
func writeEscapedString(b *bytes.Buffer, s string) {
	b.WriteByte('"')
	last := 0
	lenS := len(s)
	for i := 0; i < lenS; i++ {
		c := s[i]
		if c == '"' || c == '\\' || c == '\n' || c == '\r' || c == '\t' {
			if i > last {
				b.WriteString(s[last:i])
			}
			b.WriteByte('\\')
			switch c {
			case '\n':
				b.WriteByte('n')
			case '\r':
				b.WriteByte('r')
			case '\t':
				b.WriteByte('t')
			default:
				b.WriteByte(c)
			}
			last = i + 1
		}
	}
	if last < lenS {
		b.WriteString(s[last:])
	}
	b.WriteByte('"')
}
