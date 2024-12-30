package terminal

import (
	"io"
	"os"
	"sync"
)

type Config struct {
	Prompt      string // eg: root > <command>
	Stdin       io.ReadCloser
	Stdout      io.Writer
	Stderr      io.Writer
	IsTerminal  bool // use interactive or not
	FuncMakeRaw func() error
	FuncExitRaw func() error
	once        sync.Once
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
		rm := new(RawMode)
		if config.FuncMakeRaw == nil {
			config.FuncMakeRaw = rm.Enter
		}
		if config.FuncExitRaw == nil {
			config.FuncExitRaw = rm.Exit
		}
	})
}
