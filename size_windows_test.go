// +build windows

package tsize

import (
	"testing"
	"unsafe"
)

type mockProc struct {
	r1, r2  uintptr
	lastErr error

	process func(a ...uintptr)
}

func (p *mockProc) Call(a ...uintptr) (uintptr, uintptr, error) {
	if p.process != nil {
		p.process(a...)
	}
	return p.r1, p.r2, p.lastErr
}

var (
	fakeGetConsoleMode = &mockProc{
		r1:      1,
		lastErr: nil,
		process: func(a ...uintptr) {
			b := (*uint32)(unsafe.Pointer(a[1]))
			*b = 7
		},
	}
	fakeSetConsoleMode = &mockProc{
		r1:      1,
		lastErr: nil,
	}

	// Always returns a windowBufferSizeEvent
	fakeReadConsoleInput = &mockProc{
		r1:      1,
		lastErr: nil,
		process: func(a ...uintptr) {
			irs := (*[8]inputRecord)(unsafe.Pointer(a[1]))
			count := (*uint32)(unsafe.Pointer(a[3]))

			(*irs)[0].eventType = windowBufferSizeEvent
			*count = 1
		},
	}

	origGetConsoleScreenBufferInfo proc
)

func fakeSize(s Size) {
	getConsoleScreenBufferInfo = &mockProc{
		r1:      1,
		lastErr: nil,
		process: func(a ...uintptr) {
			csbi := (*consoleScreenBufferInfo)(unsafe.Pointer(a[1]))
			csbi.size.x = int16(s.Width)
			csbi.size.y = int16(s.Height)
		},
	}
}

func unFakeSize() {
	getConsoleScreenBufferInfo = origGetConsoleScreenBufferInfo
}

func init() {
	// Mock some windows functions
	getConsoleMode = fakeGetConsoleMode
	setConsoleMode = fakeSetConsoleMode
	readConsoleInput = fakeReadConsoleInput
	origGetConsoleScreenBufferInfo = getConsoleScreenBufferInfo
}

func triggerFakeResize() {
}


func TestGetTerminalSize(t *testing.T) {
	value := Size{11, 12}
	fakeSize(value)
	defer unFakeSize()

	s, err := getTerminalSize(nil)

	if err != nil {
		t.Fatal("getTerminalSize should not error with:", err)
	}

	if value != s {
		t.Fatal("getTerminalSize should return the faked values")
	}
}
