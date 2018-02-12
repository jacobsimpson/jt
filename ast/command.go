package ast

import (
	"fmt"
	"os"
	"strings"
)

type Command interface {
	Execute(environment map[string]string)
	AddParameter(parameter Expression)
}

type command struct {
	name       string
	parameters []Expression
}

func (c *command) Execute(environment map[string]string) {
	fmt.Printf("there should be some generic function handler here ...\n")
}

func (c *command) AddParameter(parameter Expression) {
	c.parameters = append(c.parameters, parameter)
}

func NewPrintCommand(parameters []Expression) Command {
	return &printCommand{
		parameters: parameters,
		newline:    false,
	}
}

func NewPrintlnCommand(parameters []Expression) Command {
	return &printCommand{
		parameters: parameters,
		newline:    true,
	}
}

type printCommand struct {
	parameters []Expression
	newline    bool
}

func (c *printCommand) Execute(environment map[string]string) {
	formats := []string{}
	values := []interface{}{}
	for _, p := range c.parameters {
		formats = append(formats, "%s")
		v, err := p.Evaluate(environment)
		if err != nil {
			// TODO This is not real error handling. This should propagate up
			// the stack.
			fmt.Fprintf(os.Stderr, "could not evaluate parameter %s: %v", p, err)
			return
		}
		values = append(values, v)
	}
	format := strings.Join(formats, " ")
	if c.newline {
		format = format + "\n"
	}
	fmt.Printf(format, values...)
}

func (c *printCommand) AddParameter(parameter Expression) {
	c.parameters = append(c.parameters, parameter)
}
