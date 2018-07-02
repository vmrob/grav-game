package game

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func (r *Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X <= r.X+r.W &&
		p.Y >= r.Y && p.Y <= r.Y+r.H
}
