package main

import (
	"testing"

	"github.com/glycerine/tcell_old_hacked_up/termbox"
)

func TestMXModeInitialCharacter(t *testing.T) {
	// Create a mock event for the 's' character
	ev := &termbox.Event{
		Type: termbox.EventKey,
		Ch:   's',
		Key:  0,
	}
	
	// Test the initial content logic (same as in mx_mode.on_key)
	initialContent := ""
	if ev.Ch != 0 && ev.Ch >= 32 { // Same logic as in on_key
		initialContent = string(ev.Ch)
	}
	
	if initialContent != "s" {
		t.Errorf("Expected initial content to be 's', got '%s'", initialContent)
	}
}

func TestMXModeNonPrintableCharacter(t *testing.T) {
	// Test that non-printable characters don't get included
	testCases := []struct {
		ch       rune
		key      termbox.Key
		expected string
		desc     string
	}{
		{0, termbox.KeyEnter, "", "Enter key should not be captured"},
		{0, termbox.KeyEsc, "", "Escape key should not be captured"},
		{0, termbox.KeyBackspace, "", "Backspace should not be captured"},
		{' ', 0, " ", "Space character should be captured"},
		{'a', 0, "a", "Regular character should be captured"},
		{'S', 0, "S", "Capital character should be captured"},
		{'1', 0, "1", "Number character should be captured"},
		{'\t', 0, "", "Tab character should not be captured (< 32)"},
		{'\n', 0, "", "Newline character should not be captured (< 32)"},
	}
	
	for _, tc := range testCases {
		// Test the same logic as in mx_mode.on_key
		initialContent := ""
		if tc.ch != 0 && tc.ch >= 32 {
			initialContent = string(tc.ch)
		}
		
		if initialContent != tc.expected {
			t.Errorf("%s: expected '%s', got '%s'", tc.desc, tc.expected, initialContent)
		}
	}
}

func TestMXModeCharacterRange(t *testing.T) {
	// Test the printable character range (32-126 are standard printable ASCII)
	printableChars := []rune{'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?', '@',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'[', '\\', ']', '^', '_', '`',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'{', '|', '}', '~'}
	
	for _, ch := range printableChars {
		initialContent := ""
		if ch != 0 && ch >= 32 {
			initialContent = string(ch)
		}
		
		expected := string(ch)
		if initialContent != expected {
			t.Errorf("Printable character '%c' (code %d) should be captured, got '%s'", ch, ch, initialContent)
		}
	}
}