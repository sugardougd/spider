package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sugardougd/spider"
	"github.com/sugardougd/spider/example/commands"
)

func main() {
	config := &spider.TCPConfig{
		Config: &spider.Config{
			Name:        "spider",
			Description: "spider is a tool to list and diagnose Go processes",
			Prompt:      "spider >",
			Welcome:     "type 'help' for more information",
			ExecutedHook: func(ctx *spider.Context, err error) {
				fmt.Printf("Executed command: %s\r\n", ctx.Command.Name)
			},
		},
		Address: ":8080",
	}
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunTCP(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
