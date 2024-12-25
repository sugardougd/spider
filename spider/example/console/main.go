package main

import (
	"fmt"
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
		"test -a",
		"test --aa",
		"test -a=false",
		"test --aa=false",
		"test subtest",
		"test -a subtest -b",
		"test --aa subtest --bb",
		"test -a=false subtest -b=false",
		"test --aa=false subtest --bb=false",
		//"err",
		//"exit",
	}
	for _, c := range cmd {
		if err := s.RunCommand(c); err != nil {
			fmt.Println("Run fail: ", err)
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
				Long:    "aa",
				Help:    "test flag",
				Usage:   "test usage",
				Default: false,
			})
		},
		Run: func(context *spider.Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println()
			return nil
		},
	}
	subCmd := &spider.Command{
		Name:        "subtest",
		Description: "use 'test [command]' for command help",
		Usage:       "use 'test [command]' for command usage",
		Flags: func(flags *spider.Flags) {
			flags.Bool(&spider.Flag{
				Short:   "b",
				Long:    "bb",
				Help:    "test flag",
				Usage:   "test usage",
				Default: false,
			})
		},
		Run: func(context *spider.Context) error {
			context.Spider.Println(context.String())
			context.Spider.Println()
			return nil
		},
	}
	cmd.AddChildren(subCmd)
	return cmd
}
