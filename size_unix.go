// +build linux

package terminal_size

import (
	"os"
	"os/signal"
	"unsafe"

	"golang.org/x/sys/unix"
)

type winsize struct {
	rows uint16
	cols uint16
	x    uint16
	y    uint16
}

func getTerminalSize(fp *os.File) (width int, height int, err error) {
	ws := winsize{}

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		fp.Fd(),
		uintptr(unix.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))

	if errno != 0 {
		err = errno
	}

	if err == nil {
		width = int(ws.cols)
		height = int(ws.rows)
	}

	return
}

func getTerminalSizeChanges(fp *os.File, sc chan Size, done chan struct{}) (error) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, unix.SIGWINCH)
	go func() {
		for {
			select {
			case <-ch:
				s := Size{}
				var err error
				s.Width, s.Height, err = getTerminalSize(fp)
				if err == nil {
					sc <- s
				}
			case <-done:
				signal.Reset(unix.SIGWINCH)
				close(ch)
				return
			}
		}
	}()

	return nil
}
