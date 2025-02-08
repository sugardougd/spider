package commands

import "spider/spider"

func TestCommand() *spider.Command {
	cmd := &spider.Command{
		Name:        "test",
		Description: "test command",
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
		Description: "subtest command",
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
	cmd.AddCommand(subCmd)
	return cmd
}

func Test1Command() *spider.Command {
	cmd := &spider.Command{
		Name:        "test1",
		Description: "test1 command",
		Flags: func(flags *spider.Flags) {
			flags.Bool(&spider.Flag{
				Short:   "a",
				Long:    "aaa",
				Help:    "aaa flag",
				Usage:   "aaa usage",
				Require: false,
				Default: false,
			})
		},
		Run: func(context *spider.Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}
	return cmd
}
