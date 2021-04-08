package main

import (
	"fmt"
	"testing"

	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input      string
		want       *ast.Program
		wantErrors parser.ErrorLister
	}{
		{
			"%1>9",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.GT_Operator,
						Right:    ast.NewIntegerValue("9", 9),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"%1<0x03",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.LT_Operator,
						Right:    ast.NewIntegerValue("0x03", 3),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			" %1 == 0x03     ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.EQ_Operator,
						Right:    ast.NewIntegerValue("0x03", 3),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			" %99 == 0b0110     ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%99"),
						Operator: ast.EQ_Operator,
						Right:    ast.NewIntegerValue("0b0110", 6),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			" %19 == 0b01_10     ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%19"),
						Operator: ast.EQ_Operator,
						Right:    ast.NewIntegerValue("0b01_10", 6),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"%2 <= 0b00_00_10_00",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%2"),
						Operator: ast.LE_Operator,
						Right:    ast.NewIntegerValue("0b00_00_10_00", 8),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
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
			nil,
		},
		{
			" %1   ==  2014-09-12T ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.EQ_Operator,
						Right:    mustNewDateTimeValue(t, "2014-09-12T"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			" >=0o723 ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.GE_Operator,
						Right:    ast.NewIntegerValue("0o723", 467),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			" %1   >=  13.45 ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%1"),
						Operator: ast.GE_Operator,
						Right:    mustNewDoubleFromString(t, "13.45"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"%0 == /things/ { print(%0) }",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "things"),
					},
					newPrintBlock(),
				},
			}},
			nil,
		},
		{
			"%0 == /things/ { notarealfunc(%2) }",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "things"),
					},
					&ast.Block{
						Commands: []*ast.Command{
							&ast.Command{
								Name:       "notarealfunc",
								Parameters: []ast.Expression{ast.NewVarValue("%2")},
							},
						},
					},
				},
			}},
			nil,
		},
		{
			"/things/ { print(%2) }",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "things"),
					},
					&ast.Block{
						Commands: []*ast.Command{
							&ast.Command{
								Name:       "print",
								Parameters: []ast.Expression{ast.NewVarValue("%2")},
							},
						},
					},
				},
			}},
			nil,
		},
		{
			" %9 == -3     ",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%9"),
						Operator: ast.EQ_Operator,
						Right:    ast.NewIntegerValue("-3", -3),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"/things/ { print(%2[3:7]) }",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "things"),
					},
					&ast.Block{
						Commands: []*ast.Command{
							&ast.Command{
								Name: "print",
								Parameters: []ast.Expression{
									&ast.RangeExpression{
										ast.NewVarValue("%2"),
										func(i int) *int { return &i }(3),
										func(i int) *int { return &i }(7),
									},
								},
							},
						},
					},
				},
			}},
			nil,
		},
		{
			"/things/ { print(%2[-3:]) }",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "things"),
					},
					&ast.Block{
						Commands: []*ast.Command{
							&ast.Command{
								Name: "print",
								Parameters: []ast.Expression{
									&ast.RangeExpression{
										ast.NewVarValue("%2"),
										func(i int) *int { return &i }(-3),
										nil,
									},
								},
							},
						},
					},
				},
			}},
			nil,
		},
		{
			"1.0 < %3 <= 2.4",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.AndComparison{
						&ast.Comparison{
							Left:     mustNewDoubleFromString(t, "1.0"),
							Operator: ast.LT_Operator,
							Right:    ast.NewVarValue("%3"),
						},
						&ast.Comparison{
							Left:     ast.NewVarValue("%3"),
							Operator: ast.LE_Operator,
							Right:    mustNewDoubleFromString(t, "2.4"),
						},
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"3.0 >= %4 > 2.4",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.AndComparison{
						&ast.Comparison{
							Left:     mustNewDoubleFromString(t, "3.0"),
							Operator: ast.GE_Operator,
							Right:    ast.NewVarValue("%4"),
						},
						&ast.Comparison{
							Left:     ast.NewVarValue("%4"),
							Operator: ast.GT_Operator,
							Right:    mustNewDoubleFromString(t, "2.4"),
						},
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"3.0 >= %4 < 2.4",
			nil,
			parser.NewErrorLister([]error{parser.NewParserError(
				fmt.Errorf("can not build a ternary boolean expression out of >= and < comparisons"),
				1,
				1,
				0,
				"test:1:1 (0): rule three_term_boolean_expression",
			)}),
		},
		//{
		//	" %3 == +6786     ",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%3"),
		//				Operator: ast.EQ_Operator,
		//				Right:    ast.NewIntegerValue("+6786", 6786),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//		nil,
		//},
		{
			"<9",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.LT_Operator,
						Right:    ast.NewIntegerValue("9", 9),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"==/this/",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "this"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"/this/",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.EQ_Operator,
						Right:    mustNewRegexpValue(t, "this"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"<2020-01-01T",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.LT_Operator,
						Right:    mustNewDateTimeValue(t, "2020-01-01T"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			">2020-01-01T",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%0"),
						Operator: ast.GT_Operator,
						Right:    mustNewDateTimeValue(t, "2020-01-01T"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		//{
		//	"%2 == today",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    ast.NewIntegerValue("9", 9),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//nil,
		//},
		//{
		//	"%2 == yesterday",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    ast.NewIntegerValue("9", 9),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//nil,
		//},
		//{
		//	"%2 == tomorrow",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    ast.NewIntegerValue("9", 9),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//nil,
		//},
		{
			`%3 == "this is the thing"`,
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%3"),
						Operator: ast.EQ_Operator,
						Right:    ast.NewStringValue(`"this is the thing"`),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		{
			"%3 == 2.4",
			&ast.Program{[]*ast.Rule{
				&ast.Rule{
					&ast.Comparison{
						Left:     ast.NewVarValue("%3"),
						Operator: ast.EQ_Operator,
						Right:    mustNewDoubleFromString(t, "2.4"),
					},
					ast.NewPrintlnBlock(),
				},
			}},
			nil,
		},
		//{
		//	"%3 in {1, 3, 5}",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left:     ast.NewVarValue("%0"),
		//				Operator: ast.LT_Operator,
		//				Right:    ast.NewIntegerValue("9", 9),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//nil,
		//},
		// This test case requires that a rule comparison be updated to support
		// arbitrary expressions on LHS and RHS.
		//{
		//	"%3[-4:] == '.txt'",
		//	&ast.Program{[]*ast.Rule{
		//		&ast.Rule{
		//			&ast.Comparison{
		//				Left: &ast.RangeExpression{
		//					ast.NewVarValue("%2"),
		//					func(i int) *int { return &i }(3),
		//					func(i int) *int { return &i }(7),
		//				},
		//				Operator: ast.LT_Operator,
		//				Right:    ast.NewStringValue("'.txt'"),
		//			},
		//			ast.NewPrintlnBlock(),
		//		},
		//	}},
		//nil,
		//},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert := assert.New(t)

			got, err := parse(test.input)

			assert.Equal(test.wantErrors, err)
			assert.Equal(test.want, got)
		})
	}
}

func mustNewRegexpValue(t *testing.T, value string) ast.Value {
	v, err := ast.NewRegexpValue(value)
	if err != nil {
		t.Fatalf("Unable to convert %q to a value", value)
	}
	return v
}

func mustNewDateTimeValue(t *testing.T, value string) ast.Value {
	v, err := ast.NewDateTimeValue(value)
	if err != nil {
		t.Fatalf("Unable to convert %q to a value", value)
	}
	return v
}

func mustNewDoubleFromString(t *testing.T, value string) ast.Value {
	v, err := ast.NewDoubleFromString(value)
	if err != nil {
		t.Fatalf("Unable to convert %q to a value", value)
	}
	return v
}

func newPrintBlock() *ast.Block {
	return &ast.Block{
		[]*ast.Command{ast.NewPrintCommand([]ast.Expression{ast.NewVarValue("%0")})},
	}
}
