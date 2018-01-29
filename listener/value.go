package listener

import (
	"fmt"
	"regexp"
	"time"
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

//
// A Value implementation to hold a date/time.
//
type datetimeValue struct {
	raw   string
	value time.Time
}

func NewDateTimeValue(s string) (Value, error) {
	date, err := parseDateLiteral(s)
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

var dateLiteralFormats = []string{
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02T15",
	"2006-01-02T",
	"20060102T15:04:05.000Z",
	"20060102T15:04:05",
	"20060102T15:04",
	"20060102T15",
	"20060102T",
}

func parseDateLiteral(str string) (time.Time, error) {
	for _, layout := range dateLiteralFormats {
		if t, err := time.Parse(layout, str); err == nil {
			return t, nil
		} else {
		}
	}
	return time.Time{}, fmt.Errorf("Unable to convert %q to a date", str)
}
