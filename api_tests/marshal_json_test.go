package test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	Bar interface{}
}

func (f Foo) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f.Bar)
	return buf.Bytes(), err
}

// Standard Encoder has trailing newline.
func TestEncodeMarshalJSON(t *testing.T) {

	foo := Foo{
		Bar: 123,
	}
	should := require.New(t)
	var buf, stdbuf bytes.Buffer
	enc := gosafejson.ConfigCompatibleWithStandardLibrary.NewEncoder(&buf)
	err := enc.Encode(foo)
	should.Nil(err)
	stdenc := json.NewEncoder(&stdbuf)
	err = stdenc.Encode(foo)
	should.Nil(err)
	should.Equal(stdbuf.Bytes(), buf.Bytes())
}
