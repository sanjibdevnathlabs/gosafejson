package gosafejson

import (
	"fmt"
	"os"
	"strings"
)

func ExampleMarshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleUnmarshal() {
	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

func ExampleConfigFastest_Marshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	stream := ConfigFastest.BorrowStream(nil)
	defer ConfigFastest.ReturnStream(stream)
	stream.WriteVal(group)
	if stream.Error != nil {
		fmt.Println("error:", stream.Error)
	}
	_, _ = os.Stdout.Write(stream.Buffer())
	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleConfigFastest_Unmarshal() {
	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	iter := ConfigFastest.BorrowIterator(jsonBlob)
	defer ConfigFastest.ReturnIterator(iter)
	iter.ReadVal(&animals)
	if iter.Error != nil {
		fmt.Println("error:", iter.Error)
	}
	fmt.Printf("%+v", animals)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

func ExampleGet() {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	fmt.Printf("%s", Get(val, "Colors", 0).ToString())
	// Output:
	// Crimson
}

func ExampleMyKey() {
	hello := MyKey("hello")
	output, err := Marshal(map[*MyKey]string{&hello: "world"})
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println(string(output))
	obj := map[*MyKey]string{}
	err = Unmarshal(output, &obj)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}
	for k, v := range obj {
		fmt.Println(*k, v)
	}
	// Output:
	// {"Hello":"world"}
	// Hel world
}

type MyKey string

func (m *MyKey) MarshalText() ([]byte, error) {
	return []byte(strings.ReplaceAll(string(*m), "h", "H")), nil
}

func (m *MyKey) UnmarshalText(text []byte) error {
	*m = MyKey(text[:3])
	return nil
}

// This example demonstrates using safe unmarshalling to handle type mismatches gracefully.
func ExampleConfigSafe_Unmarshal() {
	// JSON with type mismatches: age should be an int, metadata should be a map
	var jsonBlob = []byte(`{
		"user_id": "user123",
		"email": "user@example.com",
		"age": "thirty",
		"metadata": ["item1", "item2"],
		"tags": "not-an-array"
	}`)

	type UserProfile struct {
		UserID   string                 `json:"user_id"`
		Email    string                 `json:"email"`
		Age      int                    `json:"age"`
		Metadata map[string]interface{} `json:"metadata"`
		Tags     []string               `json:"tags"`
	}

	// Using standard unmarshalling
	var profile1 UserProfile
	err1 := ConfigCompatibleWithStandardLibrary.Unmarshal(jsonBlob, &profile1)
	fmt.Println("Standard unmarshalling:")
	fmt.Println("Error:", err1 != nil)

	// Using safe unmarshalling
	var profile2 UserProfile
	err2 := ConfigSafe.Unmarshal(jsonBlob, &profile2)
	fmt.Println("\nSafe unmarshalling:")
	fmt.Println("Error:", err2 != nil)

	// Check if it's a composite error
	if compErr, ok := err2.(*CompositeError); ok {
		fmt.Printf("Found %d errors\n", len(compErr.Errors))
	}

	// Even with errors, successfully parsed fields are available
	fmt.Println("Successfully parsed fields:")
	fmt.Printf("UserID: %s\n", profile2.UserID)
	fmt.Printf("Email: %s\n", profile2.Email)

	// Output:
	// Standard unmarshalling:
	// Error: true
	//
	// Safe unmarshalling:
	// Error: true
	// Found 3 errors
	// Successfully parsed fields:
	// UserID: user123
	// Email: user@example.com
}
