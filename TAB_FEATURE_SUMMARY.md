# Tab Rendering Control Feature

## Overview
This feature allows users to control how tabs are rendered in gemacs, changing the default from 8 spaces to 4 spaces and providing runtime configuration.

## Changes Made

### 1. **Default Tab Size Changed**
- Changed default tab stop length from 8 to 4 spaces
- Defined `default_tabstop_length = 4` constant

### 2. **Configurable Tab Size**
- Added `tabstop_length int` field to `gemacs` struct
- Initialized in `new_gemacs()` with default value
- Runtime configurable via command interface

### 3. **Updated Tab Calculation Functions**
- Modified `rune_advance_len()` to accept `tabstop_length` parameter
- Updated `vlen()` function to accept `tabstop_length` parameter  
- Updated `find_closest_offsets()` method to accept `tabstop_length` parameter
- Updated `voffset()` and `voffset_coffset()` methods to accept `tabstop_length` parameter

### 4. **Updated All Call Sites**
- View rendering in `draw_line()` now uses `v.g.tabstop_length`
- All cursor positioning functions now pass tab size from gemacs instance
- Fill region functionality updated to respect tab settings

### 5. **User Interface**
- Added **Ctrl-x t** command to set tab size interactively
- Prompts user for tab size (1-32)
- Shows current tab size as default
- Validates input and provides error messages
- Refreshes all views to apply new tab size immediately

## Usage

### Setting Tab Size
1. Press `Ctrl-x` to enter extended mode
2. Press `t` to set tab size
3. Enter desired tab size (1-32) and press Enter
4. Tab size is applied immediately to all views

### Current Features
- **Default**: 4 spaces per tab (changed from 8)
- **Range**: 1-32 spaces per tab
- **Scope**: Global setting affects all buffers and views
- **Persistence**: Setting lasts for current session
- **Validation**: Invalid inputs are rejected with error messages

## Technical Details

### Function Signatures Updated
```go
// Old signatures
func rune_advance_len(r rune, pos int) int
func vlen(data []byte, pos int) int
func (l *line) find_closest_offsets(voffset int) (bo, co, vo int)
func (c *cursor_location) voffset() (vo int)
func (c *cursor_location) voffset_coffset() (vo, co int)

// New signatures  
func rune_advance_len(r rune, pos int, tabstop_length int) int
func vlen(data []byte, pos int, tabstop_length int) int
func (l *line) find_closest_offsets(voffset int, tabstop_length int) (bo, co, vo int)
func (c *cursor_location) voffset(tabstop_length int) (vo int)
func (c *cursor_location) voffset_coffset(tabstop_length int) (vo, co int)
```

### Architecture
- Tab size is stored in the main `gemacs` struct
- All view operations access tab size via `v.g.tabstop_length`
- Functions that need tab calculation receive it as parameter
- Immediate refresh ensures changes are visible right away

## Testing
- Comprehensive unit tests for tab size calculations
- Tests verify correct tabstop positioning at various positions
- Tests ensure default tab size is properly set
- All existing tests continue to pass

This feature provides a much more modern default (4 spaces) while maintaining full configurability for user preferences.