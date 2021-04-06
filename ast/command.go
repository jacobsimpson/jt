package ast

import (
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Name       string
	Parameters []Expression
}

func (c *Command) Execute(environment map[string]string) {
	switch c.Name {
	case "println", "print":
		formats := []string{}
		values := []interface{}{}
		for _, p := range c.Parameters {
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
		if c.Name == "println" {
			format = format + "\n"
		}
		fmt.Printf(format, values...)
	}
}

func (c *Command) AddParameter(parameter Expression) {
	c.Parameters = append(c.Parameters, parameter)
}

func NewPrintCommand(parameters []Expression) *Command {
	return &Command{
		Name:       "print",
		Parameters: parameters,
	}
}

func NewPrintlnCommand(parameters []Expression) *Command {
	return &Command{
		Name:       "println",
		Parameters: parameters,
	}
}
