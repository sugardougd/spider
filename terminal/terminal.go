package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Terminal struct {
	config     *Config
	line       *ReadLine
	History    *History
	stopSignal chan struct{}
}

// New Terminal
func New(config *Config) *Terminal {
	config.Init()
	ter := Terminal{
		config:     config,
		line:       NewReadLine(128),
		History:    NewHistory(5),
		stopSignal: make(chan struct{}),
	}
	go func() {
		ter.ioloop()
	}()
	ter.WritePrompt()
	return &ter
}

func NewConsole(prompt string) *Terminal {
	config := Config{
		Prompt: prompt,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return New(&config)
}

// Readline returns a line of input from the terminal.
func (t *Terminal) Readline() <-chan string {
	return t.line.Line()
}

func (t *Terminal) Write(buf []byte) (int, error) {
	return t.config.Stdout.Write(buf)
}

func (t *Terminal) Stdin() io.ReadCloser {
	return t.config.Stdin
}

func (t *Terminal) Stdout() io.Writer {
	return t.config.Stdout
}

func (t *Terminal) Stderr() io.Writer {
	return t.config.Stderr
}

func (t *Terminal) WritePrompt() (int, error) {
	buf := []byte(t.config.Prompt)
	return t.Write(buf)
}

func (t *Terminal) WriteLine(buf []byte) (int, error) {
	line := make([]byte, len(buf)+len(crlf))
	copy(line, buf)
	copy(line[len(buf):], crlf)
	return t.Write(line)
}

// WriteWithCRLF writes buf to w but replaces all occurrences of \n with \r\n.
func (t *Terminal) WriteWithCRLF(buf []byte) (n int, err error) {
	for len(buf) > 0 {
		i := bytes.IndexByte(buf, '\n')
		todo := len(buf)
		if i >= 0 {
			todo = i
		}

		var nn int
		nn, err = t.Write(buf[:todo])
		n += nn
		if err != nil {
			return n, err
		}
		buf = buf[todo:]

		if i >= 0 {
			if _, err = t.Write(crlf); err != nil {
				return n, err
			}
			n++
			buf = buf[1:]
		}
	}
	return n, nil
}

func (t *Terminal) Stop() error {
	fmt.Println("Stop Terminal")
	t.config.Stdin.Close()
	close(t.stopSignal)
	close(t.line.stopSignal)
	return nil
}
func (t *Terminal) Done() <-chan struct{} {
	return t.stopSignal
}

func (t *Terminal) isRunning() bool {
	select {
	case <-t.stopSignal:
		return false
	default:
		return true
	}
}

func (t *Terminal) ioloop() {
	defer t.line.stop()
	buf := bufio.NewReader(t.config.Stdin)
	for t.isRunning() {
		t.config.FuncMakeRaw()
		b, err := buf.ReadByte()
		if err != nil {
			break
		}
		t.line.enterKey(b)
	}
	t.config.FuncExitRaw()
}
