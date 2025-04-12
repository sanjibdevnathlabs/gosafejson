package any_tests

import (
	"testing"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

func Test_wrap_and_valuetype_everything(t *testing.T) {
	should := require.New(t)
	var i interface{}
	any := gosafejson.Get([]byte("123"))
	// default of number type is float64
	i = float64(123)
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(int8(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	//  get interface is not int8 interface
	// i = int8(10)
	// should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(int16(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	//i = int16(10)
	//should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(int32(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	i = int32(10)
	should.Equal(i, any.GetInterface())
	any = gosafejson.Wrap(int64(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	i = int64(10)
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(uint(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	// not equal
	//i = uint(10)
	//should.Equal(i, any.GetInterface())
	any = gosafejson.Wrap(uint8(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	// not equal
	// i = uint8(10)
	// should.Equal(i, any.GetInterface())
	any = gosafejson.Wrap(uint16(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	any = gosafejson.Wrap(uint32(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	i = uint32(10)
	should.Equal(i, any.GetInterface())
	any = gosafejson.Wrap(uint64(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	i = uint64(10)
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(float32(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	// not equal
	//i = float32(10)
	//should.Equal(i, any.GetInterface())
	any = gosafejson.Wrap(float64(10))
	should.Equal(any.ValueType(), gosafejson.NumberValue)
	should.Equal(any.LastError(), nil)
	i = float64(10)
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(true)
	should.Equal(any.ValueType(), gosafejson.BoolValue)
	should.Equal(any.LastError(), nil)
	i = true
	should.Equal(i, any.GetInterface())
	any = gosafejson.Wrap(false)
	should.Equal(any.ValueType(), gosafejson.BoolValue)
	should.Equal(any.LastError(), nil)
	i = false
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(nil)
	should.Equal(any.ValueType(), gosafejson.NilValue)
	should.Equal(any.LastError(), nil)
	i = nil
	should.Equal(i, any.GetInterface())

	stream := gosafejson.NewStream(gosafejson.ConfigDefault, nil, 32)
	any.WriteTo(stream)
	should.Equal("null", string(stream.Buffer()))
	should.Equal(any.LastError(), nil)

	any = gosafejson.Wrap(struct{ age int }{age: 1})
	should.Equal(any.ValueType(), gosafejson.ObjectValue)
	should.Equal(any.LastError(), nil)
	i = struct{ age int }{age: 1}
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap(map[string]interface{}{"abc": 1})
	should.Equal(any.ValueType(), gosafejson.ObjectValue)
	should.Equal(any.LastError(), nil)
	i = map[string]interface{}{"abc": 1}
	should.Equal(i, any.GetInterface())

	any = gosafejson.Wrap("abc")
	i = "abc"
	should.Equal(i, any.GetInterface())
	should.Equal(nil, any.LastError())

}
