package listener

import (
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/antlrgen"
	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/debug"
)

func NewASTVisitor() antlrgen.ProgramVisitor {
	return &astVisitor{}
}

type astVisitor struct{}

func (v *astVisitor) Visit(tree antlr.ParseTree) interface{} {
	debug.Debug("astVisitor.Visit")
	return tree.Accept(v)
}

func (v *astVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	debug.Debug("astVisitor.VisitChildren")
	result := []interface{}{}
	for _, child := range node.GetChildren() {
		if _, ok := child.(antlr.TerminalNode); ok {
			// Skip TerminalNodes. They will be handled in the Visit method of
			// the node that contains the TerminalNode.
		} else if c, ok := child.(antlr.ParseTree); ok {
			r := c.Accept(v)
			if _, ok := r.(error); ok {
				return r
			}
			result = append(result, r)
		}
	}
	return result
}

func (v *astVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	debug.Debug("astVisitor.VisitTerminal node=%v", node)
	return nil
}

func (v *astVisitor) VisitErrorNode(node antlr.ErrorNode) interface{} {
	debug.Debug("astVisitor.VisitErrorNode")
	return nil
}

func (v *astVisitor) VisitProgram(ctx *antlrgen.ProgramContext) interface{} {
	debug.Debug("astVisitor.VisitProgram")
	result := v.VisitChildren(ctx)
	if err := getError(result); err != nil {
		return err
	}
	rules := []*ast.Rule{}
	for _, r := range result.([]interface{}) {
		if rule, ok := r.(*ast.Rule); ok {
			rules = append(rules, rule)
		} else {
			return fmt.Errorf("expecting rule, found %v", r)
		}
	}
	return &ast.Program{rules}
}

func (v *astVisitor) VisitProcessingRule(ctx *antlrgen.ProcessingRuleContext) interface{} {
	debug.Debug("astVisitor.VisitProcessingRule")

	var selection ast.Expression
	if ctx.Selection() != nil {
		r := ctx.Selection().Accept(v)
		if err := getError(r); err != nil {
			return err
		}
		selection = r.(ast.Expression)
	}

	var block *ast.Block
	if ctx.Block() != nil {
		r := ctx.Block().Accept(v)
		if err := getError(r); err != nil {
			return err
		}
		block = r.(*ast.Block)
	} else {
		block = ast.NewPrintlnBlock()
	}
	debug.Debug("selection = %v, block = %v", selection, block)
	return &ast.Rule{selection, block}
}

func (v *astVisitor) VisitSelection(ctx *antlrgen.SelectionContext) interface{} {
	debug.Debug("astVisitor.VisitSelection")
	if ctx.REGULAR_EXPRESSION() != nil {
		regexpString := ctx.REGULAR_EXPRESSION().GetSymbol().GetText()
		regexpString = regexpString[1 : len(regexpString)-1]
		value, err := ast.NewRegexpValue(regexpString)
		if err != nil {
			return &ParsingError{
				msg:    fmt.Sprintf("could not parse regular expression %q: %v", regexpString, err),
				line:   ctx.REGULAR_EXPRESSION().GetSymbol().GetLine(),
				column: ctx.REGULAR_EXPRESSION().GetSymbol().GetColumn(),
			}
		}
		return &ast.Comparison{
			Left:     ast.NewVarValue("%0"),
			Operator: ast.EQ_Operator,
			Right:    value,
		}
	}
	r := ctx.Expression().Accept(v)
	if err := getError(r); err != nil {
		return err
	}
	return r.(ast.Expression)
}

func (v *astVisitor) VisitExpression(ctx *antlrgen.ExpressionContext) interface{} {
	debug.Debug("astVisitor.VisitExpression")
	if ctx.GetParen() != nil {
		r := ctx.GetParen().Accept(v)
		if err := getError(r); err != nil {
			return err
		}
		return r
	} else if ctx.GetNegative() != nil {
		r := ctx.GetNegative().Accept(v)
		if err := getError(r); err != nil {
			return err
		}
		return ast.NewNegativeExpression(r.(ast.Expression))
		//} else if ctx.GetLeft() != nil {
		//	l := ctx.GetLeft().Accept(v)
		//	if err := getError(r); err != nil {
		//		return err
		//	}
		//	r := ctx.GetRight().Accept(v)
		//	if err := getError(r); err != nil {
		//		return err
		//	}
		//	o := ctx.GetOp().Accept(v)
		//	if err := getError(r); err != nil {
		//		return err
		//	}
		//	return &ast.Comparison{
		//		Left:     l.(ast.Expression),
		//		Operator: o.(ast.Operator),
		//		Right:    r.(ast.Expression),
		//	}
	} else if ctx.Comparison() != nil {
		return ctx.Comparison().Accept(v)
	}
	return v.VisitChildren(ctx)
}

func (v *astVisitor) VisitComparison(ctx *antlrgen.ComparisonContext) interface{} {
	debug.Debug("astVisitor.VisitComparison")
	l := ctx.GetLeft().Accept(v)
	if err := getError(l); err != nil {
		return err
	}
	r := ctx.GetRight().Accept(v)
	if err := getError(r); err != nil {
		return err
	}
	o := ctx.GetOp().Accept(v)
	if err := getError(o); err != nil {
		return err
	}
	return &ast.Comparison{
		Left:     l.(ast.Value),
		Operator: o.(ast.Operator),
		Right:    r.(ast.Value),
	}
}

func (v *astVisitor) VisitValue(ctx *antlrgen.ValueContext) interface{} {
	debug.Debug("astVisitor.VisitValue")
	if ctx.COLUMN() != nil {
		return ast.NewVarValue(ctx.COLUMN().GetSymbol().GetText())
	} else if ctx.REGULAR_EXPRESSION() != nil {
		regexpString := ctx.REGULAR_EXPRESSION().GetSymbol().GetText()
		regexpString = regexpString[1 : len(regexpString)-1]
		value, err := ast.NewRegexpValue(regexpString)
		if err != nil {
			return &ParsingError{
				msg:    fmt.Sprintf("could not parse regular expression %q: %v", regexpString, err),
				line:   ctx.REGULAR_EXPRESSION().GetSymbol().GetLine(),
				column: ctx.REGULAR_EXPRESSION().GetSymbol().GetColumn(),
			}
		}
		return value
	} else if ctx.STRING() != nil {
		return ast.NewStringValue(ctx.STRING().GetSymbol().GetText())
	} else if ctx.INTEGER() != nil {
		symbol := ctx.INTEGER().GetSymbol()
		value, err := ast.NewIntegerValue(symbol.GetText())
		if err != nil {
			return &ParsingError{
				msg:    err.Error(),
				line:   symbol.GetLine(),
				column: symbol.GetColumn(),
			}
		}
		return value
	} else if ctx.HEX_INTEGER() != nil {
		symbol := ctx.HEX_INTEGER().GetSymbol()
		value, err := ast.NewIntegerValueFromHexString(symbol.GetText())
		if err != nil {
			return &ParsingError{
				msg:    err.Error(),
				line:   symbol.GetLine(),
				column: symbol.GetColumn(),
			}
		}
		return value
	} else if ctx.BINARY_INTEGER() != nil {
		symbol := ctx.BINARY_INTEGER().GetSymbol()
		value, err := ast.NewIntegerValueFromBinaryString(symbol.GetText())
		if err != nil {
			return &ParsingError{
				msg:    err.Error(),
				line:   symbol.GetLine(),
				column: symbol.GetColumn(),
			}
		}
		return value
	} else if ctx.DOUBLE() != nil {
		symbol := ctx.DOUBLE().GetSymbol()
		value, err := ast.NewDoubleFromString(symbol.GetText())
		if err != nil {
			return &ParsingError{
				msg:    err.Error(),
				line:   symbol.GetLine(),
				column: symbol.GetColumn(),
			}
		}
		return value
	} else if ctx.DATE_TIME() != nil {
		symbol := ctx.DATE_TIME().GetSymbol()
		value, err := ast.NewDateTimeValue(symbol.GetText())
		if err != nil {
			return &ParsingError{
				msg:    err.Error(),
				line:   symbol.GetLine(),
				column: symbol.GetColumn(),
			}
		}
		return value
	}
	return nil
}

func (v *astVisitor) VisitComparator(ctx *antlrgen.ComparatorContext) interface{} {
	debug.Debug("astVisitor.VisitComparator")
	if ctx.LT() != nil {
		return ast.LT_Operator
	} else if ctx.LE() != nil {
		return ast.LE_Operator
	} else if ctx.EQ() != nil {
		return ast.EQ_Operator
	} else if ctx.NE() != nil {
		return ast.NE_Operator
	} else if ctx.GE() != nil {
		return ast.GE_Operator
	} else if ctx.GT() != nil {
		return ast.GT_Operator
	}
	return fmt.Errorf("unknown operator")
}

func (v *astVisitor) VisitBinary(ctx *antlrgen.BinaryContext) interface{} {
	debug.Debug("astVisitor.VisitBinary")
	return v.VisitChildren(ctx)
}

func (v *astVisitor) VisitBoolean(ctx *antlrgen.BooleanContext) interface{} {
	debug.Debug("astVisitor.VisitBoolean")
	return v.VisitChildren(ctx)
}

func (v *astVisitor) VisitBlock(ctx *antlrgen.BlockContext) interface{} {
	debug.Debug("astVisitor.VisitBlock")
	children := v.VisitChildren(ctx)
	if err := getError(children); err != nil {
		return err
	}
	commands := []*ast.Command{}
	for _, c := range children.([]interface{}) {
		command := c.(*ast.Command)
		commands = append(commands, command)
	}
	block := &ast.Block{Commands: commands}
	return block
}

func (v *astVisitor) VisitCommand(ctx *antlrgen.CommandContext) interface{} {
	debug.Debug("astVisitor.VisitCommand")
	parameters := []ast.Expression{}
	if ctx.ParameterList() != nil {
		r := ctx.ParameterList().Accept(v)
		if err, ok := r.(error); ok {
			return err
		}
		for _, e := range r.([]interface{}) {
			parameters = append(parameters, e.(ast.Expression))
		}
	}
	symbol := ctx.IDENTIFIER().GetSymbol()
	switch symbol.GetText() {
	case "print":
		return ast.NewPrintCommand(parameters)
	case "println":
		return ast.NewPrintlnCommand(parameters)
	default:
		return &ParsingError{
			msg:    fmt.Sprintf("unknown function %q", symbol.GetText()),
			line:   symbol.GetLine(),
			column: symbol.GetColumn(),
		}
	}
	return &ParsingError{
		msg:    fmt.Sprintf("illegal state"),
		line:   symbol.GetLine(),
		column: symbol.GetColumn(),
	}
}

func (v *astVisitor) VisitParameterList(ctx *antlrgen.ParameterListContext) interface{} {
	debug.Debug("astVisitor.VisitParameterList")
	children := v.VisitChildren(ctx)
	for _, c := range children.([]interface{}) {
		if err, ok := c.(error); ok {
			return err
		}
	}
	return children
}

func (v *astVisitor) VisitVariable(ctx *antlrgen.VariableContext) interface{} {
	debug.Debug("astVisitor.VisitVariable")
	var expression ast.Expression
	if ctx.COLUMN() != nil {
		expression = ast.NewVarValue(ctx.COLUMN().GetSymbol().GetText())
	} else if ctx.IDENTIFIER() != nil {
		expression = ast.NewVarValue(ctx.IDENTIFIER().GetSymbol().GetText())
	}
	for _, c := range v.VisitChildren(ctx).([]interface{}) {
		if e, ok := c.(*ast.RangeExpression); ok {
			e.SetExpression(expression)
			expression = e
		}
	}
	debug.Debug("astVisitor.VisitVariable: resulting expression = %q", expression)
	return expression
}

func (v *astVisitor) VisitSlice(ctx *antlrgen.SliceContext) interface{} {
	debug.Debug("astVisitor.VisitSlice")
	var start, end *int
	if ctx.GetLeft() != nil {
		t := ctx.GetLeft().GetText()
		if n, err := strconv.Atoi(t); err != nil {
			return fmt.Errorf("unable to convert %q to integer", t)
		} else {
			start = &n
		}
	}
	if ctx.GetRight() != nil {
		t := ctx.GetRight().GetText()
		if n, err := strconv.Atoi(t); err != nil {
			return fmt.Errorf("unable to convert %q to integer", t)
		} else {
			end = &n
		}
	}
	e := ast.NewRangeExpression(nil, start, end)
	debug.Debug("VisitSlice: range expression = %s", e)
	return e
}

func getError(p interface{}) error {
	if err, ok := p.(error); ok {
		return err
	}
	if a, ok := p.([]interface{}); ok {
		for _, o := range a {
			if err, ok := o.(error); ok {
				return err
			}
		}
	}
	return nil
}
