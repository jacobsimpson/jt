package listener

import (
	"regexp"
)

type ValueType int

const (
	StringValue ValueType = iota
	RegexpValue
	IntegerValue
	DateTimeValue
	UnknownValue
)

type Value interface {
	Type() ValueType
	Raw() string
	Value() interface{}
	String() string
}

//
// A Value implementation to hold a variable.
//
type varValue struct {
	name string
}

func (v *varValue) Type() ValueType {
	return UnknownValue
}

func (v *varValue) Raw() string {
	return v.name
}

func (v *varValue) Value() interface{} {
	return v.name
}

func (v *varValue) String() string {
	return v.name
}

//
// A Value implementation to hold a regular expression.
//
type regexpValue struct {
	raw string
	re  *regexp.Regexp
}

func NewRegexpValue(regexpString string) (Value, error) {
	re, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}
	return &regexpValue{
		raw: regexpString,
		re:  re,
	}, nil
}

func (v *regexpValue) Type() ValueType {
	return RegexpValue
}

func (v *regexpValue) Raw() string {
	return v.raw
}

func (v *regexpValue) Value() interface{} {
	return v.re
}

func (v *regexpValue) String() string {
	return v.raw
}

//
// A Value implementation to hold a string.
//
type stringValue struct {
	raw   string
	value string
}

func NewStringValue(s string) Value {
	return &stringValue{
		raw:   s,
		value: s,
	}
}

func (v *stringValue) Type() ValueType {
	return StringValue
}

func (v *stringValue) Raw() string {
	return v.raw
}

func (v *stringValue) Value() interface{} {
	return v.value
}

func (v *stringValue) String() string {
	return v.value
}
