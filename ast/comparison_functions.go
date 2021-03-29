package ast

import (
	"strconv"
	"strings"

	"github.com/jacobsimpson/jt/datetime"
	"github.com/jacobsimpson/jt/debug"
	"github.com/shopspring/decimal"
)

func lt(environment map[string]string, left, right Value) bool {
	switch l := left.(type) {
	case *VarValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeGTAny(r, resolveVar(environment, left.Value()))
		case *IntegerValue:
			return integerGTAny(r, resolveVar(environment, left.Value()))
		case *DoubleValue:
			return doubleGTAny(r, resolveVar(environment, left.Value()))
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeLTAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *AnyValue:
			return integerLTAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	default:
		return false
	}
	return false
}

func le(environment map[string]string, left, right Value) bool {
	switch l := left.(type) {
	case *VarValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeGTAny(r, resolveVar(environment, left.Value())) ||
				dateTimeEQAny(r, resolveVar(environment, left.Value()))
		case *IntegerValue:
			return integerGTAny(r, resolveVar(environment, left.Value())) ||
				integerEQAny(r, resolveVar(environment, left.Value()))
		case *DoubleValue:
			return doubleGTAny(r, resolveVar(environment, left.Value())) ||
				doubleEQAny(r, resolveVar(environment, left.Value()))
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *VarValue:
			return dateTimeLTAny(l, resolveVar(environment, right.Value())) ||
				dateTimeEQAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *VarValue:
			return integerLTAny(l, resolveVar(environment, right.Value())) ||
				integerEQAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	default:
		return false
	}
	return false
}

func eq(environment map[string]string, left, right Value) bool {
	switch l := left.(type) {
	case *RegexpValue:
		switch r := right.(type) {
		case *StringValue:
			return compareStringEQRegexp(r, l)
		case *RegexpValue:
			return regexpEQUnknown(environment, l, right.Value())
		default:
			return false
		}
	case *StringValue:
		switch r := right.(type) {
		case *RegexpValue:
			return compareStringEQRegexp(l, r)
		case *StringValue:
			return compareStringEQString(l, r)
		default:
			// Error: "Can not compare %s to %s using %s",
			//     left.Type(), right.Type(), operator
			return false
		}
	case *VarValue:
		switch r := right.(type) {
		case *RegexpValue:
			return regexpEQUnknown(environment, r, left.Value())
		case *DateTimeValue:
			return dateTimeEQAny(r, resolveVar(environment, left.Value()))
		case *IntegerValue:
			return integerEQAny(r, resolveVar(environment, left.Value()))
		case *DoubleValue:
			return doubleEQAny(r, resolveVar(environment, left.Value()))
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeEQAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *IntegerValue:
			return integerEQAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	default:
		return false
	}
	return false
}

func ne(environment map[string]string, left, right Value) bool {
	return !eq(environment, left, right)
}

func ge(environment map[string]string, left, right Value) bool {
	switch l := left.(type) {
	case *VarValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeLTAny(r, resolveVar(environment, left.Value())) ||
				dateTimeEQAny(r, resolveVar(environment, left.Value()))
		case *IntegerValue:
			return integerLTAny(r, resolveVar(environment, left.Value())) ||
				integerEQAny(r, resolveVar(environment, left.Value()))
		case *DoubleValue:
			return doubleLTAny(r, resolveVar(environment, left.Value())) ||
				doubleEQAny(r, resolveVar(environment, left.Value()))
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeGTAny(l, resolveVar(environment, right.Value())) ||
				dateTimeEQAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *AnyValue:
			return integerGTAny(l, resolveVar(environment, right.Value())) ||
				integerEQAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	default:
		return false
	}
	return false
}

func gt(environment map[string]string, left, right Value) bool {
	switch l := left.(type) {
	case *VarValue:
		switch r := right.(type) {
		case *DateTimeValue:
			return dateTimeLTAny(r, resolveVar(environment, left.Value()))
		case *IntegerValue:
			return integerLTAny(r, resolveVar(environment, left.Value()))
		case *DoubleValue:
			return doubleLTAny(r, resolveVar(environment, left.Value()))
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeGTAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *AnyValue:
			return integerGTAny(l, resolveVar(environment, right.Value()))
		default:
			return false
		}
	default:
		return false
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

func regexpEQUnknown(environment map[string]string, re *RegexpValue, s interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := s.(string)
	sv := environment[varName]
	rev := re.re
	debug.Info("comparing %s (%s) to %s", varName, sv, rev)
	return rev.MatchString(sv)
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

func resolveVar(environment map[string]string, v interface{}) *AnyValue {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	return &AnyValue{val}
}
