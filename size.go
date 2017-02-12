// Get terminal size. Supports Linux and Windows.

package tsize

import (
	"errors"
	"os"

	isatty "github.com/mattn/go-isatty"
)

// Terminal size in columns and rows as Width and Height, respectively.
type Size struct {
	Width  int
	Height int
}

// Error to return if the given file to FgetSize isn't a terminal
var NotATerminal = errors.New("Given file is not a terminal")

// Get the current terminal size.
func GetSize() (s Size, err error) {
	return FgetSize(os.Stdout)
}

// Get the terminal size of a given os.File.
func FgetSize(fp *os.File) (s Size, err error) {
	if fp == nil || !isatty.IsTerminal(fp.Fd()) {
		err = NotATerminal
		return
	}

	s.Width, s.Height, err = getTerminalSize(fp)
	return
}

// Listens to terminal size changes
type SizeChanger struct {
	Change <-chan Size

	done chan struct{}
}

// Stop listening to terminal size changes
func (sc *SizeChanger) Close() (err error) {
	if sc.done != nil {
		close(sc.done)
		sc.done = nil
		sc.Change = nil
	}

	return
}

// Create a new size change listener
func NewSizeChanger() (sc *SizeChanger, err error) {
	sc = &SizeChanger{}

	sizechan := make(chan Size, 1)
	sc.Change = sizechan
	sc.done = make(chan struct{})

	err = getTerminalSizeChanges(sizechan, sc.done)
	if err != nil {
		close(sizechan)
		close(sc.done)
		sc = &SizeChanger{}
		return
	}

	return
}
