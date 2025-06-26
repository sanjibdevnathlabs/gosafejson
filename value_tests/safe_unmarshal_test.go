package test

import (
	"fmt"
	"testing"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

type TestMismatchStruct struct {
	ID   *string                `json:"id"`
	Data map[string]interface{} `json:"data"`
	Name *string                `json:"name"`
}

type ComplexStruct struct {
	ID        *int                `json:"id"`
	Count     float64             `json:"count"`
	IsActive  bool                `json:"is_active"`
	Tags      []string            `json:"tags"`
	Metadata  map[string]string   `json:"metadata"`
	UserInfo  map[string]int      `json:"user_info"`
	Relations map[string][]string `json:"relations"`
}

func Test_SafeUnmarshal(t *testing.T) {
	should := require.New(t)

	// Test case with type mismatch - expecting map but getting array
	jsonStr := `{"id":"12345", "data": [{"a":"b"}, {"c":"d"}], "name": "Random"}`

	// Standard unmarshal should fail
	var standardResult TestMismatchStruct
	err := gosafejson.Unmarshal([]byte(jsonStr), &standardResult)
	should.Error(err)
	t.Logf("Standard unmarshal error: %v", err)

	// Safe unmarshal should continue and parse what it can
	var safeResult TestMismatchStruct
	err = gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	should.Error(err) // Should still return error but continue processing
	t.Logf("Safe unmarshal error: %v", err)
	t.Logf("ID: %p, Name: %p, Data: %v", safeResult.ID, safeResult.Name, safeResult.Data)

	// ID and Name should be parsed correctly
	should.NotNil(safeResult.ID)
	should.NotNil(safeResult.Name)
	should.Equal("12345", *safeResult.ID)
	should.Equal("Random", *safeResult.Name)

	// Test with valid JSON to ensure normal operation works
	validJsonStr := `{"id":"12345", "data": {"key1":"value1", "key2":"value2"}, "name": "Random"}`
	var validResult TestMismatchStruct
	err = gosafejson.ConfigSafe.Unmarshal([]byte(validJsonStr), &validResult)
	should.NoError(err)
	t.Logf("Valid JSON unmarshaling result: ID=%p, Name=%p, Data=%v", validResult.ID, validResult.Name, validResult.Data)
}

func Test_SafeUnmarshalMultipleMismatches(t *testing.T) {
	should := require.New(t)

	// JSON with multiple type mismatches
	jsonStr := `{
		"id": "not-an-int",
		"count": "not-a-float",
		"is_active": 123,
		"tags": {"not": "an-array"},
		"metadata": ["not", "a", "map"],
		"user_info": {"a": "not-an-int"},
		"relations": {"key": "not-an-array"}
	}`

	// Standard unmarshal should fail early
	var standardResult ComplexStruct
	err := gosafejson.Unmarshal([]byte(jsonStr), &standardResult)
	should.Error(err)
	t.Logf("Standard unmarshal error: %v", err)

	// Safe unmarshal should collect multiple errors
	var safeResult ComplexStruct
	err = gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	should.Error(err)
	t.Logf("Safe unmarshal error: %v", err)

	// Check that it's a CompositeError
	compositeErr, ok := err.(*gosafejson.CompositeError)
	should.True(ok, "Expected CompositeError")
	t.Logf("Error type: %T", err)

	// Should have multiple errors collected
	should.True(len(compositeErr.Errors) > 1, "Should have multiple errors")
}

// Test safe unmarshal with primitive type mismatches
func Test_SafeUnmarshalPrimitiveTypes(t *testing.T) {
	should := require.New(t)

	type PrimitiveStruct struct {
		IntField    int     `json:"int_field"`
		StringField string  `json:"string_field"`
		FloatField  float64 `json:"float_field"`
		BoolField   bool    `json:"bool_field"`
	}

	// JSON with all primitive type mismatches
	jsonStr := `{
		"int_field": "not-an-int",
		"string_field": 123,
		"float_field": true,
		"bool_field": "not-a-bool"
	}`

	var safeResult PrimitiveStruct
	err := gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	t.Logf("Primitive types result: %+v, error: %v", safeResult, err)

	// Safe unmarshal should succeed by skipping invalid values
	should.NoError(err, "Safe unmarshal should succeed by skipping invalid values")

	// All fields should remain at their zero values since all had type mismatches
	should.Equal(0, safeResult.IntField)
	should.Equal("", safeResult.StringField)
	should.Equal(0.0, safeResult.FloatField)
	should.Equal(false, safeResult.BoolField)
}

// Test safe unmarshal with nested structures
func Test_SafeUnmarshalNestedStructures(t *testing.T) {
	should := require.New(t)

	type NestedStruct struct {
		Value string `json:"value"`
	}

	type ParentStruct struct {
		ID     int          `json:"id"`
		Nested NestedStruct `json:"nested"`
		Items  []string     `json:"items"`
	}

	// JSON with nested structure type mismatch
	jsonStr := `{
		"id": 123,
		"nested": "should-be-object",
		"items": {"should": "be-array"}
	}`

	var safeResult ParentStruct
	err := gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	should.Error(err)

	// ID should be parsed correctly
	should.Equal(123, safeResult.ID)

	// Nested should remain at zero value
	should.Equal("", safeResult.Nested.Value)

	// Items should remain empty
	should.Empty(safeResult.Items)
}

// Test safe unmarshal with arrays and slices
func Test_SafeUnmarshalArraysAndSlices(t *testing.T) {
	should := require.New(t)

	type ArrayStruct struct {
		Numbers []int         `json:"numbers"`
		Strings []string      `json:"strings"`
		Mixed   []interface{} `json:"mixed"`
	}

	// JSON with array type mismatches
	jsonStr := `{
		"numbers": "not-an-array",
		"strings": 123,
		"mixed": {"not": "an-array"}
	}`

	var safeResult ArrayStruct
	err := gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	should.Error(err)

	// All arrays should remain empty
	should.Empty(safeResult.Numbers)
	should.Empty(safeResult.Strings)
	should.Empty(safeResult.Mixed)
}

// Test safe unmarshal with maps
func Test_SafeUnmarshalMaps(t *testing.T) {
	should := require.New(t)

	type MapStruct struct {
		StringMap map[string]string      `json:"string_map"`
		IntMap    map[string]int         `json:"int_map"`
		AnyMap    map[string]interface{} `json:"any_map"`
	}

	// JSON with map type mismatches
	jsonStr := `{
		"string_map": ["not", "a", "map"],
		"int_map": "not-a-map",
		"any_map": 123
	}`

	var safeResult MapStruct
	err := gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	should.Error(err)

	// All maps should remain empty/nil
	should.Empty(safeResult.StringMap)
	should.Empty(safeResult.IntMap)
	should.Empty(safeResult.AnyMap)
}

// Test CompositeError functionality
func Test_CompositeError(t *testing.T) {
	should := require.New(t)

	// Create a CompositeError manually
	err1 := fmt.Errorf("first error")
	err2 := fmt.Errorf("second error")

	composite := &gosafejson.CompositeError{
		Errors: []error{err1, err2},
	}

	// Test Error() method
	errStr := composite.Error()
	should.Contains(errStr, "2 errors occurred")
	should.Contains(errStr, "first error")
	should.Contains(errStr, "second error")

	// Test that it implements error interface
	var err error = composite
	should.NotNil(err)
}

// Test safe unmarshal with partial success - good values parsed, bad values skipped
func Test_SafeUnmarshalPartialSuccess(t *testing.T) {
	should := require.New(t)

	type MixedStruct struct {
		GoodInt    int    `json:"good_int"`
		BadInt     int    `json:"bad_int"`
		GoodString string `json:"good_string"`
		BadString  string `json:"bad_string"`
		GoodBool   bool   `json:"good_bool"`
		BadBool    bool   `json:"bad_bool"`
	}

	// JSON with some good and some bad values
	jsonStr := `{
		"good_int": 42,
		"bad_int": "not-an-int",
		"good_string": "hello",
		"bad_string": 123,
		"good_bool": true,
		"bad_bool": "not-a-bool"
	}`

	var safeResult MixedStruct
	err := gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &safeResult)
	t.Logf("Mixed types result: %+v, error: %v", safeResult, err)

	// Safe unmarshal should succeed by parsing good values and skipping bad ones
	should.NoError(err, "Safe unmarshal should succeed with partial success")

	// Good values should be parsed correctly
	should.Equal(42, safeResult.GoodInt)
	should.Equal("hello", safeResult.GoodString)
	should.Equal(true, safeResult.GoodBool)

	// Bad values should remain at zero values
	should.Equal(0, safeResult.BadInt)
	should.Equal("", safeResult.BadString)
	should.Equal(false, safeResult.BadBool)
}
