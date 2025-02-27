package spider

import (
	"context"
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

func RunConsole(config *Config, commands *Commands, ctx context.Context) error {
	fd := int(os.Stdout.Fd())
	raw, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, raw)

	s := New(config)
	s.AddCommands(commands)
	terminal := term.NewTerminal(&ReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, config.Prompt)

	// register width change
	onWindowChanged(fd, func(width, height int) {
		s.SetSize(width, height)
	})

	if err = s.RunWithTerminal(terminal, ctx); err != nil {
		return err
	}
	return nil
}

func onWindowChanged(fd int, windowChanged func(width, height int)) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				break
			}
			if width, height, err := term.GetSize(fd); err == nil {
				windowChanged(width, height)
			}
		}
	}()
}
