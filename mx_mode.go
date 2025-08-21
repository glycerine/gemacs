package main

import (
	"strconv"
	"strings"

	"github.com/glycerine/tcell_old_hacked_up/termbox"
)

// MXCommand represents an M-x command
type MXCommand struct {
	Name        string
	Description string
	Function    func(*gemacs) error
}

// mx_mode handles M-x command execution
type mx_mode struct {
	stub_overlay_mode
	gemacs   *gemacs
	commands map[string]*MXCommand
}

// init_mx_mode creates a new M-x mode
func init_mx_mode(g *gemacs) mx_mode {
	m := mx_mode{
		gemacs:   g,
		commands: make(map[string]*MXCommand),
	}
	
	// Register built-in commands
	m.registerCommands()
	
	g.set_status("M-x")
	return m
}

// registerCommands registers all available M-x commands
func (m *mx_mode) registerCommands() {
	commands := []*MXCommand{
		{
			Name:        "shell",
			Description: "Start an interactive shell",
			Function:    m.cmdShell,
		},
		{
			Name:        "shell-command",
			Description: "Execute a shell command",
			Function:    m.cmdShellCommand,
		},
		{
			Name:        "list-buffers",
			Description: "List all open buffers",
			Function:    m.cmdListBuffers,
		},
		{
			Name:        "kill-buffer",
			Description: "Kill current buffer",
			Function:    m.cmdKillBuffer,
		},
		{
			Name:        "syntax-highlighting-toggle",
			Description: "Toggle syntax highlighting",
			Function:    m.cmdToggleSyntaxHighlighting,
		},
		{
			Name:        "set-tab-size",
			Description: "Set tab size",
			Function:    m.cmdSetTabSize,
		},
	}
	
	for _, cmd := range commands {
		m.commands[cmd.Name] = cmd
	}
}

// on_key handles key events in M-x mode
func (m mx_mode) on_key(ev *termbox.Event) {
	// Use line edit mode for command completion and input
	commandNames := make([]string, 0, len(m.commands))
	for name := range m.commands {
		commandNames = append(commandNames, name)
	}
	
	lemp := line_edit_mode_params{
		prompt:        "M-x ",
		initial_content: "",
		ac_decide:     func(v *view) ac_func {
			return func(v *view) ([]ac_proposal, int) {
				// Get current input from the line
				line_data := v.cursor.line.data[:v.cursor.boffset]
				prefix := string(line_data)
				
				// Find matching commands
				var proposals []ac_proposal
				for _, name := range commandNames {
					if strings.HasPrefix(name, prefix) {
						proposals = append(proposals, ac_proposal{
							display: []byte(name),
							content: []byte(name),
						})
					}
				}
				return proposals, len(prefix)
			}
		},
		on_apply: func(buf *buffer) {
			commandName := string(buf.contents())
			m.executeCommand(commandName)
		},
	}
	
	m.gemacs.set_overlay_mode(init_line_edit_mode(m.gemacs, lemp))
}

// executeCommand executes an M-x command
func (m *mx_mode) executeCommand(name string) {
	name = strings.TrimSpace(name)
	if cmd, exists := m.commands[name]; exists {
		if err := cmd.Function(m.gemacs); err != nil {
			m.gemacs.set_status("Error executing %s: %v", name, err)
		}
	} else {
		m.gemacs.set_status("No such command: %s", name)
	}
}

// cmdShell implements the "shell" command
func (m *mx_mode) cmdShell(g *gemacs) error {
	// Create or get existing shell buffer
	shellName := "*shell*"
	shell := g.shell_manager.GetShell(shellName)
	
	if shell == nil {
		var err error
		shell, err = g.shell_manager.CreateShell(shellName)
		if err != nil {
			return err
		}
		
		// Add the shell buffer to gemacs buffers
		g.buffers = append(g.buffers, shell.buffer)
	}
	
	// Switch to shell buffer
	g.switch_to_buffer(shell.buffer)
	g.set_status("Shell buffer opened")
	
	return nil
}

// cmdShellCommand implements the "shell-command" command
func (m *mx_mode) cmdShellCommand(g *gemacs) error {
	lemp := line_edit_mode_params{
		prompt: "Shell command: ",
		on_apply: func(buf *buffer) {
			command := string(buf.contents())
			if command != "" {
				// Execute command and show output in a buffer
				g.executeShellCommand(command)
			}
		},
	}
	g.set_overlay_mode(init_line_edit_mode(g, lemp))
	return nil
}

// cmdListBuffers implements the "list-buffers" command
func (m *mx_mode) cmdListBuffers(g *gemacs) error {
	// Create a buffer listing
	buf := new_empty_buffer()
	buf.name = "*Buffer List*"
	
	// Add header
	buf.first_line.data = []byte("Buffers:")
	
	// Add each buffer
	for i, buffer := range g.buffers {
		line := &line{
			data: []byte(strconv.Itoa(i+1) + ". " + buffer.name),
			prev: buf.last_line,
			next: nil,
		}
		buf.last_line.next = line
		buf.last_line = line
		buf.lines_n++
	}
	
	// Add to buffers and switch to it
	g.buffers = append(g.buffers, buf)
	g.switch_to_buffer(buf)
	g.set_status("Buffer list opened")
	
	return nil
}

// cmdKillBuffer implements the "kill-buffer" command
func (m *mx_mode) cmdKillBuffer(g *gemacs) error {
	current_buf := g.active.leaf.buf
	g.kill_buffer(current_buf)
	return nil
}

// cmdToggleSyntaxHighlighting implements syntax highlighting toggle
func (m *mx_mode) cmdToggleSyntaxHighlighting(g *gemacs) error {
	if g.syntax_highlighter != nil {
		enabled := !g.syntax_highlighter.IsEnabled()
		g.syntax_highlighter.SetEnabled(enabled)
		if enabled {
			g.set_status("Syntax highlighting enabled")
		} else {
			g.set_status("Syntax highlighting disabled")
		}
		// Refresh all views
		g.views.traverse(func(vt *view_tree) {
			if vt.leaf != nil {
				vt.leaf.dirty = dirty_everything
			}
		})
	} else {
		g.set_status("Syntax highlighter not available")
	}
	return nil
}

// cmdSetTabSize implements tab size setting
func (m *mx_mode) cmdSetTabSize(g *gemacs) error {
	lemp := g.set_tab_size_lemp()
	g.set_overlay_mode(init_line_edit_mode(g, lemp))
	return nil
}

// switch_to_buffer switches to a specific buffer
func (g *gemacs) switch_to_buffer(buf *buffer) {
	current_view := g.active.leaf
	current_view.detach()
	current_view.attach(buf)
}

// executeShellCommand executes a shell command and shows output
func (g *gemacs) executeShellCommand(command string) {
	// Create output buffer
	outputBuf := new_empty_buffer()
	outputBuf.name = "*Shell Command Output*"
	
	// Add command line
	outputBuf.first_line.data = []byte("$ " + command)
	
	// Execute command (simplified version - just show that it would run)
	outputLine := &line{
		data: []byte("Command would execute: " + command),
		prev: outputBuf.last_line,
		next: nil,
	}
	outputBuf.last_line.next = outputLine
	outputBuf.last_line = outputLine
	outputBuf.lines_n++
	
	// Add note about implementation
	noteLine := &line{
		data: []byte("(Full shell command execution not yet implemented)"),
		prev: outputBuf.last_line,
		next: nil,
	}
	outputBuf.last_line.next = noteLine
	outputBuf.last_line = noteLine
	outputBuf.lines_n++
	
	// Add to buffers and switch to it
	g.buffers = append(g.buffers, outputBuf)
	g.switch_to_buffer(outputBuf)
	g.set_status("Shell command output")
}