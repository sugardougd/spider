package terminal

import (
	"errors"
)

const (
	KeyLineStart = 1
	KeyBackward  = 2
	KeyInterrupt = 3
	KeyDelete    = 4
	KeyLineEnd   = 5
	KeyForward   = 6
	KeyBell      = 7
	KeyCtrlH     = 8
	KeyTab       = 9
	KeyCtrlJ     = 10
	KeyKill      = 11
	KeyCtrlL     = 12
	KeyEnter     = 13
	KeyNext      = 14
	KeyPrev      = 16
	KeyBckSearch = 18
	KeyFwdSearch = 19
	KeyTranspose = 20
	KeyCtrlU     = 21
	KeyCtrlW     = 23
	KeyCtrlY     = 25
	KeyCtrlZ     = 26
	KeyEsc       = 27
	KeyEscapeEx  = 91
	KeyBackspace = 127
	keyLeft      = -1
	keyRight     = -2
)

var (
	crlf          = []byte{'\r', '\n'}
	pasteStart    = []byte{KeyEsc, '[', '2', '0', '0', '~'}
	pasteEnd      = []byte{KeyEsc, '[', '2', '0', '1', '~'}
	ErrNotRunning = errors.New("NOT_RUNNING")
)

type RawMode struct {
	state *State
}

func (r *RawMode) Enter() (err error) {
	r.state, err = MakeRaw(GetStdin())
	return err
}

func (r *RawMode) Exit() error {
	if r.state == nil {
		return nil
	}
	return Restore(GetStdin(), r.state)
}

func Restore(fd int, state *State) error {
	err := restoreTerm(fd, state)
	if err != nil {
		// errno 0 means everything is ok :)
		if err.Error() == "errno 0" {
			return nil
		} else {
			return err
		}
	}
	return nil
}
