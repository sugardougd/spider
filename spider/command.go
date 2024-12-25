package spider

import (
	"fmt"
)

type CommandsFunc func() *Commands

type Commands struct {
	list []*Command
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

func (commands *Commands) Parse(args []string) (command *Command, flagValues FlagValues, err error) {
	for len(args) > 0 {
		// find command
		cur := commands.Find(args[0])
		if cur == nil {
			err = fmt.Errorf("illagel command '%s'", args[0])
			return
		}
		command = cur
		args = args[1:]
		// parse flags
		args, err = command.flags.parse(command, args, flagValues)
		if err != nil {
			return
		}
	}
	return
}

type Command struct {
	Name        string                    // name
	Aliases     []string                  // name aliases.
	Description string                    // description message for the command.
	Usage       string                    // define how to use the command.Sample: start [OPTIONS] CONTAINER [CONTAINER...]
	Flags       func(f *Flags)            // define all command flags within this function.
	Args        func(c *Command, a *Args) // define all command arguments within this function.
	Run         func(c *Context) error    // function to execute for the command.
	flags       Flags
	args        Args
	parent      *Command
	children    Commands
	builtin     bool // Whenever this is a build-in command not added by the user.
}

func (command *Command) AddChildren(cmd *Command) error {
	err := command.children.Add(cmd)
	if err == nil {
		cmd.parent = command
	}
	return err
}

func (command *Command) RegisterFlags(flag func(c *Command, f *Flags) error) error {
	return flag(command, &command.flags)
}

func (command *Command) registerFlags() {
	if command.Flags != nil {
		command.Flags(&command.flags)
	}
}

func (command *Command) registerArgs() {
	if command.Args != nil {
		command.Args(command, &command.args)
	}
}

func (command *Command) fullName() string {
	if command.parent != nil {
		return command.parent.fullName() + "." + command.Name
	}
	return command.Name
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
