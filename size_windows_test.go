// +build windows

package terminal_size

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

			// spew.Dump(a)
			// fmt.Println("JEJEE!")
			// b := &a[1]
			b := (*uint32)(unsafe.Pointer(a[1]))
			// *a[1] = 7
			*b = 7
			// spew.Dump(a)
		},
	}
	fakeSetConsoleMode = &mockProc{
		r1:      1,
		lastErr: nil,
	}
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

	// Start mocking
	getConsoleMode = fakeGetConsoleMode
	setConsoleMode = fakeSetConsoleMode
	readConsoleInput = fakeReadConsoleInput
	// getConsoleScreenBufferInfo proc = kernel32.NewProc("GetConsoleScreenBufferInfo")
	// getConsoleMode             proc = kernel32.NewProc("GetConsoleMode")
	// setConsoleMode             proc = kernel32.NewProc("SetConsoleMode")
	// readConsoleInput           proc = kernel32.NewProc("ReadConsoleInputW")
}

func triggerFakeResize() {
}
