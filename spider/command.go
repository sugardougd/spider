package spider

import (
	"fmt"
)

type CommandsFunc func() *Commands

type Commands struct {
	list []*Command
}

func NewCommands(command ...*Command) *Commands {
	commands := Commands{}
	for _, c := range command {
		commands.Add(c)
	}
	return &commands
}

func (commands *Commands) Add(command *Command) error {
	err := command.validate()
	if err != nil {
		return err
	}
	command.registerFlags()
	command.registerArgs()
	commands.list = append(commands.list, command)
	return nil
}

func (commands *Commands) Find(name string) *Command {
	for _, command := range commands.list {
		if command.Name == name {
			return command
		}
		for _, a := range command.Aliases {
			if a == name {
				return command
			}
		}
	}
	return nil
}

func (commands *Commands) parse(args []string, required bool, flagValues FlagValues) (command *Command, remaining []string, err error) {
	cur := command
	for len(args) > 0 {
		// find command
		if cur == nil {
			cur = commands.Find(args[0])
		} else {
			cur = cur.FindChildren(args[0])
		}
		if cur == nil {
			break
		}
		command = cur
		args = args[1:]
		// parse flags
		args, err = command.flags.parse(command, args, required, flagValues)
		if err != nil {
			break
		}
	}
	remaining = args
	return
}

type Command struct {
	Name        string                       // name
	Aliases     []string                     // name aliases.
	Description string                       // description message for the command.
	Usage       string                       // define how to use the command.Sample: start [OPTIONS] CONTAINER [CONTAINER...]
	Flags       func(flags *Flags)           // define all command flags within this function.
	Args        func(args *Args)             // define all command arguments within this function.
	Run         func(context *Context) error // function to execute for the command.
	Parent      *Command
	Children    Commands
	flags       Flags
	args        Args
	builtin     bool // Whenever this is a build-in command not added by the user.
}

func (command *Command) AddChildren(cmd *Command) error {
	err := command.Children.Add(cmd)
	if err == nil {
		cmd.Parent = command
	}
	return err
}

func (command *Command) FindChildren(name string) *Command {
	return command.Children.Find(name)
}

func (command *Command) RegisterFlags(flag func(c *Command, f *Flags) error) error {
	return flag(command, &command.flags)
}

func (command *Command) FullName() string {
	if command.Parent != nil {
		return command.Parent.FullName() + "." + command.Name
	}
	return command.Name
}

func (command *Command) registerFlags() {
	if command.Flags != nil {
		command.Flags(&command.flags)
	}
}

func (command *Command) registerArgs() {
	if command.Args != nil {
		command.Args(&command.args)
	}
}

func (command *Command) validate() error {
	if len(command.Name) == 0 {
		return fmt.Errorf("empty command name")
	} else if command.Name[0] == '-' {
		return fmt.Errorf("invalid command name '%s' must not start with a '-'", command.Name)
	} else if len(command.Description) == 0 {
		return fmt.Errorf("empty command description")
	}
	return nil
}
