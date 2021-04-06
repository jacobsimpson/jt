package main

import (
	"testing"

	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/pparser"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input string
		want  *ast.Program
	}{
		{
			"%1>9",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.GT_Operator,
						Right:    mustNewIntegerValue(t, "9"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
		},
		{
			"%1<0x03",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.LT_Operator,
						Right:    mustNewIntegerValue(t, "03"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
		},
		{
			" %1 == 0x03     ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.EQ_Operator,
						Right:    mustNewIntegerValue(t, "03"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
		},
		{
			" %0   ==  /things/ ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "things"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
		},
		//{
		//	"%0 == /things/ { print(%0) }",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.EQ_Operator,
		//				Right:    mustNewRegexpValue(t, "things"),
		//			},
		//			newPrintBlock(),
		//		},
		//	}},
		//},
		//{
		//	"<9",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"==/this/",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"/this/",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"<2020-01-01T",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	">2020-01-01T",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%2 == today",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%2 == yesterday",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%2 == tomorrow",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%3 == 'this is the thing'",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%3 == 2.4",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%3 in {1, 3, 5}",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
		//{
		//	"%3[-4] == '.txt'",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    mustNewIntegerValue(t, "9"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			got, err := parse(test.input)

			assert.NoError(err)
			assert.Equal(test.want, got)

			got, err = pparser.ParseString(test.input)

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

func mustNewRegexpValue(t *testing.T, value string) ast.Value {
	v, err := ast.NewRegexpValue(value)
	if err != nil {
		t.Fatalf("Unable to convert %q to a value", value)
	}
	return v
}

func newPrintBlock() *ast.Block {
	return &ast.Block{
		[]ast.Command{ast.NewPrintCommand([]ast.Expression{ast.NewVariableExpression("%0")})},
	}
}
