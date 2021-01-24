package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCodec(t *testing.T) {
	RegisterCodec("testCodec", nil)

	codec := GetCodec("testCodec")
	assert.Equal(t, codec, nil)
}


func TestDefaultCodec_Encode_Decode(t *testing.T) {
	str := "hello,world!"
	bytes, err := DefaultCodec.Encode([]byte(str))
	assert.Nil(t, err)
	bytes, err = DefaultCodec.Decode(bytes)
	assert.Nil(t, err)
	assert.Equal(t, str, string(bytes))
}
