package listener

import (
	"fmt"
	"regexp"
	"time"

	"github.com/jacobsimpson/jt/debug"
)

func lt(environment map[string]string, left, right Value) bool {
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

var dateCoercionFormats = []string{
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

func dateTimeEQUnknown(environment map[string]string, dt interface{}, v interface{}) bool {
	// TODO: This is going to crash hard if the variable doesn't exist.
	varName := v.(string)
	val := environment[varName]
	datetime := dt.(time.Time)
	debug.Info("comparing %s (%s) to %s", varName, val, datetime)
	coerced, err := coerceDateTime(val)
	if err != nil {
		return false
	}
	return datetime.Equal(*coerced)
}

func coerceDateTime(str string) (*time.Time, error) {
	for _, layout := range dateLiteralFormats {
		if t, err := time.Parse(layout, str); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("Unable to convert %q to a date", str)
}
