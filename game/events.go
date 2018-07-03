package game

import (
	"math"
	"math/rand"
)

func randomPointInRect(r Rect) Point {
	return Point{
		X: rand.Float64()*r.W + r.X,
		Y: rand.Float64()*r.H + r.Y,
	}
}

func orbitVector(p Point, b *Body) Vector {
	v := math.Sqrt(gravitationalConstant * b.Mass / distance(p, b.Position))
	return Vector{p.Y, -p.X}.WithMagnitude(v).Add(b.Velocity)
}

func ThreatSpawnEvent(u *Universe) func() {
	return func() {
		var largestBody *Body
		for _, b := range u.Bodies() {
			if largestBody == nil || b.Mass > largestBody.Mass {
				largestBody = b
			}
		}

		p := randomPointInRect(u.Bounds())
		m := float64(PlayerStartMass)
		v := Vector{0, 0}

		if largestBody != nil {
			m = rand.Float64() * math.Min(largestBody.Mass*2, PlayerStartMass*10)
			v = orbitVector(p, largestBody)
		}

		b := Body{
			Position: p,
			Mass:     m,
			Velocity: v,
		}
		u.AddBody(&b)
	}
}

func FoodSpawnEvent(u *Universe) func() {
	return func() {
		var largestBody *Body
		for _, b := range u.Bodies() {
			if largestBody == nil || b.Mass > largestBody.Mass {
				largestBody = b
			}
		}

		p := randomPointInRect(u.Bounds())
		m := rand.Float64()*PlayerStartMass*0.4 + PlayerStartMass*0.1
		v := Vector{0, 0}

		if largestBody != nil {
			v = orbitVector(p, largestBody)
		}

		b := Body{
			Position: p,
			Mass:     m,
			Velocity: v,
		}
		u.AddBody(&b)
	}
}
