package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Terminal struct {
	config     *Config
	line       *Line
	stopSignal chan struct{}
}

type Line struct {
	buf        []byte
	pos        int
	line       chan string
	stopSignal chan struct{}
}

// New Terminal
func New(config *Config) *Terminal {
	ter := Terminal{
		config: config,
		line: &Line{
			buf:        make([]byte, 0, 128),
			pos:        0,
			line:       make(chan string),
			stopSignal: make(chan struct{}),
		},
		stopSignal: make(chan struct{}),
	}
	go func() {
		ter.ioloop()
	}()
	ter.WritePrompt()
	return &ter
}

// Readline returns a line of input from the terminal.
func (t *Terminal) Readline() (string, error) {
	if t.IsRunning() {
		return t.line.Line()
	}
	return "", ErrNotRunning
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

func (t *Terminal) IsRunning() bool {
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
	for t.IsRunning() {
		b, err := buf.ReadByte()
		if err != nil {
			break
		}
		t.line.enterKey(b)
	}
}

func (l *Line) enterKey(r byte) {
	switch r {
	case KeyTab:
		l.typeTab(r)
	case KeyDelete:
		l.typeDelete(r)
	case KeyCtrlJ, KeyEnter:
		l.typeEnter(r)
	default:
		l.appendRune(r)
	}
}

func (l *Line) typeTab(b byte) {
	fmt.Println("type tab: ", b)
}

func (l *Line) typeLeft(b byte) {
	fmt.Println("type left: ", b)
}

func (l *Line) typeRight(b byte) {
	fmt.Println("type right: ", b)
}

func (l *Line) typeDelete(b byte) {
	if 0 < l.pos && l.pos <= len(l.buf) {
		l.pos--
		buf := make([]byte, len(l.buf)-1)
		copy(buf, l.buf[:l.pos])
		copy(buf[l.pos:], l.buf[l.pos+1:])
		l.buf = buf
	}
}

func (l *Line) typeEnter(b byte) {
	line := string(l.buf)
	l.buf = l.buf[:0]
	select {
	case l.line <- line:
	case <-l.stopSignal:
	}
}

func (l *Line) appendRune(b byte) {
	//fmt.Println("type: ", b)
	if l.pos < 0 || l.pos == len(l.buf) || l.pos > len(l.buf) {
		l.buf = append(l.buf, b)
		l.pos = len(l.buf)
	} else {
		buf := make([]byte, len(l.buf)+1)
		copy(buf, l.buf[:l.pos])
		buf[l.pos] = b
		copy(buf[l.pos+1:], l.buf[l.pos:])
		l.buf = buf
		l.pos++
	}

}

func (l *Line) Line() (string, error) {
	select {
	case <-l.stopSignal:
		return "", ErrNotRunning
	case line := <-l.line:
		return line, nil
	}
}
func (l *Line) stop() error {
	fmt.Println("stop line")
	close(l.line)
	l.buf = l.buf[0:]
	return nil
}
