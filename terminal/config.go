package terminal

import "io"

type Config struct {
	Prompt     string // eg: root > <command>
	Stdin      io.ReadCloser
	Stdout     io.Writer
	Stderr     io.Writer
	IsTerminal bool // use interactive or not
}
