package main

import (
	"fmt"
	"spider"
)

func main() {
	RunSSH()
}

func RunGrumble() {
	/*spider.Run(&spider.Config{
		Port: 4322,
		Commands: func(commands *grumble.Commands) {
			commands.Add(&grumble.Command{
				Name: "time",
				Help: "displays time now",
				Run: func(c *grumble.Context) error {
					c.App.Println(time.Now())
					return nil
				},
			})
		},
	})*/
}

func RunSSH() {
	cfg := &spider.Config{
		Name:        "spider",
		Description: "spider is a tool to list and diagnose Go processes",
		SSHConfig: spider.SSHConfig{
			TCPConfig:    spider.TCPConfig{Address: ":4322"},
			NoClientAuth: false,
			Banner:       "welcome to spider!\r\n",
			PrivateFile:  "ssh/spider",
			PasswordFunc: func(user, password string) bool {
				return "admin" == user && "admin" == password
			},
		},
	}
	if err := spider.Start(cfg); err != nil {
		fmt.Printf("Start() error = %v", err)
	}
}
