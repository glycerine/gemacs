package tulib

import (
	"bytes"
	"github.com/gdamore/tcell/termbox"
)

func KeyToString(key termbox.Key, ch rune, mod termbox.Modifier) string {
	var buf bytes.Buffer
	if mod&termbox.ModAlt != 0 {
		buf.WriteString("M-")
	}

	switch key {
	case termbox.KeyF1:
		buf.WriteString("<f1>")
	case termbox.KeyF2:
		buf.WriteString("<f2>")
	case termbox.KeyF3:
		buf.WriteString("<f3>")
	case termbox.KeyF4:
		buf.WriteString("<f4>")
	case termbox.KeyF5:
		buf.WriteString("<f5>")
	case termbox.KeyF6:
		buf.WriteString("<f6>")
	case termbox.KeyF7:
		buf.WriteString("<f7>")
	case termbox.KeyF8:
		buf.WriteString("<f8>")
	case termbox.KeyF9:
		buf.WriteString("<f9>")
	case termbox.KeyF10:
		buf.WriteString("<f10>")
	case termbox.KeyF11:
		buf.WriteString("<f11>")
	case termbox.KeyF12:
		buf.WriteString("<f12>")
	case termbox.KeyInsert:
		buf.WriteString("<insert>")
	case termbox.KeyDelete:
		buf.WriteString("<delete>")
	case termbox.KeyHome:
		buf.WriteString("<home>")
	case termbox.KeyEnd:
		buf.WriteString("<end>")
	case termbox.KeyPgup:
		buf.WriteString("<pgup>")
	case termbox.KeyPgdn:
		buf.WriteString("<pgdn>")
	case termbox.KeyArrowUp:
		buf.WriteString("<up>")
	case termbox.KeyArrowDown:
		buf.WriteString("<down>")
	case termbox.KeyArrowLeft:
		buf.WriteString("<left>")
	case termbox.KeyArrowRight:
		buf.WriteString("<right>")
	case termbox.KeyCtrlSpace:
		if ch == 0 {
			buf.WriteString("C-<space>")
		}
	case termbox.KeyCtrlA:
		buf.WriteString("C-a")
	case termbox.KeyCtrlB:
		buf.WriteString("C-b")
	case termbox.KeyCtrlC:
		buf.WriteString("C-c")
	case termbox.KeyCtrlD:
		buf.WriteString("C-d")
	case termbox.KeyCtrlE:
		buf.WriteString("C-e")
	case termbox.KeyCtrlF:
		buf.WriteString("C-f")
	case termbox.KeyCtrlG:
		buf.WriteString("C-g")
	case termbox.KeyBackspace:
		buf.WriteString("<backspace>")
	case termbox.KeyTab:
		buf.WriteString("<tab>")
	case termbox.KeyCtrlJ:
		buf.WriteString("C-j")
	case termbox.KeyCtrlK:
		buf.WriteString("C-k")
	case termbox.KeyCtrlL:
		buf.WriteString("C-l")
	case termbox.KeyEnter:
		buf.WriteString("<enter>")
	case termbox.KeyCtrlN:
		buf.WriteString("C-n")
	case termbox.KeyCtrlO:
		buf.WriteString("C-o")
	case termbox.KeyCtrlP:
		buf.WriteString("C-p")
	case termbox.KeyCtrlQ:
		buf.WriteString("C-q")
	case termbox.KeyCtrlR:
		buf.WriteString("C-r")
	case termbox.KeyCtrlS:
		buf.WriteString("C-s")
	case termbox.KeyCtrlT:
		buf.WriteString("C-t")
	case termbox.KeyCtrlU:
		buf.WriteString("C-u")
	case termbox.KeyCtrlV:
		buf.WriteString("C-v")
	case termbox.KeyCtrlW:
		buf.WriteString("C-w")
	case termbox.KeyCtrlX:
		buf.WriteString("C-x")
	case termbox.KeyCtrlY:
		buf.WriteString("C-y")
	case termbox.KeyCtrlZ:
		buf.WriteString("C-z")
	case termbox.KeyCtrlLsqBracket:
		buf.WriteString("C-[")
	case termbox.KeyCtrlBackslash:
		buf.WriteString("C-\\")
	case termbox.KeyCtrlRsqBracket:
		buf.WriteString("C-]")
	case termbox.KeyCtrl6:
		buf.WriteString("C-6")
	case termbox.KeyCtrlUnderscore:
		buf.WriteString("C-/")
	case termbox.KeySpace:
		buf.WriteString("<space>")
	case termbox.KeyBackspace2:
		buf.WriteString("<backspace2>")
	}

	if ch != 0 {
		buf.WriteRune(ch)
	}
	return buf.String()
}
