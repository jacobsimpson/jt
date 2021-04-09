package ast

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_lt(t *testing.T) {
	tests := []struct {
		environment *Environment
		left        Value
		right       Value
		want        bool
	}{
		{
			&Environment{},
			&AnyValue{"123"},
			&DateTimeValue{
				"2010-10-11T06:45",
				time.Date(2010, 10, 11, 6, 45, 0, 0, time.Now().Location()),
			},
			false,
		},
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
