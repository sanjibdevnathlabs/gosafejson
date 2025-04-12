package misc_tests

import (
	"bytes"
	"io"
	"testing"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

func Test_read_null(t *testing.T) {
	should := require.New(t)
	iter := gosafejson.ParseString(gosafejson.ConfigDefault, `null`)
	should.True(iter.ReadNil())
	iter = gosafejson.ParseString(gosafejson.ConfigDefault, `null`)
	should.Nil(iter.Read())
	iter = gosafejson.ParseString(gosafejson.ConfigDefault, `navy`)
	iter.Read()
	should.True(iter.Error != nil && iter.Error != io.EOF)
	iter = gosafejson.ParseString(gosafejson.ConfigDefault, `navy`)
	iter.ReadNil()
	should.True(iter.Error != nil && iter.Error != io.EOF)
}

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := gosafejson.NewStream(gosafejson.ConfigDefault, buf, 4096)
	stream.WriteNil()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("null", buf.String())
}

func Test_decode_null_object_field(t *testing.T) {
	should := require.New(t)
	iter := gosafejson.ParseString(gosafejson.ConfigDefault, `[null,"a"]`)
	iter.ReadArray()
	if iter.ReadObject() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
	type TestObject struct {
		Field string
	}
	objs := []TestObject{}
	should.Nil(gosafejson.UnmarshalFromString("[null]", &objs))
	should.Len(objs, 1)
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := gosafejson.ParseString(gosafejson.ConfigDefault, `[null,"a"]`)
	should.True(iter.ReadArray())
	should.True(iter.ReadNil())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := gosafejson.ParseString(gosafejson.ConfigDefault, `[null,"a"]`)
	should.True(iter.ReadArray())
	should.Equal("", iter.ReadString())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_skip(t *testing.T) {
	iter := gosafejson.ParseString(gosafejson.ConfigDefault, `[null,"a"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}
