package main

import (
	"strings"
	"testing"
	"time"
)

func TestShellManager(t *testing.T) {
	sm := NewShellManager()
	
	// Test creating a shell
	shell, err := sm.CreateShell("*test-shell*")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}
	
	if shell == nil {
		t.Fatal("Shell is nil")
	}
	
	if shell.name != "*test-shell*" {
		t.Errorf("Expected shell name '*test-shell*', got '%s'", shell.name)
	}
	
	if !shell.running {
		t.Error("Shell should be running")
	}
	
	// Test getting existing shell
	shell2 := sm.GetShell("*test-shell*")
	if shell2 != shell {
		t.Error("GetShell should return the same shell instance")
	}
	
	// Test closing shell
	err = sm.CloseShell("*test-shell*")
	if err != nil {
		t.Errorf("Failed to close shell: %v", err)
	}
	
	// Give a moment for cleanup
	time.Sleep(100 * time.Millisecond)
	
	if shell.running {
		t.Error("Shell should not be running after close")
	}
}

func TestShellBuffer(t *testing.T) {
	sm := NewShellManager()
	shell, err := sm.CreateShell("*test-shell-2*")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}
	defer sm.CloseShell("*test-shell-2*")
	
	// Test shell buffer properties
	if shell.buffer == nil {
		t.Fatal("Shell buffer is nil")
	}
	
	if shell.lines_n < 1 {
		t.Error("Shell should have at least one line")
	}
	
	// Test command history
	initialHistoryLen := len(shell.history)
	
	err = shell.SendCommand("echo hello")
	if err != nil {
		t.Errorf("Failed to send command: %v", err)
	}
	
	if len(shell.history) != initialHistoryLen+1 {
		t.Error("Command should be added to history")
	}
	
	if shell.history[len(shell.history)-1] != "echo hello" {
		t.Error("Last command in history should be 'echo hello'")
	}
	
	// Test history navigation
	upCmd := shell.GetHistoryUp()
	if upCmd != "echo hello" {
		t.Errorf("History up should return 'echo hello', got '%s'", upCmd)
	}
	
	downCmd := shell.GetHistoryDown()
	if downCmd != "" {
		t.Errorf("History down should return empty string, got '%s'", downCmd)
	}
}

func TestMXCommands(t *testing.T) {
	InitTestScreen()
	g := new_gemacs([]string{})
	
	// Test M-x mode creation
	mx := init_mx_mode(g)
	
	// Test that commands are registered
	if len(mx.commands) == 0 {
		t.Error("No M-x commands registered")
	}
	
	// Test shell command exists
	if _, exists := mx.commands["shell"]; !exists {
		t.Error("Shell command not registered")
	}
	
	// Test shell command execution
	err := mx.cmdShell(g)
	if err != nil {
		t.Errorf("Shell command failed: %v", err)
	}
	
	// Check that shell buffer was created
	shellFound := false
	for _, buf := range g.buffers {
		if buf.name == "*shell*" {
			shellFound = true
			break
		}
	}
	
	if !shellFound {
		t.Error("Shell buffer not created")
	}
	
	// Test buffer identification
	for _, buf := range g.buffers {
		if buf.name == "*shell*" {
			if !IsShellBuffer(buf) {
				t.Error("Shell buffer not identified correctly")
			}
			break
		}
	}
}

func TestShellBufferIdentification(t *testing.T) {
	// Test shell buffer identification
	shell_buf := new_empty_buffer()
	shell_buf.name = "*shell*"
	
	if !IsShellBuffer(shell_buf) {
		t.Error("Buffer with '*shell*' name should be identified as shell buffer")
	}
	
	regular_buf := new_empty_buffer()
	regular_buf.name = "test.go"
	
	if IsShellBuffer(regular_buf) {
		t.Error("Regular buffer should not be identified as shell buffer")
	}
	
	another_shell := new_empty_buffer()
	another_shell.name = "*shell-2*"
	
	if !IsShellBuffer(another_shell) {
		t.Error("Buffer with '*shell-*' name should be identified as shell buffer")
	}
}

func TestMXAutoCompletion(t *testing.T) {
	InitTestScreen()
	g := new_gemacs([]string{})
	mx := init_mx_mode(g)
	
	// Test that shell command can be found by autocomplete
	commandNames := make([]string, 0, len(mx.commands))
	for name := range mx.commands {
		commandNames = append(commandNames, name)
	}
	
	// Check that "shell" is in the command list
	found := false
	for _, name := range commandNames {
		if name == "shell" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("'shell' command not found in command list")
	}
	
	// Test prefix matching
	shellCommands := []string{}
	for _, name := range commandNames {
		if strings.HasPrefix(name, "shell") {
			shellCommands = append(shellCommands, name)
		}
	}
	
	if len(shellCommands) < 1 {
		t.Error("Should find at least one command starting with 'shell'")
	}
}