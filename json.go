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

// FormatToJSONMap converts an error to a JSON-compatible map representation. Handles both erax and standard Go errors.
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

// FormatToJSONString converts an error to a JSON string representation.
//
// Uses a sync.Pool for buffer reuse to improve performance.
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

// FromJSONMap reconstructs an erax error from a JSON map representation.
//
// Use this to deserialize errors previously serialized with FormatToJSONMap.
func FromJSONMap(m map[string]any) error {
	if m == nil || len(m) == 0 {
		return nil
	}

	return mapToError(m)
}
