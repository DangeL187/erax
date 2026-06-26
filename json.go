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

	if e, isErax := asErax(err); isErax {
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

		if len(e.errs) > 0 {
			buf.WriteString(`,"cause":`)

			if len(e.errs) == 1 {
				writeErrorJSON(buf, e.errs[0])
			} else {
				buf.WriteByte('[')
				for i, ue := range e.errs {
					if i > 0 {
						buf.WriteByte(',')
					}
					writeErrorJSON(buf, ue)
				}
				buf.WriteByte(']')
			}
		}
	}

	buf.WriteByte('}')
}

func writeEscapedString(b *bytes.Buffer, s string) {
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"', '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteByte('\\')
			b.WriteByte('n')
		case '\r':
			b.WriteByte('\\')
			b.WriteByte('r')
		case '\t':
			b.WriteByte('\\')
			b.WriteByte('t')
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('"')
}
