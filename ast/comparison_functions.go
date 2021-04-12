package ast

import (
	"strconv"
	"strings"
	"time"

	"github.com/jacobsimpson/jt/datetime"
	"github.com/shopspring/decimal"
)

func lt(environment *Environment, left, right Expression) bool {
	left = resolveVar(environment, left)
	right = resolveVar(environment, right)

	switch l := left.(type) {
	case *AnyValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeGTAny(r, l)
		case *IntegerValue:
			return integerGTAny(r, l)
		case *DoubleValue:
			return doubleGTAny(r, l)
		case *StringValue:
			return stringGTAny(r, l)
		case *AnyValue:
			return anyGTAny(r, l)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeLTAny(l, r)
		}
	case *DoubleValue:
		switch r := right.(type) {
		case *AnyValue:
			return doubleLTAny(l, r)
		case *DoubleValue:
			return doubleLTDouble(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerLTAny(l, r)
		case *DoubleValue:
			return integerLTDouble(l, r)
		}
	case *StringValue:
		switch r := right.(type) {
		case *AnyValue:
			return stringLTAny(l, r)
		case *StringValue:
			return stringLTString(l, r)
		}
	}
	return false
}

func le(environment *Environment, left, right Expression) bool {
	left = resolveVar(environment, left)
	right = resolveVar(environment, right)

	switch l := left.(type) {
	case *AnyValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeGTAny(r, l) || dateTimeEQAny(r, l)
		case *IntegerValue:
			return integerGTAny(r, l) || integerEQAny(r, l)
		case *DoubleValue:
			return doubleGTAny(r, l) || doubleEQAny(r, l)
		case *StringValue:
			return stringGTAny(r, l) || stringEQAny(r, l)
		case *AnyValue:
			return anyGTAny(r, l) || anyEQAny(r, l)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeLTAny(l, r) || dateTimeEQAny(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerLTAny(l, r) || integerEQAny(l, r)
		}
	case *StringValue:
		switch r := right.(type) {
		case *AnyValue:
			return stringLTAny(l, r) || stringEQAny(l, r)
		}
	}
	return false
}

func eq(environment *Environment, left, right Expression) bool {
	left = resolveVar(environment, left)
	right = resolveVar(environment, right)

	switch l := left.(type) {
	case *RegexpValue:
		switch r := right.(type) {
		case *StringValue:
			return compareStringEQRegexp(r, l)
		case *AnyValue:
			return regexpEQAny(l, r)
		}
	case *AnyValue:
		switch r := right.(type) {
		case *RegexpValue:
			return regexpEQAny(r, l)
		case *DateTimeValue:
			return dateTimeEQAny(r, l)
		case *IntegerValue:
			return integerEQAny(r, l)
		case *DoubleValue:
			return doubleEQAny(r, l)
		case *StringValue:
			return stringEQAny(r, l)
		case *AnyValue:
			return anyEQAny(r, l)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeEQAny(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerEQAny(l, r)
		}
	case *StringValue:
		switch r := right.(type) {
		case *AnyValue:
			return stringEQAny(l, r)
		case *RegexpValue:
			return compareStringEQRegexp(l, r)
		case *StringValue:
			return compareStringEQString(l, r)
		}
	}
	return false
}

func ne(environment *Environment, left, right Expression) bool {
	return !eq(environment, left, right)
}

func ge(environment *Environment, left, right Expression) bool {
	left = resolveVar(environment, left)
	right = resolveVar(environment, right)

	switch l := left.(type) {
	case *AnyValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeLTAny(r, l) || dateTimeEQAny(r, l)
		case *IntegerValue:
			return integerLTAny(r, l) || integerEQAny(r, l)
		case *DoubleValue:
			return doubleLTAny(r, l) || doubleEQAny(r, l)
		case *StringValue:
			return stringLTAny(r, l) || stringEQAny(r, l)
		case *AnyValue:
			return anyLTAny(r, l) || anyEQAny(r, l)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeGTAny(l, r) || dateTimeEQAny(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerGTAny(l, r) || integerEQAny(l, r)
		}
	case *StringValue:
		switch r := right.(type) {
		case *AnyValue:
			return stringGTAny(l, r) || stringEQAny(l, r)
		}
	}
	return false
}

func gt(environment *Environment, left, right Expression) bool {
	left = resolveVar(environment, left)
	right = resolveVar(environment, right)

	switch l := left.(type) {
	case *AnyValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeLTAny(r, l)
		case *IntegerValue:
			return integerLTAny(r, l)
		case *DoubleValue:
			return doubleLTAny(r, l)
		case *StringValue:
			return stringLTAny(r, l)
		case *AnyValue:
			return anyGTAny(l, r)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeGTAny(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerGTAny(l, r)
		}
	case *StringValue:
		switch r := right.(type) {
		case *AnyValue:
			return stringGTAny(l, r)
		}
	}
	return false
}

func compareStringEQRegexp(s *StringValue, re *RegexpValue) bool {
	sv := s.value
	rev := re.re
	return rev.MatchString(sv)
}

func compareStringEQString(left *StringValue, right *StringValue) bool {
	return left.value == right.value
}

func regexpEQAny(re *RegexpValue, val *AnyValue) bool {
	rev := re.re
	return rev.MatchString(val.raw)
}

func dateTimeEQAny(dtValue *DateTimeValue, val *AnyValue) bool {
	dt := dtValue.value
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val.raw)
	if err != nil {
		return false
	}
	return dt.Equal(coerced)
}

func dateTimeLTAny(dtValue *DateTimeValue, val *AnyValue) bool {
	dt := dtValue.value
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val.raw)
	if err != nil {
		return false
	}
	return dt.Before(coerced)
}

func dateTimeGTAny(dtValue *DateTimeValue, val *AnyValue) bool {
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val.raw)
	if err != nil {
		return false
	}
	return dtValue.value.After(coerced)
}

func integerEQAny(iValue *IntegerValue, val *AnyValue) bool {
	i := iValue.value
	parsed, err := parseInt(val.raw)
	if err != nil {
		return false
	}
	return i == parsed
}

func integerLTAny(iValue *IntegerValue, val *AnyValue) bool {
	i := iValue.value
	parsed, err := parseInt(val.raw)
	if err != nil {
		return false
	}
	return i < parsed
}

func integerLTDouble(lhs *IntegerValue, rhs *DoubleValue) bool {
	return decimal.NewFromInt(lhs.value).LessThan(*rhs.value)
}

func integerGTAny(iValue *IntegerValue, val *AnyValue) bool {
	i := iValue.value
	parsed, err := parseInt(val.raw)
	if err != nil {
		return false
	}
	return i > parsed
}

func doubleEQAny(dValue *DoubleValue, val *AnyValue) bool {
	d := dValue.value
	parsed, err := decimal.NewFromString(val.raw)
	if err != nil {
		return false
	}
	return d.Equal(parsed)
}

func doubleLTAny(dValue *DoubleValue, val *AnyValue) bool {
	d := dValue.value
	parsed, err := decimal.NewFromString(val.raw)
	if err != nil {
		parsedInt, err := parseInt(val.raw)
		if err != nil {
			return false
		}
		parsed = decimal.New(parsedInt, 0)
	}
	return d.LessThan(parsed)
}

func doubleGTAny(dValue *DoubleValue, val *AnyValue) bool {
	d := dValue.value
	parsed, err := decimal.NewFromString(val.raw)
	if err != nil {
		parsedInt, err := parseInt(val.raw)
		if err != nil {
			return false
		}
		parsed = decimal.New(parsedInt, 0)
	}
	return d.GreaterThan(parsed)
}

func doubleLTDouble(lhs *DoubleValue, rhs *DoubleValue) bool {
	return lhs.value.LessThan(*rhs.value)
}

func stringEQAny(lValue *StringValue, val *AnyValue) bool {
	return lValue.value == val.raw
}

func stringLTAny(lValue *StringValue, val *AnyValue) bool {
	return lValue.value < val.raw
}

func stringGTAny(lValue *StringValue, val *AnyValue) bool {
	return lValue.value > val.raw
}

func stringLTString(lhs *StringValue, rhs *StringValue) bool {
	return lhs.value < rhs.value
}

func anyGTAny(lValue *AnyValue, val *AnyValue) bool {
	return lValue.raw > val.raw
}

func anyEQAny(lValue *AnyValue, val *AnyValue) bool {
	return lValue.raw == val.raw
}

func anyLTAny(lValue *AnyValue, val *AnyValue) bool {
	return lValue.raw < val.raw
}

func parseInt(s string) (int64, error) {
	s = strings.Map(func(r rune) rune {
		if r == '_' {
			return -1
		}
		return r
	}, s)
	if len(s) > 2 {
		if strings.HasPrefix(s, "0x") {
			return strconv.ParseInt(s, 0, 64)
		}
	}
	if len(s) > 2 {
		if strings.HasPrefix(s, "0b") {
			return strconv.ParseInt(s[2:], 2, 64)
		}
	}
	return strconv.ParseInt(s, 10, 64)
}

func resolveVar(environment *Environment, v Expression) Expression {
	// TODO: This is going to crash hard if the variable doesn't exist.
	switch vr := v.(type) {
	case *VarValue:
		return environment.Resolve(vr)
	case *KeywordValue:
		switch vr.value {
		case "yesterday":
			t := time.Now().Add(-24 * time.Hour)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
			return &DateTimeValue{
				raw:   t.String(),
				value: t,
			}
		case "today":
			t := time.Now()
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
			return &DateTimeValue{
				raw:   t.String(),
				value: t,
			}
		case "now":
			t := time.Now()
			return &DateTimeValue{
				raw:   t.String(),
				value: t,
			}
		case "tomorrow":
			t := time.Now().Add(24 * time.Hour)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
			return &DateTimeValue{
				raw:   t.String(),
				value: t,
			}
		}
	}
	return v
}
