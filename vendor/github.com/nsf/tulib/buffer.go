package tulib

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/termbox"
	//"github.com/nsf/termbox-go"
	"unicode/utf8"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

type Buffer struct {
	Screen tcell.Screen
	//Cells  []termbox.Cell
	Rect
}

func NewBuffer(w, h int) Buffer {
	// jea: just ignore w, h

	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	e = s.Init()
	if e != nil {
		panic(e)
	}
	w2, h2 := s.Size()
	/*	if w2 != w {
			panic(fmt.Sprintf("w2 != w. w2=%v, w=%v", w2, w))
		}
		if h2 != h {
			panic(fmt.Sprintf("h2 != h"))
		}
	*/
	return Buffer{
		Screen: s,
		Rect:   Rect{0, 0, w2, h2},
	}
}

func TermboxBuffer() Buffer {

	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	e = s.Init()
	if e != nil {
		panic(e)
	}
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
	if x < 0 || x >= this.Width {
		return
	}
	if y < 0 || y >= this.Height {
		return
	}
	this.Screen.SetContent(x, y, proto.Ch, nil, 0)
}

/*
// Gets a pointer to the cell at specified position or nil if it's out
//  of range.
func (this *Buffer) Get(x, y int) *termbox.Cell {
	if x < 0 || x >= this.Width {
		return nil
	}
	if y < 0 || y >= this.Height {
		return nil
	}
	off := this.Width*y + x
	return &this.Cells[off]

	mainc, combc, style, width := this.Screen.GetContent(x, y)
}
*/

// Resizes the Buffer, buffer contents are invalid after the resize.
func (this *Buffer) Resize(nw, nh int) {

	b := TermboxBuffer()
	this.Screen = b.Screen
	this.Rect = b.Rect
	/*
		this.Width = nw
		this.Height = nh

		this.Screen.Resize()
		nsize := nw * nh
		if nsize <= cap(this.Cells) {
			this.Cells = this.Cells[:nsize]
		} else {
			this.Cells = make([]termbox.Cell, nsize)
		}
	*/
}

func (this *Buffer) Blit(dstr Rect, srcx, srcy int, src *Buffer) {
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
	stride := this.Width
	off := this.Width*dest.Y + dest.X
	for y := 0; y < dest.Height; y++ {
		for x := 0; x < dest.Width; x++ {
			//this.Cells[off+x] = proto
			this.Screen.SetContent(x, y, proto.Ch, nil, 0)
		}
		off += stride
	}
}

// draws from left to right, 'off' is the beginning position
// (DrawLabel uses that method)
func (this *Buffer) draw_n_first_runes(off, n int, params *LabelParams, text []byte, destX, destY int) {

	beg := off
	for n > 0 {
		r, size := utf8.DecodeRune(text)

		this.Screen.SetContent(destX+(off-beg), destY, r, nil, 0)
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

	i := 0
	for n > 0 {
		r, size := utf8.DecodeLastRune(text)

		this.Screen.SetContent(destX-i, destY, r, nil, 0)
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
			this.Screen.SetContent(dest.X+dest.Width/2, dest.Y, ellipsis.Ch, nil, 0)
			//this.Cells[off+dest.Width/2] = ellipsis
		} else {
			switch params.Align {
			case AlignLeft:
				this.Screen.SetContent(dest.X+dest.Width-1, dest.Y, ellipsis.Ch, nil, 0)
				//this.Cells[off+dest.Width-1] = ellipsis
			case AlignCenter:
				//this.Cells[off] = ellipsis
				//this.Cells[off+dest.Width-1] = ellipsis
				this.Screen.SetContent(dest.X, dest.Y, ellipsis.Ch, nil, 0)
				this.Screen.SetContent(dest.X+dest.Width-1, dest.Y, ellipsis.Ch, nil, 0)

				n--
			case AlignRight:
				//this.Cells[off] = ellipsis
				this.Screen.SetContent(dest.X, dest.Y, ellipsis.Ch, nil, 0)

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
