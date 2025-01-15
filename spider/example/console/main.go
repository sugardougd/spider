package main

import (
	"fmt"
	"os"
	"spider/spider"
)

func main() {
	RunConsole()
}

func RunConsole() {
	s := spider.New(&spider.Config{
		Name:        "spider",
		Description: "spider is a tool to list and diagnose Go processes",
		Prompt:      "spider > ",
	}, spider.NewCommands(testCommand()))
	if err := s.RunConsole(); err != nil {
		fmt.Println("run spider fail", err.Error())
	}
}

func RunCommand() {
	s := spider.New(&spider.Config{
		Name:        "spider",
		Description: "spider is a tool to list and diagnose Go processes",
		Prompt:      "spider > ",
	}, spider.NewCommands(testCommand()))
	cmd := []string{
		//"help",
		"help test",
		//"test",
		//"test -a -b 1 -c 1.1 -d 1 false 2 3 4",
		//"test --aaa --bbb 1 --ccc 1.1 --ddd 1 true 2 3 4",
		//"test -a=false -b=1 -c=1.1 -d=1 true 2",
		//"test --aaa=false --bbb=1 --ccc=1.1 --ddd=1 true",
		//"test subtest false 2 3 4",
		//"test subtest -a -b 1 -c 1.1 -d 1 false 2 3 4",
		//"test subtest --aaa --bbb 1 --ccc 1.1 --ddd 1 true 2 3 4",
		//"test subtest -a=false -b=1 -c=1.1 -d=1 true 2",
		//"test subtest --aaa=false --bbb=1 --ccc=1.1 --ddd=1 true",
		//"test --aaa=false --bbb=1 --ccc=1.1 --ddd=1 subtest --aaa=false --bbb=1 --ccc=1.1 --ddd=1",
		//"err",
		//"exit",
	}
	for _, c := range cmd {
		if err := s.RunCommand(c); err != nil {
			fmt.Println("Run fail: ", err)
			os.Exit(0)
		}
	}
}

func testCommand() *spider.Command {
	cmd := &spider.Command{
		Name:        "test",
		Description: "use 'test [command]' for command help",
		Usage:       "use 'test [command]' for command usage",
		Flags: func(flags *spider.Flags) {
			flags.Bool(&spider.Flag{
				Short:   "a",
				Long:    "aaa",
				Help:    "aaa flag",
				Usage:   "aaa usage",
				Require: false,
				Default: false,
			})
			flags.Int(&spider.Flag{
				Short:   "b",
				Long:    "bbb",
				Help:    "bbb flag",
				Usage:   "bbb usage",
				Default: 99,
			})
			flags.Float64(&spider.Flag{
				Short:   "c",
				Long:    "ccc",
				Help:    "ccc flag",
				Usage:   "ccc usage",
				Default: 99.99,
			})
			flags.Float64(&spider.Flag{
				Short:   "d",
				Long:    "ddd",
				Help:    "ddd flag",
				Usage:   "ddd usage",
				Default: "hello spider",
			})
		},
		Args: func(args *spider.Args) {
			args.Bool(&spider.Arg{
				Name:    "e",
				Help:    "e argument",
				Default: false,
				Require: false,
			})
			args.Int(&spider.Arg{
				Name:    "f",
				Help:    "f argument",
				Default: 1,
				Require: false,
			})
			args.Float64(&spider.Arg{
				Name:    "g",
				Help:    "g argument",
				Default: 1.1,
				Require: false,
			})
			args.String(&spider.Arg{
				Name:    "h",
				Help:    "h argument",
				Default: "1",
				Require: false,
			})
		},
		Run: func(context *spider.Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}
	subCmd := &spider.Command{
		Name:        "subtest",
		Description: "use 'test [command]' for command help",
		Usage:       "use 'test [command]' for command usage",
		Flags: func(flags *spider.Flags) {
			flags.Bool(&spider.Flag{
				Short:   "a",
				Long:    "aaa",
				Help:    "aaa flag",
				Usage:   "aaa usage",
				Require: false,
				Default: false,
			})
			flags.Int(&spider.Flag{
				Short:   "b",
				Long:    "bbb",
				Help:    "bbb flag",
				Usage:   "bbb usage",
				Default: 99,
			})
			flags.Float64(&spider.Flag{
				Short:   "c",
				Long:    "ccc",
				Help:    "ccc flag",
				Usage:   "ccc usage",
				Default: 99.99,
			})
			flags.Float64(&spider.Flag{
				Short:   "d",
				Long:    "ddd",
				Help:    "ddd flag",
				Usage:   "ddd usage",
				Default: "hello spider",
			})
		},
		Args: func(args *spider.Args) {
			args.Bool(&spider.Arg{
				Name:    "e",
				Help:    "e argument",
				Default: false,
				Require: false,
			})
			args.Int(&spider.Arg{
				Name:    "f",
				Help:    "f argument",
				Default: 1,
				Require: false,
			})
			args.Float64(&spider.Arg{
				Name:    "g",
				Help:    "g argument",
				Default: 1.1,
				Require: false,
			})
			args.String(&spider.Arg{
				Name:    "h",
				Help:    "h argument",
				Default: "1",
				Require: false,
			})
		},
		Run: func(context *spider.Context) error {
			return context.Command.Parent.Run(context)
		},
	}
	cmd.AddChildren(subCmd)
	return cmd
}
