package ast

// Block represents a block of commands in a program.
type Block struct {
	Commands []Command
}

// NewPrintlnBlock is a convenience method for a block with a single println
// statement that prints the complete line.
func NewPrintlnBlock() *Block {
	return &Block{
		Commands: []Command{
			&printCommand{
				parameters: []Expression{NewVariableExpression("%0")},
				newline:    true,
			},
		},
	}
}

func (b *Block) Execute(environment map[string]string) {
	for _, command := range b.Commands {
		command.Execute(environment)
	}
}

func (b *Block) LastCommand() Command {
	return b.Commands[len(b.Commands)-1]
}
