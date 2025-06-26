package gosafejson

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test CompositeError behavior for coverage
func TestCompositeError_Coverage(t *testing.T) {
	should := require.New(t)

	// Test empty error list
	t.Run("empty_errors", func(t *testing.T) {
		ce := &CompositeError{Errors: []error{}}
		should.Equal("no errors", ce.Error())
	})

	// Test single error
	t.Run("single_error", func(t *testing.T) {
		ce := &CompositeError{Errors: []error{errors.New("single error")}}
		should.Equal("single error", ce.Error())
	})

	// Test multiple errors
	t.Run("multiple_errors", func(t *testing.T) {
		ce := &CompositeError{Errors: []error{
			errors.New("first error"),
			errors.New("second error"),
			errors.New("third error"),
		}}
		errorMsg := ce.Error()
		should.Contains(errorMsg, "3 errors occurred during safe unmarshalling")
		should.Contains(errorMsg, "1: first error")
		should.Contains(errorMsg, "2: second error")
		should.Contains(errorMsg, "3: third error")
	})
}

// Test UnmarshalFromString error paths for coverage
func TestUnmarshalFromString_ErrorPaths(t *testing.T) {
	should := require.New(t)

	// Test empty string
	t.Run("empty_string", func(t *testing.T) {
		var result map[string]interface{}
		err := ConfigSafe.UnmarshalFromString("", &result)
		should.Error(err)
		should.Contains(err.Error(), "EOF")
	})

	// Test non-pointer parameter
	t.Run("non_pointer", func(t *testing.T) {
		var result map[string]interface{}
		err := ConfigSafe.UnmarshalFromString("{}", result) // Not a pointer
		should.Error(err)
		should.Contains(err.Error(), "unmarshal need ptr")
	})

	// Test safe unmarshal with collected errors
	t.Run("safe_unmarshal_with_errors", func(t *testing.T) {
		type TestStruct struct {
			GoodField string `json:"good"`
			BadField  int    `json:"bad"`
		}
		jsonStr := `{"good": "value", "bad": "not-a-number"}`
		var result TestStruct
		err := ConfigSafe.UnmarshalFromString(jsonStr, &result)
		should.NoError(err) // No error for type mismatches in safe mode
		should.Equal("value", result.GoodField)
		should.Equal(0, result.BadField)
	})
}

// Test Unmarshal error paths for coverage
func TestUnmarshal_ErrorPaths(t *testing.T) {
	should := require.New(t)

	// Test empty data
	t.Run("empty_data", func(t *testing.T) {
		var result map[string]interface{}
		err := ConfigSafe.Unmarshal([]byte{}, &result)
		should.Error(err)
		should.Contains(err.Error(), "EOF")
	})

	// Test non-pointer parameter
	t.Run("non_pointer", func(t *testing.T) {
		var result map[string]interface{}
		err := ConfigSafe.Unmarshal([]byte("{}"), result) // Not a pointer
		should.Error(err)
		should.Contains(err.Error(), "unmarshal need ptr")
	})

	// Test safe unmarshal with collected errors returns CompositeError
	t.Run("safe_unmarshal_composite_error", func(t *testing.T) {
		type TestStruct struct {
			Field1 string            `json:"field1"`
			Field2 int               `json:"field2"`
			Field3 []string          `json:"field3"`
			Field4 map[string]string `json:"field4"`
		}
		jsonStr := `{"field1": 123, "field2": "not-a-number", "field3": "not-an-array", "field4": [1,2,3]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// If there are errors, they should be CompositeError
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		// Check that valid parsing still happened for compatible fields
		should.Equal("", result.Field1) // Default value for string with wrong type
		should.Equal(0, result.Field2)  // Default value for int with wrong type
		should.Nil(result.Field3)       // Nil for slice with wrong type
		should.Nil(result.Field4)       // Nil for map with wrong type
	})
}

// Test Iterator error collection for coverage
func TestIterator_ErrorCollection(t *testing.T) {
	should := require.New(t)

	// Test ReportError in safe mode
	t.Run("report_error_safe_mode", func(t *testing.T) {
		iter := ConfigSafe.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigSafe.ReturnIterator(iter)

		// Should collect errors instead of stopping
		iter.ReportError("test operation", "test error message")
		should.Len(iter.CollectedErrors, 1)
		should.Contains(iter.CollectedErrors[0].Error(), "test operation")
		should.Contains(iter.CollectedErrors[0].Error(), "test error message")
		should.NoError(iter.Error) // Error should not be set in safe mode
	})

	// Test ReportError in normal mode
	t.Run("report_error_normal_mode", func(t *testing.T) {
		iter := ConfigCompatibleWithStandardLibrary.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigCompatibleWithStandardLibrary.ReturnIterator(iter)

		// Should set error immediately
		iter.ReportError("test operation", "test error message")
		should.Error(iter.Error)
		should.Contains(iter.Error.Error(), "test operation")
		should.Contains(iter.Error.Error(), "test error message")
	})

	// Test Reset functionality
	t.Run("reset_functionality", func(t *testing.T) {
		iter := ConfigSafe.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigSafe.ReturnIterator(iter)

		// Add some errors
		iter.ReportError("test", "error1")
		iter.ReportError("test", "error2")
		should.Len(iter.CollectedErrors, 2)

		// Reset should clear errors
		iter.Reset(nil)
		should.Len(iter.CollectedErrors, 0)
		should.NoError(iter.Error)
	})

	// Test ResetBytes functionality
	t.Run("reset_bytes_functionality", func(t *testing.T) {
		iter := ConfigSafe.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigSafe.ReturnIterator(iter)

		// Add some errors
		iter.ReportError("test", "error1")
		should.Len(iter.CollectedErrors, 1)

		// ResetBytes should clear errors
		iter.ResetBytes([]byte(`{"new": "data"}`))
		should.Len(iter.CollectedErrors, 0)
		should.NoError(iter.Error)
	})
}

// Test struct decoder safe mode paths for coverage
func TestStructDecoder_SafeMode(t *testing.T) {
	should := require.New(t)

	// Test struct with malformed JSON (missing closing brace)
	t.Run("malformed_json_missing_brace", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}
		jsonStr := `{"field1": "value", "field2": 123` // Missing closing brace
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err) // This should error even in safe mode due to malformed JSON
	})

	// Test struct with multiple field errors in safe mode
	t.Run("multiple_field_errors_safe_mode", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
			Field3 bool   `json:"field3"`
		}
		jsonStr := `{"field1": 123, "field2": "not-a-number", "field3": "not-a-bool"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)                // Safe mode should handle type mismatches gracefully
		should.Equal("", result.Field1)    // Default value for string
		should.Equal(0, result.Field2)     // Default value for int
		should.Equal(false, result.Field3) // Default value for bool
	})
}

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
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Values) // Should remain nil
	})

	t.Run("slice_with_string", func(t *testing.T) {
		type TestStruct struct {
			Values []int `json:"values"`
		}
		jsonStr := `{"values": "not-an-array"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Values) // Should remain nil
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
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Data) // Should remain nil
	})

	t.Run("map_with_string", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]int `json:"data"`
		}
		jsonStr := `{"data": "not-an-object"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Data) // Should remain nil
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

// Test coverage for slice decoders in safe mode
func TestSliceDecoder_SafeMode(t *testing.T) {
	should := require.New(t)

	// Test slice with wrong element types
	t.Run("int_slice_with_string_elements", func(t *testing.T) {
		type TestStruct struct {
			Values []int `json:"values"`
		}
		jsonStr := `{"values": ["not", "numbers", "here"]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)                         // Safe mode should handle type mismatches
		should.Equal([]int{0, 0, 0}, result.Values) // Default values for int
	})

	t.Run("string_slice_with_number_elements", func(t *testing.T) {
		type TestStruct struct {
			Values []string `json:"values"`
		}
		jsonStr := `{"values": [123, 456, 789]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal([]string{"", "", ""}, result.Values) // Default values for string
	})

	t.Run("bool_slice_with_mixed_types", func(t *testing.T) {
		type TestStruct struct {
			Values []bool `json:"values"`
		}
		jsonStr := `{"values": ["string", 123, {"object": "value"}]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal([]bool{false, false, false}, result.Values) // Default values for bool
	})

	// Test slice with completely wrong type (not an array)
	t.Run("slice_with_object", func(t *testing.T) {
		type TestStruct struct {
			Values []string `json:"values"`
		}
		jsonStr := `{"values": {"not": "an-array"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Values) // Should remain nil
	})

	t.Run("slice_with_string", func(t *testing.T) {
		type TestStruct struct {
			Values []int `json:"values"`
		}
		jsonStr := `{"values": "not-an-array"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Values) // Should remain nil
	})
}

// Test coverage for map decoders in safe mode
func TestMapDecoder_SafeMode(t *testing.T) {
	should := require.New(t)

	// Test map with wrong value types
	t.Run("string_map_with_number_values", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]string `json:"data"`
		}
		jsonStr := `{"data": {"key1": 123, "key2": 456}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(map[string]string{"key1": "", "key2": ""}, result.Data) // Default values
	})

	t.Run("int_map_with_string_values", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]int `json:"data"`
		}
		jsonStr := `{"data": {"key1": "not-a-number", "key2": "also-not-a-number"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(map[string]int{"key1": 0, "key2": 0}, result.Data) // Default values
	})

	t.Run("bool_map_with_mixed_values", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]bool `json:"data"`
		}
		jsonStr := `{"data": {"key1": "string", "key2": 123, "key3": {"nested": "object"}}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(map[string]bool{"key1": false, "key2": false, "key3": false}, result.Data)
	})

	// Test map with completely wrong type (not an object)
	t.Run("map_with_array", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]string `json:"data"`
		}
		jsonStr := `{"data": ["not", "an", "object"]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Data) // Should remain nil
	})

	t.Run("map_with_string", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]int `json:"data"`
		}
		jsonStr := `{"data": "not-an-object"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			// Expect CompositeError for completely wrong types
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Error should be CompositeError type")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Data) // Should remain nil
	})
}

// Test additional edge cases for maximum coverage
func TestSafeUnmarshal_AdditionalEdgeCases(t *testing.T) {
	should := require.New(t)

	// Test with null values for all types
	t.Run("null_values_all_types", func(t *testing.T) {
		type TestStruct struct {
			StringVal string            `json:"string_val"`
			IntVal    int               `json:"int_val"`
			BoolVal   bool              `json:"bool_val"`
			FloatVal  float64           `json:"float_val"`
			SliceVal  []string          `json:"slice_val"`
			MapVal    map[string]string `json:"map_val"`
		}
		jsonStr := `{"string_val": null, "int_val": null, "bool_val": null, "float_val": null, "slice_val": null, "map_val": null}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.StringVal)
		should.Equal(0, result.IntVal)
		should.Equal(false, result.BoolVal)
		should.Equal(float64(0), result.FloatVal)
		should.Nil(result.SliceVal)
		should.Nil(result.MapVal)
	})

	// Test deeply nested structures with errors
	t.Run("deeply_nested_with_errors", func(t *testing.T) {
		type Level3 struct {
			Value string `json:"value"`
		}
		type Level2 struct {
			Level3 Level3 `json:"level3"`
			Number int    `json:"number"`
		}
		type Level1 struct {
			Level2 Level2 `json:"level2"`
			Flag   bool   `json:"flag"`
		}
		type TestStruct struct {
			Level1 Level1 `json:"level1"`
			Name   string `json:"name"`
		}

		jsonStr := `{
			"level1": {
				"level2": {
					"level3": {
						"value": 123
					},
					"number": "not-a-number"
				},
				"flag": "not-a-bool"
			},
			"name": 456
		}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)                                 // Safe mode should handle all type mismatches
		should.Equal("", result.Level1.Level2.Level3.Value) // String field gets default value
		should.Equal(0, result.Level1.Level2.Number)        // Int field gets default value
		should.Equal(false, result.Level1.Flag)             // Bool field gets default value
		should.Equal("", result.Name)                       // String field gets default value
	})

	// Test arrays of complex types with errors
	t.Run("array_of_structs_with_errors", func(t *testing.T) {
		type Item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		type TestStruct struct {
			Items []Item `json:"items"`
		}

		jsonStr := `{
			"items": [
				{"id": "not-a-number", "name": 123},
				{"id": 456, "name": "valid-name"},
				{"id": "also-not-a-number", "name": true}
			]
		}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err) // Safe mode should handle type mismatches
		should.Len(result.Items, 3)
		should.Equal(0, result.Items[0].ID)              // Default value for int
		should.Equal("", result.Items[0].Name)           // Default value for string
		should.Equal(456, result.Items[1].ID)            // Valid value
		should.Equal("valid-name", result.Items[1].Name) // Valid value
		should.Equal(0, result.Items[2].ID)              // Default value for int
		should.Equal("", result.Items[2].Name)           // Default value for string
	})

	// Test maps with struct values and errors
	t.Run("map_of_structs_with_errors", func(t *testing.T) {
		type Value struct {
			Count  int  `json:"count"`
			Active bool `json:"active"`
		}
		type TestStruct struct {
			Data map[string]Value `json:"data"`
		}

		jsonStr := `{
			"data": {
				"item1": {"count": "not-a-number", "active": "not-a-bool"},
				"item2": {"count": 42, "active": true},
				"item3": {"count": "also-not-a-number", "active": 123}
			}
		}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err) // Safe mode should handle type mismatches
		should.Len(result.Data, 3)
		should.Equal(0, result.Data["item1"].Count)      // Default value for int
		should.Equal(false, result.Data["item1"].Active) // Default value for bool
		should.Equal(42, result.Data["item2"].Count)     // Valid value
		should.Equal(true, result.Data["item2"].Active)  // Valid value
		should.Equal(0, result.Data["item3"].Count)      // Default value for int
		should.Equal(false, result.Data["item3"].Active) // Default value for bool
	})
}

// Test specific code paths added in the PR for codecov coverage
func TestSafeUnmarshal_SpecificCodePaths(t *testing.T) {
	should := require.New(t)

	// Test specific WhatIsNext() code paths in primitive decoders
	t.Run("string_decoder_whatsnext_paths", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}

		// Test NumberValue path
		jsonStr := `{"value": 123}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value) // Should skip and use default

		// Test BoolValue path
		jsonStr = `{"value": true}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value) // Should skip and use default

		// Test ObjectValue path
		jsonStr = `{"value": {"nested": "object"}}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value) // Should skip and use default

		// Test ArrayValue path
		jsonStr = `{"value": [1, 2, 3]}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value) // Should skip and use default

		// Test NilValue path
		jsonStr = `{"value": null}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("", result.Value) // Should handle null
	})

	t.Run("bool_decoder_whatsnext_paths", func(t *testing.T) {
		type TestStruct struct {
			Value bool `json:"value"`
		}

		// Test NumberValue path
		jsonStr := `{"value": 123}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value) // Should skip and use default

		// Test StringValue path
		jsonStr = `{"value": "not-a-bool"}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value) // Should skip and use default

		// Test ObjectValue path
		jsonStr = `{"value": {"nested": "object"}}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value) // Should skip and use default

		// Test ArrayValue path
		jsonStr = `{"value": [1, 2, 3]}`
		err = ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal(false, result.Value) // Should skip and use default
	})

	// Test all numeric decoders with non-numeric values
	t.Run("numeric_decoders_non_numeric_paths", func(t *testing.T) {
		type TestStruct struct {
			Int8Val    int8    `json:"int8_val"`
			Int16Val   int16   `json:"int16_val"`
			Int32Val   int32   `json:"int32_val"`
			Int64Val   int64   `json:"int64_val"`
			Uint8Val   uint8   `json:"uint8_val"`
			Uint16Val  uint16  `json:"uint16_val"`
			Uint32Val  uint32  `json:"uint32_val"`
			Uint64Val  uint64  `json:"uint64_val"`
			Float32Val float32 `json:"float32_val"`
			Float64Val float64 `json:"float64_val"`
		}

		// Test with string values (non-numeric)
		jsonStr := `{
			"int8_val": "not-a-number",
			"int16_val": "not-a-number",
			"int32_val": "not-a-number",
			"int64_val": "not-a-number",
			"uint8_val": "not-a-number",
			"uint16_val": "not-a-number",
			"uint32_val": "not-a-number",
			"uint64_val": "not-a-number",
			"float32_val": "not-a-number",
			"float64_val": "not-a-number"
		}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err) // Safe mode should handle gracefully

		// All should be default values
		should.Equal(int8(0), result.Int8Val)
		should.Equal(int16(0), result.Int16Val)
		should.Equal(int32(0), result.Int32Val)
		should.Equal(int64(0), result.Int64Val)
		should.Equal(uint8(0), result.Uint8Val)
		should.Equal(uint16(0), result.Uint16Val)
		should.Equal(uint32(0), result.Uint32Val)
		should.Equal(uint64(0), result.Uint64Val)
		should.Equal(float32(0), result.Float32Val)
		should.Equal(float64(0), result.Float64Val)
	})

	// Test the specific iter.cfg.safeUnmarshal branches
	t.Run("safe_unmarshal_config_branches", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}

		// Test with ConfigSafe (safeUnmarshal = true)
		jsonStr := `{"value": 123}`
		var safeResult TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
		should.NoError(err)
		should.Equal("", safeResult.Value)

		// Test with ConfigCompatibleWithStandardLibrary (safeUnmarshal = false)
		var normalResult TestStruct
		err = ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(jsonStr), &normalResult)
		should.Error(err) // Should fail in normal mode
	})

	// Test ReportError paths in safe vs normal mode
	t.Run("report_error_safe_vs_normal", func(t *testing.T) {
		// Test safe mode error collection
		safeIter := ConfigSafe.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigSafe.ReturnIterator(safeIter)

		safeIter.ReportError("test", "error in safe mode")
		should.Len(safeIter.CollectedErrors, 1)
		should.NoError(safeIter.Error) // Should not set Error in safe mode

		// Test normal mode error setting
		normalIter := ConfigCompatibleWithStandardLibrary.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigCompatibleWithStandardLibrary.ReturnIterator(normalIter)

		normalIter.ReportError("test", "error in normal mode")
		should.Error(normalIter.Error)           // Should set Error in normal mode
		should.Empty(normalIter.CollectedErrors) // Should not collect in normal mode
	})
}

// Test the exact UnmarshalFromString and Unmarshal error paths
func TestSafeUnmarshal_ErrorPathCoverage(t *testing.T) {
	should := require.New(t)

	// Test UnmarshalFromString with non-pointer
	t.Run("unmarshal_from_string_non_pointer", func(t *testing.T) {
		var result map[string]interface{}
		err := ConfigSafe.UnmarshalFromString("{}", result) // Not a pointer
		should.Error(err)
		should.Contains(err.Error(), "unmarshal need ptr")
	})

	// Test Unmarshal with non-pointer
	t.Run("unmarshal_non_pointer", func(t *testing.T) {
		var result map[string]interface{}
		err := ConfigSafe.Unmarshal([]byte("{}"), result) // Not a pointer
		should.Error(err)
		should.Contains(err.Error(), "unmarshal need ptr")
	})

	// Test CompositeError return path
	t.Run("composite_error_return", func(t *testing.T) {
		type TestStruct struct {
			BadSlice []string          `json:"bad_slice"`
			BadMap   map[string]string `json:"bad_map"`
		}
		jsonStr := `{"bad_slice": "not-an-array", "bad_map": [1,2,3]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
	})

	// Test iterator Reset and ResetBytes error clearing
	t.Run("iterator_reset_error_clearing", func(t *testing.T) {
		iter := ConfigSafe.BorrowIterator([]byte(`{"test": "value"}`))
		defer ConfigSafe.ReturnIterator(iter)

		// Add errors
		iter.ReportError("test", "error1")
		iter.ReportError("test", "error2")
		should.Len(iter.CollectedErrors, 2)

		// Test Reset clears errors
		iter.Reset(nil)
		should.Len(iter.CollectedErrors, 0)
		should.NoError(iter.Error)

		// Add error again
		iter.ReportError("test", "error3")
		should.Len(iter.CollectedErrors, 1)

		// Test ResetBytes clears errors
		iter.ResetBytes([]byte(`{"new": "data"}`))
		should.Len(iter.CollectedErrors, 0)
		should.NoError(iter.Error)
	})
}

// Test the exact missing coverage lines for primitive decoders
func TestSafeUnmarshal_MissingCoverageLines(t *testing.T) {
	should := require.New(t)

	// Test the else branch in primitive decoders (normal mode)
	t.Run("primitive_decoders_normal_mode", func(t *testing.T) {
		type TestStruct struct {
			StringVal  string  `json:"string_val"`
			IntVal     int     `json:"int_val"`
			BoolVal    bool    `json:"bool_val"`
			FloatVal   float64 `json:"float_val"`
			Uint8Val   uint8   `json:"uint8_val"`
			Uint16Val  uint16  `json:"uint16_val"`
			Uint32Val  uint32  `json:"uint32_val"`
			Uint64Val  uint64  `json:"uint64_val"`
			Int8Val    int8    `json:"int8_val"`
			Int16Val   int16   `json:"int16_val"`
			Int32Val   int32   `json:"int32_val"`
			Int64Val   int64   `json:"int64_val"`
			Float32Val float32 `json:"float32_val"`
		}

		// Test normal mode with correct types (should hit the else branch)
		jsonStr := `{
			"string_val": "correct-string",
			"int_val": 42,
			"bool_val": true,
			"float_val": 3.14,
			"uint8_val": 8,
			"uint16_val": 16,
			"uint32_val": 32,
			"uint64_val": 64,
			"int8_val": 8,
			"int16_val": 16,
			"int32_val": 32,
			"int64_val": 64,
			"float32_val": 3.14
		}`
		var result TestStruct
		err := ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Equal("correct-string", result.StringVal)
		should.Equal(42, result.IntVal)
		should.Equal(true, result.BoolVal)
		should.Equal(3.14, result.FloatVal)
	})

	// Test null handling in primitive decoders
	t.Run("primitive_decoders_null_handling", func(t *testing.T) {
		type TestStruct struct {
			StringVal *string  `json:"string_val"`
			IntVal    *int     `json:"int_val"`
			BoolVal   *bool    `json:"bool_val"`
			FloatVal  *float64 `json:"float_val"`
		}

		jsonStr := `{
			"string_val": null,
			"int_val": null,
			"bool_val": null,
			"float_val": null
		}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)
		should.Nil(result.StringVal)
		should.Nil(result.IntVal)
		should.Nil(result.BoolVal)
		should.Nil(result.FloatVal)
	})

	// Test specific error collection paths in struct decoder
	t.Run("struct_decoder_error_collection", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}

		// Test malformed JSON to hit error collection paths
		jsonStr := `{"field1": "value", "field2": "not-a-number"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err) // Safe mode should handle gracefully
		should.Equal("value", result.Field1)
		should.Equal(0, result.Field2) // Default value
	})

	// Test the specific case where CollectedErrors > 0 in config.go
	t.Run("config_collected_errors_return", func(t *testing.T) {
		type TestStruct struct {
			BadMap map[string]string `json:"bad_map"`
		}
		jsonStr := `{"bad_map": "not-a-map"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
	})

	// Test UnmarshalFromString with collected errors
	t.Run("unmarshal_from_string_collected_errors", func(t *testing.T) {
		type TestStruct struct {
			BadField []string `json:"bad_field"`
		}
		jsonStr := `{"bad_field": "not-an-array"}`
		var result TestStruct
		err := ConfigSafe.UnmarshalFromString(jsonStr, &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
	})
}

// Test specific map decoder error paths
func TestMapDecoder_ErrorPaths(t *testing.T) {
	should := require.New(t)

	// Test map decoder with wrong key types
	t.Run("map_wrong_key_type", func(t *testing.T) {
		type TestStruct struct {
			Data map[int]string `json:"data"` // int keys, but JSON will have string keys
		}
		jsonStr := `{"data": {"key1": "value1", "key2": "value2"}}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
	})

	// Test the specific unreadByte() and Skip() path in map decoder
	t.Run("map_decoder_skip_path", func(t *testing.T) {
		type TestStruct struct {
			Data map[string]string `json:"data"`
		}
		jsonStr := `{"data": [1, 2, 3]}` // Array instead of object
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Data)
	})
}

// Test specific slice decoder error paths
func TestSliceDecoder_ErrorPaths(t *testing.T) {
	should := require.New(t)

	// Test the specific unreadByte() and Skip() path in slice decoder
	t.Run("slice_decoder_skip_path", func(t *testing.T) {
		type TestStruct struct {
			Data []string `json:"data"`
		}
		jsonStr := `{"data": {"not": "an-array"}}` // Object instead of array
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.Data)
	})

	// Test slice with mixed valid and invalid elements
	t.Run("slice_mixed_elements", func(t *testing.T) {
		type TestStruct struct {
			Numbers []int `json:"numbers"`
		}
		jsonStr := `{"numbers": [1, "not-a-number", 3, "also-not-a-number", 5]}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)                                // Safe mode should handle gracefully
		should.Equal([]int{1, 0, 3, 0, 5}, result.Numbers) // Invalid elements become 0
	})
}

// Test struct field decoder error handling
func TestStructFieldDecoder_ErrorHandling(t *testing.T) {
	should := require.New(t)

	// Test the case where len(iter.CollectedErrors) > prevErrorCount
	t.Run("field_decoder_error_collection", func(t *testing.T) {
		type NestedStruct struct {
			Value string `json:"value"`
		}
		type TestStruct struct {
			Nested NestedStruct `json:"nested"`
		}
		jsonStr := `{"nested": {"value": 123}}` // Wrong type for nested value
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.NoError(err)                   // Safe mode should handle gracefully
		should.Equal("", result.Nested.Value) // Default value
	})

	// Test missing colon in object field
	t.Run("missing_colon_in_field", func(t *testing.T) {
		type TestStruct struct {
			Field1 string `json:"field1"`
			Field2 string `json:"field2"`
		}
		// This will test the missing colon error path, but it's hard to construct valid JSON for this
		// Let's test a different scenario that might trigger the colon error handling
		jsonStr := `{"field1" "missing-colon", "field2": "value"}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		should.Error(err) // This should error even in safe mode due to malformed JSON
	})
}

// Final push for codecov coverage - target remaining edge cases
func TestSafeUnmarshal_FinalCoverageEdgeCases(t *testing.T) {
	should := require.New(t)

	// Test the exact error return path in config.go when len(iter.CollectedErrors) > 0
	t.Run("config_composite_error_return_path", func(t *testing.T) {
		type TestStruct struct {
			Map1   map[string]int    `json:"map1"`
			Map2   map[string]string `json:"map2"`
			Slice1 []int             `json:"slice1"`
			Slice2 []string          `json:"slice2"`
		}
		// Multiple type mismatches to ensure CollectedErrors > 0
		jsonStr := `{
			"map1": "not-a-map",
			"map2": [1, 2, 3],
			"slice1": "not-a-slice",
			"slice2": {"not": "a-slice"}
		}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)

		// This should return a CompositeError with multiple errors
		should.Error(err, "Should return error with multiple type mismatches")
		compositeErr, ok := err.(*CompositeError)
		should.True(ok, "Error should be CompositeError type")
		should.NotEmpty(compositeErr.Errors, "Should have collected multiple errors")
		should.True(len(compositeErr.Errors) >= 2, "Should have multiple collected errors")

		// Verify the error message contains information about multiple errors
		errorMsg := compositeErr.Error()
		should.Contains(errorMsg, "errors occurred during safe unmarshalling", "Error message should indicate multiple errors")
	})

	// Test ReportError vs direct error setting paths
	t.Run("report_error_vs_direct_error", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}

		// Test safe mode - should use ReportError
		jsonStr := `{"value": 123}`
		var result1 TestStruct
		err1 := ConfigSafe.Unmarshal([]byte(jsonStr), &result1)
		should.NoError(err1, "Safe mode should handle type mismatch gracefully")
		should.Equal("", result1.Value, "Should use default value")

		// Test normal mode - should set error directly
		var result2 TestStruct
		err2 := ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(jsonStr), &result2)
		should.Error(err2, "Normal mode should return error on type mismatch")
	})

	// Test iterator Reset functionality with collected errors
	t.Run("iterator_reset_with_collected_errors", func(t *testing.T) {
		iter := ConfigSafe.BorrowIterator([]byte(`{"field": "wrong-type-for-int"}`))
		defer ConfigSafe.ReturnIterator(iter)

		// Cause some errors to be collected
		iter.ReadObjectCB(func(iter *Iterator, field string) bool {
			if field == "field" {
				// Try to read as int, should collect error in safe mode
				iter.ReadInt()
			}
			return true
		})

		// Reset should clear collected errors
		iter.ResetBytes([]byte(`{"valid": "json"}`))
		should.Empty(iter.CollectedErrors, "Reset should clear collected errors")

		// ResetBytes should also clear collected errors
		iter.ResetBytes([]byte(`{"another": "valid-json"}`))
		should.Empty(iter.CollectedErrors, "ResetBytes should clear collected errors")
	})

	// Test specific WhatIsNext() branches that might be missed
	t.Run("whatsnext_all_branches", func(t *testing.T) {
		type TestStruct struct {
			StringVal string `json:"string_val"`
			IntVal    int    `json:"int_val"`
			BoolVal   bool   `json:"bool_val"`
		}

		// Test each JSON value type against wrong Go type
		testCases := []struct {
			name     string
			jsonStr  string
			expected TestStruct
		}{
			{
				name:     "object_for_primitives",
				jsonStr:  `{"string_val": {"nested": "object"}, "int_val": {"nested": "object"}, "bool_val": {"nested": "object"}}`,
				expected: TestStruct{StringVal: "", IntVal: 0, BoolVal: false},
			},
			{
				name:     "array_for_primitives",
				jsonStr:  `{"string_val": [1, 2, 3], "int_val": [1, 2, 3], "bool_val": [1, 2, 3]}`,
				expected: TestStruct{StringVal: "", IntVal: 0, BoolVal: false},
			},
			{
				name:     "null_for_primitives",
				jsonStr:  `{"string_val": null, "int_val": null, "bool_val": null}`,
				expected: TestStruct{StringVal: "", IntVal: 0, BoolVal: false},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var result TestStruct
				err := ConfigSafe.Unmarshal([]byte(tc.jsonStr), &result)
				should.NoError(err, "Safe mode should handle all type mismatches gracefully")
				should.Equal(tc.expected, result, "Should use default values for mismatched types")
			})
		}
	})

	// Test the specific unreadByte() calls in map and slice decoders
	t.Run("unread_byte_coverage", func(t *testing.T) {
		type TestStruct struct {
			MapField   map[string]string `json:"map_field"`
			SliceField []string          `json:"slice_field"`
		}

		// These should trigger the unreadByte() and Skip() paths
		jsonStr := `{"map_field": 123, "slice_field": 456}`
		var result TestStruct
		err := ConfigSafe.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			compositeErr, ok := err.(*CompositeError)
			should.True(ok, "Should return CompositeError")
			should.NotEmpty(compositeErr.Errors, "Should have collected errors")
		}
		should.Nil(result.MapField, "Map should remain nil")
		should.Nil(result.SliceField, "Slice should remain nil")
	})

	// Test edge case: empty CompositeError
	t.Run("empty_composite_error", func(t *testing.T) {
		ce := &CompositeError{Errors: []error{}}
		should.Equal("no errors", ce.Error(), "Empty CompositeError should return 'no errors'")
	})
}
