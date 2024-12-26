package spider

import (
	"fmt"
	"io"
	"os"
	"spider/terminal"
	"strings"
)

type Spider struct {
	config     *Config
	commands   Commands
	stopSignal chan struct{}
	terminal   *terminal.Terminal
}

func New(config *Config, commands *Commands) *Spider {
	s := Spider{
		config:     config,
		stopSignal: make(chan struct{}),
	}
	// Add general builtin commands. help exit
	s.AddCommand(&Command{
		Name:        "help",
		Description: "use 'help [command]' for command help",
		Usage:       "use 'help [command]' for command usage",
		builtin:     true,
		Run: func(c *Context) error {
			c.Spider.Printf("exec '%s'\n", c.Command.Name)
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
	defer s.terminal.Stop()
	//TODO
	return nil
}

func (s *Spider) RunWithTerminal(terminal *terminal.Terminal) error {
	s.terminal = terminal
	return s.Run()
}

func (s *Spider) RunCommand(cmd string) error {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return fmt.Errorf("illagel command '%s'", cmd)
	}
	flagValues := make(FlagValues)
	command, args, err := s.commands.Parse(args, flagValues)
	if err != nil {
		return err
	}
	// parse args
	argValues := make(ArgValues)
	args, err = command.args.parse(args, argValues)
	if err != nil {
		return err
	}
	context := &Context{
		Spider:     s,
		Command:    command,
		CommandStr: cmd,
		flagValues: flagValues,
		argValues:  argValues,
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
	return s.commands.Add(cmd)
}

func (s *Spider) Stop() error {
	// TODO
	close(s.stopSignal)
	return nil
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
