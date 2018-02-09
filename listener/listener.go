package listener

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/parser"
)

type InterpreterListener struct {
	Rules       []ast.Rule
	currentRule ast.Rule
	Errors      []*ParsingError
}

func NewInterpreterListener() *InterpreterListener {
	return &InterpreterListener{}
}

func (l *InterpreterListener) FoundErrors() bool {
	return len(l.Errors) > 0
}

func (l *InterpreterListener) VisitTerminal(node antlr.TerminalNode) {}
func (l *InterpreterListener) VisitErrorNode(node antlr.ErrorNode)   {}

func (l *InterpreterListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}
func (l *InterpreterListener) ExitEveryRule(ctx antlr.ParserRuleContext)  {}

func (l *InterpreterListener) EnterProgram(ctx *parser.ProgramContext) {}
func (l *InterpreterListener) ExitProgram(ctx *parser.ProgramContext)  {}

func (l *InterpreterListener) EnterProcessingRule(ctx *parser.ProcessingRuleContext) {
	l.currentRule = ast.NewRule(ast.NewPrintlnBlock())
}
func (l *InterpreterListener) ExitProcessingRule(ctx *parser.ProcessingRuleContext) {
	if l.currentRule != nil {
		l.Rules = append(l.Rules, l.currentRule)
		l.currentRule = nil
	}
}

func (l *InterpreterListener) EnterSelection(ctx *parser.SelectionContext) {
	if ctx.REGULAR_EXPRESSION() != nil {
		regexpString := ctx.REGULAR_EXPRESSION().GetSymbol().GetText()
		regexpString = regexpString[1 : len(regexpString)-1]
		value, err := ast.NewRegexpValue(regexpString)
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    fmt.Sprintf("could not parse regular expression %q: %v", regexpString, err),
				line:   ctx.REGULAR_EXPRESSION().GetSymbol().GetLine(),
				column: ctx.REGULAR_EXPRESSION().GetSymbol().GetColumn(),
			})
			return
		}
		l.currentRule.SetSelection(&ast.Comparison{
			Left:     ast.NewVarValue("%0"),
			Operator: ast.EQ_Operator,
			Right:    value,
		})
	}
}
func (l *InterpreterListener) ExitSelection(ctx *parser.SelectionContext) {}

func (l *InterpreterListener) EnterComparison(ctx *parser.ComparisonContext) {
	l.currentRule.SetSelection(&ast.Comparison{})
}
func (l *InterpreterListener) ExitComparison(ctx *parser.ComparisonContext) {}

func (l *InterpreterListener) EnterComparator(ctx *parser.ComparatorContext) {
	cmp := l.currentRule.Selection().(*ast.Comparison)
	if ctx.LT() != nil {
		cmp.Operator = ast.LT_Operator
	} else if ctx.LE() != nil {
		cmp.Operator = ast.LE_Operator
	} else if ctx.EQ() != nil {
		cmp.Operator = ast.EQ_Operator
	} else if ctx.NE() != nil {
		cmp.Operator = ast.NE_Operator
	} else if ctx.GE() != nil {
		cmp.Operator = ast.GE_Operator
	} else if ctx.GT() != nil {
		cmp.Operator = ast.GT_Operator
	}
}
func (l *InterpreterListener) ExitComparator(ctx *parser.ComparatorContext) {}

func (l *InterpreterListener) EnterValue(ctx *parser.ValueContext) {
	var value ast.Value
	if ctx.COLUMN() != nil {
		value = ast.NewVarValue(ctx.COLUMN().GetSymbol().GetText())
	} else if ctx.REGULAR_EXPRESSION() != nil {
		regexpString := ctx.REGULAR_EXPRESSION().GetSymbol().GetText()
		regexpString = regexpString[1 : len(regexpString)-1]
		var err error
		value, err = ast.NewRegexpValue(regexpString)
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    fmt.Sprintf("could not parse regular expression %q: %v", regexpString, err),
				line:   ctx.REGULAR_EXPRESSION().GetSymbol().GetLine(),
				column: ctx.REGULAR_EXPRESSION().GetSymbol().GetColumn(),
			})
			return
		}
	} else if ctx.STRING() != nil {
		value = ast.NewStringValue(ctx.STRING().GetSymbol().GetText())
	} else if ctx.INTEGER() != nil {
		var err error
		value, err = ast.NewIntegerValue(ctx.INTEGER().GetSymbol().GetText())
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    err.Error(),
				line:   ctx.DATE_TIME().GetSymbol().GetLine(),
				column: ctx.DATE_TIME().GetSymbol().GetColumn(),
			})
			return
		}
	} else if ctx.HEX_INTEGER() != nil {
	} else if ctx.BINARY_INTEGER() != nil {
	} else if ctx.DECIMAL() != nil {
	} else if ctx.DATE_TIME() != nil {
		var err error
		value, err = ast.NewDateTimeValue(ctx.DATE_TIME().GetSymbol().GetText())
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    err.Error(),
				line:   ctx.DATE_TIME().GetSymbol().GetLine(),
				column: ctx.DATE_TIME().GetSymbol().GetColumn(),
			})
			return
		}
	}
	cmp := l.currentRule.Selection().(*ast.Comparison)
	if cmp.Left == nil {
		cmp.Left = value
	} else {
		cmp.Right = value
	}
}
func (l *InterpreterListener) ExitValue(ctx *parser.ValueContext) {}

func (l *InterpreterListener) EnterBlock(ctx *parser.BlockContext) {
	l.currentRule.SetBlock(ast.NewBlock())
}
func (l *InterpreterListener) ExitBlock(ctx *parser.BlockContext) {}

func (l *InterpreterListener) EnterCommand(ctx *parser.CommandContext) {
	symbol := ctx.IDENTIFIER().GetSymbol()
	switch symbol.GetText() {
	case "print":
		l.currentRule.Block().AddCommand(ast.NewPrintCommand())
	case "println":
		l.currentRule.Block().AddCommand(ast.NewPrintlnCommand())
	default:
		l.Errors = append(l.Errors, &ParsingError{
			msg:    fmt.Sprintf("unknown function %q", symbol.GetText()),
			line:   symbol.GetLine(),
			column: symbol.GetColumn(),
		})
	}
}
func (l *InterpreterListener) ExitCommand(ctx *parser.CommandContext) {}

func (l *InterpreterListener) EnterParameterList(ctx *parser.ParameterListContext) {
	if len(l.Errors) > 0 {
		return
	}
	currentCommand := l.currentRule.Block().LastCommand()
	for _, c := range ctx.AllCOLUMN() {
		currentCommand.AddParameter(c.GetSymbol().GetText())
	}
}
func (l *InterpreterListener) ExitParameterList(ctx *parser.ParameterListContext) {}

func (l *InterpreterListener) EnterBinary(ctx *parser.BinaryContext) {}
func (l *InterpreterListener) ExitBinary(ctx *parser.BinaryContext)  {}

func (l *InterpreterListener) EnterBoolean(ctx *parser.BooleanContext) {}
func (l *InterpreterListener) ExitBoolean(ctx *parser.BooleanContext)  {}

func (l *InterpreterListener) EnterExpression(ctx *parser.ExpressionContext) {}
func (l *InterpreterListener) ExitExpression(ctx *parser.ExpressionContext)  {}
