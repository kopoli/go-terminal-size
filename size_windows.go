// +build windows

package terminal_size

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/sys/windows"
)

var (
	kernel32                   = windows.NewLazySystemDLL("kernel32")
	getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	getConsoleMode             = kernel32.NewProc("GetConsoleMode")
	setConsoleMode             = kernel32.NewProc("SetConsoleMode")
	readConsoleInput           = kernel32.NewProc("ReadConsoleInputW")
)

type coord struct {
	x int16
	y int16
}

type smallRect struct {
	left   int16
	top    int16
	right  int16
	bottom int16
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        uint16
	window            smallRect
	maximumWindowSize coord
}

// Console modes
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033.aspx
const (
	enableProcessedInput uint32 = 1 << iota
	enableLineInput
	enableEchoInput
	enableWindowInput
	enableMouseInput
	enableInsertMode
	enableQuickEditMode
	enableExtendedFlags

	enableVirtualTerminalInput uint32 = 0x0200
)

type windowBufferSizeRecord struct {
	size coord
}

const (
	windowBufferSizeEvent uint16 = 0x0004
)

// INPUT_RECORD is defined in https://msdn.microsoft.com/en-us/library/windows/desktop/ms683499(v=vs.85).aspx
type inputRecord struct {
	eventType uint16
	win       windowBufferSizeRecord

	// Largest sub-struct in the union is the KEY_EVENT_RECORD with 4+2+2+2+2+4=16 bytes
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms684166(v=vs.85).aspx

	// Make the inputRecord be of the same size as INPUT_RECORD. We are
	// only interested on the WINDOW_BUFFER_SIZE_RECORDs, so buffer this with 16-4=12 bytes
	buf [12]byte
}

func getTerminalSize(fp *os.File) (width int, height int, err error) {
	csbi := consoleScreenBufferInfo{}
	ret, _, err := getConsoleScreenBufferInfo.Call(uintptr(windows.Handle(fp.Fd())),
		uintptr(unsafe.Pointer(&csbi)))

	fmt.Println("Ret on", ret, "ja err", err)

	fmt.Println(csbi)

	if ret == 0 {
		return
	}

	err = nil
	width = int(csbi.size.x)
	height = int(csbi.size.y)

	// Debug
	// var oldmode uint32
	// ret, _ , err = getConsoleMode.Call(uintptr(windows.Handle(fp.Fd())),
	// 	uintptr(unsafe.Pointer(&oldmode)))

	// fmt.Println("Old mode is", oldmode, "Ret", ret)

	return
}

// changes can be read with https://msdn.microsoft.com/en-us/library/windows/desktop/ms685035.aspx

func getTerminalSizeChanges(sc chan Size, done chan struct{}) (err error) {

	var oldmode, newmode uint32

	// Get
	handle := uintptr(windows.Handle(os.Stdin.Fd()))
	ret, _, err := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&oldmode)))

	fmt.Println("Old mode is", oldmode, "Ret", ret, "err", err)

	if ret == 0 {
		return
	}

	newmode = oldmode | enableWindowInput

	ret, _, err = setConsoleMode.Call(handle, uintptr(newmode))

	fmt.Println("new mode setting Ret", ret, "err", err, "newmode", newmode)

	if ret == 0 {
		return
	}

	go func() {
		var irs [8]inputRecord
		var count uint32

		for {
			ret, _, err := readConsoleInput.Call(handle,
				uintptr(unsafe.Pointer(&irs)),
				uintptr(len(irs)),
				uintptr(unsafe.Pointer(&count)),
			)

			fmt.Println("ret consolein", ret, "err", err, "count", count, "len", len(irs))

			if ret != 0 {
				spew.Dump("returned ", irs[:count])

				var i uint32
				for i = 0; i < count; i++ {
					if irs[i].eventType == windowBufferSizeEvent {
						var s Size

						// Getting the terminal size through Stdout gives the proper values.
						s.Width, s.Height, err = getTerminalSize(os.Stdout)
						spew.Dump(s)
						if err == nil {
							sc <- s
						}
						break
					}
				}
			}

			select {
			case <-done:
				getConsoleMode.Call(handle, uintptr(oldmode))
				return
			default:
			}
		}
	}()

	return nil
}
