package main

import (
	"context"
	"fmt"
	"github.com/sugardougd/spider"
	"github.com/sugardougd/spider/example/commands"
	"time"
)

func main() {
	config := spider.NewConfig(
		spider.ConfigName("spider"),
		spider.ConfigDescription("spider is a tool to list and diagnose Go processes"),
		spider.ConfigPrompt("spider > "),
		spider.ConfigAddress(":8080"),
		spider.ConfigExecutedHook(func(ctx *spider.Context, err error) {
			fmt.Printf("Executed command: %s\r\n", ctx.Command.Name)
		}))
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunTCP(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
