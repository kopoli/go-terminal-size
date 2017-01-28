package terminal_size

import (
	"testing"
	"time"
)

func TestGetSize(t *testing.T) {

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
	sc, err := NewSizeChanger()

	if err != nil {
		t.Fatal("Creating SizeChanger failed with", err)
	}

	triggerFakeResize()
	select {
	case s := <-sc.Change:
		if s.Width == 0 || s.Height == 0 {
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
