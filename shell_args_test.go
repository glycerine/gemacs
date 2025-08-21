package main

import (
	"reflect"
	"testing"
)

func TestGetShellArgs(t *testing.T) {
	tests := []struct {
		shellCmd string
		expected []string
		description string
	}{
		{"/bin/bash", []string{"-i", "+m"}, "bash should have job control disabled"},
		{"/usr/bin/bash", []string{"-i", "+m"}, "bash with full path should have job control disabled"},
		{"/bin/zsh", []string{"-i"}, "zsh should use interactive mode"},
		{"/usr/local/bin/fish", []string{"-i"}, "fish should use interactive mode"},
		{"/bin/sh", []string{"-i"}, "sh should use interactive mode"},
		{"/usr/bin/tcsh", []string{"-i"}, "unknown shell should default to interactive mode"},
		{"bash", []string{"-i", "+m"}, "bash without path should work"},
		{"zsh", []string{"-i"}, "zsh without path should work"},
	}
	
	for _, test := range tests {
		result := getShellArgs(test.shellCmd)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("%s: expected %v, got %v", test.description, test.expected, result)
		}
	}
}

func TestShellJobControlFix(t *testing.T) {
	// Test that bash gets the +m flag to disable job control
	bashArgs := getShellArgs("/bin/bash")
	
	foundJobControlDisable := false
	for _, arg := range bashArgs {
		if arg == "+m" {
			foundJobControlDisable = true
			break
		}
	}
	
	if !foundJobControlDisable {
		t.Error("bash should have +m flag to disable job control and prevent warnings")
	}
	
	// Test that other shells don't get the +m flag
	zshArgs := getShellArgs("/bin/zsh")
	for _, arg := range zshArgs {
		if arg == "+m" {
			t.Error("zsh should not have +m flag")
			break
		}
	}
}