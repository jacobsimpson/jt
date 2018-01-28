package listener

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/parser"
)

type Block interface {
	Execute(environment map[string]string, line string, lineNumber int)
}

type printlnBlock struct{}

func (b *printlnBlock) Execute(environment map[string]string, line string, lineNumber int) {
	fmt.Println(line)
}

func NewPrintlnBlock() Block {
	return &printlnBlock{}
}

type InterpreterListener struct {
	Rules       []*rule
	currentRule *rule
}

func NewInterpreterListener() *InterpreterListener {
	return &InterpreterListener{}
}

func (l *InterpreterListener) VisitTerminal(node antlr.TerminalNode) {}
func (l *InterpreterListener) VisitErrorNode(node antlr.ErrorNode)   {}

func (l *InterpreterListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}
func (l *InterpreterListener) ExitEveryRule(ctx antlr.ParserRuleContext)  {}

func (l *InterpreterListener) EnterProgram(c *parser.ProgramContext) {}
func (l *InterpreterListener) ExitProgram(c *parser.ProgramContext)  {}

func (l *InterpreterListener) EnterProcessingRule(c *parser.ProcessingRuleContext) {
	l.currentRule = &rule{
		block: NewPrintlnBlock(),
	}
}
func (l *InterpreterListener) ExitProcessingRule(c *parser.ProcessingRuleContext) {
	if l.currentRule != nil {
		l.Rules = append(l.Rules, l.currentRule)
		l.currentRule = nil
	}
}

func (l *InterpreterListener) EnterSelection(c *parser.SelectionContext) {
	if c.REGULAR_EXPRESSION() != nil {
		re := c.REGULAR_EXPRESSION().GetSymbol().GetText()
		re = re[1 : len(re)-1]
		// TODO: error handling. Can't return it here, have to capture somehow.
		selection, _ := NewRegexpMatcher(re)
		l.currentRule.selection = selection
	}
}
func (s *InterpreterListener) ExitSelection(ctx *parser.SelectionContext) {}

func (l *InterpreterListener) EnterValue(c *parser.ValueContext) {
	if c.COLUMN() != nil {
	} else if c.REGULAR_EXPRESSION() != nil {
		re := c.REGULAR_EXPRESSION().GetSymbol().GetText()
		re = re[1 : len(re)-1]
		// TODO: error handling. Can't return it here, have to capture somehow.
		selection, _ := NewRegexpMatcher(re)
		l.currentRule.selection = selection
	} else if c.STRING() != nil {
	} else if c.DATE_TIME != nil {
	} else if c.INTEGER != nil {
	} else if c.HEX_INTEGER != nil {
	} else if c.BINARY_INTEGER != nil {
	}
}
func (s *InterpreterListener) ExitValue(ctx *parser.ValueContext) {}

func (l *InterpreterListener) EnterBlock(c *parser.BlockContext) {}
func (l *InterpreterListener) ExitBlock(c *parser.BlockContext)  {}

func (l *InterpreterListener) EnterCommand(c *parser.CommandContext) {}
func (l *InterpreterListener) ExitCommand(c *parser.CommandContext)  {}

func (l *InterpreterListener) EnterBinary(c *parser.BinaryContext)  {}
func (s *InterpreterListener) ExitBinary(ctx *parser.BinaryContext) {}

func (l *InterpreterListener) EnterBoolean(c *parser.BooleanContext)  {}
func (s *InterpreterListener) ExitBoolean(ctx *parser.BooleanContext) {}

func (l *InterpreterListener) EnterComparator(c *parser.ComparatorContext)  {}
func (s *InterpreterListener) ExitComparator(ctx *parser.ComparatorContext) {}

func (l *InterpreterListener) EnterExpression(c *parser.ExpressionContext)  {}
func (s *InterpreterListener) ExitExpression(ctx *parser.ExpressionContext) {}

func (l *InterpreterListener) EnterParameterList(c *parser.ParameterListContext) {}
func (l *InterpreterListener) ExitParameterList(c *parser.ParameterListContext)  {}

func (l *InterpreterListener) EnterComparison(c *parser.ComparisonContext) {}
func (l *InterpreterListener) ExitComparison(c *parser.ComparisonContext)  {}
