//go:build darwin

// Package terminal provides support functions for dealing with terminals, as
// commonly found on UNIX systems.
//
// Putting a terminal into raw mode is the most common requirement:
//
//	oldState, err := terminal.MakeRaw(0)
//	if err != nil {
//	        panic(err)
//	}
//	defer terminal.Restore(0, oldState)
package terminal

import (
	"syscall"
	"unsafe"
)

type Termios syscall.Termios

// State contains the state of a terminal.
type State struct {
	termios Termios
}

func getTermios(fd int) (*Termios, error) {
	termios := new(Termios)
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCGETA, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return termios, nil
}

func setTermios(fd int, termios *Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSETA, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) {
	var oldState State

	if termios, err := getTermios(fd); err != nil {
		return nil, err
	} else {
		oldState.termios = *termios
	}

	newState := oldState.termios
	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	// newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8

	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	return &oldState, setTermios(fd, &newState)
}

func Restore(fd int, state *State) error {
	err := restore(fd, state)
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

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func restore(fd int, state *State) error {
	return setTermios(fd, &state.termios)
}
