//go:build go1.18
// +build go1.18

package jsoniter

import (
	"testing"
)

// FuzzUnmarshal targets the main Unmarshal function with ConfigCompatibleWithStandardLibrary.
// It aims to find inputs that cause panics or unexpected errors when unmarshalling
// JSON (that jsoniter considers valid) into a generic interface{}.
func FuzzUnmarshal(f *testing.F) {
	// Seed corpus with various valid JSON structures
	f.Add([]byte(`{}`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`"hello"`))
	f.Add([]byte(`123`))
	f.Add([]byte(`true`))
	f.Add([]byte(`false`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"key": "value", "array": [1, null, true]}`))
	f.Add([]byte(`[{"a":"b"}, 123.45, "test"]`))

	json := ConfigCompatibleWithStandardLibrary

	f.Fuzz(func(t *testing.T, data []byte) {
		// Check if the input is considered valid by jsoniter first.
		isValid := json.Valid(data)

		var v interface{}
		err := json.Unmarshal(data, &v)

		// If jsoniter thought the data was valid, Unmarshal should not error
		// when decoding into an empty interface.
		// Panics are caught automatically by the fuzzing engine.
		if isValid && err != nil {
			t.Errorf("Unmarshal returned error for supposedly valid JSON: %v\ndata: %s", err, string(data))
		}

		// Optional: Consider comparing with encoding/json.Valid behavior, but note
		// that jsoniter might have slightly different validation rules.
		// stdValid := stdjson.Valid(data)
		// if stdValid != isValid { ... }
	})
}
