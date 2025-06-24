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
		spider.ConfigInteractive(true),
		spider.ConfigAddress(":8080"),
		spider.ConfigNoClientAuth(false),
		spider.ConfigBanner("welcome to spider!\n"),
		spider.ConfigPrivateFile("spider/ssh/spider"),
		spider.ConfigPasswordValidator(passwordValidator))
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunSSH(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}

func passwordValidator(user, password string) bool {
	return "admin" == user && "admin" == password
}
