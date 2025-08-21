package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/termbox"
)

// Note: Using termbox for screen management, so no global screen needed

//----------------------------------------------------------------------------
// tcell-based replacements for tulib functionality
//----------------------------------------------------------------------------

// Rect represents a rectangular area
type Rect struct {
	X, Y          int
	Width, Height int
}

// Intersection returns the intersection of two rectangles
func (r Rect) Intersection(other Rect) Rect {
	x1 := max(r.X, other.X)
	y1 := max(r.Y, other.Y)
	x2 := min(r.X+r.Width, other.X+other.Width)
	y2 := min(r.Y+r.Height, other.Y+other.Height)
	
	if x2 <= x1 || y2 <= y1 {
		return Rect{0, 0, 0, 0} // No intersection
	}
	
	return Rect{x1, y1, x2 - x1, y2 - y1}
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LabelParams holds styling parameters for text labels
type LabelParams struct {
	Fg       termbox.Attribute
	Bg       termbox.Attribute
	Ellipsis rune
}

// DefaultLabelParams provides default label styling
var DefaultLabelParams = LabelParams{
	Fg:       termbox.ColorDefault,
	Bg:       termbox.ColorDefault,
	Ellipsis: '~',
}

// ScreenBuffer wraps a tcell Screen and provides buffered operations
type ScreenBuffer struct {
	Screen tcell.Screen
	Width  int
	Height int
	Rect   Rect
	// Internal buffer to store content before rendering
	cells  [][]tcell.Style
	chars  [][]rune
}

// NewScreenBuffer creates a new screen buffer with the given dimensions
func NewScreenBuffer(screen tcell.Screen) *ScreenBuffer {
	w, h := 80, 24 // Default size, will be resized as needed
	if screen != nil {
		w, h = screen.Size()
	}
	sb := &ScreenBuffer{
		Screen: screen,
		Width:  w,
		Height: h,
		Rect:   Rect{0, 0, w, h},
	}
	sb.initBuffer()
	return sb
}

// initBuffer initializes the internal buffer arrays
func (sb *ScreenBuffer) initBuffer() {
	sb.cells = make([][]tcell.Style, sb.Height)
	sb.chars = make([][]rune, sb.Height)
	for y := 0; y < sb.Height; y++ {
		sb.cells[y] = make([]tcell.Style, sb.Width)
		sb.chars[y] = make([]rune, sb.Width)
		for x := 0; x < sb.Width; x++ {
			sb.cells[y][x] = tcell.StyleDefault
			sb.chars[y][x] = ' '
		}
	}
}

// InitTestScreen initializes termbox for testing
func InitTestScreen() {
	// For tests, we can initialize termbox with simulation mode if needed
	// For now, just return since tests don't need actual screen
}

// TermboxBuffer creates a new screen buffer from the current termbox screen
func TermboxBuffer() *ScreenBuffer {
	w, h := termbox.Size()
	sb := &ScreenBuffer{
		Screen: nil, // We don't need the screen reference for termbox mode
		Width:  w,
		Height: h,
		Rect:   Rect{0, 0, w, h},
	}
	sb.initBuffer()
	return sb
}

// SetContent sets content at the given position with style in the buffer
func (sb *ScreenBuffer) SetContent(x, y int, primary rune, combining []rune, style tcell.Style) {
	if x >= 0 && x < sb.Width && y >= 0 && y < sb.Height {
		sb.chars[y][x] = primary
		sb.cells[y][x] = style
		// Note: we ignore combining characters for simplicity
	}
}

// Fill fills a rectangular area with the given character and style
func (sb *ScreenBuffer) Fill(r Rect, ch rune, style tcell.Style) {
	for y := r.Y; y < r.Y+r.Height && y < sb.Height; y++ {
		for x := r.X; x < r.X+r.Width && x < sb.Width; x++ {
			sb.SetContent(x, y, ch, nil, style)
		}
	}
}

// DrawLabel draws text within a rectangle with the given parameters
func (sb *ScreenBuffer) DrawLabel(r Rect, lp *LabelParams, text string) {
	style := tcell.StyleDefault
	if lp.Fg != termbox.ColorDefault {
		style = style.Foreground(tcell.PaletteColor(int(lp.Fg) - 1))
	}
	if lp.Bg != termbox.ColorDefault {
		style = style.Background(tcell.PaletteColor(int(lp.Bg) - 1))
	}

	runes := []rune(text)
	x := r.X
	for _, ch := range runes {
		if x >= r.X+r.Width {
			// Truncate with ellipsis if text is too long
			if r.X+r.Width-1 >= r.X {
				sb.SetContent(r.X+r.Width-1, r.Y, lp.Ellipsis, nil, style)
			}
			break
		}
		if x >= 0 && x < sb.Width && r.Y >= 0 && r.Y < sb.Height {
			sb.SetContent(x, r.Y, ch, nil, style)
		}
		x++
	}
}

// Blit copies content from another screen buffer to this one
func (sb *ScreenBuffer) Blit(r Rect, srcX, srcY int, src *ScreenBuffer) {
	for y := 0; y < r.Height && y+r.Y < sb.Height; y++ {
		for x := 0; x < r.Width && x+r.X < sb.Width; x++ {
			srcYPos := srcY + y
			srcXPos := srcX + x
			if srcXPos < src.Width && srcYPos < src.Height && srcXPos >= 0 && srcYPos >= 0 {
				destX := r.X + x
				destY := r.Y + y
				if destX >= 0 && destY >= 0 && destX < sb.Width && destY < sb.Height {
					sb.chars[destY][destX] = src.chars[srcYPos][srcXPos]
					sb.cells[destY][destX] = src.cells[srcYPos][srcXPos]
				}
			}
		}
	}
}

// Resize changes the size of the screen buffer
func (sb *ScreenBuffer) Resize(w, h int) {
	sb.Width = w
	sb.Height = h
	sb.Rect = Rect{0, 0, w, h}
	sb.initBuffer()
}

// Flush renders the buffer content to the termbox screen
func (sb *ScreenBuffer) Flush() {
	for y := 0; y < sb.Height; y++ {
		for x := 0; x < sb.Width; x++ {
			// Convert tcell.Style to termbox attributes
			style := sb.cells[y][x]
			fg, bg, _ := style.Decompose()
			
			// Convert colors to termbox attributes (simplified)
			var tbFg, tbBg termbox.Attribute
			if fg == tcell.ColorDefault {
				tbFg = termbox.ColorDefault
			} else {
				// For now, just use basic color mapping
				tbFg = termbox.Attribute(fg + 1)
			}
			if bg == tcell.ColorDefault {
				tbBg = termbox.ColorDefault
			} else {
				tbBg = termbox.Attribute(bg + 1)
			}
			
			termbox.SetCell(x, y, sb.chars[y][x], tbFg, tbBg)
		}
	}
}

// Set is an alias for SetContent to maintain compatibility
func (sb *ScreenBuffer) Set(x, y int, cell termbox.Cell) {
	style := MakeStyle(cell.Fg, cell.Bg)
	sb.SetContent(x, y, cell.Ch, nil, style)
}

// KeyToString converts key information to a string representation
func KeyToString(key termbox.Key, ch rune, mod termbox.Modifier) string {
	if key == 0 { // KeyRune equivalent
		if mod&termbox.ModAlt != 0 {
			return fmt.Sprintf("M-%c", ch)
		}
		return string(ch)
	}

	var keyStr string
	switch key {
	case termbox.KeyF1:
		keyStr = "F1"
	case termbox.KeyF2:
		keyStr = "F2"
	case termbox.KeyF3:
		keyStr = "F3"
	case termbox.KeyF4:
		keyStr = "F4"
	case termbox.KeyF5:
		keyStr = "F5"
	case termbox.KeyF6:
		keyStr = "F6"
	case termbox.KeyF7:
		keyStr = "F7"
	case termbox.KeyF8:
		keyStr = "F8"
	case termbox.KeyF9:
		keyStr = "F9"
	case termbox.KeyF10:
		keyStr = "F10"
	case termbox.KeyF11:
		keyStr = "F11"
	case termbox.KeyF12:
		keyStr = "F12"
	case termbox.KeyInsert:
		keyStr = "Insert"
	case termbox.KeyDelete:
		keyStr = "Delete"
	case termbox.KeyHome:
		keyStr = "Home"
	case termbox.KeyEnd:
		keyStr = "End"
	case termbox.KeyArrowUp:
		keyStr = "Up"
	case termbox.KeyArrowDown:
		keyStr = "Down"
	case termbox.KeyArrowLeft:
		keyStr = "Left"
	case termbox.KeyArrowRight:
		keyStr = "Right"
	case termbox.KeyPgup:
		keyStr = "PgUp"
	case termbox.KeyPgdn:
		keyStr = "PgDn"
	case termbox.KeyEnter:
		keyStr = "Enter"
	case termbox.KeyEsc:
		keyStr = "Esc"
	case termbox.KeyTab:
		keyStr = "Tab"
	case termbox.KeyBackspace:
		keyStr = "Backspace"
	case termbox.KeySpace:
		keyStr = "Space"
	default:
		keyStr = fmt.Sprintf("Key(%d)", int(key))
	}

	if mod&termbox.ModAlt != 0 {
		keyStr = "M-" + keyStr
	}

	return keyStr
}

// Helper function to create a tcell Style from termbox attributes
func MakeStyle(fg, bg termbox.Attribute) tcell.Style {
	style := tcell.StyleDefault
	if fg != termbox.ColorDefault {
		style = style.Foreground(tcell.PaletteColor(int(fg) - 1))
	}
	if bg != termbox.ColorDefault {
		style = style.Background(tcell.PaletteColor(int(bg) - 1))
	}
	return style
}

// PollEvent polls for events from the tcell screen and converts them to termbox events
func PollEvent() termbox.Event {
	if globalScreen == nil {
		panic("Global screen not initialized")
	}
	
	ev := globalScreen.PollEvent()
	
	switch tev := ev.(type) {
	case *tcell.EventKey:
		key := tev.Key()
		ch := tev.Rune()
		mod := tev.Modifiers()
		
		// Convert tcell key to termbox key
		var tbKey termbox.Key
		if key == tcell.KeyRune {
			tbKey = termbox.Key(0) // KeyRune equivalent
		} else {
			tbKey = termbox.Key(key)
		}
		
		// Convert tcell modifiers to termbox modifiers
		var tbMod termbox.Modifier
		if mod&tcell.ModAlt != 0 {
			tbMod |= termbox.ModAlt
		}
		
		return termbox.Event{
			Type: termbox.EventKey,
			Key:  tbKey,
			Ch:   ch,
			Mod:  tbMod,
		}
		
	case *tcell.EventResize:
		w, h := tev.Size()
		return termbox.Event{
			Type:   termbox.EventResize,
			Width:  w,
			Height: h,
		}
		
	default:
		// For other events, return a none event
		return termbox.Event{
			Type: termbox.EventNone,
		}
	}
}