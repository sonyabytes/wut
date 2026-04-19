package ui

import (
	"bytes"
	"strings"
	"testing"
)

// ANSI helpers write a fixed byte pattern; these tests lock the contract so a
// "harmless rewrite" can't silently disable the cursor hide/show or change
// the dim color to something else.

func TestDimWrapsWithReset(t *testing.T) {
	got := dim("hi")
	if !strings.HasPrefix(got, "\x1b[2m") {
		t.Errorf("dim missing SGR 2 prefix: %q", got)
	}
	if !strings.HasSuffix(got, "\x1b[0m") {
		t.Errorf("dim missing reset suffix: %q", got)
	}
	if !strings.Contains(got, "hi") {
		t.Errorf("dim dropped content: %q", got)
	}
}

func TestHideShowCursor(t *testing.T) {
	var buf bytes.Buffer
	hideCursor(&buf)
	if buf.String() != "\x1b[?25l" {
		t.Errorf("hideCursor = %q, want \\x1b[?25l", buf.String())
	}
	buf.Reset()
	showCursor(&buf)
	if buf.String() != "\x1b[?25h" {
		t.Errorf("showCursor = %q, want \\x1b[?25h", buf.String())
	}
}

// Select's public error contract: with no options, it refuses to draw rather
// than opening a tty. Exercises the validation branch without needing a tty.
func TestSelectRejectsEmptyOptions(t *testing.T) {
	if _, err := Select("t", nil, 0); err == nil {
		t.Fatal("Select(nil opts) should error")
	}
}
