package main

import (
	"github.com/gdamore/tcell/termbox"
)

type redo_mode struct {
	stub_overlay_mode
	gemacs *gemacs
}

func init_redo_mode(gemacs *gemacs) redo_mode {
	r := redo_mode{gemacs: gemacs}
	return r
}

func (r redo_mode) on_key(ev *termbox.Event) {
	g := r.gemacs
	v := g.active.leaf
	if ev.Mod == 0 && ev.Key == termbox.KeyCtrlSlash {
		v.on_vcommand(vcommand_redo, 0)
		return
	}

	g.set_overlay_mode(nil)
	g.on_key(ev)
}
