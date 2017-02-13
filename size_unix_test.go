// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package tsize

import (
	"golang.org/x/sys/unix"
)

func triggerFakeResize() {
	unix.Kill(unix.Getpid(), unix.SIGWINCH)
}

func fakeSize(s Size) {
}

func unFakeSize() {
}
