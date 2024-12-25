package main

import (
	"fmt"
	"spider/grumble"
	"spider_v1"
	"time"
)

func main() {
	RunConsole()
	//RunTCP()
	//RunSSH()
}

func RunConsole() {
	cfg := &spider_v1.ConsoleConfig{
		Config: &grumble.Config{
			Name:        "spider_v1",
			Prompt:      "spider_v1 > ",
			Description: "spider_v1 is a tool to list and diagnose Go processes",
		},
		Commands: func(commands *grumble.Commands) *grumble.Commands {
			return spider_v1.NoyaCommands()
		},
	}
	if err := spider_v1.RunConsole(cfg); err != nil {
		fmt.Printf("Start() error = %v", err)
	}
}

func RunTCP() {
	cfg := &spider_v1.TCPConfig{
		Address: ":4322",
		Config: &grumble.Config{
			Name:        "spider_v1",
			Prompt:      "spider_v1 > ",
			Description: "spider_v1 is a tool to list and diagnose Go processes",
			FuncIsTerminal: func() bool {
				return true
			},
		},
		Commands: func(commands *grumble.Commands) *grumble.Commands {
			return timeCommand(commands)
		},
	}
	if err := spider_v1.RunTCP(cfg); err != nil {
		fmt.Printf("Start() error = %v", err)
	}
}

func RunSSH() {
	cfg := &spider_v1.SSHConfig{
		Address:      ":4322",
		NoClientAuth: true,
		Banner:       "welcome to spider_v1!\r\n",
		PrivateFile:  "ssh/spider_v1",
		PasswordFunc: func(user, password string) bool {
			return "admin" == user && "admin" == password
		},
		Config: &grumble.Config{
			Name:        "spider_v1",
			Prompt:      "spider_v1 > ",
			Description: "spider_v1 is a tool to list and diagnose Go processes",
		},
		Commands: func(commands *grumble.Commands) *grumble.Commands {
			return timeCommand(commands)
		},
	}
	if err := spider_v1.RunSSH(cfg); err != nil {
		fmt.Printf("Start() error = %v", err)
	}
}

func timeCommand(commands *grumble.Commands) *grumble.Commands {
	tc := grumble.Commands{}
	commands.Add(&grumble.Command{
		Name: "time",
		Help: "displays time now",
		Run: func(c *grumble.Context) error {
			c.App.Println(time.Now())
			return nil
		},
	})
	return &tc
}
