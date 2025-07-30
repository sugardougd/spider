package main

import (
	"context"
	"fmt"
	"github.com/sugardougd/spider"
	"github.com/sugardougd/spider/example/commands"
	"time"
)

func main() {
	config := spider.NewConsoleConfig("spider", "spider is a tool to list and diagnose Go processes",
		"spider >", "type 'help' for more information", nil)
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunConsole(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
