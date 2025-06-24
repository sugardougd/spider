package main

import (
	"context"
	"fmt"
	spider2 "spider"
	"spider/example/commands"
	"time"
)

func main() {
	config := spider2.NewConfig(
		spider2.ConfigName("spider"),
		spider2.ConfigDescription("spider is a tool to list and diagnose Go processes"),
		spider2.ConfigInteractive(true),
		spider2.ConfigPrompt("spider > "))
	commands := spider2.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider2.RunConsole(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
