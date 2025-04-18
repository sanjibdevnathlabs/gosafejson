package test

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

func Test_customize_type_decoder(t *testing.T) {
	t.Skip()
	gosafejson.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *gosafejson.Iterator) {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", iter.ReadString(), time.UTC)
		if err != nil {
			iter.Error = err
			return
		}
		*((*time.Time)(ptr)) = t
	})
	//defer gosafejson.ConfigDefault.(*frozenConfig).cleanDecoders()
	val := time.Time{}
	err := gosafejson.Unmarshal([]byte(`"2016-12-05 08:43:28"`), &val)
	if err != nil {
		t.Fatal(err)
	}
	year, month, day := val.Date()
	if year != 2016 || month != 12 || day != 5 {
		t.Fatal(val)
	}
}

func Test_customize_byte_array_encoder(t *testing.T) {
	t.Skip()
	//gosafejson.ConfigDefault.(*frozenConfig).cleanEncoders()
	should := require.New(t)
	gosafejson.RegisterTypeEncoderFunc("[]uint8", func(ptr unsafe.Pointer, stream *gosafejson.Stream) {
		t := *((*[]byte)(ptr))
		stream.WriteString(string(t))
	}, nil)
	//defer gosafejson.ConfigDefault.(*frozenConfig).cleanEncoders()
	val := []byte("abc")
	str, err := gosafejson.MarshalToString(val)
	should.Nil(err)
	should.Equal(`"abc"`, str)
}

type CustomEncoderAttachmentTestStruct struct {
	Value int32 `json:"value"`
}

type CustomEncoderAttachmentTestStructEncoder struct{}

func (c *CustomEncoderAttachmentTestStructEncoder) Encode(ptr unsafe.Pointer, stream *gosafejson.Stream) {
	attachVal, ok := stream.Attachment.(int)
	stream.WriteRaw(`"`)
	stream.WriteRaw(fmt.Sprintf("%t %d", ok, attachVal))
	stream.WriteRaw(`"`)
}

func (c *CustomEncoderAttachmentTestStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func Test_custom_encoder_attachment(t *testing.T) {

	gosafejson.RegisterTypeEncoder("test.CustomEncoderAttachmentTestStruct", &CustomEncoderAttachmentTestStructEncoder{})
	expectedValue := 17
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := gosafejson.NewStream(gosafejson.Config{SortMapKeys: true}.Froze(), buf, 4096)
	stream.Attachment = expectedValue
	val := map[string]CustomEncoderAttachmentTestStruct{"a": {}}
	stream.WriteVal(val)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("{\"a\":\"true 17\"}", buf.String())
}

type Tom struct {
	Field1 string
}

func Test_customize_field_decoder(t *testing.T) {
	defer gosafejson.TestingOnlyCleanDecoders()
	gosafejson.RegisterFieldDecoderFunc(reflect.TypeOf(Tom{}).String(), "Field1", func(ptr unsafe.Pointer, iter *gosafejson.Iterator) {
		*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
	})
	tom := Tom{}
	err := gosafejson.Unmarshal([]byte(`{"Field1": 100}`), &tom)
	if err != nil {
		t.Fatal(err)
	}
	should := require.New(t)
	should.Equal("100", tom.Field1)
}

func Test_recursive_empty_interface_customization(t *testing.T) {
	t.Skip()
	var obj interface{}
	gosafejson.RegisterTypeDecoderFunc("interface {}", func(ptr unsafe.Pointer, iter *gosafejson.Iterator) {
		switch iter.WhatIsNext() {
		case gosafejson.NumberValue:
			*(*interface{})(ptr) = iter.ReadInt64()
		default:
			*(*interface{})(ptr) = iter.Read()
		}
	})
	should := require.New(t)
	err := gosafejson.Unmarshal([]byte("[100]"), &obj)
	should.Nil(err)
	should.Equal([]interface{}{int64(100)}, obj)
}

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}

func Test_read_custom_interface(t *testing.T) {
	t.Skip()
	should := require.New(t)
	var val MyInterface
	gosafejson.RegisterTypeDecoderFunc("gosafejson.MyInterface", func(ptr unsafe.Pointer, iter *gosafejson.Iterator) {
		*((*MyInterface)(ptr)) = MyString(iter.ReadString())
	})
	err := gosafejson.UnmarshalFromString(`"hello"`, &val)
	should.Nil(err)
	should.Equal("hello", val.Hello())
}

const flow1 = `
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}`

const flow2 = `
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
{"A":"hello"}
`

type (
	Type1 struct {
		A string
	}

	Type2 struct {
		A string
	}
)

func (t *Type2) UnmarshalJSON(data []byte) error {
	return nil
}

func (t *Type2) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func TestType1NoFinalLF(t *testing.T) {
	reader := bytes.NewReader([]byte(flow1))
	dec := gosafejson.NewDecoder(reader)

	i := 0
	for dec.More() {
		data := &Type1{}
		if err := dec.Decode(data); err != nil {
			t.Errorf("at %v got %v", i, err)
		}
		i++
	}
}

func TestType1FinalLF(t *testing.T) {
	reader := bytes.NewReader([]byte(flow2))
	dec := gosafejson.NewDecoder(reader)

	i := 0
	for dec.More() {
		data := &Type1{}
		if err := dec.Decode(data); err != nil {
			t.Errorf("at %v got %v", i, err)
		}
		i++
	}
}

func TestType2NoFinalLF(t *testing.T) {
	reader := bytes.NewReader([]byte(flow1))
	dec := gosafejson.NewDecoder(reader)

	i := 0
	for dec.More() {
		data := &Type2{}
		if err := dec.Decode(data); err != nil {
			t.Errorf("at %v got %v", i, err)
		}
		i++
	}
}

func TestType2FinalLF(t *testing.T) {
	reader := bytes.NewReader([]byte(flow2))
	dec := gosafejson.NewDecoder(reader)

	i := 0
	for dec.More() {
		data := &Type2{}
		if err := dec.Decode(data); err != nil {
			t.Errorf("at %v got %v", i, err)
		}
		i++
	}
}
