package spider

import (
	"context"
	"golang.org/x/term"
	"os"
)

func RunConsole(config *Config, commands *Commands, ctx context.Context) error {
	fd := int(os.Stdout.Fd())
	raw, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, raw)

	s := New(config, commands)
	go s.RunWithTerminal(term.NewTerminal(&ReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, config.Prompt))

	select {
	case <-ctx.Done():
		s.Stop()
	case <-s.stopSignal:
		break
	}
	return nil
}
