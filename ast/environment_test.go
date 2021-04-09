package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		environment *Environment
		variable    *VarValue
		want        Expression
	}{
		{
			"get positive valid column from environment",
			&Environment{
				Row: &Row{1, []string{"whole line 8", "whole", "line", "7"}},
			},
			NewVarValue("%2").(*VarValue),
			&AnyValue{"line"},
		},
		{
			"get negative valid column from environment",
			&Environment{
				Row: &Row{10, []string{"whole line 8", "whole", "line", "7"}},
			},
			NewVarValue("%-1").(*VarValue),
			&AnyValue{"7"},
		},
		{
			"get positive invalid column from environment",
			&Environment{
				Row: &Row{10, []string{"whole line 8", "whole", "line", "7"}},
			},
			NewVarValue("%6").(*VarValue),
			&AnyValue{""},
		},
		{
			"get negative invalid column from environment",
			&Environment{
				Row: &Row{10, []string{"whole line 8", "whole", "line", "7"}},
			},
			NewVarValue("%-6").(*VarValue),
			&AnyValue{""},
		},
		{
			"get whole line from environment",
			&Environment{
				Row: &Row{10, []string{"whole line 8", "whole", "line", "7"}},
			},
			NewVarValue("%0").(*VarValue),
			&AnyValue{"whole line 8"},
		},
		{
			"get matching variable from environment",
			&Environment{
				Variables: map[string]Value{"varname": &VarValue{"matching value"}},
			},
			NewVarValue("varname").(*VarValue),
			&VarValue{"matching value"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got := test.environment.Resolve(test.variable)

			assert.Equal(test.want, got)
		})
	}
}
