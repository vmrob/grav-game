package game

import (
	"math"
)

func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

func distanceSquared(p1, p2 Point) float64 {
	return math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2)
}
