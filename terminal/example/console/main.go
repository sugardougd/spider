package main

import (
	"fmt"
	"os"
	"spider/terminal"
	"strings"
)

func main() {
	TestTerminal()
}

func TestTerminal() {
	fd := int(os.Stdin.Fd())
	raw, err := terminal.MakeRaw(fd)
	//raw, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	defer terminal.Restore(fd, raw)
	//defer term.Restore(fd, raw)
	ter := terminal.NewConsole("spider$ ")

	for {
		select {
		case line := <-ter.Readline():
			command := strings.TrimSpace(line)
			if _, err := fmt.Fprintln(ter, "type: "+command); err != nil {
				fmt.Printf("Write fail %v", err)
			}
		case <-ter.Done():
			break
		}
	}
}
