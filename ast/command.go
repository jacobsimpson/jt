package ast

import (
	"fmt"
	"strings"
)

type Command struct {
	Name       string
	Parameters []Expression
}

func (c *Command) Execute(environment map[string]string) error {
	switch c.Name {
	case "println", "print":
		formats := []string{}
		values := []interface{}{}
		for _, p := range c.Parameters {
			formats = append(formats, "%s")
			v, err := p.Evaluate(environment)
			if err != nil {
				return fmt.Errorf("could not evaluate parameter %s: %v", p, err)
			}
			values = append(values, v)
		}
		format := strings.Join(formats, " ")
		if c.Name == "println" {
			format = format + "\n"
		}
		fmt.Printf(format, values...)
	default:
		return fmt.Errorf("unknown function %q: 1:11", c.Name)
	}
	return nil
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
