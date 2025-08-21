package main

import (
	"github.com/gdamore/tcell/v2"
	"strconv"
)

//----------------------------------------------------------------------------
// extended mode
//----------------------------------------------------------------------------

type extended_mode struct {
	stub_overlay_mode
	gemacs *gemacs
}

func init_extended_mode(gemacs *gemacs) extended_mode {
	e := extended_mode{gemacs: gemacs}
	e.gemacs.set_status("C-x")
	return e
}

func (e extended_mode) on_key(ev *tcell.EventKey) {
	//pp("extended_mode on_key, Ch='%v', ev='%#v'", string(ev.Rune()), ev)
	g := e.gemacs
	v := g.active.leaf
	b := v.buf

	switch ev.Key() {
	case tcell.KeyCtrlC:
		if g.has_unsaved_buffers() {
			g.set_overlay_mode(init_key_press_mode(
				g,
				map[rune]func(){
					'y': func() {
							g.quitflag = true
						},
					'n': func() {},
				},
				0,
				"Modified buffers exist; exit anyway? (y or n)",
			))
			return
		} else {
			g.quitflag = true
		}
	case tcell.KeyCtrlX:
		v.on_vcommand(vcommand_swap_cursor_and_mark, 0)
	case tcell.KeyCtrlV:
		g.set_overlay_mode(init_view_op_mode(g))
		return
	case tcell.KeyCtrlW:
		g.set_overlay_mode(init_line_edit_mode(g,
			g.save_as_buffer_lemp(true)))
		return
	case tcell.KeyCtrlA:
		v.on_vcommand(vcommand_autocompl_init, 0)
	case tcell.KeyCtrlU:
		v.on_vcommand(vcommand_region_to_upper, 0)
	case tcell.KeyCtrlL:
		v.on_vcommand(vcommand_region_to_lower, 0)
	case tcell.KeyCtrlF:
		g.set_overlay_mode(init_line_edit_mode(g, g.open_buffer_lemp()))
		return
	case tcell.KeyCtrlS:
		g.save_active_buffer(false)
		return
	case tcell.KeyCtrlSlash:
		g.active.leaf.on_vcommand(vcommand_redo, 0)
		g.set_overlay_mode(init_redo_mode(g))
		return
	case tcell.KeyCtrlR:
		if !v.buf.is_mark_set() {
			v.ctx.set_status("The mark is not set now, so there is no region")
			break
		}
		g.set_overlay_mode(init_line_edit_mode(g, g.search_and_replace_lemp1()))
		return
	default:
		switch ev.Rune() {
		case '0':
			g.kill_active_view()
		case '1':
			g.kill_all_views_but_active()
		case '2':
			g.split_vertically()
		case '3':			g.split_horizontally()
		case 'o':
			next := g.active.nextInCycle()
			if next != nil && next.leaf != nil {
				g.active.leaf.deactivate()
				g.active = next
				g.active.leaf.activate()
			}
		case 'b':
			g.set_overlay_mode(init_line_edit_mode(g, g.switch_buffer_lemp()))
			return
		case '(': 
			g.set_status("Defining keyboard macro...")
			g.recording = true
			g.keymacros = g.keymacros[:0]
		case ')':
			g.stop_recording()
		case 'e':
			g.stop_recording()
			if len(g.keymacros) > 0 {
				g.set_overlay_mode(init_macro_repeat_mode(g))
				return
			}
		case '>':
			g.set_overlay_mode(init_region_indent_mode(g, 1))
			return
		case '<':
			g.set_overlay_mode(init_region_indent_mode(g, -1))
			return
		case 'k':
			if !b.synced_with_disk() {
				g.set_overlay_mode(init_key_press_mode(
						g,
						map[rune]func(){
								'y': func() {
										g.kill_buffer(b)
								},
								'n': func() {},
							},
							0,
							"Buffer "+b.name+" modified; kill anyway? (y or n)",
						))
				return
			} else {
				g.kill_buffer(b)
			}
		case 'S':
			if ev.Mod()&tcell.ModAlt != 0 {
				g.set_overlay_mode(init_line_edit_mode(g,
						g.save_as_buffer_lemp(true)))
				return
			}
			g.save_active_buffer(true)
			return
		case 's':
			if ev.Mod()&tcell.ModAlt != 0 {
				g.set_overlay_mode(init_line_edit_mode(g,
						g.save_as_buffer_lemp(false)))
				return
			}
		case '=':
			var r rune
			if v.cursor.eol() {
				r = '\n'
			} else {
				r, _ = v.cursor.rune_under()
			}
			cursor_ex := make_cursor_location_ex(v.cursor)
			g.set_status("Char: %s (dec: %d, oct: %s, hex: %s), Cursor offset: %d bytes",
				strconv.QuoteRune(r), r,
				strconv.FormatInt(int64(r), 8),
				strconv.FormatInt(int64(r), 16),
				cursor_ex.abs_boffset)
		case '!':
			g.set_overlay_mode(init_line_edit_mode(g, g.filter_region_lemp()))
			return
		case 'h':
			// Toggle syntax highlighting
			if g.syntax_highlighter != nil {
				enabled := !g.syntax_highlighter.IsEnabled()
				g.syntax_highlighter.SetEnabled(enabled)
				if enabled {
					g.set_status("Syntax highlighting enabled")
				} else {
					g.set_status("Syntax highlighting disabled")
				}
				// Refresh all views
				g.views.traverse(func(vt *view_tree) {
					if vt.leaf != nil {
						vt.leaf.dirty = dirty_everything
					}
				})
			} else {
				g.set_status("Syntax highlighter not available")
			}
		case 't':
			// Set tab size
			g.set_overlay_mode(init_line_edit_mode(g, g.set_tab_size_lemp()))
			return
		default:
			goto undefined
		}
	}

	g.set_overlay_mode(nil)
	return
undefined:
	g.set_status("C-x %s is undefined", KeyToString(ev))
	g.set_overlay_mode(nil)
}
