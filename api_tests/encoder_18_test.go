//go:build go1.8
// +build go1.8

package test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/sanjibdevnathlabs/gosafejson"
	"github.com/stretchr/testify/require"
)

func Test_new_encoder(t *testing.T) {
	should := require.New(t)
	buf1 := &bytes.Buffer{}
	encoder1 := json.NewEncoder(buf1)
	encoder1.SetEscapeHTML(false)
	err := encoder1.Encode([]int{1})
	should.Nil(err)
	should.Equal("[1]\n", buf1.String())
	buf2 := &bytes.Buffer{}
	encoder2 := gosafejson.NewEncoder(buf2)
	encoder2.SetEscapeHTML(false)
	err = encoder2.Encode([]int{1})
	should.Nil(err)
	should.Equal("[1]\n", buf2.String())
}

func Test_string_encode_with_std_without_html_escape(t *testing.T) {
	api := gosafejson.Config{EscapeHTML: false}.Froze()
	should := require.New(t)
	for i := 0; i < utf8.RuneSelf; i++ {
		input := string([]byte{byte(i)})
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(input)
		should.Nil(err)
		stdOutput := buf.String()
		stdOutput = stdOutput[:len(stdOutput)-1]
		jsoniterOutputBytes, err := api.Marshal(input)
		should.Nil(err)
		jsoniterOutput := string(jsoniterOutputBytes)
		// Normalize standard library output to handle differences in control character escaping
		// Standard lib uses \b, \f while jsoniter uses \u0008, \u000c
		normalizedStdOutput := strings.ReplaceAll(stdOutput, "\\b", "\\u0008")
		normalizedStdOutput = strings.ReplaceAll(normalizedStdOutput, "\\f", "\\u000c")
		should.Equal(normalizedStdOutput, jsoniterOutput)
	}
}
