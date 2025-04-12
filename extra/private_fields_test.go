package extra

import (
	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_private_fields(t *testing.T) {
	type TestObject struct {
		field1 string
	}
	SupportPrivateFields()
	should := require.New(t)
	obj := TestObject{}
	should.Nil(gosafejson.UnmarshalFromString(`{"field1":"Hello"}`, &obj))
	should.Equal("Hello", obj.field1)
}
