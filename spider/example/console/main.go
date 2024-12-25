package main

import (
	"fmt"
	"os"
	"spider/spider"
)

func main() {
	s := spider.New(&spider.Config{
		Name:        "",
		Description: "",
		Prompt:      "",
	}, spider.NewCommands(testCommand()))
	cmd := []string{
		//"help",
		"test ",
		"test -a -b 1 -c 1.1 -d 1",
		"test --aaa --bbb 1 --ccc 1.1 --ddd 1",
		"test -a=false -b=1 -c=1.1 -d=1",
		"test --aaa=false --bbb=1 --ccc=1.1 --ddd=1",
		"test subtest",
		"test subtest -a -b 1 -c 1.1 -d 1",
		"test subtest --aaa --bbb 1 --ccc 1.1 --ddd 1",
		"test subtest -a=false -b=1 -c=1.1 -d=1",
		"test subtest --aaa=false --bbb=1 --ccc=1.1 --ddd=1",
		"test --aaa=false --bbb=1 --ccc=1.1 --ddd=1 subtest --aaa=false --bbb=1 --ccc=1.1 --ddd=1",
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
				Help:    "test flag",
				Usage:   "test usage",
				Default: false,
			})
			flags.Int(&spider.Flag{
				Short:   "b",
				Long:    "bbb",
				Help:    "test flag",
				Usage:   "test usage",
				Default: 99,
			})
			flags.Float64(&spider.Flag{
				Short:   "c",
				Long:    "ccc",
				Help:    "test flag",
				Usage:   "test usage",
				Default: 99.99,
			})
			flags.Float64(&spider.Flag{
				Short:   "d",
				Long:    "ddd",
				Help:    "test flag",
				Usage:   "test usage",
				Default: "hello spider",
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
				Help:    "test flag",
				Usage:   "test usage",
				Default: false,
			})
			flags.Int(&spider.Flag{
				Short:   "b",
				Long:    "bbb",
				Help:    "test flag",
				Usage:   "test usage",
				Default: 99,
			})
			flags.Float64(&spider.Flag{
				Short:   "c",
				Long:    "ccc",
				Help:    "test flag",
				Usage:   "test usage",
				Default: 99.99,
			})
			flags.Float64(&spider.Flag{
				Short:   "d",
				Long:    "ddd",
				Help:    "test flag",
				Usage:   "test usage",
				Default: "hello spider",
			})
		},
		Run: func(context *spider.Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println("------------------------------------------------------------")
			return nil
		},
	}
	cmd.AddChildren(subCmd)
	return cmd
}
