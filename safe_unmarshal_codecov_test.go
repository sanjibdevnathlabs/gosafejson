package gosafejson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test coverage for all primitive type safe unmarshal decoders
func TestSafeUnmarshalCoverage_AllPrimitiveTypes(t *testing.T) {
	should := require.New(t)

	// Test uint8 with various wrong types
	t.Run("uint8_with_string", func(t *testing.T) {
		type TestStruct struct {
			Value uint8 `json:"value"`
		}
		jsonStr := `{"value": "not-a-number"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(uint8(0), result.Value)
	})

	t.Run("uint8_with_bool", func(t *testing.T) {
		type TestStruct struct {
			Value uint8 `json:"value"`
		}
		jsonStr := `{"value": true}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(uint8(0), result.Value)
	})

	t.Run("uint8_with_object", func(t *testing.T) {
		type TestStruct struct {
			Value uint8 `json:"value"`
		}
		jsonStr := `{"value": {"not": "a-number"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(uint8(0), result.Value)
	})

	// Test uint16 with various wrong types
	t.Run("uint16_with_array", func(t *testing.T) {
		type TestStruct struct {
			Value uint16 `json:"value"`
		}
		jsonStr := `{"value": [1, 2, 3]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(uint16(0), result.Value)
	})

	// Test uint32 with various wrong types
	t.Run("uint32_with_null", func(t *testing.T) {
		type TestStruct struct {
			Value uint32 `json:"value"`
		}
		jsonStr := `{"value": null}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(uint32(0), result.Value)
	})

	// Test uint64 with various wrong types
	t.Run("uint64_with_string", func(t *testing.T) {
		type TestStruct struct {
			Value uint64 `json:"value"`
		}
		jsonStr := `{"value": "not-a-number"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(uint64(0), result.Value)
	})

	// Test int8 with various wrong types
	t.Run("int8_with_string", func(t *testing.T) {
		type TestStruct struct {
			Value int8 `json:"value"`
		}
		jsonStr := `{"value": "not-a-number"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err) // String to int mismatch should be handled gracefully
		should.Equal(int8(0), result.Value)
	})

	// Test int16 with various wrong types
	t.Run("int16_with_string", func(t *testing.T) {
		type TestStruct struct {
			Value int16 `json:"value"`
		}
		jsonStr := `{"value": "not-a-number"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(int16(0), result.Value)
	})

	// Test int32 with various wrong types
	t.Run("int32_with_bool", func(t *testing.T) {
		type TestStruct struct {
			Value int32 `json:"value"`
		}
		jsonStr := `{"value": false}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(int32(0), result.Value)
	})

	// Test int64 with various wrong types
	t.Run("int64_with_array", func(t *testing.T) {
		type TestStruct struct {
			Value int64 `json:"value"`
		}
		jsonStr := `{"value": []}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(int64(0), result.Value)
	})

	// Test float32 with various wrong types
	t.Run("float32_with_string", func(t *testing.T) {
		type TestStruct struct {
			Value float32 `json:"value"`
		}
		jsonStr := `{"value": "not-a-float"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(float32(0), result.Value)
	})

	// Test float64 with various wrong types
	t.Run("float64_with_object", func(t *testing.T) {
		type TestStruct struct {
			Value float64 `json:"value"`
		}
		jsonStr := `{"value": {"not": "a-float"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(float64(0), result.Value)
	})

	// Test string with various wrong types
	t.Run("string_with_number", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}
		jsonStr := `{"value": 12345}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value)
	})

	t.Run("string_with_bool", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}
		jsonStr := `{"value": true}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value)
	})

	t.Run("string_with_array", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}
		jsonStr := `{"value": [1, 2, 3]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value)
	})

	// Test bool with various wrong types
	t.Run("bool_with_string", func(t *testing.T) {
		type TestStruct struct {
			Value bool `json:"value"`
		}
		jsonStr := `{"value": "not-a-bool"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value)
	})

	t.Run("bool_with_number", func(t *testing.T) {
		type TestStruct struct {
			Value bool `json:"value"`
		}
		jsonStr := `{"value": 123}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value)
	})

	t.Run("bool_with_object", func(t *testing.T) {
		type TestStruct struct {
			Value bool `json:"value"`
		}
		jsonStr := `{"value": {"not": "a-bool"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value)
	})
}

// Test coverage for slice and map safe unmarshal with errors
func TestSafeUnmarshalCoverage_CollectionTypes(t *testing.T) {
	should := require.New(t)

	// Test slice with various wrong types - should return CompositeError
	t.Run("slice_with_object", func(t *testing.T) {
		type TestStruct struct {
			Values []string `json:"values"`
		}
		jsonStr := `{"values": {"not": "an-array"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Values)

		// Verify it's a CompositeError
		compositeErr, ok := err.(*CompositeError)
		should.True(ok)
		should.Len(compositeErr.Errors, 1)
	})

	t.Run("slice_with_string", func(t *testing.T) {
		type TestStruct struct {
			Values []int `json:"values"`
		}
		jsonStr := `{"values": "not-an-array"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Values)
	})

	t.Run("slice_with_number", func(t *testing.T) {
		type TestStruct struct {
			Values []string `json:"values"`
		}
		jsonStr := `{"values": 12345}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Values)
	})

	// Test map with various wrong types - should return CompositeError
	t.Run("map_with_array", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]string `json:"data"`
		}
		jsonStr := `{"data": ["not", "an", "object"]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Data)

		// Verify it's a CompositeError
		compositeErr, ok := err.(*CompositeError)
		should.True(ok)
		should.Len(compositeErr.Errors, 1)
	})

	t.Run("map_with_string", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]int `json:"data"`
		}
		jsonStr := `{"data": "not-an-object"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Data)
	})

	t.Run("map_with_number", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]bool `json:"data"`
		}
		jsonStr := `{"data": 12345}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Data)
	})

	t.Run("map_with_bool", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]interface{} `json:"data"`
		}
		jsonStr := `{"data": true}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)
		should.Nil(result.Data)
	})
}

// Test coverage for complex nested structures with multiple type mismatches
func TestSafeUnmarshalCoverage_ComplexStructures(t *testing.T) {
	should := require.New(t)

	t.Run("deeply_nested_with_multiple_errors", func(t *testing.T) {
		type NestedStruct struct {
			ID       uint64            `json:"id"`
			Name     string            `json:"name"`
			Active   bool              `json:"active"`
			Score    float32           `json:"score"`
			Tags     []string          `json:"tags"`
			Metadata map[string]string `json:"metadata"`
		}

		type TestStruct struct {
			Count    int                    `json:"count"`
			Items    []NestedStruct         `json:"items"`
			Settings map[string]interface{} `json:"settings"`
		}

		// JSON with multiple type mismatches at different levels
		jsonStr := `{
			"count": "not-a-number",
			"items": "not-an-array",
			"settings": [
				{"not": "an-object"}
			]
		}`

		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)

		// Verify it's a CompositeError with multiple errors
		compositeErr, ok := err.(*CompositeError)
		should.True(ok)
		should.True(len(compositeErr.Errors) >= 2) // Should have multiple errors

		// Check that good default values are set
		should.Equal(0, result.Count)
		should.Nil(result.Items)
		should.Nil(result.Settings)
	})

	t.Run("mixed_success_and_failure", func(t *testing.T) {
		type TestStruct struct {
			GoodString string            `json:"good_string"`
			BadString  string            `json:"bad_string"`
			GoodInt    int               `json:"good_int"`
			BadInt     int               `json:"bad_int"`
			GoodSlice  []string          `json:"good_slice"`
			BadSlice   []string          `json:"bad_slice"`
			GoodMap    map[string]string `json:"good_map"`
			BadMap     map[string]string `json:"bad_map"`
		}

		jsonStr := `{
			"good_string": "hello",
			"bad_string": 12345,
			"good_int": 42,
			"bad_int": "not-a-number",
			"good_slice": ["a", "b", "c"],
			"bad_slice": {"not": "an-array"},
			"good_map": {"key": "value"},
			"bad_map": ["not", "an", "object"]
		}`

		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)

		// Verify it's a CompositeError
		compositeErr, ok := err.(*CompositeError)
		should.True(ok, "Expected CompositeError, got: %T", err)
		should.True(len(compositeErr.Errors) >= 2, "Expected at least 2 errors, got: %d", len(compositeErr.Errors)) // Should have multiple errors

		// Check that good values were parsed correctly
		should.Equal("hello", result.GoodString)
		should.Equal(42, result.GoodInt)
		should.Equal([]string{"a", "b", "c"}, result.GoodSlice)
		should.Equal(map[string]string{"key": "value"}, result.GoodMap)

		// Check that bad values have zero values
		should.Equal("", result.BadString)
		should.Equal(0, result.BadInt)
		should.Nil(result.BadSlice)
		should.Nil(result.BadMap)
	})
}

// Test coverage for edge cases and error paths
func TestSafeUnmarshalCoverage_EdgeCases(t *testing.T) {
	should := require.New(t)

	t.Run("all_primitive_types_with_null", func(t *testing.T) {
		type TestStruct struct {
			UInt8Value   uint8   `json:"uint8_value"`
			UInt16Value  uint16  `json:"uint16_value"`
			UInt32Value  uint32  `json:"uint32_value"`
			UInt64Value  uint64  `json:"uint64_value"`
			Int8Value    int8    `json:"int8_value"`
			Int16Value   int16   `json:"int16_value"`
			Int32Value   int32   `json:"int32_value"`
			Int64Value   int64   `json:"int64_value"`
			Float32Value float32 `json:"float32_value"`
			Float64Value float64 `json:"float64_value"`
			StringValue  string  `json:"string_value"`
			BoolValue    bool    `json:"bool_value"`
		}

		jsonStr := `{
			"uint8_value": null,
			"uint16_value": null,
			"uint32_value": null,
			"uint64_value": null,
			"int8_value": null,
			"int16_value": null,
			"int32_value": null,
			"int64_value": null,
			"float32_value": null,
			"float64_value": null,
			"string_value": null,
			"bool_value": null
		}`

		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err) // null values should be handled gracefully

		// All values should be zero values
		should.Equal(uint8(0), result.UInt8Value)
		should.Equal(uint16(0), result.UInt16Value)
		should.Equal(uint32(0), result.UInt32Value)
		should.Equal(uint64(0), result.UInt64Value)
		should.Equal(int8(0), result.Int8Value)
		should.Equal(int16(0), result.Int16Value)
		should.Equal(int32(0), result.Int32Value)
		should.Equal(int64(0), result.Int64Value)
		should.Equal(float32(0), result.Float32Value)
		should.Equal(float64(0), result.Float64Value)
		should.Equal("", result.StringValue)
		should.Equal(false, result.BoolValue)
	})

	t.Run("extremely_wrong_types", func(t *testing.T) {
		type TestStruct struct {
			Number uint64            `json:"number"`
			Text   string            `json:"text"`
			Flag   bool              `json:"flag"`
			List   []int             `json:"list"`
			Dict   map[string]string `json:"dict"`
		}

		// Use completely wrong types for everything
		jsonStr := `{
			"number": {"deeply": {"nested": {"object": "value"}}},
			"text": [[[["deeply", "nested"], "array"]]],
			"flag": {"another": "object"},
			"list": "definitely-not-an-array",
			"dict": 3.14159
		}`

		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err)

		// Verify it's a CompositeError with multiple errors
		compositeErr, ok := err.(*CompositeError)
		should.True(ok)
		should.True(len(compositeErr.Errors) >= 2) // Should have multiple errors

		// All values should be zero values
		should.Equal(uint64(0), result.Number)
		should.Equal("", result.Text)
		should.Equal(false, result.Flag)
		should.Nil(result.List)
		should.Nil(result.Dict)
	})
}
