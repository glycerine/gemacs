package main

import (
	"github.com/gdamore/tcell/termbox"
	"github.com/nsf/tulib"
	"strings"
	"unicode/utf8"
)

//----------------------------------------------------------------------------
// line edit mode
//----------------------------------------------------------------------------

type line_edit_mode struct {
	stub_overlay_mode
	line_edit_mode_params
	gemacs   *gemacs
	linebuf  *buffer
	lineview *view
	prompt   []byte
	prompt_w int
}

type line_edit_mode_params struct {
	on_apply        func(buffer *buffer)
	on_cancel       func()
	ac_decide       ac_decide_func
	prompt          string
	initial_content string
	init_autocompl  bool
}

func (l *line_edit_mode) exit() {
	if l.on_cancel != nil {
		l.on_cancel()
	}
}

func (l *line_edit_mode) on_key(ev *termbox.Event) {
	pp("line_edit_mode on_key, Ch='%v', ev.Key = '%#v'", string(ev.Ch), ev)
	switch ev.Key {
	case termbox.KeyEnter, termbox.KeyCtrlJ:
		pp("enter")
		if l.lineview.ac != nil {
			pp("l.lineview.ac != nil")
			l.lineview.on_key(ev)
			if !l.init_autocompl {
				break
			}
		}

		// reset overlay mode earlier so that 'on_apply' can
		// override it
		l.gemacs.set_overlay_mode(nil)
		if l.on_apply != nil {
			pp("l.on_apply != nil")
			l.on_apply(l.linebuf)
		}
	case termbox.KeyTab:
		pp("KeyTab.") //  l.lineview='%#v'", l.lineview)
		if l.lineview.ac == nil {
			l.lineview.on_vcommand(vcommand_autocompl_init, 0)
		} else {
			// jea, 1st time above, vs 2nd time below:
			l.lineview.on_vcommand(vcommand_autocompl_tab, 0)
		}
	default:
		pp("default")
		l.lineview.on_key(ev)
	}
}

func (l *line_edit_mode) resize(ev *termbox.Event) {
	w, h := ev.Width-l.prompt_w-1, 1
	if w < 1 || ev.Height < 1 {
		return
	}
	l.lineview.resize(w, h)
}

func (l *line_edit_mode) draw() {
	ui := l.gemacs.uibuf
	view := l.lineview

	// update label
	prompt_r := tulib.Rect{
		0, ui.Height - 1,
		l.prompt_w + 1, 1,
	}
	ui.Fill(prompt_r, termbox.Cell{
		Fg: termbox.ColorDefault,
		Bg: termbox.ColorDefault,
		Ch: ' ',
	})
	lp := default_label_params
	lp.Fg = termbox.ColorCyan
	ui.DrawLabel(prompt_r, &lp, l.prompt)

	// update line view
	view.resize(ui.Width-l.prompt_w-1, 1)
	view.draw()
	line_r := tulib.Rect{
		l.prompt_w + 1, ui.Height - 1,
		view.uibuf.Width, view.uibuf.Height,
	}
	ui.Blit(line_r, 0, 0, view.uibuf) // unnamed being written here.
	if view.ac == nil {
		return
	}

	// draw autocompletion
	proposals := view.ac.actual_proposals()
	if len(proposals) > 0 {
		cx, cy := view.cursor_position_for(view.ac.origin)
		view.ac.draw_onto(ui, line_r.X+cx, line_r.Y+cy)
	}
}

func (l *line_edit_mode) needs_cursor() bool {
	return true
}

func (l *line_edit_mode) cursor_position() (int, int) {
	y := l.gemacs.uibuf.Height - 1
	x := l.prompt_w + 1
	lx, ly := l.lineview.cursor_position()
	return x + lx, y + ly
}

func init_line_edit_mode(gemacs *gemacs, p line_edit_mode_params) *line_edit_mode {
	l := new(line_edit_mode)
	l.gemacs = gemacs
	l.line_edit_mode_params = p

	l.linebuf, _ = new_buffer(strings.NewReader(p.initial_content))
	l.lineview = new_view(gemacs.view_context(), l.linebuf, gemacs)
	l.lineview.oneline = true          // enable one line mode
	l.lineview.ac_decide = p.ac_decide // override ac_decide function
	l.prompt = []byte(p.prompt)
	l.prompt_w = utf8.RuneCount(l.prompt)
	l.lineview.resize(l.gemacs.uibuf.Width-l.prompt_w-1, 1)
	l.lineview.on_vcommand(vcommand_move_cursor_end_of_line, 0)
	if l.init_autocompl {
		l.lineview.on_vcommand(vcommand_autocompl_init, 0)
	}
	return l
}
