package test

import (
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
	// Test case with a type mismatch (data should be map but is array)
	jsonStr := `{"id":"12345", "data": [{"a":"b"}, {"c":"d"}], "name": "Random"}`

	// Standard unmarshal should fail
	var details1 TestMismatchStruct
	err1 := gosafejson.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(jsonStr, &details1)
	require.Error(t, err1, "Standard unmarshal should fail")
	t.Logf("Standard unmarshal error: %v", err1)

	// Safe unmarshal should return errors but still decode valid fields
	var details2 TestMismatchStruct
	err2 := gosafejson.ConfigSafe.UnmarshalFromString(jsonStr, &details2)
	require.Error(t, err2, "Safe unmarshal should return error")
	t.Logf("Safe unmarshal error: %v", err2)
	t.Logf("ID: %v, Name: %v, Data: %v", details2.ID, details2.Name, details2.Data)

	// ID should be set correctly (it's before the error)
	require.NotNil(t, details2.ID, "ID should not be nil")
	require.Equal(t, "12345", *details2.ID, "ID should have correct value")

	// Test with a successful case
	jsonStr2 := `{"id":"67890", "data": {"key1": "value1", "key2": "value2"}, "name": "Random2"}`

	var details3 TestMismatchStruct
	err3 := gosafejson.ConfigSafe.UnmarshalFromString(jsonStr2, &details3)
	t.Logf("Valid JSON unmarshaling result: ID=%v, Name=%v, Data=%v", details3.ID, details3.Name, details3.Data)
	if err3 != nil {
		t.Logf("Unexpected error with valid JSON: %v", err3)
	}
	require.NoError(t, err3, "Safe unmarshal should not return error for valid JSON")
	require.NotNil(t, details3.ID, "ID should be unmarshalled correctly")
	require.Equal(t, "67890", *details3.ID, "ID should have correct value")
	require.NotNil(t, details3.Data, "Data should be unmarshalled correctly")
	if details3.Data != nil {
		value, ok := details3.Data["key1"]
		require.True(t, ok, "Data should contain 'key1'")
		require.Equal(t, "value1", value, "Data should contain correct values")
	}
	require.NotNil(t, details3.Name, "Name should be unmarshalled correctly")
	require.Equal(t, "Random2", *details3.Name, "Name should have correct value")
}

func Test_SafeUnmarshalMultipleMismatches(t *testing.T) {
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

	// Standard unmarshal should fail on first mismatch
	var obj1 ComplexStruct
	err1 := gosafejson.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(jsonStr, &obj1)
	require.Error(t, err1, "Standard unmarshal should fail")
	t.Logf("Standard unmarshal error: %v", err1)

	// Safe unmarshal should continue and collect all errors
	var obj2 ComplexStruct
	err2 := gosafejson.ConfigSafe.UnmarshalFromString(jsonStr, &obj2)
	require.Error(t, err2, "Safe unmarshal should return error")
	t.Logf("Safe unmarshal error: %v", err2)

	// Print error type
	t.Logf("Error type: %T", err2)

	// Check if it's a CompositeError
	compErr, ok := err2.(*gosafejson.CompositeError)
	require.True(t, ok, "Error should be of type CompositeError")
	require.NotEmpty(t, compErr.Errors, "CompositeError should contain errors")
}
