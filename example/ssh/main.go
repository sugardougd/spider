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
		spider2.ConfigPrompt("spider > "),
		spider2.ConfigInteractive(true),
		spider2.ConfigAddress(":8080"),
		spider2.ConfigNoClientAuth(false),
		spider2.ConfigBanner("welcome to spider!\n"),
		spider2.ConfigPrivateFile("spider/ssh/spider"),
		spider2.ConfigPasswordValidator(passwordValidator))
	commands := spider2.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider2.RunSSH(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}

func passwordValidator(user, password string) bool {
	return "admin" == user && "admin" == password
}
