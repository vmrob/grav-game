package game

type Rect struct {
	X int
	Y int
	W int
	H int
}

func (r *Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X <= r.X+r.W &&
		p.Y >= r.Y && p.Y <= r.Y+r.H
}
