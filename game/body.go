package game

import (
	"math"
	"time"
)

type Body struct {
	Position           Point
	Mass               float64
	Radius             float64
	Static             bool
	Velocity           Vector
	GravitationalForce Vector
	AdditionalForce    Vector
	NetForce           Vector
}

func (b *Body) Step(d time.Duration) {
	b.updateRadius()
	b.updateNetForce(d)
	b.updateVelocity(d)
	b.updatePosition(d)
}

func (b *Body) Decay(pct float64) {
	qty := b.Mass * pct
	if b.Mass >= minDecayMass {
		b.Mass -= math.Min(qty, b.Mass)
	}
}

func (b *Body) ForceDecay(pct float64) {
	// TODO: this should be cleaned up
	qty := b.Mass * pct
	if b.Mass >= minDecayMass {
		b.Mass -= math.Min(qty, b.Mass)
	} else {
		b.Mass -= math.Min(
			math.Max(minDecayMass, minDecayMassForced),
			b.Mass)
	}
}

func (b *Body) CollidesWith(other *Body) bool {
	return b.DistanceTo(other.Position) < b.Radius+other.Radius
}

func (b *Body) DistanceTo(p Point) float64 {
	return math.Sqrt(math.Pow(b.Position.X-p.X, 2) + math.Pow(b.Position.Y-p.Y, 2))
}

func (b *Body) MergeWith(other *Body) {
	// don't conserve velocity against static bodies
	if !b.Static && !other.Static {
		b.Velocity = Vector{
			X: (b.Velocity.X*b.Mass + other.Velocity.X*other.Mass) / (b.Mass + other.Mass),
			Y: (b.Velocity.Y*b.Mass + other.Velocity.Y*other.Mass) / (b.Mass + other.Mass),
		}
	}
	// conservation of position?
	b.Position.X = (b.Position.X*b.Mass + other.Position.X*other.Mass) / (b.Mass + other.Mass)
	b.Position.Y = (b.Position.Y*b.Mass + other.Position.Y*other.Mass) / (b.Mass + other.Mass)
	b.Mass += other.Mass
	other.Mass = 0
}

func (b *Body) GravitationalForceTo(other *Body) Vector {
	if b.Static || other.Static {
		return Vector{}
	}
	force := gravitationalConstant * b.Mass * other.Mass / math.Pow(b.DistanceTo(other.Position), 2)
	return b.Position.VectorTo(other.Position).WithMagnitude(force)
}

func (b *Body) updateRadius() {
	b.Radius = math.Cbrt((b.Mass * 3) / (4 * math.Pi))
}

func (b *Body) updateNetForce(d time.Duration) {
	b.NetForce = b.GravitationalForce.Add(b.AdditionalForce)
}

func (b *Body) updateVelocity(d time.Duration) {
	if b.Static {
		return
	}
	if b.Mass == 0 {
		panic("updateVelocity called on zero-mass body")
	}
	b.Velocity.X += b.NetForce.X / b.Mass * d.Seconds()
	b.Velocity.Y += b.NetForce.Y / b.Mass * d.Seconds()
}

func (b *Body) updatePosition(d time.Duration) {
	if b.Static {
		return
	}
	b.Position.X += b.Velocity.X * d.Seconds()
	b.Position.Y += b.Velocity.Y * d.Seconds()
}
