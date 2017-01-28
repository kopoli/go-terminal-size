// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package terminal_size

import (
	"golang.org/x/sys/unix"
)

func triggerFakeResize() {
	unix.Kill(unix.Getpid(), unix.SIGWINCH)
}
