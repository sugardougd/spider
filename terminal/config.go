package terminal

import (
	"io"
	"os"
	"sync"
)

type Config struct {
	Prompt         string // eg: root > <command>
	Stdin          io.ReadCloser
	Stdout         io.Writer
	Stderr         io.Writer
	FuncIsTerminal func() bool // use interactive or not
	FuncMakeRaw    func() error
	FuncExitRaw    func() error
	FuncGetWidth   func() int
	once           sync.Once
}

func (config *Config) Init() {
	config.once.Do(func() {
		if config.Stdin == nil {
			config.Stdin = os.Stdin
		}
		if config.Stdout == nil {
			config.Stdout = os.Stdout
		}
		if config.Stderr == nil {
			config.Stderr = os.Stderr
		}
		if config.FuncIsTerminal == nil {
			config.FuncIsTerminal = DefaultIsTerminal
		}
		rm := new(RawMode)
		if config.FuncMakeRaw == nil {
			config.FuncMakeRaw = rm.Enter
		}
		if config.FuncExitRaw == nil {
			config.FuncExitRaw = rm.Exit
		}
		if config.FuncGetWidth == nil {
			config.FuncGetWidth = GetScreenWidth
		}
	})
}
