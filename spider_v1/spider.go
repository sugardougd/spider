package spider_v1

import (
	"os"
	"spider/grumble"
)

type CommandsSet func(commands *grumble.Commands) *grumble.Commands

func newGrumble(cfg *grumble.Config, commands CommandsSet) *grumble.App {
	app := grumble.New(cfg)

	// add to App
	for _, c := range spiderCommands().All() {
		app.AddCommand(c)
	}

	// add to App
	if commands != nil {
		for _, c := range commands(app.Commands()).All() {
			app.AddCommand(c)
		}
	}
	return app
}

func spiderCommands() *grumble.Commands {
	commands := grumble.Commands{}
	sc := &grumble.Command{
		Name: "spider_v1",
		Help: "a tool to list and diagnose Go processes",
		Run: func(c *grumble.Context) error {
			c.App.PrintCommandHelp(c.Command, true)
			return nil
		},
	}
	// stats
	stats := &grumble.Command{
		Name: "stats",
		Help: "displays process info",
		Run: func(c *grumble.Context) error {
			c.App.Println("PPID:", os.Getppid())
			c.App.Println("PID:", os.Getpid())
			return nil
		},
	}
	sc.AddCommand(stats)
	commands.Add(sc)
	return &commands
}
