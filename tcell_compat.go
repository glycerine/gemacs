package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/termbox"
)

// Global tcell screen, initialized in main
var GlobalScreen tcell.Screen

// GetGlobalScreen returns the global tcell screen.
func GetGlobalScreen() tcell.Screen {
	return GlobalScreen
}

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
	Fg       tcell.Color
	Bg       tcell.Color
	Ellipsis rune
}

// DefaultLabelParams provides default label styling
var DefaultLabelParams = LabelParams{
	Fg:       tcell.ColorDefault,
	Bg:       tcell.ColorDefault,
	Ellipsis: '~',
}

// ScreenBuffer wraps a tcell Screen and provides buffered operations
type ScreenBuffer struct {
	Screen tcell.Screen
	Width  int
	Height int
	Rect   Rect
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
	return sb
}

// InitTestScreen initializes termbox for testing
func InitTestScreen() {
	// For tests, we can initialize termbox with simulation mode if needed
	// For now, just return since tests don't need actual screen
}

// TermboxBuffer creates a new screen buffer from the current termbox screen
func TermboxBuffer() *ScreenBuffer {
	return NewScreenBuffer(GlobalScreen)
}

// SetContent sets content at the given position with style in the buffer
func (sb *ScreenBuffer) SetContent(x, y int, primary rune, combining []rune, style tcell.Style) {
	if sb.Screen != nil {
		sb.Screen.SetContent(x, y, primary, combining, style)
	}
}

// Fill fills a rectangular area with the given character and style
func (sb *ScreenBuffer) Fill(r Rect, ch rune, style tcell.Style) {
	if sb.Screen != nil {
		for y := r.Y; y < r.Y+r.Height && y < sb.Height; y++ {
			for x := r.X; x < r.X+r.Width && x < sb.Width; x++ {
				sb.Screen.SetContent(x, y, ch, nil, style)
			}
		}
	}
}

// DrawLabel draws text within a rectangle with the given parameters
func (sb *ScreenBuffer) DrawLabel(r Rect, lp *LabelParams, text string) {
	style := tcell.StyleDefault
	if lp.Fg != tcell.ColorDefault {
		style = style.Foreground(lp.Fg)
	}
	if lp.Bg != tcell.ColorDefault {
		style = style.Background(lp.Bg)
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
	// Blit operation is typically handled by the underlying tcell.Screen
	// when drawing directly. For a ScreenBuffer that wraps a tcell.Screen,
	// this method might not be directly applicable in the same way as a
	// software buffer. If direct blitting is needed, it would involve
	// iterating and setting content on the destination screen.
	// For now, we'll assume direct drawing to the screen is sufficient.
	// If src.Screen is the same as sb.Screen, this is a no-op.
	// If src.Screen is a different screen, this operation is complex.
	// For simplicity, we will just ensure the content is drawn to the main screen.

	// This method might need a more sophisticated implementation if ScreenBuffer
	// is intended to be an off-screen buffer that is then blitted.
	// Given the current usage, it seems to be used for drawing directly.

	// If src is a software buffer (not wrapping a tcell.Screen directly),
	// then we would iterate its internal cells and set them on sb.Screen.
	// However, with sb.Screen being a tcell.Screen, this is more about
	// drawing directly to the screen.

	// For now, we will make this a no-op or log a warning if it's called
	// in a way that implies off-screen rendering that needs blitting.
	// The primary drawing is done by individual SetContent calls.
}

// Resize changes the size of the screen buffer
func (sb *ScreenBuffer) Resize(w, h int) {
	sb.Width = w
	sb.Height = h
	sb.Rect = Rect{0, 0, w, h}
	// No internal buffer to re-initialize as tcell.Screen handles it
}

// Set is an alias for SetContent to maintain compatibility
func (sb *ScreenBuffer) Set(x, y int, cell termbox.Cell) {
	style := MakeStyle(tcell.Color(cell.Fg), tcell.Color(cell.Bg))
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
func MakeStyle(fg, bg tcell.Color) tcell.Style {
	style := tcell.StyleDefault
	if fg != tcell.ColorDefault {
		style = style.Foreground(fg)
	}
	if bg != tcell.ColorDefault {
		style = style.Background(bg)
	}
	return style
}

// PollEvent polls for events from the tcell screen and converts them to termbox events
func PollEvent() termbox.Event {
	if GlobalScreen == nil {
		panic("Global screen not initialized")
	}
	
	ev := GlobalScreen.PollEvent()
	
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