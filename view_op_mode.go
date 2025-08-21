package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/termbox"
)

//----------------------------------------------------------------------------
// view op mode
//----------------------------------------------------------------------------

type view_op_mode struct {
	stub_overlay_mode
	gemacs *gemacs
}

const view_names = `1234567890abcdefgijlmnpqrstuwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`

var view_op_mode_name = []byte("view operations mode")

func init_view_op_mode(gemacs *gemacs) view_op_mode {
	GlobalScreen.HideCursor()
	v := view_op_mode{gemacs: gemacs}
	return v
}

func (v view_op_mode) draw() {
	g := v.gemacs
	r := g.uibuf.Rect
	r.Y = r.Height - 1
	r.Height = 1
	fillStyle := MakeStyle(tcell.ColorDefault, tcell.ColorDefault)
	g.uibuf.Fill(r, ' ', fillStyle)
	lp := default_label_params
	lp.Fg = tcell.ColorYellow
	g.uibuf.DrawLabel(r, &lp, string(view_op_mode_name))

	// draw views names
	name := 0
	g.views.traverse(func(leaf *view_tree) {
		if name >= len(view_names) {
			return
		}
		bg := tcell.ColorBlue
		if leaf == g.active {
			bg = tcell.ColorRed
		}
		r := leaf.Rect
		r.Width = 3
		r.Height = 1
		x := r.X + 1
		y := r.Y
		bgStyle := MakeStyle(tcell.ColorDefault, bg)
		g.uibuf.Fill(r, ' ', bgStyle)
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(bg).Bold(true)
		g.uibuf.SetContent(x, y, rune(view_names[name]), nil, style)
		name++
	})

	// draw splitters
	r = g.active.Rect
	var x, y int

	// horizontal ----------------------
	hr := r
	hr.X += (r.Width - 1) / 2
	hr.Width = 1
	hr.Height = 3
	hrStyle := MakeStyle(tcell.ColorWhite, tcell.ColorRed)
	g.uibuf.Fill(hr, '|', hrStyle)

	x = hr.X
	y = hr.Y + 1
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed).Bold(true)
	g.uibuf.SetContent(x, y, 'h', nil, style)

	// vertical ----------------------
	vr := r
	vr.Y += (r.Height - 1) / 2
	vr.Height = 1
	vr.Width = 5
	vrStyle := MakeStyle(tcell.ColorWhite, tcell.ColorRed)
	g.uibuf.Fill(vr, '-', vrStyle)

	x = vr.X + 2
	y = vr.Y
	style = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed).Bold(true)
	g.uibuf.SetContent(x, y, 'v', nil, style)
}

func (v view_op_mode) select_name(ch rune) *view_tree {
	g := v.gemacs
	sel := (*view_tree)(nil)
	name := 0
	g.views.traverse(func(leaf *view_tree) {
		if name >= len(view_names) {
			return
		}
		if rune(view_names[name]) == ch {
			sel = leaf
		}
		name++
	})

	return sel
}

func (v view_op_mode) needs_cursor() bool {
	return true
}

func (v view_op_mode) on_key(ev *termbox.Event) {
	g := v.gemacs
	if ev.Ch != 0 {
		leaf := v.select_name(ev.Ch)
		if leaf != nil {
			g.active.leaf.deactivate()
			g.active = leaf
			g.active.leaf.activate()
			return
		}

		switch ev.Ch {
		case 'h':
			g.split_horizontally()
			return
		case 'v':
			g.split_vertically()
			return
		case 'k':
			g.kill_active_view()
			return
		}
	}

	switch ev.Key {
	case termbox.KeyCtrlN, termbox.KeyArrowDown:
		node := g.active.nearest_vsplit()
		if node != nil {
			node.step_resize(1)
		}
		return
	case termbox.KeyCtrlP, termbox.KeyArrowUp:
		node := g.active.nearest_vsplit()
		if node != nil {
			node.step_resize(-1)
		}
		return
	case termbox.KeyCtrlF, termbox.KeyArrowRight:
		node := g.active.nearest_hsplit()
		if node != nil {
			node.step_resize(1)
		}
		return
	case termbox.KeyCtrlB, termbox.KeyArrowLeft:
		node := g.active.nearest_hsplit()
		if node != nil {
			node.step_resize(-1)
		}
		return
	}

	g.set_overlay_mode(nil)
}