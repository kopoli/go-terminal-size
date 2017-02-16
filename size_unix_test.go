// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package tsize

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func genFakeSyscall(s Size) func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {

	return func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
		ws := (*winsize)(unsafe.Pointer(a3))
		(*ws).cols = uint16(s.Width)
		(*ws).rows = uint16(s.Height)

		return 0, 0, 0
	}
}

func triggerFakeResize() {
	unix.Kill(unix.Getpid(), unix.SIGWINCH)
}

func fakeSize(s Size) {
	unixSyscall = genFakeSyscall(s)
}

func unFakeSize() {
	unixSyscall = unix.Syscall
}
