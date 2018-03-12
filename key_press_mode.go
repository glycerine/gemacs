package main

import (
	"github.com/gdamore/tcell/termbox"
)

type key_press_mode struct {
	stub_overlay_mode
	gemacs  *gemacs
	actions map[rune]func()
	def     rune
	prompt  string
}

func init_key_press_mode(gemacs *gemacs, actions map[rune]func(), def rune, prompt string) *key_press_mode {
	k := new(key_press_mode)
	k.gemacs = gemacs
	k.actions = actions
	k.def = def
	k.prompt = prompt
	k.gemacs.set_status(prompt)
	return k
}

func (k *key_press_mode) on_key(ev *termbox.Event) {
	if ev.Mod != 0 {
		return
	}

	ch := ev.Ch
	if ev.Key == termbox.KeyEnter || ev.Key == termbox.KeyCtrlJ {
		ch = k.def
	}

	action, ok := k.actions[ch]
	if ok {
		action()
		k.gemacs.set_overlay_mode(nil)
	} else {
		k.gemacs.set_status(k.prompt)
	}
}
