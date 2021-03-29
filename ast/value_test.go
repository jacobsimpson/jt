package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntegerValueFromBinaryString(t *testing.T) {
	assert := assert.New(t)

	v, err := NewIntegerValueFromBinaryString("0b1000")

	assert.NoError(err)
	assert.Equal(v, &IntegerValue{raw: "1000", value: 8})
}
