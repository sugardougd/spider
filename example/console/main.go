package main

import (
	"context"
	"fmt"
	"github.com/sugardougd/spider"
	"github.com/sugardougd/spider/example/commands"
	"time"
)

func main() {
	config :=
		spider.NewConfig(
			spider.ConfigName("spider"),
			spider.ConfigDescription("spider is a tool to list and diagnose Go processes"),
			spider.ConfigInteractive(true),
			spider.ConfigPrompt("spider > "),
			spider.ConfigWelcome("type 'help' for more information"))
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunConsole(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
