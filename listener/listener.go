package listener

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/parser"
)

type Block interface {
	Execute(environment map[string]string)
}

func NewPrintlnBlock() *block {
	return &block{
		commands: []Command{
			&printCommand{
				parameters: []string{"%0"},
				newline:    true,
			},
		},
	}
}

type block struct {
	commands []Command
}

func (b *block) Execute(environment map[string]string) {
	for _, command := range b.commands {
		command.Execute(environment)
	}
}

type Command interface {
	Execute(environment map[string]string)
	AddParameter(parameter string)
}

type command struct {
	name       string
	parameters []string
}

func (c *command) Execute(environment map[string]string) {
	fmt.Printf("there should be some generic function handler here ...\n")
}

func (c *command) AddParameter(parameter string) {
	c.parameters = append(c.parameters, parameter)
}

type printCommand struct {
	parameters []string
	newline    bool
}

func (c *printCommand) Execute(environment map[string]string) {
	formats := []string{}
	values := []interface{}{}
	for _, p := range c.parameters {
		formats = append(formats, "%s")
		values = append(values, environment[p])
	}
	format := strings.Join(formats, " ")
	if c.newline {
		format = format + "\n"
	}
	fmt.Printf(format, values...)
}

func (c *printCommand) AddParameter(parameter string) {
	c.parameters = append(c.parameters, parameter)
}

type InterpreterListener struct {
	Rules       []*rule
	currentRule *rule
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
	l.currentRule = &rule{
		block: NewPrintlnBlock(),
	}
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
		value, err := NewRegexpValue(regexpString)
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    fmt.Sprintf("could not parse regular expression %q: %v", regexpString, err),
				line:   ctx.REGULAR_EXPRESSION().GetSymbol().GetLine(),
				column: ctx.REGULAR_EXPRESSION().GetSymbol().GetColumn(),
			})
			return
		}
		l.currentRule.selection = &comparison{
			left:     &varValue{name: "%0"},
			operator: EQ_Operator,
			right:    value,
		}
	}
}
func (l *InterpreterListener) ExitSelection(ctx *parser.SelectionContext) {}

func (l *InterpreterListener) EnterComparison(ctx *parser.ComparisonContext) {
	l.currentRule.selection = &comparison{}
}
func (l *InterpreterListener) ExitComparison(ctx *parser.ComparisonContext) {}

func (l *InterpreterListener) EnterComparator(ctx *parser.ComparatorContext) {
	cmp := l.currentRule.selection.(*comparison)
	if ctx.LT() != nil {
		cmp.operator = LT_Operator
	} else if ctx.LE() != nil {
		cmp.operator = LE_Operator
	} else if ctx.EQ() != nil {
		cmp.operator = EQ_Operator
	} else if ctx.NE() != nil {
		cmp.operator = NE_Operator
	} else if ctx.GE() != nil {
		cmp.operator = GE_Operator
	} else if ctx.GT() != nil {
		cmp.operator = GT_Operator
	}
}
func (l *InterpreterListener) ExitComparator(ctx *parser.ComparatorContext) {}

func (l *InterpreterListener) EnterValue(ctx *parser.ValueContext) {
	var value Value
	if ctx.COLUMN() != nil {
		value = &varValue{
			name: ctx.COLUMN().GetSymbol().GetText(),
		}
	} else if ctx.REGULAR_EXPRESSION() != nil {
		regexpString := ctx.REGULAR_EXPRESSION().GetSymbol().GetText()
		regexpString = regexpString[1 : len(regexpString)-1]
		var err error
		value, err = NewRegexpValue(regexpString)
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    fmt.Sprintf("could not parse regular expression %q: %v", regexpString, err),
				line:   ctx.REGULAR_EXPRESSION().GetSymbol().GetLine(),
				column: ctx.REGULAR_EXPRESSION().GetSymbol().GetColumn(),
			})
			return
		}
	} else if ctx.STRING() != nil {
		value = NewStringValue(ctx.STRING().GetSymbol().GetText())
	} else if ctx.INTEGER() != nil {
		var err error
		value, err = NewIntegerValue(ctx.INTEGER().GetSymbol().GetText())
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
		value, err = NewDateTimeValue(ctx.DATE_TIME().GetSymbol().GetText())
		if err != nil {
			l.Errors = append(l.Errors, &ParsingError{
				msg:    err.Error(),
				line:   ctx.DATE_TIME().GetSymbol().GetLine(),
				column: ctx.DATE_TIME().GetSymbol().GetColumn(),
			})
			return
		}
	}
	cmp := l.currentRule.selection.(*comparison)
	if cmp.left == nil {
		cmp.left = value
	} else {
		cmp.right = value
	}
}
func (l *InterpreterListener) ExitValue(ctx *parser.ValueContext) {}

func (l *InterpreterListener) EnterBlock(ctx *parser.BlockContext) {
	l.currentRule.block = &block{}
}
func (l *InterpreterListener) ExitBlock(ctx *parser.BlockContext) {}

func (l *InterpreterListener) EnterCommand(ctx *parser.CommandContext) {
	l.currentRule.block.commands = append(l.currentRule.block.commands,
		&printCommand{
			parameters: []string{},
			newline:    true,
		})
}
func (l *InterpreterListener) ExitCommand(ctx *parser.CommandContext) {}

func (l *InterpreterListener) EnterParameterList(ctx *parser.ParameterListContext) {
	currentCommand := l.currentRule.block.commands[len(l.currentRule.block.commands)-1]
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
