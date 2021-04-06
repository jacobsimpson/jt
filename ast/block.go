package ast

import "fmt"

// Block represents a block of commands in a program.
type Block struct {
	Commands []*Command
}

// NewPrintlnBlock is a convenience method for a block with a single println
// statement that prints the complete line.
func NewPrintlnBlock() *Block {
	return &Block{
		Commands: []*Command{
			&Command{
				Name:       "println",
				Parameters: []Expression{NewVarValue("%0")},
			},
		},
	}
}

func (b *Block) Execute(environment map[string]string) error {
	for _, command := range b.Commands {
		if err := command.Execute(environment); err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) LastCommand() *Command {
	return b.Commands[len(b.Commands)-1]
}

func (b *Block) String() string {
	return fmt.Sprintf("Block[%+v]", b.Commands)
}
