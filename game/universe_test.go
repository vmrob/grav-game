package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddBody(t *testing.T) {
	u := NewUniverse(Rect{X: 0, Y: 0, W: 100, H: 100})

	expectedPosition := Point{X: 50, Y: 50}
	expectedMass := 1000.0

	id := u.AddBody(&Body{
		Position: expectedPosition,
		Mass:     expectedMass,
	})

	b := u.GetBody(id)

	assert.Equal(t, b.Position, expectedPosition)
	assert.Equal(t, b.Mass, expectedMass)
}

func TestRemoveBody(t *testing.T) {
	u := NewUniverse(Rect{X: 0, Y: 0, W: 100, H: 100})

	id := u.AddBody(&Body{
		Position: Point{X: 50, Y: 50},
		Mass:     1000.0,
	})

	u.RemoveBody(id)

	b := u.GetBody(id)

	assert.Nil(t, b)
}

func TestInBoundsDecay(t *testing.T) {
	u := NewUniverse(Rect{X: 0, Y: 0, W: 100, H: 100})

	startingMass := playerStartMass * 2.0
	b := Body{
		Position: Point{X: 50, Y: 50},
		Mass:     startingMass,
	}

	u.AddBody(&b)
	u.decayBodies()

	assert.True(t, b.Mass < startingMass)
}

func TestOutOfBoundsDecay(t *testing.T) {
	u := NewUniverse(Rect{X: 0, Y: 0, W: 100, H: 100})

	startingMass := playerStartMass * 2.0
	b := Body{
		Position: Point{X: 150, Y: 150},
		Mass:     startingMass,
	}

	u.AddBody(&b)
	u.decayBodies()

	assert.True(t, b.Mass < startingMass)
}
