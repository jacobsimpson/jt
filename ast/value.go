package ast

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jacobsimpson/jt/datetime"
	"github.com/shopspring/decimal"
)

type Value interface {
	Raw() string
	Value() interface{}
	String() string
}

//
// A Value implementation to hold a variable.
//
func NewVarValue(name string) Value {
	return &VarValue{
		name: name,
	}
}

type VarValue struct {
	name string
}

func (v *VarValue) Raw() string {
	return v.name
}

func (v *VarValue) Value() interface{} {
	return v.name
}

func (v *VarValue) String() string {
	return v.name
}

//
// A Value implementation to hold a regular expression.
//
type RegexpValue struct {
	raw string
	re  *regexp.Regexp
}

func NewRegexpValue(regexpString string) (Value, error) {
	re, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}
	return &RegexpValue{
		raw: regexpString,
		re:  re,
	}, nil
}

func (v *RegexpValue) Raw() string {
	return v.raw
}

func (v *RegexpValue) Value() interface{} {
	return v.re
}

func (v *RegexpValue) String() string {
	return v.raw
}

//
// A Value implementation to hold a string.
//
type StringValue struct {
	raw   string
	value string
}

func NewStringValue(s string) Value {
	return &StringValue{
		raw:   s,
		value: s,
	}
}

func (v *StringValue) Raw() string {
	return v.raw
}

func (v *StringValue) Value() interface{} {
	return v.value
}

func (v *StringValue) String() string {
	return v.value
}

//
// A Value implementation to hold a date/time.
//
type DateTimeValue struct {
	raw   string
	value time.Time
}

func NewDateTimeValue(s string) (Value, error) {
	date, err := datetime.ParseDateTime(datetime.LiteralFormats, s)
	if err != nil {
		return nil, err
	}
	return &DateTimeValue{
		raw:   s,
		value: date,
	}, nil
}

func (v *DateTimeValue) Raw() string {
	return v.raw
}

func (v *DateTimeValue) Value() interface{} {
	return v.value
}

func (v *DateTimeValue) String() string {
	return v.value.String()
}

//
// A Value implementation to hold a integer.
//
type IntegerValue struct {
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
	return &IntegerValue{
		raw:   s,
		value: i,
	}, nil
}

func (v *IntegerValue) Raw() string {
	return v.raw
}

func (v *IntegerValue) Value() interface{} {
	return v.value
}

func (v *IntegerValue) String() string {
	return fmt.Sprintf("%d", v.value)
}

//
// A Value implementation to hold a double.
//
type DoubleValue struct {
	raw   string
	value *decimal.Decimal
}

func NewDoubleFromString(s string) (Value, error) {
	// '_' characters are allowed in decimal representations to improve
	// readability, but they have no other purpose and are stripped here to
	// allow parsing.
	s = strings.Map(func(r rune) rune {
		if r == '_' {
			return -1
		}
		return r
	}, s)
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &DoubleValue{
		raw:   s,
		value: &d,
	}, nil
}

func (v *DoubleValue) Raw() string {
	return v.raw
}

func (v *DoubleValue) Value() interface{} {
	return v.value
}

func (v *DoubleValue) String() string {
	return v.value.String()
}

//
// AnyValue is a Value implementation to hold a value that is, as yet,
// typeless. A value taken from a column is typeless. It is a string, but it
// wasn't specified as a string by the programmer, so coercing it to a
// different type, depending on parseability, is legal.
//
type AnyValue struct {
	raw string
}

func (v *AnyValue) Raw() string {
	return v.raw
}

func (v *AnyValue) Value() interface{} {
	return v.raw
}

func (v *AnyValue) String() string {
	return v.raw
}
