package util

import "testing"

func TestStyleLine(t *testing.T) {
	line := "test"
	styled := StyleLine(line)
	if styled != line {
		t.Errorf("expected %s, got %s", line, styled)
	}
}

func TestStyleLine_Brew(t *testing.T) {
	line := "brew"
	styled := StyleLine(line)
	if styled != brewStyle.Render(line) {
		t.Errorf("expected %s, got %s", brewStyle.Render(line), styled)
	}
}
	
func TestStyleLine_Pacman(t *testing.T) {
	line := "pacman"
	styled := StyleLine(line)
	if styled != pacmanStyle.Render(line) {
		t.Errorf("expected %s, got %s", pacmanStyle.Render(line), styled)
	}
}
