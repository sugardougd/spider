package main

import (
	"fmt"
	"spider"
	"spider/grumble"
	"strconv"
	"strings"
	"time"
)

func main() {
	//padding(9999, 4)
	//padding(0, 4)
	//padding(1, 4)
	//padding(10, 4)
	//padding(19, 4)
	//padding(199, 4)
	RunConsole()
	//RunTCP()
	//RunSSH()
}

func padding(number, width int) {
	// 将数字转换为字符串
	numStr := strconv.Itoa(number)

	// 计算左右填充的空格数
	padding := (width - len(numStr)) / 2
	leftPadding := strings.Repeat(" ", padding)
	rightPadding := strings.Repeat(" ", width-len(numStr)-padding)
	// 如果宽度是奇数，确保字符串真正居中
	if width%2 == 1 {
		if len(numStr)%2 == 0 {
			rightPadding += " "
		} else {
			leftPadding += " "
		}
	}

	// 最终居中对齐的字符串
	centered := leftPadding + numStr + rightPadding
	fmt.Println("Centered Number:", centered)
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
