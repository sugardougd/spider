package spider

import (
	"context"
	"os"
	"syscall"

	"golang.org/x/term"
)

func RunConsole(ctx context.Context, config *Config, commands *Commands) error {
	fd := int(syscall.Stdin)
	raw, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, raw)

	config.Interactive = true
	s := New(config, commands)
	terminal := term.NewTerminal(&ReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, config.Prompt)

	// register width change
	registerWindowChange(fd, s.onWindowChanged)

	if err = s.RunWithTerminal(ctx, terminal); err != nil {
		return err
	}
	return nil
}
