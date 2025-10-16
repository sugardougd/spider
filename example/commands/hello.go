package commands

import (
	"github.com/sugardougd/spider"
	"time"
)

func HelloCommand() *spider.Command {
	hello := &spider.Command{
		Name:        "hello",
		Description: "hello tools",
		Usage:       "hello...",
		Run: func(c *spider.Context) error {
			c.Spider.Println("Hello")
			return nil
		},
	}
	// date
	hello.AddCommand(&spider.Command{
		Name:        "date",
		Description: "hello date tools",
		Usage:       "hello date",
		Flags:       noyaCommandFlags,
		Run: func(c *spider.Context) error {
			c.Spider.Println("Hello", time.Now().Format("2006-01-02"))
			return nil
		},
	})
	// time
	hello.AddCommand(&spider.Command{
		Name:        "time",
		Description: "hello time tools",
		Usage:       "hello time",
		Flags:       noyaCommandFlags,
		Run: func(c *spider.Context) error {
			c.Spider.Println("Hello", time.Now().Format("15:04:05"))
			return nil
		},
	})
	// date-time
	hello.AddCommand(&spider.Command{
		Name:        "date-time",
		Description: "hello date-time tools",
		Usage:       "hello date-time",
		Flags:       noyaCommandFlags,
		Run: func(c *spider.Context) error {
			c.Spider.Println("Hello", time.Now().Format("2006-01-02 15:04:05"))
			return nil
		},
	})
	return hello
}
