//go:build darwin || dragonfly || freebsd || (linux && !appengine) || netbsd || openbsd || solaris

package spider

import (
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

func registerWindowChange(fd int, f func(width, height int)) {
	windowChanged := f
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				break
			}
			if width, height, err := term.GetSize(fd); err == nil {
				windowChanged(width, height)
			}
		}
	}()
}
