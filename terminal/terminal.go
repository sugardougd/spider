package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Terminal struct {
	prompt     string
	reader     io.Reader
	writer     io.Writer
	buf        *RuneBuffer
	line       chan string
	history    *History
	stopSignal chan struct{}
}

// New Terminal
func New(prompt string, reader io.Reader, writer io.Writer, interactive bool) *Terminal {
	ter := Terminal{
		prompt:     prompt,
		reader:     reader,
		writer:     writer,
		buf:        NewRuneBuffer(prompt, writer, interactive),
		line:       make(chan string),
		history:    NewHistory(5),
		stopSignal: make(chan struct{}),
	}
	go func() {
		ter.ioloop()
	}()
	return &ter
}

func NewConsole(prompt string) *Terminal {
	return New(prompt, os.Stdin, os.Stdout, true)
}

// Readline returns a line of input from the terminal.
func (t *Terminal) Readline() <-chan string {
	return t.line
}

func (t *Terminal) Write(buf []byte) (n int, err error) {
	n, err = t.writer.Write(buf)
	t.WritePrompt()
	return
}

func (t *Terminal) WritePrompt() (int, error) {
	buf := []byte(t.prompt)
	return t.writer.Write(buf)
}

func (t *Terminal) Stop() error {
	fmt.Println("Stop Terminal")
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
	t.WritePrompt()
	buf := bufio.NewReader(t.reader)
	var isEscape, isEscapeEx bool
	for t.isRunning() {
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
			t.typeEnter(r)
		case CharEsc: // 27 -> ESC
			isEscape = true
		default:
			t.buf.WriteRune(r)
		}
	}
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
	//if pre, ok := t.history.Prev(); ok {
	//	fmt.Println("pre:", pre)
	//t.buf = []rune(pre)
	//t.pos = len(t.buf)
	//}
}

func (t *Terminal) typeNext(r rune) {
	fmt.Println("type next: ", r)
	//if next, ok := t.history.Next(); ok {
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

func (t *Terminal) typeEnter(r rune) {
	t.buf.MoveToLineEnd()
	t.buf.WriteRune(r)
	line := string(t.buf.Reset())
	line = line[:len(line)-1] // trim \n
	t.history.Push(line)
	select {
	case t.line <- line:
	case <-t.stopSignal:
	}
}
