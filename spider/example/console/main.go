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
	}, func() *spider.Commands {
		commands := spider.Commands{}
		commands.Add(&spider.Command{
			Name:        "test",
			Description: "use 'test [command]' for command help",
			Usage:       "use 'test [command]' for command usage",
			Flags: func(flags *spider.Flags) {
				flags.Bool(&spider.Flag{
					Short:   "h",
					Long:    "help",
					Usage:   "use help",
					Default: flags,
				})
			},
			Run: func(c *spider.Context) error {
				c.Spider.Printf("exec '%s'\n", c.Command.Name)
				return nil
			},
		})
		return &commands
	})
	cmd := []string{"help", "test", "err", "exit"}
	for _, c := range cmd {
		if err := s.RunCommand(c); err != nil {
			fmt.Println("Run fail: ", err)
		}
	}
}
