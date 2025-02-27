package spider

import (
	"fmt"
)

func helpCommand() *Command {
	return &Command{
		Name:        "help",
		Description: "use 'help [command]' for command help",
		Usage:       "use 'help [command]' for command usage",
		builtin:     true,
		Args: func(args *Args) {
			args.StringList(&Arg{
				Name: "command",
				Help: "the name of the command",
			})
		},
		Run: func(c *Context) error {
			args, err := c.ArgValues.StringList("command")
			if err != nil {
				return err
			}
			if len(args) == 0 {
				c.Spider.PrintHelp()
				return nil
			}
			flagValues := make(FlagValues)
			command, _, err := c.Spider.Commands.parse(args, false, flagValues)
			if err != nil {
				return err
			} else if command == nil {
				return fmt.Errorf("command not found")
			}
			c.Spider.PrintCommandHelp(command)
			return nil
		},
	}
}

func spiderCommand() *Command {
	cmd := &Command{
		Name:        "spider",
		Description: "spider tool",
		Usage:       "spider <command> [flags] [args]",
		builtin:     true,
		Run: func(c *Context) error {
			c.Spider.PrintCommandHelp(c.Command)
			return nil
		},
	}
	cmd.AddCommand(&Command{
		Name:        "gps",
		Description: "displays process info",
		Usage:       "spider gps",
		builtin:     true,
		Run:         spiderGps,
	})
	cmd.AddCommand(&Command{
		Name:        "memory",
		Description: "displays process memory info",
		Usage:       "spider memory",
		builtin:     true,
		Run:         spiderGps,
	})
	return cmd
}

func exitCommand() *Command {
	return &Command{
		Name:        "exit",
		Description: "exit the shell",
		builtin:     true,
		Run: func(c *Context) error {
			return c.Stop()
		},
	}
}
