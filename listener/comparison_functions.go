package listener

import (
	"regexp"
	"time"

	"github.com/jacobsimpson/jt/datetime"
	"github.com/jacobsimpson/jt/debug"
)

func lt(environment map[string]string, left, right Value) bool {
	switch left.Type() {
	case UnknownValue:
		switch right.Type() {
		case DateTimeValue:
			return dateTimeGTUnknown(environment, right.Value(), left.Value())
		default:
			return false
		}
	case DateTimeValue:
		switch right.Type() {
		case UnknownValue:
			return dateTimeLTUnknown(environment, left.Value(), right.Value())
		default:
			return false
		}
	default:
		return false
	}
	return false
}

func le(environment map[string]string, left, right Value) bool {
	return false
}

func eq(environment map[string]string, left, right Value) bool {
	switch left.Type() {
	case RegexpValue:
		switch right.Type() {
		case StringValue:
			return compareStringEQRegexp(right.Value(), left.Value())
		case RegexpValue:
			return regexpEQUnknown(environment, left.Value(), right.Value())
		default:
			return false
		}
	case StringValue:
		switch right.Type() {
		case RegexpValue:
			return compareStringEQRegexp(left.Value(), right.Value())
		case StringValue:
			return compareStringEQString(left.Value(), right.Value())
		default:
			// Error: "Can not compare %s to %s using %s",
			//     left.Type(), right.Type(), operator
			return false
		}
	case UnknownValue:
		switch right.Type() {
		case RegexpValue:
			return regexpEQUnknown(environment, right.Value(), left.Value())
		case DateTimeValue:
			return dateTimeEQUnknown(environment, right.Value(), left.Value())
		default:
			return false
		}
	case DateTimeValue:
		switch right.Type() {
		case UnknownValue:
			return dateTimeEQUnknown(environment, left.Value(), right.Value())
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
	return false
}

func gt(environment map[string]string, left, right Value) bool {
	switch left.Type() {
	case UnknownValue:
		switch right.Type() {
		case DateTimeValue:
			return dateTimeLTUnknown(environment, right.Value(), left.Value())
		default:
			return false
		}
	case DateTimeValue:
		switch right.Type() {
		case UnknownValue:
			return dateTimeGTUnknown(environment, left.Value(), right.Value())
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

func dateTimeEQUnknown(environment map[string]string, dtValue interface{}, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	dt := dtValue.(*time.Time)
	debug.Info("comparing %s (%s) to %s", varName, val, dt)
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val)
	if err != nil {
		return false
	}
	return dt.Equal(*coerced)
}

func dateTimeLTUnknown(environment map[string]string, dtValue interface{}, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	dt := dtValue.(*time.Time)
	debug.Info("comparing %s (%s) to %s", varName, val, dt)
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val)
	if err != nil {
		return false
	}
	return dt.Before(*coerced)
}

func dateTimeGTUnknown(environment map[string]string, dtValue interface{}, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	dt := dtValue.(*time.Time)
	debug.Info("comparing %s (%s) to %s", varName, val, dt)
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, val)
	if err != nil {
		return false
	}
	return dt.After(*coerced)
}
