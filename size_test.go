package terminal_size

import "testing"

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
