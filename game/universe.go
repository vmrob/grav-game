package game

import (
	"math/rand"
	"sort"
	"time"
)

type Universe struct {
	bounds Rect
	bodies map[BodyId]*Body
	nextId BodyId
	events chan func()
}

func NewUniverse(bounds Rect) *Universe {
	return &Universe{
		bounds: bounds,
		bodies: make(map[BodyId]*Body),
		events: make(chan func(), 1000), // TODO: this isn't too scalable
	}
}

func (u *Universe) Bounds() Rect {
	return u.bounds
}

func (u *Universe) Bodies() map[BodyId]*Body {
	return u.bodies
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

func (u *Universe) consumeAvailableEvents() {
	for {
		select {
		case f := <-u.events:
			f()
		default:
			return
		}
	}
}

func (u *Universe) Step(d time.Duration) {
	u.consumeAvailableEvents()
	u.decayBodies()
	u.checkCollisions()
	u.applyForces()

	rankings := make([]*Body, 0, len(u.bodies))
	for _, b := range u.bodies {
		b.Step(d)
		rankings = append(rankings, b)
	}
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Mass > rankings[j].Mass
	})

	if len(rankings) >= 100 {
		majorThreshold := rankings[len(rankings)/100].Mass * 3
		for _, b := range rankings[:len(rankings)/100] {
			if b.MajorName == "" && b.Mass >= majorThreshold {
				b.MajorName = u.NewMajorName()
			}
		}
	}

	minor := len(rankings) / 100
	for _, b := range rankings[:minor] {
		if b.MinorName == "" {
			b.MinorName = u.NewMinorName()
		}
	}
}

func (u *Universe) decayBodies() {
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
					other.MergeWith(body)
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

func (u *Universe) AddEvent(f func()) {
	u.events <- f
}

var phoneticAlphabet = []string{
	"Alfa", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf", "Hotel", "India", "Juliett",
	"Kilo", "Lima", "Mike", "November", "Oscar", "Papa", "Quebec", "Romeo", "Sierra", "Tango",
	"Uniform", "Victor", "Whiskey", "X-Ray", "Yankee", "Zulu",
}

var alphaNumeric = "QWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func (u *Universe) NewMinorName() string {
	for {
		name := phoneticAlphabet[rand.Intn(len(phoneticAlphabet))] + " "
		for i := 0; i < 5; i++ {
			name += string(rune(alphaNumeric[rand.Intn(len(alphaNumeric))]))
		}
		inUse := false
		for _, body := range u.bodies {
			if body.MinorName == name {
				inUse = true
				break
			}
		}
		if !inUse {
			return name
		}
	}
}

var majorNames = []string{
	"Zeus", "Hera", "Poseidon", "Demeter", "Ares", "Athena", "Apollo", "Artemis", "Hephaestus",
	"Aphrodite", "Hermes", "Dionysus", "Hades", "Hypnos", "Nike", "Janus", "Nemesis", "Iris",
	"Hecate", "Tyche",
}

func (u *Universe) NewMajorName() string {
	indices := rand.Perm(len(majorNames))
	for _, i := range indices {
		name := majorNames[i]
		inUse := false
		for _, body := range u.bodies {
			if body.MajorName == name {
				inUse = true
				break
			}
		}
		if !inUse {
			return name
		}
	}
	return ""
}
