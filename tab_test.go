package main

import (
	"testing"
)

func TestTabSize(t *testing.T) {
	// Test rune_advance_len with different tab sizes
	tests := []struct {
		rune        rune
		pos         int
		tabsize     int
		expected    int
		description string
	}{
		{'\t', 0, 4, 4, "Tab at position 0 with tabsize 4"},
		{'\t', 1, 4, 3, "Tab at position 1 with tabsize 4"},
		{'\t', 2, 4, 2, "Tab at position 2 with tabsize 4"},
		{'\t', 3, 4, 1, "Tab at position 3 with tabsize 4"},
		{'\t', 4, 4, 4, "Tab at position 4 with tabsize 4"},
		{'\t', 0, 8, 8, "Tab at position 0 with tabsize 8"},
		{'\t', 1, 8, 7, "Tab at position 1 with tabsize 8"},
		{'\t', 4, 8, 4, "Tab at position 4 with tabsize 8"},
		{'a', 0, 4, 1, "Regular character should always be 1"},
		{'a', 5, 8, 1, "Regular character should always be 1"},
	}

	for _, test := range tests {
		result := rune_advance_len(test.rune, test.pos, test.tabsize)
		if result != test.expected {
			t.Errorf("%s: expected %d, got %d", test.description, test.expected, result)
		}
	}
}

func TestGemacsTabSizeDefault(t *testing.T) {
	// Test that new gemacs instance has correct default tab size
	InitTestScreen()
	g := new_gemacs([]string{})
	if g.tabstop_length != default_tabstop_length {
		t.Errorf("Expected default tab size %d, got %d", default_tabstop_length, g.tabstop_length)
	}
	if g.tabstop_length != 4 {
		t.Errorf("Expected default tab size to be 4, got %d", g.tabstop_length)
	}
}

func TestVlen(t *testing.T) {
	// Test vlen function with different tab sizes
	tests := []struct {
		input       string
		pos         int
		tabsize     int
		expected    int
		description string
	}{
		{"\t", 0, 4, 4, "Single tab with tabsize 4"},
		{"\t", 0, 8, 8, "Single tab with tabsize 8"},
		{"a\t", 0, 4, 4, "Character + tab with tabsize 4 (advances to next tabstop)"},
		{"ab\t", 0, 4, 4, "Two chars + tab with tabsize 4"},
		{"abc\t", 0, 4, 4, "Three chars + tab with tabsize 4"},
		{"abcd\t", 0, 4, 8, "Four chars + tab with tabsize 4"},
		{"hello", 0, 4, 5, "Regular text should be same as length"},
	}

	for _, test := range tests {
		result := vlen([]byte(test.input), test.pos, test.tabsize)
		if result != test.expected {
			t.Errorf("%s: expected %d, got %d", test.description, test.expected, result)
		}
	}
}