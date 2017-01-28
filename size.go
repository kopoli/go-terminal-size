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
