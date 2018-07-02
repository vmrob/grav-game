package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: test other decay scenarios

func TestDecay(t *testing.T) {
	b := Body{Mass: 100000}
	b.Decay(0.10)
	assert.Equal(t, float64(90000), b.Mass)
}

func TestForceDecay(t *testing.T) {
	b := Body{Mass: 100000}
	b.ForceDecay(0.10)
	assert.Equal(t, float64(90000), b.Mass)
}

func TestCollidesWith(t *testing.T) {
	b1 := Body{
		Position: Point{0, 0},
		Radius:   100,
	}
	b2 := Body{
		Position: Point{0, 0},
		Radius:   1,
	}

	assert.True(t, b1.CollidesWith(&b2))
	assert.True(t, b2.CollidesWith(&b1))

	b3 := Body{
		Position: Point{101, 0},
		Radius:   1,
	}

	assert.False(t, b1.CollidesWith(&b3))
	assert.False(t, b3.CollidesWith(&b1))

	b4 := Body{
		Position: Point{100, 0},
		Radius:   1,
	}

	assert.True(t, b1.CollidesWith(&b4))
	assert.True(t, b4.CollidesWith(&b1))
	assert.True(t, b3.CollidesWith(&b4))
}

func TestDistanceTo(t *testing.T) {
	b := Body{
		Position: Point{0, 0},
		Radius:   1,
	}

	assert.Equal(t, b.DistanceTo(Point{0, 0}), float64(0))
	assert.Equal(t, b.DistanceTo(Point{1, 0}), float64(1))
	assert.Equal(t, b.DistanceTo(Point{2, 0}), float64(2))
	assert.Equal(t, b.DistanceTo(Point{3, 4}), float64(5))
}

func TestMergeWith(t *testing.T) {
	b1 := Body{
		Velocity: Vector{10, 0},
		Position: Point{10, 0},
		Mass:     9,
	}

	b2 := Body{
		Velocity: Vector{0, 0},
		Position: Point{0, 0},
		Mass:     1,
	}

	b1.MergeWith(&b2)

	assert.Equal(t, b1.Mass, float64(10))
	assert.Equal(t, b1.Position, Point{9, 0})
	assert.Equal(t, b1.Velocity, Vector{9, 0})

	assert.Equal(t, b2.Mass, float64(0))
}

func TestGravitationalForceTo(t *testing.T) {
	b1 := Body{
		Mass:     10,
		Position: Point{0, 0},
	}

	b2 := Body{
		Mass:     100,
		Position: Point{0, 1},
	}

	assert.Equal(t, Vector{0, 100000}, b1.GravitationalForceTo(&b2))
	assert.Equal(t, Vector{0, -100000}, b2.GravitationalForceTo(&b1))

	b3 := Body{
		Static:   true,
		Mass:     100,
		Position: Point{0, 1},
	}

	assert.Equal(t, Vector{0, 0}, b1.GravitationalForceTo(&b3))
	assert.Equal(t, Vector{0, 0}, b3.GravitationalForceTo(&b1))
}
