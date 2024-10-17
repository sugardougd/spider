package main

import (
	"fmt"
	"spider"
	"spider/grumble"
	"time"
)

func main() {
	RunConsole()
	//RunTCP()
	//RunSSH()
}

func RunConsole() {
	cfg := &spider.ConsoleConfig{
		Config: &grumble.Config{
			Name:        "spider",
			Prompt:      "spider > ",
			Description: "spider is a tool to list and diagnose Go processes",
		},
		Commands: func(commands *grumble.Commands) *grumble.Commands {
			return spider.NoyaCommands()
		},
	}
	if err := spider.RunConsole(cfg); err != nil {
		fmt.Printf("Start() error = %v", err)
	}
}

func RunTCP() {
	cfg := &spider.TCPConfig{
		Address: ":4322",
		Config: &grumble.Config{
			Name:        "spider",
			Prompt:      "spider > ",
			Description: "spider is a tool to list and diagnose Go processes",
			FuncIsTerminal: func() bool {
				return true
			},
		},
		Commands: func(commands *grumble.Commands) *grumble.Commands {
			return timeCommand(commands)
		},
	}
	if err := spider.RunTCP(cfg); err != nil {
		fmt.Printf("Start() error = %v", err)
	}
}

func RunSSH() {
	cfg := &spider.SSHConfig{
		Address:      ":4322",
		NoClientAuth: true,
		Banner:       "welcome to spider!\r\n",
		PrivateFile:  "ssh/spider",
		PasswordFunc: func(user, password string) bool {
			return "admin" == user && "admin" == password
		},
		Config: &grumble.Config{
			Name:        "spider",
			Prompt:      "spider > ",
			Description: "spider is a tool to list and diagnose Go processes",
		},
		Commands: func(commands *grumble.Commands) *grumble.Commands {
			return timeCommand(commands)
		},
	}
	if err := spider.RunSSH(cfg); err != nil {
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
