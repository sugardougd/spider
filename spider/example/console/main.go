package main

import (
	"context"
	"fmt"
	"spider/spider"
	"spider/spider/example/commands"
	"time"
)

func main() {
	config := spider.NewConfig(
		spider.ConfigName("spider"),
		spider.ConfigDescription("spider is a tool to list and diagnose Go processes"),
		spider.ConfigInteractive(true),
		spider.ConfigPrompt("spider > "))
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunConsole(config, commands, ctx); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
