package game

import (
	"math"
	"time"
)

type Body struct {
	Position Point
	Mass     float64
}

func (b *Body) Step(d time.Duration) {

}

func (b *Body) Decay(pct float64) {
	qty := b.Mass * pct
	if b.Mass >= minDecayMass {
		b.Mass -= math.Min(qty, b.Mass)
	}
}

func (b *Body) ForceDecay(pct float64) {
	qty := b.Mass * pct
	if b.Mass >= minDecayMass {
		b.Mass -= math.Min(qty, b.Mass)
	} else {
		b.Mass -= math.Min(
			math.Max(minDecayMass, minDecayMassForced),
			b.Mass)
	}
}
