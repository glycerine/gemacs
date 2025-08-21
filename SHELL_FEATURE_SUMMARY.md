# Shell Feature Implementation

## Overview
This feature implements emacs-style shell functionality, allowing users to run interactive shells inside gemacs using the `M-x shell` command, just like in emacs.

## Features Implemented

### 1. **M-x Command System**
- Added `Alt-x` (M-x) key binding to access extended commands
- Implemented command completion and autocomplete
- Integrated with existing line edit mode for consistent UX

### 2. **Shell Buffer Management**
- **ShellManager**: Manages multiple shell instances
- **ShellBuffer**: Special buffer type for shell interaction
- **Process Management**: Starts and manages shell subprocesses
- **Cross-platform**: Works on Unix-like systems and Windows

### 3. **Interactive Shell Features**
- **Real-time I/O**: Captures shell output and displays in buffer
- **Command History**: Up/down arrow navigation through command history
- **Prompt Detection**: Handles shell prompts appropriately
- **Multiple Shells**: Support for multiple shell buffers

### 4. **Integration with gemacs**
- **Buffer System**: Shell buffers integrate seamlessly with existing buffer management
- **View Handling**: Special key handling for shell buffers
- **Status Line**: Shows shell command input in status line
- **Cleanup**: Proper cleanup of shell processes on exit

## Usage

### Starting a Shell
1. Press `Alt-x` (M-x) to enter command mode
2. Type `shell` and press Enter
3. A new shell buffer `*shell*` will open
4. Shell starts automatically with your default shell

### Interacting with Shell
- **Enter commands**: Press Enter in shell buffer to enter command mode
- **Type command**: Type your command in the status line
- **Execute**: Press Enter to send command to shell
- **History**: Use Up/Down arrows to navigate command history
- **Cancel**: Press Ctrl-G to cancel command input

### Navigation in Shell Buffer
- **Read-only**: Shell output is read-only, preventing accidental edits
- **Scroll**: Use normal navigation keys (arrows, Page Up/Down, etc.)
- **Cursor**: Move cursor to read shell output

## Available M-x Commands

| Command | Description |
|---------|-------------|
| `shell` | Start an interactive shell |
| `shell-command` | Execute a single shell command |
| `list-buffers` | List all open buffers |
| `kill-buffer` | Kill current buffer |
| `syntax-highlighting-toggle` | Toggle syntax highlighting |
| `set-tab-size` | Set tab size |

## Technical Architecture

### Core Components

1. **ShellManager** (`shell.go`)
   - Manages shell instances
   - Creates and tracks shell buffers
   - Handles cleanup

2. **ShellBuffer** (`shell.go`)
   - Extends regular buffer with shell functionality
   - Manages shell process lifecycle
   - Handles I/O with shell subprocess

3. **MX Mode** (`mx_mode.go`)
   - Implements M-x command system
   - Provides command completion
   - Executes commands

4. **Shell Mode** (`shell.go`)
   - Handles shell command input
   - Manages command history
   - Provides shell-specific key bindings

### Process Management
- **Cross-platform**: Detects OS and uses appropriate shell (bash/cmd)
- **Environment**: Inherits user environment variables
- **Interactive**: Starts shell in interactive mode
- **Cleanup**: Graceful termination with fallback to force kill

### I/O Handling
- **Asynchronous**: Separate goroutines for stdout/stderr reading
- **Real-time**: Output appears immediately in buffer
- **Thread-safe**: Mutex protection for concurrent access
- **Buffered**: Uses Go's bufio for efficient I/O

## Example Workflows

### Basic Shell Usage
```
1. Alt-x
2. Type: shell
3. Press Enter
4. Shell buffer opens with prompt
5. Press Enter to start typing commands
6. Type: ls -la
7. Press Enter to execute
8. Output appears in buffer
```

### Command History
```
1. In shell buffer, press Enter
2. Type: echo "first command"
3. Press Enter
4. Press Enter again for new command
5. Press Up arrow - shows "echo "first command""
6. Press Down arrow - clears input
```

### Multiple Shells
- Each `M-x shell` creates a new shell if none exists
- Existing shell buffer is reused if already open
- Future enhancement could support multiple numbered shells

## Testing
- **Unit Tests**: Comprehensive test suite covering all components
- **Shell Management**: Tests shell creation, execution, cleanup
- **M-x System**: Tests command registration and execution
- **Buffer Integration**: Tests shell buffer identification and handling
- **Process Lifecycle**: Tests shell startup and termination

## Platform Support
- **Unix/Linux**: Uses `/bin/bash` or `$SHELL` environment variable
- **macOS**: Uses user's default shell
- **Windows**: Uses `cmd` command prompt
- **Interactive Mode**: Starts shells in interactive mode for proper prompt handling

## Future Enhancements
- **Terminal Emulation**: More complete terminal emulation (colors, escape sequences)
- **Multiple Shells**: Support for multiple numbered shell buffers
- **Shell Integration**: Better integration with shell features (completion, etc.)
- **Custom Shells**: Configuration for custom shell commands
- **Background Execution**: Support for background process management

This implementation provides a solid foundation for shell interaction within gemacs, closely mimicking the emacs shell experience while integrating seamlessly with gemacs's existing architecture.