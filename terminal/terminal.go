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
	buf        *RuneBuffer
	line       chan string
	History    *History
	stopSignal chan struct{}
}

// New Terminal
func New(config *Config) *Terminal {
	config.Init()
	ter := Terminal{
		config:     config,
		buf:        NewRuneBuffer(config),
		line:       make(chan string),
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
	return t.line
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

func (t *Terminal) Write(buf []byte) (int, error) {
	return t.config.Stdout.Write(buf)
}

func (t *Terminal) WriteLine(buf []byte) (int, error) {
	line := make([]byte, len(buf)+len(crlf))
	copy(line, buf)
	copy(line[len(buf):], crlf)
	return t.Write(line)
}

func (t *Terminal) WritePrompt() (int, error) {
	buf := []byte(t.config.Prompt)
	return t.Write(buf)
}

func (t *Terminal) WriteWithPrompt(buf []byte) (int, error) {
	prompt := []byte(t.config.Prompt)
	data := make([]byte, len(prompt)+len(buf))
	copy(data, prompt)
	copy(data[len(prompt):], buf)
	return t.Write(data)
}

func (t *Terminal) RefreshWithPrompt(buf []byte) (int, error) {
	t.Clean()
	return t.WriteWithPrompt(buf)
}

func (t *Terminal) Clean() {

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
	buf := bufio.NewReader(t.config.Stdin)
	var isEscape, isEscapeEx bool
	for t.isRunning() {
		t.config.FuncMakeRaw()
		r, _, err := buf.ReadRune()
		if err != nil {
			break
		}
		//fmt.Println("[DEBUG]type:", r)
		if isEscape {
			isEscape = false
			if r == CharEscapeEx { // 91 -> [
				isEscapeEx = true
			}
			continue
		}
		if isEscapeEx {
			isEscapeEx = false
			switch r {
			case 'A':
				r = CharPrev
			case 'B':
				r = CharNext
			case 'C':
				r = CharForward
			case 'D':
				r = CharBackward
			case 'H':
				r = CharLineStart
			case 'F':
				r = CharLineEnd
			default:
				fmt.Println("un-support EscapeEx:", r)
			}
		}
		switch r {
		case CharBackward:
			t.typeBackward()
		case CharForward:
			t.typeForward()
		case CharPrev:
			t.typePrev(r)
		case CharNext:
			t.typeNext(r)
		case CharTab:
			t.typeTab(r)
		case CharDelete, CharBackspace:
			t.typeDelete()
		case CharCtrlJ, CharEnter:
			t.typeEnter()
		case CharEsc: // 27 -> ESC
			isEscape = true
		default:
			t.buf.WriteRune(r)
		}
	}
	t.config.FuncExitRaw()
}

func (t *Terminal) typeTab(r rune) {
	fmt.Println("type tab: ", r)
}

func (t *Terminal) typeForward() {
	t.buf.MoveForward()
}

func (t *Terminal) typeBackward() {
	t.buf.MoveBackward()
}

func (t *Terminal) typePrev(r rune) {
	fmt.Println("type prev: ", r)
	//if pre, ok := t.History.Prev(); ok {
	//	fmt.Println("pre:", pre)
	//t.buf = []rune(pre)
	//t.pos = len(t.buf)
	//}
}

func (t *Terminal) typeNext(r rune) {
	fmt.Println("type next: ", r)
	//if next, ok := t.History.Next(); ok {
	//t.buf = []rune(next)
	//t.pos = len(t.buf)
	//}
}

func (t *Terminal) typeLeft(r rune) {
	fmt.Println("type left: ", r)
}

func (t *Terminal) typeRight(r rune) {
	fmt.Println("type right: ", r)
}

func (t *Terminal) typeDelete() {
	t.buf.Backspace()
}

func (t *Terminal) typeEnter() {
	t.buf.MoveToLineEnd()
	line := string(t.buf.Reset())
	t.History.Push(line)
	select {
	case t.line <- line:
	case <-t.stopSignal:
	}
}
