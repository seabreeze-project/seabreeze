package shell

type CommandExecutor func(args []string) error

type Command struct {
	Name    string
	Desc    string
	Long    string
	Aliases []string
	Run     CommandExecutor
}

func (c *Command) Execute(args []string) error {
	return c.Run(args)
}
