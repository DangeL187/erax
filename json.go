package erax

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

func FormatToJSONMap(err error) map[string]any {
	if err == nil {
		return nil
	}

	if next, isErax := asErax(err); isErax {
		return errorToMap(next)
	}

	return map[string]any{
		"message": err.Error(),
	}
}

func FormatToJSONString(err error) string {
	if err == nil {
		return "{}"
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	writeErrorJSON(buf, err)
	res := buf.String()

	if buf.Cap() <= 16384 {
		bufferPool.Put(buf)
	}

	return res
}

func FromJSONMap(m map[string]any) error {
	if m == nil || len(m) == 0 {
		return nil
	}

	return mapToError(m)
}

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
