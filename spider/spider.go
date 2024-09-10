package spider

import (
	"os"
	"spider/grumble"
)

type CommandsSet func(*grumble.Commands)

func NewGrumble(cfg *grumble.Config, commands CommandsSet) *grumble.App {
	app := grumble.New(cfg)

	// stats
	app.AddCommand(&grumble.Command{
		Name: "stats",
		Help: "displays process info",
		Run: func(c *grumble.Context) error {
			c.App.Println("PID:", os.Getpid())
			c.App.Println("PPID:", os.Getppid())
			return nil
		},
	})

	//
	if commands != nil {
		commands(app.Commands())
	}
	return app
}
