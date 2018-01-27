package main

import (
	"fmt"
	"regexp"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/parser"
)

type ValueType int

const (
	StringValue ValueType = iota
	RegularExpressionValue
	IntegerValue
	DateTimeValue
)

type Value interface {
	Type() ValueType
	Raw() string
}

type Selection interface {
	Evaluate(line string, lineNumber int) bool
	String() string
}

type regularExpressionMatcher struct {
	regexString string
	re          *regexp.Regexp
}

func NewRegularExpressionMatcher(regexString string) (Selection, error) {
	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}
	return &regularExpressionMatcher{
		regexString: regexString,
		re:          re,
	}, nil
}

func (m *regularExpressionMatcher) Evaluate(line string, lineNumber int) bool {
	return m.re.MatchString(line)
}

func (m *regularExpressionMatcher) String() string {
	return m.regexString
}

type Block interface {
	Execute(line string, lineNumber int)
}

type printlnBlock struct{}

func (b *printlnBlock) Execute(line string, lineNumber int) {
	fmt.Println(line)
}

func NewPrintlnBlock() Block {
	return &printlnBlock{}
}

type Expression struct {
	Selection Selection
	Block     Block
}

type InterpreterListener struct {
	Expressions       []*Expression
	currentExpression *Expression
}

func NewInterpreterListener() *InterpreterListener {
	return &InterpreterListener{}
}

func (l *InterpreterListener) VisitTerminal(node antlr.TerminalNode)      {}
func (l *InterpreterListener) VisitErrorNode(node antlr.ErrorNode)        {}
func (l *InterpreterListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}
func (l *InterpreterListener) ExitEveryRule(ctx antlr.ParserRuleContext)  {}

// EnterProgram is called when entering the program production.
func (l *InterpreterListener) EnterProgram(c *parser.ProgramContext) {}

// EnterExpression is called when entering the expression production.
func (l *InterpreterListener) EnterExpression(c *parser.ExpressionContext) {
	l.currentExpression = &Expression{
		Block: NewPrintlnBlock(),
	}
}

func (l *InterpreterListener) EnterSelection(c *parser.SelectionContext) {
	if c.REGULAR_EXPRESSION() != nil {
		re := c.REGULAR_EXPRESSION().GetSymbol().GetText()
		re = re[1 : len(re)-1]
		// TODO: error handling. Can't return it here, have to capture somehow.
		selection, _ := NewRegularExpressionMatcher(re)
		l.currentExpression.Selection = selection
	}
}

func (l *InterpreterListener) EnterValue(c *parser.ValueContext) {
	if c.COLUMN() != nil {
	} else if c.REGULAR_EXPRESSION() != nil {
		re := c.REGULAR_EXPRESSION().GetSymbol().GetText()
		re = re[1 : len(re)-1]
		// TODO: error handling. Can't return it here, have to capture somehow.
		selection, _ := NewRegularExpressionMatcher(re)
		l.currentExpression.Selection = selection
	} else if c.STRING() != nil {
	} else if c.DATE_TIME != nil {
	} else if c.INTEGER != nil {
	} else if c.HEX_INTEGER != nil {
	} else if c.BINARY_INTEGER != nil {
	}
}

// EnterBlock is called when entering the block production.
func (l *InterpreterListener) EnterBlock(c *parser.BlockContext) {}

// EnterCommand is called when entering the command production.
func (l *InterpreterListener) EnterCommand(c *parser.CommandContext) {}

// EnterParameter_list is called when entering the parameter_list production.
func (l *InterpreterListener) EnterParameter_list(c *parser.Parameter_listContext) {}

// ExitProgram is called when exiting the program production.
func (l *InterpreterListener) ExitProgram(c *parser.ProgramContext) {}

// ExitExpression is called when exiting the expression production.
func (l *InterpreterListener) ExitExpression(c *parser.ExpressionContext) {
	if l.currentExpression != nil {
		l.Expressions = append(l.Expressions, l.currentExpression)
		l.currentExpression = nil
	}
}

func (s *InterpreterListener) ExitSelection(ctx *parser.SelectionContext) {}
func (s *InterpreterListener) ExitValue(ctx *parser.ValueContext)         {}

// ExitBlock is called when exiting the block production.
func (l *InterpreterListener) ExitBlock(c *parser.BlockContext) {}

// ExitCommand is called when exiting the command production.
func (l *InterpreterListener) ExitCommand(c *parser.CommandContext) {}

// ExitParameter_list is called when exiting the parameter_list production.
func (l *InterpreterListener) ExitParameter_list(c *parser.Parameter_listContext) {}
