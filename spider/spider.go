package spider

import (
	"fmt"
	"io"
	"os"
	"spider/terminal"
	"strings"
)

type Spider struct {
	Config     *Config
	Commands   Commands
	stopSignal chan struct{}
	terminal   *terminal.Terminal
}

func New(config *Config, commands *Commands) *Spider {
	s := Spider{
		Config:     config,
		stopSignal: make(chan struct{}),
	}
	// Add general builtin Commands. help exit
	s.AddCommand(&Command{
		Name:        "help",
		Description: "use 'help [command]' for command help",
		Usage:       "use 'help [command]' for command usage",
		builtin:     true,
		Args: func(args *Args) {
			args.StringList(&Arg{
				Name: "command",
				Help: "the name of the command",
			})
		},
		Run: func(c *Context) error {
			args, err := c.ArgValues.StringList("command")
			if err != nil {
				return err
			}
			if len(args) == 0 {
				c.Spider.PrintHelp()
				return nil
			}
			flagValues := make(FlagValues)
			command, _, err := c.Spider.Commands.parse(args, false, flagValues)
			if err != nil {
				return err
			} else if command == nil {
				return fmt.Errorf("command not found")
			}
			c.Spider.PrintCommandHelp(command)
			return nil
		},
	})
	s.AddCommand(&Command{
		Name:        "exit",
		Description: "exit the spider",
		builtin:     true,
		Run: func(c *Context) error {
			return c.Stop()
		},
	})
	for _, cmd := range commands.list {
		s.AddCommand(cmd)
	}
	return &s
}

func (s *Spider) Run() error {
	for s.IsRunning() {
		select {
		case cmd := <-s.terminal.Readline():
			if err := s.RunCommand(cmd); err != nil {
				fmt.Fprintf(s.Stderr(), "%v\n", err)
			}
			s.PrintPrompt()
		case <-s.terminal.Done():
			s.Stop()
			break
		case <-s.stopSignal:
			break
		}
	}
	return nil
}

func (s *Spider) RunWithTerminal(terminal *terminal.Terminal) error {
	s.terminal = terminal
	return s.Run()
}

func (s *Spider) RunCommand(cmd string) error {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		//return fmt.Errorf("illegal command '%s'", cmd)
		return nil
	}
	flagValues := make(FlagValues)
	command, args, err := s.Commands.parse(args, true, flagValues)
	if err != nil {
		return err
	}
	// parse args
	argValues := make(ArgValues)
	args, err = command.args.parse(args, argValues)
	if err != nil {
		return err
	}

	// Check, if values from the argument string are not consumed (and therefore invalid).
	if len(args) > 0 {
		return fmt.Errorf("invalid usage of command '%s' (unconsumed input '%s'), try 'help'", command.Name, strings.Join(args, " "))
	}

	context := &Context{
		Spider:     s,
		Command:    command,
		CommandStr: cmd,
		FlagValues: flagValues,
		ArgValues:  argValues,
	}
	if command.Run == nil {
		return fmt.Errorf("illagel command Run '%s'", command.Name)
	}
	if err = command.Run(context); err != nil {
		return err
	}
	return nil
}

func (s *Spider) AddCommand(cmd *Command) error {
	return s.Commands.Add(cmd)
}

func (s *Spider) PrintPrompt() (int, error) {
	return fmt.Fprintf(s.Stdout(), s.Config.Prompt)
}

func (s *Spider) Stop() error {
	close(s.stopSignal)
	if s.terminal != nil {
		s.terminal.Stop()
	}
	return nil
}

func (s *Spider) IsRunning() bool {
	select {
	case <-s.stopSignal:
		return false
	default:
		return true
	}
}

func (s *Spider) Stdin() io.ReadCloser {
	if s.terminal != nil {
		return s.terminal.Stdin()
	}
	return os.Stdin
}

func (s *Spider) Stdout() io.Writer {
	if s.terminal != nil {
		return s.terminal.Stdout()
	}
	return os.Stdout
}

func (s *Spider) Stderr() io.Writer {
	if s.terminal != nil {
		return s.terminal.Stderr()
	}
	return os.Stderr
}

func (s *Spider) Write(p []byte) (n int, err error) {
	return s.Stdout().Write(p)
}

func (s *Spider) Print(args ...interface{}) (int, error) {
	return fmt.Fprint(s, args...)
}

func (s *Spider) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(s, format, args...)
}

func (s *Spider) Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(s, args...)
}

func (s *Spider) PrintHelp() {
	s.Println()
	s.Println(s.Config.Description)
	s.Println()
	s.Println("Commands:")
	for _, command := range s.Commands.list {
		s.Println("\t" + command.Name + "\t" + command.Description)
	}
}

func (s *Spider) PrintCommandHelp(command *Command) {
	s.Println()
	s.Println(command.Description)
	s.Println("Usage:")
	s.Println("\t" + command.Usage)
	s.Println("Flags:")
	for _, flag := range command.flags.list {
		s.Println("\t", "-"+flag.Short, "--"+flag.Long, "\t"+flag.Help)
	}
	s.Println("Args:")
	for _, arg := range command.args.list {
		s.Println("\t", arg.Name, arg.Help)
	}
	s.Println("Sub Commands:")
	for _, sub := range command.Children.list {
		s.Println("\t" + sub.Name + "\t" + sub.Description)
	}
}
