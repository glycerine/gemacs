package main

import (
	"github.com/gdamore/tcell/termbox"
)

type macro_repeat_mode struct {
	stub_overlay_mode
	gemacs *gemacs
}

func init_macro_repeat_mode(gemacs *gemacs) macro_repeat_mode {
	m := macro_repeat_mode{gemacs: gemacs}
	gemacs.set_overlay_mode(nil)
	m.gemacs.replay_macro()
	m.gemacs.set_status("(Type e to repeat macro)")
	return m
}

func (m macro_repeat_mode) on_key(ev *termbox.Event) {
	g := m.gemacs
	if ev.Mod == 0 && ev.Ch == 'e' {
		g.set_overlay_mode(nil)
		g.replay_macro()
		g.set_overlay_mode(m)
		g.set_status("(Type e to repeat macro)")
		return
	}

	g.set_overlay_mode(nil)
	g.on_key(ev)
}
