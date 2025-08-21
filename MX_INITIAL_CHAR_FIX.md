# M-x Initial Character Fix

## Problem
When using `M-x` command:
1. User presses `Alt-x` to enter M-x mode
2. User types "s" as the first character of "shell"
3. The "s" character was lost/consumed and didn't appear in the command line
4. Only a space appeared instead of "s"

## Root Cause
The issue was in the `mx_mode.on_key()` method:
1. When M-x mode is activated, it immediately switches to line edit mode
2. The first key event (containing "s") was not passed to the line edit mode
3. The line edit mode started with empty `initial_content`
4. The "s" character was lost in the transition

## Solution
Modified `mx_mode.on_key()` to capture and preserve the initial character:

### Before:
```go
func (m mx_mode) on_key(ev *termbox.Event) {
    // ... setup code ...
    
    lemp := line_edit_mode_params{
        prompt:        "M-x ",
        initial_content: "",  // Always empty!
        // ... rest of params ...
    }
```

### After:
```go
func (m mx_mode) on_key(ev *termbox.Event) {
    // ... setup code ...
    
    // Capture the initial character if it's printable
    initialContent := ""
    if ev.Ch != 0 && ev.Ch >= 32 { // Printable character
        initialContent = string(ev.Ch)
    }
    
    lemp := line_edit_mode_params{
        prompt:        "M-x ",
        initial_content: initialContent,  // Preserves the character!
        // ... rest of params ...
    }
```

## Technical Details

### Character Detection Logic
- **Printable Check**: `ev.Ch != 0 && ev.Ch >= 32`
- **ASCII Range**: Characters 32-126 are standard printable ASCII
- **Non-printable**: Control characters, function keys, etc. are ignored
- **Safety**: Only captures actual character input, not special keys

### Printable Characters Handled
- **Letters**: a-z, A-Z
- **Numbers**: 0-9
- **Symbols**: !, @, #, $, %, etc.
- **Space**: Space character (ASCII 32)
- **Punctuation**: All standard punctuation marks

### Non-printable Characters Ignored
- **Control Keys**: Ctrl+C, Ctrl+G, etc.
- **Function Keys**: F1, F2, arrow keys, etc.
- **Special Keys**: Enter, Escape, Tab, Backspace, etc.
- **Control Characters**: Characters below ASCII 32

## User Experience Impact

### Before the Fix:
```
User: Alt-x s h e l l
Screen shows: M-x  hell    (missing 's')
```

### After the Fix:
```
User: Alt-x s h e l l  
Screen shows: M-x shell    (complete word)
```

## Testing
Added comprehensive test suite:
- `TestMXModeInitialCharacter`: Verifies 's' is captured correctly
- `TestMXModeNonPrintableCharacter`: Ensures special keys are ignored
- `TestMXModeCharacterRange`: Tests all printable ASCII characters

## Edge Cases Handled
1. **Function Keys**: F1, F2, etc. don't create spurious characters
2. **Control Sequences**: Ctrl+key combinations are ignored
3. **Arrow Keys**: Navigation keys don't interfere
4. **International Characters**: Unicode characters above 127 work correctly
5. **Empty Events**: Malformed events with no character are handled safely

## Future Considerations
- Could extend to handle multi-byte Unicode characters
- Could add support for initial character sequences
- Could implement more sophisticated input preprocessing

This fix ensures that M-x command input works exactly as users expect, with the first typed character immediately visible and processed correctly.