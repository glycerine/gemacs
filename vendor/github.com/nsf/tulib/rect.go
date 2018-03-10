package tulib

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Rect) IsValid() bool {
	return r.Width >= 0 && r.Height >= 0
}

func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

func (a Rect) FitsIn(b Rect) bool {
	return a == a.Intersection(b)
}

func (a Rect) Intersection(b Rect) Rect {
	// no intersection cases
	if a.X >= b.X+b.Width || a.Y >= b.Y+b.Height {
		return Rect{0, 0, 0, 0}
	}

	if a.X+a.Width <= b.X || a.Y+a.Height <= b.Y {
		return Rect{0, 0, 0, 0}
	}

	// adjust X or Width
	if a.X+a.Width > b.X+b.Width {
		a.Width = (b.X + b.Width) - a.X
	}

	if a.X < b.X {
		a.Width -= b.X - a.X
		a.X = b.X
	}

	// adjust Y or Height
	if a.Y+a.Height > b.Y+b.Height {
		a.Height = (b.Y + b.Height) - a.Y
	}

	if a.Y < b.Y {
		a.Height -= b.Y - a.Y
		a.Y = b.Y
	}

	return a
}
