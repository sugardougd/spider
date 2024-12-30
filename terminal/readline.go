package terminal

import "fmt"

type ReadLine struct {
	buf        []byte
	pos        int
	line       chan string
	stopSignal chan struct{}
}

func NewReadLine(size int) *ReadLine {
	l := ReadLine{
		buf:        make([]byte, 0, size),
		pos:        0,
		line:       make(chan string),
		stopSignal: make(chan struct{}),
	}
	return &l
}

func (l *ReadLine) enterKey(r byte) {
	fmt.Println("[DEBUG]type:", r)
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

func (l *ReadLine) typeTab(b byte) {
	fmt.Println("type tab: ", b)
}

func (l *ReadLine) typeLeft(b byte) {
	fmt.Println("type left: ", b)
}

func (l *ReadLine) typeRight(b byte) {
	fmt.Println("type right: ", b)
}

func (l *ReadLine) typeDelete(b byte) {
	if 0 < l.pos && l.pos <= len(l.buf) {
		l.pos--
		buf := make([]byte, len(l.buf)-1)
		copy(buf, l.buf[:l.pos])
		copy(buf[l.pos:], l.buf[l.pos+1:])
		l.buf = buf
	}
}

func (l *ReadLine) typeEnter(b byte) {
	line := string(l.buf)
	l.buf = l.buf[:0]
	select {
	case l.line <- line:
	case <-l.stopSignal:
	}
}

func (l *ReadLine) appendRune(b byte) {
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

func (l *ReadLine) Line() <-chan string {
	return l.line
}
func (l *ReadLine) stop() error {
	fmt.Println("stop line")
	close(l.line)
	l.buf = l.buf[0:]
	return nil
}
