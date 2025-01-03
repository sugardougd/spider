package terminal

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

type RuneBuffer struct {
	buf         []rune
	pos         int
	prompt      []rune
	writer      io.Writer
	width       int
	interactive bool
	sync.Mutex
}

func NewRuneBuffer(config *Config) *RuneBuffer {
	rb := RuneBuffer{
		prompt:      []rune(config.Prompt),
		writer:      config.Stdout,
		width:       config.FuncGetWidth(),
		interactive: config.FuncIsTerminal(),
	}
	return &rb
}

func (r *RuneBuffer) Reset() []rune {
	ret := runes.Copy(r.buf)
	r.buf = r.buf[:0]
	r.pos = 0
	return ret
}

func (r *RuneBuffer) WriteRune(s rune) {
	r.WriteRunes([]rune{s})
}

func (r *RuneBuffer) WriteRunes(s []rune) {
	r.Refresh(func(r *RuneBuffer) {
		tail := append(s, r.buf[r.pos:]...)
		r.buf = append(r.buf[:r.pos], tail...)
		r.pos += len(s)
	})
}

func (r *RuneBuffer) MoveBackward() {
	r.Refresh(func(r *RuneBuffer) {
		if r.pos == 0 {
			return
		}
		r.pos--
	})
}

func (r *RuneBuffer) MoveForward() {
	r.Refresh(func(r *RuneBuffer) {
		if r.pos == len(r.buf) {
			return
		}
		r.pos++
	})
}

func (r *RuneBuffer) MoveToLineEnd() {
	r.Refresh(func(r *RuneBuffer) {
		r.pos = len(r.buf)
	})
}

func (r *RuneBuffer) Backspace() {
	r.Refresh(func(r *RuneBuffer) {
		if r.pos == 0 {
			return
		}
		r.pos--
		r.buf = append(r.buf[:r.pos], r.buf[r.pos+1:]...)
	})
}

func (r *RuneBuffer) Refresh(f func(buf *RuneBuffer)) {
	r.Lock()
	defer r.Unlock()
	if !r.interactive {
		if f != nil {
			f(r)
		}
		return
	}
	r.clean()
	if f != nil {
		f(r)
	}
	r.print()
}

func (r *RuneBuffer) Clean() {
	r.Lock()
	defer r.Unlock()
	r.clean()
}

func (r *RuneBuffer) print() {
	r.writer.Write(r.output())
}

func (r *RuneBuffer) output() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(string(r.prompt))
	for _, e := range r.buf {
		buf.WriteRune(e)
	}

	// cursor position
	if len(r.buf) > r.pos {
		buf.Write(r.getBackspaceSequence())
	}
	return buf.Bytes()
}

func (r *RuneBuffer) getBackspaceSequence() []byte {
	var buf []byte
	for i := len(r.buf); i > r.pos; i-- {
		// move input to the left of one
		buf = append(buf, '\b')
	}
	return buf
}

func (r *RuneBuffer) promptLen() int {
	return runes.WidthAll(r.prompt)
}

func (r *RuneBuffer) clean() {
	buf := bufio.NewWriter(r.writer)
	buf.Write([]byte("\033[J")) // just like ^k :)
	buf.WriteString("\033[2K")
	buf.WriteString("\r")
	buf.Flush()
}
