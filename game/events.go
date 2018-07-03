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

func ThreatSpawnEvent(u *Universe) func() {
	return func() {
		maxMass := 0.0
		for _, b := range u.Bodies() {
			maxMass = math.Max(maxMass, b.Mass)
		}
		maxMass = math.Min(maxMass*2, PlayerStartMass*10)

		b := Body{
			Position: randomPointInRect(u.Bounds()),
			Mass:     rand.Float64() * maxMass,
			Velocity: Vector{
				X: float64(rand.Int31n(200) - 100),
				Y: float64(rand.Int31n(200) - 100),
			},
		}
		u.AddBody(&b)
	}
}

func FoodSpawnEvent(u *Universe) func() {
	return func() {
		b := Body{
			Position: randomPointInRect(u.Bounds()),
			Mass:     rand.Float64()*PlayerStartMass*0.4 + PlayerStartMass*0.1,
			Velocity: Vector{
				X: float64(rand.Int31n(1000) - 500),
				Y: float64(rand.Int31n(1000) - 500),
			},
		}
		u.AddBody(&b)
	}
}
