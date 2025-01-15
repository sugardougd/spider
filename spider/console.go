package spider

import (
	"fmt"
	"golang.org/x/term"
	"os"
)

func (s *Spider) RunConsole() error {
	terminal := term.NewTerminal(&ReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, s.Config.Prompt)

	fd := int(os.Stdout.Fd())
	raw, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("failed to get terminal state from MakeRaw: %s", err)
		return err
	}
	defer term.Restore(fd, raw)
	return s.runWithTerminal(terminal)
}
