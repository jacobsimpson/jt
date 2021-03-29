package ast

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_lt(t *testing.T) {
	tests := []struct {
		environment map[string]string
		left        Value
		right       Value
		want        bool
	}{
		{
			map[string]string{},
			&AnyValue{"123"},
			&DateTimeValue{
				"2010-10-11T06:45",
				time.Date(2010, 10, 11, 6, 45, 0, 0, time.Now().Location()),
			},
			false,
		},
		{
			map[string]string{"varname": "2010-10-11T05:15"},
			&VarValue{"varname"},
			&DateTimeValue{
				"2010-10-11T06:45",
				time.Date(2010, 10, 11, 6, 45, 0, 0, time.Now().Location()),
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
		environment map[string]string
		left        Value
		right       Value
		want        bool
	}{
		{
			map[string]string{"%2": "8"},
			&VarValue{"%2"},
			&IntegerValue{raw: "1000", value: 8},
			true,
		},
		{
			map[string]string{"%2": "7"},
			&VarValue{"%2"},
			&IntegerValue{raw: "1000", value: 8},
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
