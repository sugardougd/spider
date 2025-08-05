package main

import (
	"context"
	"fmt"
	"github.com/sugardougd/spider"
	"github.com/sugardougd/spider/example/commands"
	"time"
)

func main() {
	config := &spider.SSHConfig{
		Config: &spider.Config{
			Name:        "spider",
			Description: "spider is a tool to list and diagnose Go processes",
			Prompt:      "spider >",
			Welcome:     "type 'help' for more information",
		},
		Address:           ":8080",
		NoClientAuth:      false,
		PasswordValidator: passwordValidator,
		Banner:            "welcome to spider!\n",
		PrivateFile:       "ssh/spider",
	}
	commands := spider.NewCommands(commands.NoyaCommand())

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	if err := spider.RunSSH(ctx, config, commands); err != nil {
		fmt.Println("exit spider", err.Error())
	}
}

func passwordValidator(user string, password []byte) bool {
	return "admin" == user && "admin" == string(password)
}
