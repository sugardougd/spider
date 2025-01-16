package main

import (
	"context"
	"fmt"
	"spider/spider"
	"spider/spider/example/commands"
	"time"
)

func main() {
	config := &spider.TCPConfig{
		Address: ":8080",
		Config: &spider.Config{
			Name:        "spider",
			Description: "spider is a tool to list and diagnose Go processes",
			Prompt:      "spider > ",
		},
	}
	commands := spider.NewCommands(commands.TestCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	if err := spider.RunTCP(config, commands, ctx); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
