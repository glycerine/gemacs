package main

import (
	"github.com/gdamore/tcell/termbox"
)

type region_indent_mode struct {
	stub_overlay_mode
	gemacs *gemacs
}

func init_region_indent_mode(gemacs *gemacs, dir int) region_indent_mode {
	v := gemacs.active.leaf
	r := region_indent_mode{gemacs: gemacs}

	beg, end := v.line_region()
	if dir > 0 {
		v.on_vcommand(vcommand_indent_region, 0)
		end.boffset++
	} else if dir < 0 {
		v.on_vcommand(vcommand_deindent_region, 0)
	}
	v.set_tags(view_tag{
		beg_line:   beg.line_num,
		beg_offset: beg.boffset,
		end_line:   end.line_num,
		end_offset: end.boffset,
		fg:         termbox.ColorDefault,
		bg:         termbox.ColorBlue,
	})
	v.dirty = dirty_everything
	gemacs.set_status("(Type > or < to indent/deindent respectively)")
	return r
}

func (r region_indent_mode) exit() {
	v := r.gemacs.active.leaf
	v.set_tags()
	v.dirty = dirty_everything
}

func (r region_indent_mode) on_key(ev *termbox.Event) {
	g := r.gemacs
	v := g.active.leaf
	beg, end := v.line_region()
	if ev.Mod == 0 {
		switch ev.Ch {
		case '>':
			v.on_vcommand(vcommand_indent_region, 0)
			g.set_status("(Type > or < to indent/deindent respectively)")
			end.boffset++
			goto update_tag
		case '<':
			v.on_vcommand(vcommand_deindent_region, 0)
			g.set_status("(Type > or < to indent/deindent respectively)")
			goto update_tag
		}
	}

	g.set_overlay_mode(nil)
	g.on_key(ev)
	return

update_tag:
	v.set_tags(view_tag{
		beg_line:   beg.line_num,
		beg_offset: beg.boffset,
		end_line:   end.line_num,
		end_offset: end.boffset,
		fg:         termbox.ColorDefault,
		bg:         termbox.ColorBlue,
	})
}
