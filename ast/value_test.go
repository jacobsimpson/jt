package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntegerValueFromBinaryString(t *testing.T) {
	assert := assert.New(t)

	v, err := NewIntegerValueFromBinaryString("0b1000")

	assert.NoError(err)
	assert.Equal(v, &IntegerValue{raw: "0b1000", value: 8})
}

func TestAllValuesAreExpressions(t *testing.T) {
	assert := assert.New(t)

	var e Expression

	e = NewVarValue("abc")
	r, err := NewRegexpValue("/abc/")
	assert.NoError(err)
	e = r

	e = NewStringValue("abc")

	d, err := NewDateTimeValue("2006-03-04T")
	assert.NoError(err)
	e = d

	i, err := NewIntegerValueFromBinaryString("0b1")
	assert.NoError(err)
	e = i

	dbl, err := NewDoubleFromString("12.12")
	assert.NoError(err)
	e = dbl

	e = &AnyValue{"abc"}

	assert.NotNil(e)
}
