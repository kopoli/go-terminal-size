// +build windows

package terminal_size

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32                   = windows.NewLazyDLL("kernel32.dll")
	GetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
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

func getTerminalSize(fp *os.File) (width int, height int, err error) {
	csbi := consoleScreenBufferInfo{}
	ret, _, err := GetConsoleScreenBufferInfo.Call(uintptr(windows.Handle(fp.Fd())), uintptr(unsafe.Pointer(&csbi)))

	fmt.Println("Ret on", ret, "ja err", err)

	fmt.Println(csbi)

	if ret == 0 {
		return
	}

	width = int(csbi.size.x)
	height = int(csbi.size.y)

	return
}

func getTerminalSizeChanges(fp *os.File, sc chan Size, done chan struct{}) error {

	return nil
}
