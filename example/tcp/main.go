package main

import (
	"context"
	"fmt"
	"github.com/sugardougd/spider"
	"github.com/sugardougd/spider/example/commands"
	"time"
)

func main() {
	config := spider.NewTCPConfig("spider", "spider is a tool to list and diagnose Go processes",
		"spider >", "type 'help' for more information", func(ctx *spider.Context, err error) {
			fmt.Printf("Executed command: %s\r\n", ctx.Command.Name)
		}, spider.TCPConfig{Address: ":8080"})
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunTCP(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
