package main

import (
	"fmt"

	tsize "github.com/kopoli/go-terminal-size"
)

func printSize(s tsize.Size) {
	fmt.Println("Current size is", s.Width, "by", s.Height)
}

func main() {

	s, err := tsize.GetSize()
	if err != nil {
		fmt.Println("Getting terminal size failed:", err)
		return
	}
	printSize(s)

	sc, err := tsize.NewSizeListener()
	if err != nil {
		fmt.Println("initializing failed:", err)
		return
	}

	defer sc.Close()

	for {
		select {
		case s = <-sc.Change:
			printSize(s)
		}
	}
}
