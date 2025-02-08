package spider

import (
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

type Spider struct {
	Config     *Config
	Commands   Commands
	stopSignal chan struct{}
	terminal   *term.Terminal
}

func New(config *Config, commands *Commands) *Spider {
	s := Spider{
		Config:     config,
		stopSignal: make(chan struct{}),
	}
	// Add general builtin Commands. help exit
	s.AddCommand(helpCommand())
	s.AddCommand(spiderCommand())
	s.AddCommand(exitCommand())

	for _, cmd := range commands.list {
		s.AddCommand(cmd)
	}
	return &s
}

func (s *Spider) RunWithTerminal(terminal *term.Terminal) error {
	s.terminal = terminal
	s.terminal.AutoCompleteCallback = s.autoComplete
	return s.run()
}

func (s *Spider) run() error {
	for s.IsRunning() {
		select {
		case <-s.stopSignal:
			break
		default:
			cmd, err := s.terminal.ReadLine()
			if err != nil {
				s.Stop()
				break
			}
			if err := s.RunCommand(cmd); err != nil {
				s.Printf("%v\n", err)
			}
		}
	}
	return nil
}

func (s *Spider) RunCommand(cmd string) error {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return nil
	}
	flagValues := make(FlagValues)
	command, args, err := s.Commands.parse(args, true, flagValues)
	if err != nil {
		return err
	}
	if command == nil {
		return fmt.Errorf("illegal command '%s'", args[0])
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
	return command.Run(context)
}

func (s *Spider) AddCommand(cmd *Command) error {
	return s.Commands.Add(cmd)
}

func (s *Spider) Stop() error {
	close(s.stopSignal)
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

func (s *Spider) Write(p []byte) (n int, err error) {
	if s.terminal != nil {
		return s.terminal.Write(p)
	}
	return os.Stdout.Write(p)
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

func (s *Spider) SetSize(width, height int) {
	if s.terminal != nil {
		s.terminal.SetSize(width, height)
	}
}

func (s *Spider) PrintHelp() {
	s.Println()
	s.Println(s.Config.Description)
	s.Println()
	s.Println("Commands:")
	w := tabwriter.NewWriter(s, 0, 2, 4, ' ', 0)
	for _, command := range s.Commands.list {
		fmt.Fprintln(w, BLANK2+command.Name+TAB+command.Description)
	}
	w.Flush()
}

func (s *Spider) PrintCommandHelp(command *Command) {
	s.Println()
	s.Println(command.Description)
	s.Println("Usage:")
	s.Println(BLANK2 + command.Usage)
	s.Println("Flags:")
	for _, flag := range command.flags.list {
		s.Println(BLANK2, "-"+flag.Short, "--"+flag.Long, "\t"+flag.Help)
	}
	s.Println("Args:")
	for _, arg := range command.args.list {
		s.Println(BLANK2, arg.Name, arg.Help)
	}
	s.Println("Sub Commands:")
	for _, sub := range command.Children.list {
		s.PrintCommandList(TAB, sub)
	}
}

func (s *Spider) PrintCommandList(padding string, command *Command) {
	s.Println(padding + command.Name + TAB + command.Description)
	for _, sub := range command.Children.list {
		s.PrintCommandList(padding, sub)
	}
}

type ReadWriter struct {
	Reader io.Reader
	Writer io.Writer
}

func (c *ReadWriter) Read(p []byte) (n int, err error) {
	return c.Reader.Read(p)
}

func (c *ReadWriter) Write(p []byte) (n int, err error) {
	return c.Writer.Write(p)
}
