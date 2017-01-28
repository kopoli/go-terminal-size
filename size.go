package terminal_size

import (
	"errors"
	"os"

	isatty "github.com/mattn/go-isatty"
)

type Size struct {
	Width  int
	Height int
}

var NotATerminal = errors.New("Given file is not a terminal")

func GetSize() (s Size, err error) {
	return FgetSize(os.Stdout)
}

func FgetSize(fp *os.File) (s Size, err error) {
	if fp == nil || !isatty.IsTerminal(fp.Fd()) {
		err = NotATerminal
		return
	}

	s.Width, s.Height, err = getTerminalSize(fp)
	return
}

type SizeChanger struct {
	Change <-chan Size

	done chan struct{}
}

func (sc *SizeChanger) Close() (err error) {
	if sc.done != nil {
		close(sc.done)
		sc.done = nil
		sc.Change = nil
	}

	return
}

func NewSizeChanger() (sc *SizeChanger, err error) {
	fp := os.Stdout
	_, err = FgetSize(fp)
	if err != nil {
		return
	}

	sc = &SizeChanger{}

	sizechan := make(chan Size, 1)
	sc.Change = sizechan
	sc.done = make(chan struct{})

	err = getTerminalSizeChanges(fp, sizechan, sc.done)
	if err != nil {
		close(sizechan)
		close(sc.done)
		sc = &SizeChanger{}
		return
	}

	return
}
