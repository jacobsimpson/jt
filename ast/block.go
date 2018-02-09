package ast

type Block interface {
	Execute(environment map[string]string)
	AddCommand(command Command)
	LastCommand() Command
}

func NewBlock() Block {
	return &block{
		commands: []Command{},
	}
}

func NewPrintlnBlock() Block {
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

func (b *block) AddCommand(command Command) {
	b.commands = append(b.commands, command)
}

func (b *block) LastCommand() Command {
	return b.commands[len(b.commands)-1]
}
