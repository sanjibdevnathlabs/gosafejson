package test

import (
	"bytes"
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

// Standard Encoder has trailing newline.
func TestEncoderHasTrailingNewline(t *testing.T) {
	should := require.New(t)
	var buf, stdbuf bytes.Buffer
	enc := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(&buf)
	err := enc.Encode(1)
	should.Nil(err)
	stdenc := json.NewEncoder(&stdbuf)
	err = stdenc.Encode(1)
	should.Nil(err)
	should.Equal(stdbuf.Bytes(), buf.Bytes())
}
