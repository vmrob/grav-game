package game

type Point struct {
	X float64
	Y float64
}

// Returns a vector from this point to another
func (p Point) VectorTo(other Point) Vector {
	return Vector{
		X: other.X - p.X,
		Y: other.Y - p.Y,
	}
}
