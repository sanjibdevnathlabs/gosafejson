package jsoniter

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
