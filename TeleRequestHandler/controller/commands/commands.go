package commands

type Command struct {
	Name       string
	Decription string
	CommandHandler[any]
}

func NewCommand(name string, description string, handler CommandHandler[any]) *Command {
	return &Command{name, description, handler}
}
