package game

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vector) Sub(v2 Vector) Vector {
	return Vector{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v Vector) Scale(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vector) DistanceSquared(v2 Vector) float64 {
	d := v.Sub(v2)
	return d.X*d.X + d.Y*d.Y
}
