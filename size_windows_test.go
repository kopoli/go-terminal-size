// +build windows

package tsize

import "unsafe"

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
)

func init() {
	// Mock some windows functions
	getConsoleMode = fakeGetConsoleMode
	setConsoleMode = fakeSetConsoleMode
	readConsoleInput = fakeReadConsoleInput
}

func triggerFakeResize() {
}
