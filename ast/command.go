package ast

import (
	"fmt"
	"strings"
)

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

func NewPrintCommand() Command {
	return &printCommand{
		parameters: []string{},
		newline:    false,
	}
}

func NewPrintlnCommand() Command {
	return &printCommand{
		parameters: []string{},
		newline:    true,
	}
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
