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

type regexpValue struct {
	raw string
	re  *regexp.Regexp
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
