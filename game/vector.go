package game

import (
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Add(other Vector) Vector {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v Vector) Sub(other Vector) Vector {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v Vector) Scale(s float64) Vector {
	v.X *= s
	v.Y *= s
	return v
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector) WithMagnitude(m float64) Vector {
	current := v.Magnitude()
	if current == 0.0 {
		panic("cannot use WithMagnitude on a zero Vector")
	}
	return v.Scale(m / current)
}
