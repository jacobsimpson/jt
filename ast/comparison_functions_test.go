package ast

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_lt(t *testing.T) {
	//- any
	//- date/time (date, timestamp, time??)
	//- integer
	//- reals
	//- regular expressions
	//- string
	//- string range
	tests := []struct {
		left  Value
		right Value
		want  bool
	}{
		{&AnyValue{"12"} /*         */, &AnyValue{"123"} /*              */, true},
		{&AnyValue{"123"} /*        */, &AnyValue{"123"} /*              */, false},
		{&AnyValue{"124"} /*        */, &AnyValue{"123"} /*              */, false},
		{&AnyValue{"123"} /*        */, mustDateTime(t, "2010-10-11T06:45"), false},
		{&AnyValue{"2010-10-11T05:15"}, mustDateTime(t, "2010-10-11T06:45"), true},
		{&AnyValue{"2010-10-11T06:45"}, mustDateTime(t, "2010-10-11T06:45"), false},
		{&AnyValue{"2010-11-11T05:15"}, mustDateTime(t, "2010-10-11T06:45"), false},
		{&AnyValue{"aaa"} /*        */, &IntegerValue{"13", 13} /*       */, false},
		{&AnyValue{"12"} /*         */, &IntegerValue{"13", 13} /*       */, true},
		{&AnyValue{"13"} /*         */, &IntegerValue{"13", 13} /*       */, false},
		{&AnyValue{"14"} /*         */, &IntegerValue{"13", 13} /*       */, false},
		{&AnyValue{"zzz"} /*        */, mustDouble(t, "13.1") /*         */, false},
		{&AnyValue{"2"} /*          */, mustDouble(t, "13.1") /*         */, true},
		{&AnyValue{"13.0"} /*       */, mustDouble(t, "13.1") /*         */, true},
		{&AnyValue{"13.1"} /*       */, mustDouble(t, "13.1") /*         */, false},
		{&AnyValue{"13.2"} /*       */, mustDouble(t, "13.1") /*         */, false},
		{&AnyValue{"200"} /*        */, mustDouble(t, "13.1") /*         */, false},
		{&AnyValue{"abc"} /*        */, NewStringValue("'bbc'") /*       */, true},
		{&AnyValue{"bbc"} /*        */, NewStringValue("'bbc'") /*       */, false},
		{&AnyValue{"cbc"} /*        */, NewStringValue("'bbc'") /*       */, false},

		{mustDateTime(t, "2010-09-08T07:06"), &AnyValue{"123"} /*               */, false},
		{mustDateTime(t, "2010-09-08T07:06"), &AnyValue{"2010-09-09T07:06"} /*  */, true},
		{mustDateTime(t, "2010-09-08T07:06"), &AnyValue{"2010-09-08T07:06"} /*  */, false},
		{mustDateTime(t, "2010-09-08T07:06"), &AnyValue{"2010-09-07T07:06"} /*  */, false},
		{mustDateTime(t, "2010-09-08T07:06"), mustDateTime(t, "2010-09-09T07:06"), false},
		{mustDateTime(t, "2010-09-08T07:06"), mustDateTime(t, "2010-09-08T07:06"), false},
		{mustDateTime(t, "2010-09-08T07:06"), mustDateTime(t, "2010-09-07T07:06"), false},
		{mustDateTime(t, "2010-09-08T07:06"), &IntegerValue{"13", 13} /*        */, false},
		{mustDateTime(t, "2010-09-08T07:06"), mustDouble(t, "13.1") /*          */, false},
		{mustDateTime(t, "2010-09-08T07:06"), NewStringValue("'bbc'") /*        */, false},
		{mustDateTime(t, "2010-09-08T07:06"), NewStringValue("'2010-09-07T07:06'"), false},

		{&IntegerValue{"40", 40}, &AnyValue{"abc"} /*              */, false},
		{&IntegerValue{"40", 40}, &AnyValue{"50"} /*               */, true},
		{&IntegerValue{"40", 40}, &AnyValue{"40"} /*               */, false},
		{&IntegerValue{"40", 40}, &AnyValue{"30"} /*               */, false},
		{&IntegerValue{"40", 40}, mustDateTime(t, "2010-09-07T07:06"), false},
		{&IntegerValue{"40", 40}, &IntegerValue{"50", 50} /*       */, false},
		{&IntegerValue{"40", 40}, &IntegerValue{"40", 40} /*       */, false},
		{&IntegerValue{"40", 40}, &IntegerValue{"30", 30} /*       */, false},
		// This behavior is currently not implemented. But it should be.
		//{&IntegerValue{"40", 40}, mustDouble(t, "50.0") /*         */, true},
		{&IntegerValue{"40", 40}, mustDouble(t, "40.0") /*         */, false},
		{&IntegerValue{"40", 40}, mustDouble(t, "30.0") /*         */, false},
		{&IntegerValue{"40", 40}, NewStringValue("'bbc'") /*       */, false},
		{&IntegerValue{"40", 40}, NewStringValue("'50'") /*        */, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s < %s", test.left, test.right), func(t *testing.T) {
			assert := assert.New(t)

			got := lt(&Environment{}, test.left, test.right)

			assert.Equal(test.want, got)
		})
	}
}

func Test_eq(t *testing.T) {
	tests := []struct {
		environment *Environment
		left        Value
		right       Value
		want        bool
	}{
		{
			&Environment{
				Row: &Row{1, []string{"whole 8", "whole", "8"}},
			},
			&VarValue{"%2"},
			&IntegerValue{raw: "1000", value: 8},
			true,
		},
		{
			&Environment{
				Row: &Row{1, []string{"whole 8", "whole", "7"}},
			},
			&VarValue{"%2"},
			&IntegerValue{raw: "1000", value: 8},
			false,
		},
		{
			&Environment{},
			&StringValue{raw: "abcd", value: "abcd"},
			&StringValue{raw: "abcd", value: "abcd"},
			true,
		},
		{
			&Environment{},
			&StringValue{raw: "abcde", value: "abcde"},
			&StringValue{raw: "abcd", value: "abcd"},
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v == %v", test.left, test.right), func(t *testing.T) {
			assert := assert.New(t)

			got := eq(test.environment, test.left, test.right)

			assert.Equal(test.want, got)
		})
	}
}

func TestVariableResolution(t *testing.T) {
	tests := []struct {
		environment *Environment
		left        Value
		right       Value
		want        bool
	}{
		{
			&Environment{
				Variables: map[string]Value{
					"varname": &AnyValue{"2010-10-11T05:15"},
				},
			},
			&VarValue{"varname"},
			&DateTimeValue{
				"2010-10-11T06:45",
				time.Date(2010, 10, 11, 6, 45, 0, 0, time.Now().Location()),
			},
			true,
		},
		{
			&Environment{
				Variables: map[string]Value{
					"varname": &AnyValue{"12"},
				},
			},
			&VarValue{"varname"},
			&IntegerValue{"12", 12},
			false,
		},
		{
			&Environment{
				Variables: map[string]Value{
					"varname": &AnyValue{"12"},
				},
			},
			&VarValue{"varname"},
			&IntegerValue{"13", 13},
			true,
		},
		{
			&Environment{
				Variables: map[string]Value{
					"varname": &AnyValue{"12"},
				},
			},
			&VarValue{"varname"},
			&DoubleValue{
				"13.1",
				func() *decimal.Decimal {
					v := decimal.NewFromFloat(13.1)
					return &v
				}(),
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v < %v", test.left, test.right), func(t *testing.T) {
			assert := assert.New(t)

			got := lt(test.environment, test.left, test.right)

			assert.Equal(test.want, got)
		})
	}
}

func mustDateTime(t *testing.T, v string) Value {
	t.Helper()

	d, err := NewDateTimeValue(v)
	if err != nil {
		t.Fatalf("Unable to convert %q to a data/time: %+v", v, err)
	}

	return d
}

func mustDouble(t *testing.T, v string) Value {
	t.Helper()

	d, err := decimal.NewFromString(v)
	if err != nil {
		t.Fatalf("Unable to convert %q to a decimal: %+v", v, err)
	}

	return &DoubleValue{v, &d}
}
