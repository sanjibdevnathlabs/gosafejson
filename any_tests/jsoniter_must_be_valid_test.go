package any_tests

import (
	"testing"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

// if must be valid is useless, just drop this test
func Test_must_be_valid(t *testing.T) {
	should := require.New(t)
	any := gosafejson.Get([]byte("123"))
	should.Equal(any.MustBeValid().ToInt(), 123)

	any = gosafejson.Wrap(int8(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(int16(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(int32(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(int64(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(uint(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(uint8(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(uint16(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(uint32(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(uint64(10))
	should.Equal(any.MustBeValid().ToInt(), 10)

	any = gosafejson.Wrap(float32(10))
	should.Equal(any.MustBeValid().ToFloat64(), float64(10))

	any = gosafejson.Wrap(float64(10))
	should.Equal(any.MustBeValid().ToFloat64(), float64(10))

	any = gosafejson.Wrap(true)
	should.Equal(any.MustBeValid().ToFloat64(), float64(1))

	any = gosafejson.Wrap(false)
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = gosafejson.Wrap(nil)
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = gosafejson.Wrap(struct{ age int }{age: 1})
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = gosafejson.Wrap(map[string]interface{}{"abc": 1})
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = gosafejson.Wrap("abc")
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = gosafejson.Wrap([]int{})
	should.Equal(any.MustBeValid().ToFloat64(), float64(0))

	any = gosafejson.Wrap([]int{1, 2})
	should.Equal(any.MustBeValid().ToFloat64(), float64(1))
}
