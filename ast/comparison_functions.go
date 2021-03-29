package ast

import (
	"regexp"
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
			return dateTimeGTUnknown(environment, r, left.Value())
		case *IntegerValue:
			return integerGTUnknown(environment, r, left.Value())
		case *DoubleValue:
			return doubleGTUnknown(environment, r, left.Value())
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeLTUnknown(environment, l, right.Value())
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *AnyValue:
			return integerLTUnknown(environment, l, right.Value())
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
			return dateTimeGTUnknown(environment, r, left.Value()) ||
				dateTimeEQUnknown(environment, r, left.Value())
		case *IntegerValue:
			return integerGTUnknown(environment, r, left.Value()) ||
				integerEQUnknown(environment, r, left.Value())
		case *DoubleValue:
			return doubleGTUnknown(environment, r, left.Value()) ||
				doubleEQUnknown(environment, r, left.Value())
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *VarValue:
			return dateTimeLTUnknown(environment, l, right.Value()) ||
				dateTimeEQUnknown(environment, l, right.Value())
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *VarValue:
			return integerLTUnknown(environment, l, right.Value()) ||
				integerEQUnknown(environment, l, right.Value())
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
		switch right.(type) {
		case *StringValue:
			return compareStringEQRegexp(right.Value(), left.Value())
		case *RegexpValue:
			return regexpEQUnknown(environment, left.Value(), right.Value())
		default:
			return false
		}
	case *StringValue:
		switch right.(type) {
		case *RegexpValue:
			return compareStringEQRegexp(left.Value(), right.Value())
		case *StringValue:
			return compareStringEQString(left.Value(), right.Value())
		default:
			// Error: "Can not compare %s to %s using %s",
			//     left.Type(), right.Type(), operator
			return false
		}
	case *VarValue:
		switch r := right.(type) {
		case *RegexpValue:
			return regexpEQUnknown(environment, right.Value(), left.Value())
		case *DateTimeValue:
			return dateTimeEQUnknown(environment, r, left.Value())
		case *IntegerValue:
			return integerEQUnknown(environment, r, left.Value())
		case *DoubleValue:
			return doubleEQUnknown(environment, r, left.Value())
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeEQUnknown(environment, l, right.Value())
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *IntegerValue:
			return integerEQUnknown(environment, l, right.Value())
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
			return dateTimeLTUnknown(environment, r, left.Value()) ||
				dateTimeEQUnknown(environment, r, left.Value())
		case *IntegerValue:
			return integerLTUnknown(environment, r, left.Value()) ||
				integerEQUnknown(environment, r, left.Value())
		case *DoubleValue:
			return doubleLTUnknown(environment, right.Value(), left.Value()) ||
				doubleEQUnknown(environment, r, left.Value())
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeGTUnknown(environment, l, right.Value()) ||
				dateTimeEQUnknown(environment, l, right.Value())
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *AnyValue:
			return integerGTUnknown(environment, l, right.Value()) ||
				integerEQUnknown(environment, l, right.Value())
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
			return dateTimeLTUnknown(environment, r, left.Value())
		case *IntegerValue:
			return integerLTUnknown(environment, r, left.Value())
		case *DoubleValue:
			return doubleLTUnknown(environment, right.Value(), left.Value())
		default:
			return false
		}
	case *DateTimeValue:
		switch right.(type) {
		case *AnyValue:
			return dateTimeGTUnknown(environment, l, right.Value())
		default:
			return false
		}
	case *IntegerValue:
		switch right.(type) {
		case *AnyValue:
			return integerGTUnknown(environment, l, right.Value())
		default:
			return false
		}
	default:
		return false
	}
	return false
}

func compareStringEQRegexp(s, re interface{}) bool {
	sv := s.(string)
	rev := re.(*regexp.Regexp)
	return rev.MatchString(sv)
}

func compareStringEQString(leftInterface interface{}, rightInterface interface{}) bool {
	left := leftInterface.(string)
	right := rightInterface.(string)
	return left == right
}

func regexpEQUnknown(environment map[string]string, re interface{}, s interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := s.(string)
	sv := environment[varName]
	rev := re.(*regexp.Regexp)
	debug.Info("comparing %s (%s) to %s", varName, sv, rev)
	return rev.MatchString(sv)
}

func dateTimeEQUnknown(environment map[string]string, dtValue *DateTimeValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	dt := dtValue.value
	debug.Info("comparing %s (%s) to %s", varName, val, dt)
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val)
	if err != nil {
		return false
	}
	return dt.Equal(coerced)
}

func dateTimeLTUnknown(environment map[string]string, dtValue *DateTimeValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	dt := dtValue.value
	debug.Info("comparing %s (%s) to %s", varName, val, dt)
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val)
	if err != nil {
		return false
	}
	return dt.Before(coerced)
}

func dateTimeGTUnknown(environment map[string]string, dtValue *DateTimeValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	dt := dtValue.value
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val)
	if err != nil {
		return false
	}
	return dt.After(coerced)
}

func integerEQUnknown(environment map[string]string, iValue *IntegerValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	i := iValue.value
	parsed, err := parseInt(val)
	debug.Info("comparing %s (%s = %d) to %d", varName, val, parsed, i)
	if err != nil {
		return false
	}
	return i == parsed
}

func integerLTUnknown(environment map[string]string, iValue *IntegerValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	i := iValue.value
	parsed, err := parseInt(val)
	debug.Info("comparing %s (%s = %d) to %d", varName, val, parsed, i)
	if err != nil {
		return false
	}
	return i < parsed
}

func integerGTUnknown(environment map[string]string, iValue *IntegerValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	i := iValue.value
	parsed, err := parseInt(val)
	debug.Info("comparing %s (%s = %d) to %d", varName, val, parsed, i)
	if err != nil {
		return false
	}
	return i > parsed
}

func doubleEQUnknown(environment map[string]string, dValue *DoubleValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	d := dValue.value
	parsed, err := decimal.NewFromString(val)
	debug.Info("comparing %s (%s = %s) to %s", varName, val, parsed, d)
	if err != nil {
		return false
	}
	return d.Equal(parsed)
}

func doubleLTUnknown(environment map[string]string, dValue interface{}, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	d := dValue.(*decimal.Decimal)
	parsed, err := decimal.NewFromString(val)
	if err != nil {
		parsedInt, err := parseInt(val)
		if err != nil {
			return false
		}
		parsed = decimal.New(parsedInt, 0)
	}
	debug.Info("comparing %s (%s = %v) to %v", varName, val, parsed, d)
	return d.LessThan(parsed)
}

func doubleGTUnknown(environment map[string]string, dValue *DoubleValue, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	d := dValue.value
	parsed, err := decimal.NewFromString(val)
	if err != nil {
		parsedInt, err := parseInt(val)
		if err != nil {
			return false
		}
		parsed = decimal.New(parsedInt, 0)
	}
	debug.Info("comparing %s (%s = %v) to %v", varName, val, parsed, d)
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
