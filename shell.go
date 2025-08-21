package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/glycerine/tcell_old_hacked_up/termbox"
)

// ShellBuffer represents a shell buffer with process management
type ShellBuffer struct {
	*buffer              // Embed regular buffer
	cmd     *exec.Cmd    // Shell process
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	prompt  string       // Current shell prompt
	history []string     // Command history
	historyPos int       // Current position in history
	currentLine string   // Current input line
	promptLine int       // Line number where current prompt starts
	running bool         // Whether shell is running
	mutex   sync.Mutex   // Protect concurrent access
}

// ShellManager manages shell buffers
type ShellManager struct {
	shells map[string]*ShellBuffer
	mutex  sync.Mutex
}

// NewShellManager creates a new shell manager
func NewShellManager() *ShellManager {
	return &ShellManager{
		shells: make(map[string]*ShellBuffer),
	}
}

// CreateShell creates a new shell buffer
func (sm *ShellManager) CreateShell(name string) (*ShellBuffer, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	if shell, exists := sm.shells[name]; exists {
		return shell, nil
	}
	
	// Create the shell buffer
	shell := &ShellBuffer{
		buffer:      new_empty_buffer(),
		history:     make([]string, 0, 100),
		historyPos:  0,
		promptLine:  1,
		running:     false,
	}
	
	shell.name = name
	shell.buffer.name = name
	
	// Start the shell process
	if err := shell.startShell(); err != nil {
		return nil, err
	}
	
	sm.shells[name] = shell
	return shell, nil
}

// GetShell returns an existing shell buffer
func (sm *ShellManager) GetShell(name string) *ShellBuffer {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return sm.shells[name]
}

// CloseShell closes and removes a shell buffer
func (sm *ShellManager) CloseShell(name string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	if shell, exists := sm.shells[name]; exists {
		shell.Close()
		delete(sm.shells, name)
	}
	return nil
}

// getShellArgs returns appropriate arguments for different shell types
func getShellArgs(shellCmd string) []string {
	// Extract just the shell name from the path
	shellName := filepath.Base(shellCmd)
	
	switch shellName {
	case "bash":
		// Interactive mode with job control disabled to prevent warnings
		// +m disables job control monitoring
		// --norc prevents reading ~/.bashrc which might interfere
		// --noprofile prevents reading profile files
		return []string{"-i", "+m", "--norc", "--noprofile"}
	case "zsh":
		// Interactive mode, zsh handles job control better
		return []string{"-i"}
	case "fish":
		// Interactive mode for fish
		return []string{"-i"}
	case "sh":
		// Basic sh, interactive mode
		return []string{"-i"}
	default:
		// Default to interactive mode for unknown shells
		return []string{"-i"}
	}
}

// startShell starts the shell process
func (shell *ShellBuffer) startShell() error {
	// Determine shell command based on OS
	var shellCmd string
	var shellArgs []string
	
	if runtime.GOOS == "windows" {
		shellCmd = "cmd"
		shellArgs = []string{"/c"}
	} else {
		// Use user's shell or default to bash
		shellCmd = os.Getenv("SHELL")
		if shellCmd == "" {
			shellCmd = "/bin/bash"
		}
		
		// Set shell arguments based on shell type
		shellArgs = getShellArgs(shellCmd)
	}
	
	// Create the command
	shell.cmd = exec.Command(shellCmd, shellArgs...)
	
	// Set up environment
	env := os.Environ()
	// Add environment variables to improve shell experience
	env = append(env, "TERM=dumb")        // Indicate simple terminal capabilities
	env = append(env, "PAGER=cat")        // Use simple pager
	env = append(env, "EDITOR=gemacs")    // Set gemacs as editor
	shell.cmd.Env = env
	
	// Set up pipes
	var err error
	shell.stdin, err = shell.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	shell.stdout, err = shell.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	shell.stderr, err = shell.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}
	
	// Start the process
	if err := shell.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}
	
	shell.running = true
	
	// Start output readers
	go shell.readOutput()
	go shell.readErrors()
	
	// Add initial empty line for shell output
	shell.addLine("")
	
	return nil
}

// readOutput reads stdout from shell process
func (shell *ShellBuffer) readOutput() {
	scanner := bufio.NewScanner(shell.stdout)
	for scanner.Scan() {
		line := scanner.Text()
		shell.mutex.Lock()
		shell.addLine(line)
		shell.mutex.Unlock()
	}
}

// readErrors reads stderr from shell process
func (shell *ShellBuffer) readErrors() {
	scanner := bufio.NewScanner(shell.stderr)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Filter out known harmless bash warnings
		if strings.Contains(line, "no job control in this shell") {
			continue // Skip this harmless warning
		}
		
		shell.mutex.Lock()
		shell.addLine("ERROR: " + line)
		shell.mutex.Unlock()
	}
}

// addLine adds a line to the shell buffer
func (shell *ShellBuffer) addLine(text string) {
	// Create new line
	newLine := &line{
		data: []byte(text),
		prev: shell.last_line,
		next: nil,
	}
	
	// Link to buffer
	if shell.last_line != nil {
		shell.last_line.next = newLine
	} else {
		shell.first_line = newLine
	}
	shell.last_line = newLine
	shell.lines_n++
	
	// Update views
	for _, view := range shell.views {
		view.dirty = dirty_everything
		// Move cursor to end
		view.cursor.line = shell.last_line
		view.cursor.line_num = shell.lines_n
		view.cursor.boffset = len(shell.last_line.data)
	}
}

// SendCommand sends a command to the shell
func (shell *ShellBuffer) SendCommand(command string) error {
	shell.mutex.Lock()
	defer shell.mutex.Unlock()
	
	if !shell.running {
		return fmt.Errorf("shell is not running")
	}
	
	// Add to history
	if command != "" && (len(shell.history) == 0 || shell.history[len(shell.history)-1] != command) {
		shell.history = append(shell.history, command)
		if len(shell.history) > 100 {
			shell.history = shell.history[1:]
		}
	}
	shell.historyPos = len(shell.history)
	
	// Echo the command in the buffer
	shell.addLine(shell.prompt + command)
	
	// Send to shell
	_, err := shell.stdin.Write([]byte(command + "\n"))
	return err
}

// GetHistoryUp returns previous command in history
func (shell *ShellBuffer) GetHistoryUp() string {
	shell.mutex.Lock()
	defer shell.mutex.Unlock()
	
	if len(shell.history) == 0 {
		return ""
	}
	
	if shell.historyPos > 0 {
		shell.historyPos--
	}
	
	if shell.historyPos < len(shell.history) {
		return shell.history[shell.historyPos]
	}
	return ""
}

// GetHistoryDown returns next command in history
func (shell *ShellBuffer) GetHistoryDown() string {
	shell.mutex.Lock()
	defer shell.mutex.Unlock()
	
	if len(shell.history) == 0 {
		return ""
	}
	
	if shell.historyPos < len(shell.history)-1 {
		shell.historyPos++
		return shell.history[shell.historyPos]
	} else {
		shell.historyPos = len(shell.history)
		return ""
	}
}

// Close closes the shell and cleans up resources
func (shell *ShellBuffer) Close() error {
	shell.mutex.Lock()
	defer shell.mutex.Unlock()
	
	if !shell.running {
		return nil
	}
	
	shell.running = false
	
	// Close pipes
	if shell.stdin != nil {
		shell.stdin.Close()
	}
	if shell.stdout != nil {
		shell.stdout.Close()
	}
	if shell.stderr != nil {
		shell.stderr.Close()
	}
	
	// Terminate process
	if shell.cmd != nil && shell.cmd.Process != nil {
		// Try graceful termination first
		shell.cmd.Process.Signal(syscall.SIGTERM)
		
		// Wait a bit for graceful termination
		done := make(chan error, 1)
		go func() {
			done <- shell.cmd.Wait()
		}()
		
		select {
		case <-done:
			// Process terminated gracefully
		case <-time.After(2 * time.Second):
			// Force kill if not terminated
			shell.cmd.Process.Kill()
		}
	}
	
	shell.addLine("Shell closed.")
	return nil
}

// IsShellBuffer returns true if the buffer is a shell buffer
func IsShellBuffer(buf *buffer) bool {
	return strings.HasPrefix(buf.name, "*shell")
}

// ShellMode represents the shell interaction mode
type ShellMode struct {
	stub_overlay_mode
	gemacs      *gemacs
	shell       *ShellBuffer
	inputBuffer string
}

// InitShellMode creates a new shell mode for a view
func InitShellMode(g *gemacs, shell *ShellBuffer) *ShellMode {
	return &ShellMode{
		gemacs:      g,
		shell:       shell,
		inputBuffer: "",
	}
}

// on_key handles key events in shell mode
func (sm *ShellMode) on_key(ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeyEnter:
		// Send command
		if sm.inputBuffer != "" {
			sm.shell.SendCommand(sm.inputBuffer)
			sm.inputBuffer = ""
		} else {
			sm.shell.SendCommand("")
		}
		sm.gemacs.set_overlay_mode(nil)
		
	case termbox.KeyArrowUp:
		// History up
		sm.inputBuffer = sm.shell.GetHistoryUp()
		sm.updateStatus()
		
	case termbox.KeyArrowDown:
		// History down  
		sm.inputBuffer = sm.shell.GetHistoryDown()
		sm.updateStatus()
		
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		// Remove character
		if len(sm.inputBuffer) > 0 {
			sm.inputBuffer = sm.inputBuffer[:len(sm.inputBuffer)-1]
			sm.updateStatus()
		}
		
	case termbox.KeyCtrlG:
		// Cancel
		sm.gemacs.set_overlay_mode(nil)
		
	default:
		if ev.Ch != 0 {
			sm.inputBuffer += string(ev.Ch)
			sm.updateStatus()
		}
	}
}

// updateStatus updates the status line with current input
func (sm *ShellMode) updateStatus() {
	sm.gemacs.set_status("Shell> %s", sm.inputBuffer)
}

// draw draws the shell mode (status line shows input)
func (sm *ShellMode) draw() {
	sm.updateStatus()
}

// needs_cursor returns false as shell mode uses status line
func (sm *ShellMode) needs_cursor() bool {
	return false
}

// cursor_position returns cursor position (not used)
func (sm *ShellMode) cursor_position() (int, int) {
	return 0, 0
}

// on_resize handles resize events
func (sm *ShellMode) on_resize(ev *termbox.Event) {
	// Nothing special needed
}

// exit handles cleanup when exiting shell mode
func (sm *ShellMode) exit() {
	// Nothing special needed
}