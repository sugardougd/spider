package spider

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
)

type Spider struct {
	Config     *Config
	Commands   Commands
	stopSignal chan struct{}
	terminal   *term.Terminal
	mux        sync.Mutex
}

func New(config *Config, commands *Commands) *Spider {
	s := Spider{
		Config: config,
	}
	// Add general builtin Commands. help exit
	s.AddCommand(spiderCommand())
	s.AddCommand(helpCommand())
	s.AddCommand(exitCommand())
	if commands != nil {
		s.AddCommands(commands)
	}
	return &s
}

func (s *Spider) AddCommands(commands *Commands) error {
	for _, cmd := range commands.list {
		if err := s.AddCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (s *Spider) AddCommand(cmd *Command) error {
	return s.Commands.Add(cmd)
}

func (s *Spider) RunWithTerminal(ctx context.Context, terminal *term.Terminal) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.stopSignal = make(chan struct{})
	defer func() {
		s.stopSignal = nil
	}()
	s.terminal = terminal
	s.terminal.AutoCompleteCallback = s.autoComplete
	if len(s.Config.Welcome) > 0 {
		s.Println(s.Config.Welcome)
	}
	return s.run(ctx)
}

func (s *Spider) run(ctx context.Context) error {
	readline := make(chan string)

	// 另起goroutine来执行 terminal.ReadLine
	go func() {
		for {
			cmd, err := s.terminal.ReadLine()
			if err != nil {
				close(readline)
				s.stop()
				return
			}
			readline <- cmd
		}
	}()
	for {
		select {
		case <-ctx.Done():
			s.stop()
			goto exit
		case <-s.stopSignal:
			goto exit
		case cmd := <-readline:
			if err := s.RunCommand(ctx, cmd); err != nil {
				s.Printf("%v\n", err)
			}
		}
	}
exit:
	return nil
}

func (s *Spider) RunCommand(ctx context.Context, cmd string) error {
	command, flagValues, argValues, err := s.parse(cmd, true)
	if err != nil {
		return err
	}
	if command == nil {
		return nil
	}
	if command.Run == nil {
		return fmt.Errorf("illagel command Run '%s'", command.Name)
	}
	context := &Context{
		Spider:     s,
		Command:    command,
		CommandStr: cmd,
		FlagValues: flagValues,
		ArgValues:  argValues,
		Ctx:        ctx,
	}
	err = command.Run(context)
	if s.Config.ExecutedHook != nil {
		s.Config.ExecutedHook(context, err)
	}
	return err
}

func (s *Spider) parse(cmd string, required bool) (*Command, FlagValues, ArgValues, error) {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return nil, nil, nil, nil
	}
	flagValues := make(FlagValues)
	command, args, err := s.Commands.parse(args, required, flagValues)
	if err != nil {
		return nil, nil, nil, err
	}
	if command == nil {
		return nil, nil, nil, fmt.Errorf("illegal command '%s'", args[0])
	}
	helpFlag, err := flagValues.Bool(command, "help")
	if (err == nil && helpFlag) || command.Run == nil {
		s.PrintCommandHelp(command)
		return nil, nil, nil, nil
	}
	// parse args
	argValues := make(ArgValues)
	args, err = command.args.parse(args, argValues)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check, if values from the argument string are not consumed (and therefore invalid).
	if len(args) > 0 {
		return nil, nil, nil, fmt.Errorf("invalid usage of command '%s' (unconsumed input '%s'), try 'help'", command.Name, strings.Join(args, " "))
	}
	return command, flagValues, argValues, nil
}

func (s *Spider) stop() error {
	if s.stopSignal == nil {
		return fmt.Errorf("has stopped")
	}
	close(s.stopSignal)
	return nil
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

func (s *Spider) onWindowChanged(width, height int) {
	s.SetSize(width, height)
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

	s.Println()
	s.Println("Usage:")
	s.Println(BLANK2 + command.Usage)

	if len(command.flags.list) > 0 {
		s.Println()
		s.Println("Flags:")
		for _, flag := range command.flags.list {
			s.Println(BLANK2, "-"+flag.Short, "--"+flag.Long, "\t"+flag.Help)
		}
	}

	if len(command.args.list) > 0 {
		s.Println()
		s.Println("Args:")
		for _, arg := range command.args.list {
			s.Println(BLANK2, arg.Name, arg.Help)
		}
	}

	if len(command.Children.list) > 0 {
		s.Println()
		s.Println("Sub Commands:")
		for _, sub := range command.Children.list {
			s.PrintCommandList(TAB, sub)
		}
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
