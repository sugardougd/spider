package main

import (
	"context"
	"fmt"
	"spider/spider"
	"spider/spider/example/commands"
	"time"
)

func main() {
	config := &spider.SSHConfig{
		TCPConfig: &spider.TCPConfig{
			Address: ":8080",
			Config: &spider.Config{
				Name:        "spider",
				Description: "spider is a tool to list and diagnose Go processes",
				Prompt:      "spider > ",
			},
		},
		NoClientAuth: false,
		Banner:       "welcome to spider!\r\n",
		PrivateFile:  "spider/ssh/spider",
		PasswordFunc: func(user, password string) bool {
			return "admin" == user && "admin" == password
		},
	}
	commands := spider.NewCommands(commands.TestCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunSSH(config, commands, ctx); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}
