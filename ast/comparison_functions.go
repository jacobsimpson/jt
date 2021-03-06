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
		case *AnyValue:
			return anyGTAny(r, l)
		case *DateTimeValue:
			return dateTimeGTAny(r, l)
		case *DoubleValue:
			return doubleGTAny(r, l)
		case *IntegerValue:
			return integerGTAny(r, l)
		case *StringValue:
			return stringGTAny(r, l)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeLTAny(l, r)
		case *DateTimeValue:
			return dateTimeLTDateTime(l, r)
		}
	case *DoubleValue:
		switch r := right.(type) {
		case *AnyValue:
			return doubleLTAny(l, r)
		case *DoubleValue:
			return doubleLTDouble(l, r)
		case *IntegerValue:
			return doubleLTInteger(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerLTAny(l, r)
		case *DoubleValue:
			return integerLTDouble(l, r)
		case *IntegerValue:
			return integerLTInteger(l, r)
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
		case *AnyValue:
			return anyGTAny(r, l) || anyEQAny(r, l)
		case *DateTimeValue:
			return dateTimeGTAny(r, l) || dateTimeEQAny(r, l)
		case *DoubleValue:
			return doubleGTAny(r, l) || doubleEQAny(r, l)
		case *IntegerValue:
			return integerGTAny(r, l) || integerEQAny(r, l)
		case *StringValue:
			return stringGTAny(r, l) || stringEQAny(r, l)
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
	case *AnyValue:
		switch r := right.(type) {
		case *AnyValue:
			return anyEQAny(r, l)
		case *DateTimeValue:
			return dateTimeEQAny(r, l)
		case *DoubleValue:
			return doubleEQAny(r, l)
		case *IntegerValue:
			return integerEQAny(r, l)
		case *RegexpValue:
			return regexpEQAny(r, l)
		case *StringValue:
			return stringEQAny(r, l)
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
	case *RegexpValue:
		switch r := right.(type) {
		case *AnyValue:
			return regexpEQAny(l, r)
		case *StringValue:
			return compareStringEQRegexp(r, l)
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
		case *AnyValue:
			return anyLTAny(r, l) || anyEQAny(r, l)
		case *DateTimeValue:
			return dateTimeLTAny(r, l) || dateTimeEQAny(r, l)
		case *DoubleValue:
			return doubleLTAny(r, l) || doubleEQAny(r, l)
		case *IntegerValue:
			return integerLTAny(r, l) || integerEQAny(r, l)
		case *StringValue:
			return stringLTAny(r, l) || stringEQAny(r, l)
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
		case *AnyValue:
			return anyGTAny(l, r)
		case *DateTimeValue:
			return dateTimeLTAny(r, l)
		case *DoubleValue:
			return doubleLTAny(r, l)
		case *IntegerValue:
			return integerLTAny(r, l)
		case *StringValue:
			return stringLTAny(r, l)
		}
	case *DateTimeValue:
		switch r := right.(type) {
		case *AnyValue:
			return dateTimeGTAny(l, r)
		case *DateTimeValue:
			return dateTimeGTDateTime(l, r)
		}
	case *DoubleValue:
		switch r := right.(type) {
		case *AnyValue:
			return doubleGTAny(l, r)
		case *DoubleValue:
			return doubleGTDouble(l, r)
		case *IntegerValue:
			return doubleGTInteger(l, r)
		}
	case *IntegerValue:
		switch r := right.(type) {
		case *AnyValue:
			return integerGTAny(l, r)
		case *DoubleValue:
			return integerGTDouble(l, r)
		case *IntegerValue:
			return integerGTInteger(l, r)
		}
	case *StringValue:
		switch r := right.(type) {
		case *AnyValue:
			return stringGTAny(l, r)
		case *StringValue:
			return stringGTString(l, r)
		}
	}
	return false
}

func compareStringEQRegexp(lhs *StringValue, rhs *RegexpValue) bool {
	return rhs.re.MatchString(lhs.value)
}

func compareStringEQString(lhs *StringValue, rhs *StringValue) bool {
	return lhs.value == rhs.value
}

func regexpEQAny(lhs *RegexpValue, rhs *AnyValue) bool {
	rev := lhs.re
	return rev.MatchString(rhs.raw)
}

func dateTimeEQAny(lhs *DateTimeValue, rhs *AnyValue) bool {
	dt := lhs.value
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, rhs.raw)
	if err != nil {
		return false
	}
	return dt.Equal(coerced)
}

func dateTimeLTAny(lhs *DateTimeValue, rhs *AnyValue) bool {
	dt := lhs.value
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, rhs.raw)
	if err != nil {
		return false
	}
	return dt.Before(coerced)
}

func dateTimeLTDateTime(lhs *DateTimeValue, rhs *DateTimeValue) bool {
	return lhs.value.Before(rhs.value)
}

func dateTimeGTAny(lhs *DateTimeValue, rhs *AnyValue) bool {
	coerced, err := datetime.ParseDateTime(datetime.CoercionFormats, rhs.raw)
	if err != nil {
		return false
	}
	return lhs.value.After(coerced)
}

func dateTimeGTDateTime(lhs *DateTimeValue, rhs *DateTimeValue) bool {
	return lhs.value.After(rhs.value)
}

func integerEQAny(lhs *IntegerValue, rhs *AnyValue) bool {
	i := lhs.value
	parsed, err := parseInt(rhs.raw)
	if err != nil {
		return false
	}
	return i == parsed
}

func integerLTAny(lhs *IntegerValue, rhs *AnyValue) bool {
	i := lhs.value
	parsed, err := parseInt(rhs.raw)
	if err != nil {
		return false
	}
	return i < parsed
}

func integerLTDouble(lhs *IntegerValue, rhs *DoubleValue) bool {
	return decimal.NewFromInt(lhs.value).LessThan(*rhs.value)
}

func integerLTInteger(lhs *IntegerValue, rhs *IntegerValue) bool {
	return lhs.value < rhs.value
}

func integerGTAny(lhs *IntegerValue, rhs *AnyValue) bool {
	i := lhs.value
	parsed, err := parseInt(rhs.raw)
	if err != nil {
		return false
	}
	return i > parsed
}

func integerGTDouble(lhs *IntegerValue, rhs *DoubleValue) bool {
	return decimal.NewFromInt(lhs.value).GreaterThan(*rhs.value)
}

func integerGTInteger(lhs *IntegerValue, rhs *IntegerValue) bool {
	return lhs.value > rhs.value
}

func doubleEQAny(lhs *DoubleValue, rhs *AnyValue) bool {
	d := lhs.value
	parsed, err := decimal.NewFromString(rhs.raw)
	if err != nil {
		return false
	}
	return d.Equal(parsed)
}

func doubleLTAny(lhs *DoubleValue, rhs *AnyValue) bool {
	d := lhs.value
	parsed, err := decimal.NewFromString(rhs.raw)
	if err != nil {
		parsedInt, err := parseInt(rhs.raw)
		if err != nil {
			return false
		}
		parsed = decimal.New(parsedInt, 0)
	}
	return d.LessThan(parsed)
}

func doubleGTAny(lhs *DoubleValue, rhs *AnyValue) bool {
	d := lhs.value
	parsed, err := decimal.NewFromString(rhs.raw)
	if err != nil {
		parsedInt, err := parseInt(rhs.raw)
		if err != nil {
			return false
		}
		parsed = decimal.New(parsedInt, 0)
	}
	return d.GreaterThan(parsed)
}

func doubleGTDouble(lhs *DoubleValue, rhs *DoubleValue) bool {
	return lhs.value.GreaterThan(*rhs.value)
}

func doubleGTInteger(lhs *DoubleValue, rhs *IntegerValue) bool {
	return lhs.value.GreaterThan(decimal.NewFromInt(rhs.value))
}

func doubleLTDouble(lhs *DoubleValue, rhs *DoubleValue) bool {
	return lhs.value.LessThan(*rhs.value)
}

func doubleLTInteger(lhs *DoubleValue, rhs *IntegerValue) bool {
	return lhs.value.LessThan(decimal.NewFromInt(rhs.value))
}

func stringEQAny(lhs *StringValue, rhs *AnyValue) bool {
	return lhs.value == rhs.raw
}

func stringLTAny(lhs *StringValue, rhs *AnyValue) bool {
	return lhs.value < rhs.raw
}

func stringGTAny(lhs *StringValue, rhs *AnyValue) bool {
	return lhs.value > rhs.raw
}

func stringGTString(lhs *StringValue, rhs *StringValue) bool {
	return lhs.value > rhs.value
}

func stringLTString(lhs *StringValue, rhs *StringValue) bool {
	return lhs.value < rhs.value
}

func anyGTAny(lhs *AnyValue, rhs *AnyValue) bool {
	return lhs.raw > rhs.raw
}

func anyEQAny(lhs *AnyValue, rhs *AnyValue) bool {
	return lhs.raw == rhs.raw
}

func anyLTAny(lhs *AnyValue, rhs *AnyValue) bool {
	return lhs.raw < rhs.raw
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
