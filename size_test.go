package tsize

import (
	"os"
	"testing"
	"time"

	isatty "github.com/mattn/go-isatty"
)

func TestGetSize(t *testing.T) {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		t.Skip("Skipping real terminal size test as not connected to a terminal")
	}

	s, err := GetSize()

	if err != nil {
		t.Fatal("Failed with", err)
	}

	if s.Width == 0 || s.Height == 0 {
		t.Fatal("Terminal size should not be", s.Width, s.Height)
	}
}

func TestFgetSize(t *testing.T) {
	_, err := FgetSize(nil)

	if err != NotATerminal {
		t.Fatal("Should fail with NotATerminal")
	}
}

func TestSizeChanger(t *testing.T) {
	defSize := Size{10, 20}
	fakeSize(defSize)
	defer unFakeSize()

	sc, err := NewSizeChanger()

	if err != nil {
		t.Fatal("Creating SizeChanger failed with", err)
	}

	triggerFakeResize()
	select {
	case s := <-sc.Change:
		if s.Width != defSize.Width || s.Height != defSize.Height {
			t.Fatal("Terminal size should not be", s.Width, s.Height)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Resize didn't trigger")
	}

	sc.Close()
	if sc.Change != nil {
		t.Fatal("Closing should nil the Change channel")
	}
}
