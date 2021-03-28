package main

import (
	"testing"

	"github.com/jacobsimpson/jt/ast"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input string
		want  ast.Program
	}{
		{
			"%1<9",
			ast.NewProgram([]ast.Rule{
				ast.NewRule(&ast.Comparison{
					Left:  ast.NewVarValue("%1"),
					Right: mustNewIntegerValue(t, "9"),
				}, ast.NewPrintlnBlock()),
			}),
		},
		{
			"%1<0x03",
			ast.NewProgram([]ast.Rule{
				ast.NewRule(&ast.Comparison{
					Left:  ast.NewVarValue("%1"),
					Right: mustNewIntegerValue(t, "03"),
				}, ast.NewPrintlnBlock()),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			got, err := parse(test.input)

			assert.NoError(err)
			assert.Equal(test.want, got)
		})
	}
}

func mustNewIntegerValue(t *testing.T, value string) ast.Value {
	v, err := ast.NewIntegerValue(value)
	if err != nil {
		t.Fatalf("Unable to convert %q to a value", value)
	}
	return v
}

func mustNewHexIntegerValue(t *testing.T, value string) ast.Value {
	v, err := ast.NewIntegerValueFromHexString(value)
	if err != nil {
		t.Fatalf("Unable to convert %q to a value", value)
	}
	return v
}
