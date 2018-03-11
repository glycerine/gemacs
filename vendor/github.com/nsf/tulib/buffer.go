package tulib

import (
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/termbox"
	"unicode/utf8"

	"github.com/glycerine/verb"
)

var pp = verb.PP

func init() {
	f, err := os.Create("./log.godit.debug")
	if err != nil {
		panic(err)
	}
	verb.OurStdout = f
}

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

type Buffer struct {
	Screen tcell.Screen
	Rect
}

func NewBuffer(w, h int) Buffer {
	// jea: just ignore w, h
	return TermboxBuffer()
}

func TermboxBuffer() Buffer {

	termbox.Init()
	s := termbox.GetScreen()
	w, h := s.Size()

	return Buffer{
		Screen: s,
		Rect:   Rect{0, 0, w, h},
	}
}

// Fills an area which is an intersection between buffer and 'dest' with 'proto'.
func (this *Buffer) Fill(dst Rect, proto termbox.Cell) {
	this.unsafe_fill(this.Rect.Intersection(dst), proto)
}

// Sets a cell at specified position
func (this *Buffer) Set(x, y int, proto termbox.Cell) {
	//pp("top of Set(x=%v, y=%v)", x, y)
	if x < 0 || x >= this.Width {
		return
	}
	if y < 0 || y >= this.Height {
		return
	}
	st := termbox.MakeStyle(proto.Fg, proto.Bg)
	this.Screen.SetContent(x, y, proto.Ch, nil, st)
}

// Resizes the Buffer, buffer contents are invalid after the resize.
func (this *Buffer) Resize(nw, nh int) {
	//pp("top of Resize")
	b := TermboxBuffer()

	this.Screen = b.Screen
	this.Rect = b.Rect
}

func (this *Buffer) Blit(dstr Rect, srcx, srcy int, src *Buffer) {
	//pp("top of Blit")
	srcr := Rect{srcx, srcy, 0, 0}

	// first adjust 'srcr' if 'dstr' has negatives
	if dstr.X < 0 {
		srcr.X -= dstr.X
	}
	if dstr.Y < 0 {
		srcr.Y -= dstr.Y
	}

	// adjust 'dstr' against 'this.Rect', copy 'dstr' size to 'srcr'
	dstr = this.Rect.Intersection(dstr)
	srcr.Width = dstr.Width
	srcr.Height = dstr.Height

	// adjust 'srcr' against 'src.Rect', copy 'srcr' size to 'dstr'
	srcr = src.Rect.Intersection(srcr)
	dstr.Width = srcr.Width
	dstr.Height = srcr.Height

	if dstr.IsEmpty() {
		return
	}

	// blit!
	//srcstride := src.Width
	//dststride := this.Width
	linew := dstr.Width
	srcoff := src.Width*srcr.Y + srcr.X
	dstoff := this.Width*dstr.Y + dstr.X

	cp := func(x, y int) int {
		mainc, combc, style, width := src.Screen.GetContent(x, y)
		this.Screen.SetContent(x, y, mainc, combc, style)
		return width - 1
	}
	if srcoff > dstoff {
		for i := 0; i < dstr.Height; i++ {
			for j := 0; j < linew; j++ {
				j += cp(i, j)
			}
		}
	} else {
		for i := dstr.Height - 1; i >= 0; i-- {
			for j := linew - 1; j >= 0; j-- {
				cp(i, j)
			}
		}
	}
}

// Unsafe part of the fill operation, doesn't check for bounds.
func (this *Buffer) unsafe_fill(dest Rect, proto termbox.Cell) {
	//pp("unsafe fill proto='%#v', dest='%#v'", proto, dest)
	stride := this.Width
	off := this.Width*dest.Y + dest.X
	st := termbox.MakeStyle(proto.Fg, proto.Bg)
	for y := 0; y < dest.Height; y++ {
		for x := 0; x < dest.Width; x++ {
			//this.Cells[off+x] = proto
			this.Screen.SetContent(dest.X+x, dest.Y+y, proto.Ch, nil, st)
		}
		off += stride
	}
}

// draws from left to right, 'off' is the beginning position
// (DrawLabel uses that method)
func (this *Buffer) draw_n_first_runes(off, n int, params *LabelParams, text []byte, destX, destY int) {
	//pp("top of draw_n_first_runes")

	st := termbox.MakeStyle(params.Fg, params.Bg)
	beg := off
	for n > 0 {
		r, size := utf8.DecodeRune(text)

		this.Screen.SetContent(destX+(off-beg), destY, r, nil, st)
		/*
			this.Cells[off] = termbox.Cell{
				Ch: r,
				Fg: params.Fg,
				Bg: params.Bg,
			}
		*/
		text = text[size:]
		off++
		n--
	}
}

// draws from right to left, 'off' is the end position
// (DrawLabel uses that method)
func (this *Buffer) draw_n_last_runes(off, n int, params *LabelParams, text []byte, destX, destY int) {
	//pp("top of draw_n_last_runes")

	st := termbox.MakeStyle(params.Fg, params.Bg)

	i := 0
	for n > 0 {
		r, size := utf8.DecodeLastRune(text)

		this.Screen.SetContent(destX-i, destY, r, nil, st)
		/*
			this.Cells[off] = termbox.Cell{
				Ch: r,
				Fg: params.Fg,
				Bg: params.Bg,
			}
		*/
		text = text[:len(text)-size]
		off--
		n--
		i++
	}
}

type LabelParams struct {
	Fg             termbox.Attribute
	Bg             termbox.Attribute
	Align          Alignment
	Ellipsis       rune
	CenterEllipsis bool
}

var DefaultLabelParams = LabelParams{
	termbox.ColorDefault,
	termbox.ColorDefault,
	AlignLeft,
	'…',
	false,
}

func skip_n_runes(x []byte, n int) []byte {
	if n <= 0 {
		return x
	}

	for n > 0 {
		_, size := utf8.DecodeRune(x)
		x = x[size:]
		n--
	}
	return x
}

func (this *Buffer) DrawLabel(dest Rect, params *LabelParams, text []byte) {
	//pp("top of DrawLabel, text = '%s', param='%#v'", string(text), params)
	st := termbox.MakeStyle(params.Fg, params.Bg)

	if dest.Height != 1 {
		dest.Height = 1
	}

	dest = this.Rect.Intersection(dest)
	if dest.Height == 0 || dest.Width == 0 {
		return
	}

	ellipsis := termbox.Cell{Ch: params.Ellipsis, Fg: params.Fg, Bg: params.Bg}
	off := dest.Y*this.Width + dest.X
	textlen := utf8.RuneCount(text)
	n := textlen
	if n > dest.Width {
		// string doesn't fit in the dest rectangle, draw ellipsis
		n = dest.Width - 1

		// if user asks for ellipsis in the center, alignment doesn't matter
		if params.CenterEllipsis {
			this.Screen.SetContent(dest.X+dest.Width/2, dest.Y, ellipsis.Ch, nil, st)
			//this.Cells[off+dest.Width/2] = ellipsis
		} else {
			switch params.Align {
			case AlignLeft:
				this.Screen.SetContent(dest.X+dest.Width-1, dest.Y, ellipsis.Ch, nil, st)
				//this.Cells[off+dest.Width-1] = ellipsis
			case AlignCenter:
				//this.Cells[off] = ellipsis
				//this.Cells[off+dest.Width-1] = ellipsis
				this.Screen.SetContent(dest.X, dest.Y, ellipsis.Ch, nil, st)
				this.Screen.SetContent(dest.X+dest.Width-1, dest.Y, ellipsis.Ch, nil, st)

				n--
			case AlignRight:
				//this.Cells[off] = ellipsis
				this.Screen.SetContent(dest.X, dest.Y, ellipsis.Ch, nil, st)

			}
		}
	}

	if n <= 0 {
		return
	}

	if params.CenterEllipsis && textlen != n {
		firsthalf := dest.Width / 2
		secondhalf := dest.Width - 1 - firsthalf
		this.draw_n_first_runes(off, firsthalf, params, text, dest.X, dest.Y)
		off += dest.Width - 1
		this.draw_n_last_runes(off, secondhalf, params, text, dest.X+off, dest.Y)
		return
	}

	switch params.Align {
	case AlignLeft:
		this.draw_n_first_runes(off, n, params, text, dest.X, dest.Y)
	case AlignCenter:
		if textlen == n {
			off += (dest.Width - n) / 2
			this.draw_n_first_runes(off, n, params, text, dest.X+off, dest.Y)
		} else {
			off++
			mid := (textlen - n) / 2
			text = skip_n_runes(text, mid)
			this.draw_n_first_runes(off, n, params, text, dest.X+off, dest.Y)
		}
	case AlignRight:
		off += dest.Width - 1
		this.draw_n_last_runes(off, n, params, text, dest.X+off, dest.Y)
	}
}
