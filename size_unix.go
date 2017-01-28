// +build linux

package terminal_size

import (
	"os"
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

	if errno == 0 {
		width = int(ws.cols)
		height = int(ws.rows)
	} else {
		err = errno
	}

	return
}
