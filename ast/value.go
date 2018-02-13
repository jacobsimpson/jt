package ast

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jacobsimpson/jt/datetime"
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
func NewVarValue(name string) Value {
	return &varValue{
		name: name,
	}
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

//
// A Value implementation to hold a date/time.
//
type datetimeValue struct {
	raw   string
	value *time.Time
}

func NewDateTimeValue(s string) (Value, error) {
	date, err := datetime.ParseDateTime(datetime.LiteralFormats, s)
	if err != nil {
		return nil, err
	}
	return &datetimeValue{
		raw:   s,
		value: date,
	}, nil
}

func (v *datetimeValue) Type() ValueType {
	return DateTimeValue
}

func (v *datetimeValue) Raw() string {
	return v.raw
}

func (v *datetimeValue) Value() interface{} {
	return v.value
}

func (v *datetimeValue) String() string {
	return v.value.String()
}

//
// A Value implementation to hold a integer.
//
type integerValue struct {
	raw   string
	value int64
}

func NewIntegerValueFromBinaryString(s string) (Value, error) {
	s = s[2:]
	// '_' characters are allowed in integer representations to improve
	// readability, but they have no other purpose and are stripped here to
	// allow parsing.
	s = strings.Map(func(r rune) rune {
		if r == '_' {
			return -1
		}
		return r
	}, s)
	v, err := parseIntFromString(s, 2)
	return v, err
}

func NewIntegerValueFromHexString(s string) (Value, error) {
	s = s[2:]
	// '_' characters are allowed in integer representations to improve
	// readability, but they have no other purpose and are stripped here to
	// allow parsing.
	s = strings.Map(func(r rune) rune {
		if r == '_' {
			return -1
		}
		return r
	}, s)
	return parseIntFromString(s, 16)
}

func NewIntegerValue(s string) (Value, error) {
	// '_' characters are allowed in integer representations to improve
	// readability, but they have no other purpose and are stripped here to
	// allow parsing.
	s = strings.Map(func(r rune) rune {
		if r == '_' {
			return -1
		}
		return r
	}, s)
	return parseIntFromString(s, 10)
}

func parseIntFromString(s string, base int) (Value, error) {
	i, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return nil, err
	}
	return &integerValue{
		raw:   s,
		value: i,
	}, nil
}

func (v *integerValue) Type() ValueType {
	return IntegerValue
}

func (v *integerValue) Raw() string {
	return v.raw
}

func (v *integerValue) Value() interface{} {
	return v.value
}

func (v *integerValue) String() string {
	return fmt.Sprintf("%d", v.value)
}
