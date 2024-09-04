package spider

type Config struct {
	SSHConfig
	Name        string
	Description string
}

/*type Config struct {
	Host           string                       // agent host
	Port           int                          // agent port
	Commands       func(*grumble.Commands)      // Commands
	Authentication func(user, pwd string) error // Authentication
}

// authentication the questions to ask authentication
var authentication = []*survey.Question{
	{
		Name:     "user",
		Prompt:   &survey.Input{Message: "Please type your username"},
		Validate: survey.Required,
	},
	{
		Name: "password",
		Prompt: &survey.PasswordFunc{
			Message: "Please type your password",
		},
		Validate: survey.Required,
	},
}

func Run(cfg *Config) {
	if cfg.Port > 0 {
		fmt.Printf("spider port %d", cfg.Port)
	} else {
		fmt.Printf("disable spider")
		return
	}

	app := New(cfg)
	if err := app.Run(); err != nil {
		app.Printf("run spider fail. %v", err)
		app.Close()
	}

	handleFunc := func(rl *readline.Instance) {
		app := New(cfg)
		if err := app.RunWithReadline(rl); err != nil {
			app.Printf("run spider fail. %v", err)
		}
	}
	rlc := &readline.Config{}
	readline.ListenRemote("tcp", cfg.Host+":"+strconv.Itoa(cfg.Port), rlc, handleFunc)
}

func New(cfg *Config) *grumble.App {
	app := grumble.New(&grumble.Config{
		Name:        "spider",
		Description: "spider is a tool to list and diagnose Go processes.",
		Flags: func(f *grumble.Flags) {
			f.Bool("c", "chinese", false, "enable chinese languages")
		},
	})

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
	if cfg.Commands != nil {
		cfg.Commands(app.Commands())
	}
	return app
}*/
