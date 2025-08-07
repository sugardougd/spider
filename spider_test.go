package spider

import (
	"testing"
)

func TestSpider_parse(t *testing.T) {
	s := NewSpiderMock()

	wrongCmd := []string{
		"hello",
		"command -f",
		"command --flag",
		"command -o",
		"command arg",
		"command sub-command1 arg1 arg2",
		"command sub-command2",
	}
	for _, c := range wrongCmd {
		command, flagValues, argValues, err := s.parse(c, true)
		if command != nil || flagValues != nil || argValues != nil || err == nil {
			t.Fatalf("parse command fail of '%s'", c)
		}
	}

	correctCmd := []string{
		"command",
		"command -h",
		"command -o flag1",
		"command --flag1 flag1",
		"command -t 2",
		"command --flag2 2",
		"command -o flag1 -t 2",
		"command --flag1 flag1 -t 2",
		"command -o flag1 --flag2 2",
		"command --flag1 flag1 --flag2 2",

		"command sub-command1",
		"command sub-command1 -o flag1",
		"command sub-command1 --flag1 flag1",
		"command sub-command1 -t 2",
		"command sub-command1 --flag2 2",
		"command sub-command1 -o flag1 -t 2",
		"command sub-command1 --flag1 flag1 -t 2",
		"command sub-command1 -o flag1 --flag2 2",
		"command sub-command1 --flag1 flag1 --flag2 2",

		"command -o flag1 sub-command1",
		"command --flag1 flag1 sub-command1",
		"command -t 2 sub-command1",
		"command --flag2 2 sub-command1",
		"command -o flag1 -t 2 sub-command1",
		"command --flag1 flag1 -t 2 sub-command1",
		"command -o flag1 --flag2 2 sub-command1",
		"command --flag1 flag1 --flag2 2 sub-command1",

		"command -o flag1 sub-command1 -o flag1",
		"command --flag1 flag1 sub-command1 --flag1 flag1",
		"command -t 2 sub-command1 -t 2",
		"command --flag2 2 sub-command1 --flag2 2",
		"command -o flag1 -t 2 sub-command1 -o flag1 -t 2",
		"command --flag1 flag1 -t 2 sub-command1 --flag1 flag1 -t 2",
		"command -o flag1 --flag2 2 sub-command1 -o flag1 --flag2 2",
		"command --flag1 flag1 --flag2 2 sub-command1 --flag1 flag1 --flag2 2",

		"command -o flag1 -t 2 sub-command2 -o flag1 -t 2",
		"command --flag1 flag1 -t 2 sub-command2 --flag1 flag1 -t 2",
		"command -o flag1 --flag2 2 sub-command2 -o flag1 --flag2 2",
		"command --flag1 flag1 --flag2 2 sub-command2 --flag1 flag1 --flag2 2",
	}
	for _, c := range correctCmd {
		command, flagValues, argValues, err := s.parse(c, true)
		if command == nil || flagValues == nil || argValues == nil || err != nil {
			t.Fatalf("parse command fail of '%s' %v", c, err)
		}
	}

}

func testCommand() *Commands {
	command := &Command{
		Name:        "command",
		Description: "test command",
		Flags: func(flags *Flags) {
			flags.String(&Flag{
				Short:   "o",
				Long:    "flag1",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: "default flag",
			})
			flags.Int(&Flag{
				Short:   "t",
				Long:    "flag2",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: 999,
			})
		},
		Run: func(context *Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}

	command1 := &Command{
		Name:        "command1",
		Description: "test command1",
		Flags: func(flags *Flags) {
			flags.String(&Flag{
				Short:   "o",
				Long:    "flag1",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: "default flag",
			})
			flags.Int(&Flag{
				Short:   "t",
				Long:    "flag2",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: 999,
			})
		},
		Run: func(context *Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}

	subCommand := &Command{
		Name:        "sub-command",
		Description: "test command",
		Flags: func(flags *Flags) {
			flags.String(&Flag{
				Short:   "o",
				Long:    "flag1",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: "default flag",
			})
			flags.Int(&Flag{
				Short:   "t",
				Long:    "flag2",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: 999,
			})
		},
		Args: func(args *Args) {
			args.String(&Arg{
				Name:    "arg",
				Help:    "test arg",
				Default: "default arg",
				Require: false,
			})
		},
		Run: func(context *Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}

	subCommand1 := &Command{
		Name:        "sub-command1",
		Description: "test command",
		Flags: func(flags *Flags) {
			flags.String(&Flag{
				Short:   "o",
				Long:    "flag1",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: "default flag",
			})
			flags.Int(&Flag{
				Short:   "t",
				Long:    "flag2",
				Help:    "test flag",
				Usage:   "test flag",
				Require: false,
				Default: 999,
			})
		},
		Args: func(args *Args) {
			args.String(&Arg{
				Name:    "arg",
				Help:    "test arg",
				Default: "default arg",
				Require: false,
			})
		},
		Run: func(context *Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}
	subCommand2 := &Command{
		Name:        "sub-command2",
		Description: "test command",
		Flags: func(flags *Flags) {
			flags.String(&Flag{
				Short:   "o",
				Long:    "flag1",
				Help:    "test flag",
				Usage:   "test flag",
				Require: true,
				Default: "default flag",
			})
			flags.Int(&Flag{
				Short:   "t",
				Long:    "flag2",
				Help:    "test flag",
				Usage:   "test flag",
				Require: true,
				Default: 999,
			})
		},
		Args: func(args *Args) {
			args.String(&Arg{
				Name:    "arg",
				Help:    "test arg",
				Default: "default arg",
				Require: false,
			})
		},
		Run: func(context *Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}
	command.AddCommand(subCommand)
	command.AddCommand(subCommand1)
	command.AddCommand(subCommand2)

	command1.AddCommand(subCommand)
	command1.AddCommand(subCommand1)
	command1.AddCommand(subCommand2)

	return NewCommands(command, command1)
}

func NewSpiderMock() *Spider {
	config := Config{
		Name:        "spider",
		Description: "spider is a tool to list and diagnose Go processes",
		Prompt:      "spider >",
	}
	s := New(&config, testCommand())
	return s
}
