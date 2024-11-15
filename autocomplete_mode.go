package main

import (
	"github.com/glycerine/tcell_old_hacked_up/termbox"
)

type autocomplete_mode struct {
	stub_overlay_mode
	gemacs     *gemacs
	origin     cursor_location
	proposals  []ac_proposal
	prefix_len int
	current    int
}

func init_autocomplete_mode(gemacs *gemacs) *autocomplete_mode {
	view := gemacs.active.leaf

	a := new(autocomplete_mode)
	a.gemacs = gemacs
	a.origin = view.cursor
	a.proposals, a.prefix_len = local_ac(view)
	a.current = -1
	a.substitute_next()
	return a
}

func (a *autocomplete_mode) substitute_next() {
	view := a.gemacs.active.leaf
	if a.current != -1 {
		// undo previous substitution
		view.undo()
		a.gemacs.set_status("") // hide undo status message
	}

	a.current++
	if a.current >= len(a.proposals) {
		a.current = -1
		a.gemacs.set_status("No further expansions found")
		return
	}

	// create a new one
	c := view.cursor
	view.finalize_action_group()
	if a.prefix_len != 0 {
		c.move_one_word_backward()
		wlen := a.origin.boffset - c.boffset
		view.action_delete(c, wlen)
	}
	newword := clone_byte_slice(a.proposals[a.current].content)
	view.action_insert(c, newword)
	view.last_vcommand = vcommand_none
	view.dirty = dirty_everything
	c.boffset += len(newword)
	view.move_cursor_to(c)
	view.finalize_action_group()
}

func (a *autocomplete_mode) on_key(ev *termbox.Event) {
	g := a.gemacs
	if ev.Mod&termbox.ModAlt != 0 && ev.Ch == '/' {
		a.substitute_next()
		return
	}

	g.set_overlay_mode(nil)
	g.on_key(ev)
}
