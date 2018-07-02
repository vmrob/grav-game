package game

import (
	"time"
)

type BodyId = int

type Universe struct {
	bounds Rect
	bodies map[BodyId]*Body
	nextId BodyId
}

func NewUniverse(bounds Rect) *Universe {
	return &Universe{
		bounds: bounds,
		bodies: make(map[BodyId]*Body),
	}
}

func (u *Universe) AddBody(b *Body) BodyId {
	id := u.nextId
	u.nextId++
	u.bodies[id] = b
	return id
}

func (u *Universe) GetBody(id BodyId) *Body {
	return u.bodies[id]
}

func (u *Universe) RemoveBody(id BodyId) {
	delete(u.bodies, id)
}

func (u *Universe) Step(d time.Duration) {
	u.DecayBodies()
	u.checkCollisions()
	u.applyForces()
	for _, b := range u.bodies {
		b.Step(d)
	}
}

func (u *Universe) DecayBodies() {
	for _, b := range u.bodies {
		if !u.bounds.Contains(b.Position) {
			b.ForceDecay(outOfBoundsDecayPerStep)
		} else {
			b.Decay(decayPerStep)
		}
	}

	for id, b := range u.bodies {
		if b.Mass == 0 {
			u.RemoveBody(id)
		}
	}
}

func (u *Universe) checkCollisions() {
	for id, body := range u.bodies {
		for otherId, other := range u.bodies {
			if id <= otherId {
				continue
			}
			if body.CollidesWith(other) {
				if body.Mass > other.Mass {
					body.MergeWith(other)
					u.RemoveBody(otherId)
				} else {
					other.MergeWith(other)
					u.RemoveBody(id)
				}
				// TODO: this could be a lot better
				u.checkCollisions()
				return
			}
		}
	}
}

func (u *Universe) applyForces() {
	for id, body := range u.bodies {
		netForces := make([]Vector, 0, len(u.bodies))
		for otherId, other := range u.bodies {
			if id == otherId {
				continue
			}
			netForces = append(netForces, body.GravitationalForceTo(other))
		}
		body.GravitationalForce = Vector{}
		for _, v := range netForces {
			body.GravitationalForce.X += v.X
			body.GravitationalForce.Y += v.Y
		}
	}
}
