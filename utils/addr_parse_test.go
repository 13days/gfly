package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServicePath(t *testing.T) {
	_, _, err := ParseServicePath("hello.world")
	assert.NotNil(t, err)

	_, _, err = ParseServicePath("hello/world")
	assert.NotNil(t, err)

	seriveName, method, err := ParseServicePath("/hello/world")
	assert.Equal(t, seriveName, "hello")
	assert.Equal(t, method, "world")
	assert.Equal(t, err, nil)
}
