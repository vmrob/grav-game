package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	r := Rect{
		X: 0,
		Y: 0,
		W: 10,
		H: 10,
	}

	assert.True(t, r.Contains(Point{0, 0}))
	assert.True(t, r.Contains(Point{10, 10}))
	assert.True(t, r.Contains(Point{0, 5}))
	assert.True(t, r.Contains(Point{5, 0}))

	assert.False(t, r.Contains(Point{-1, 0}))
	assert.False(t, r.Contains(Point{11, 0}))
	assert.False(t, r.Contains(Point{0, -5}))
	assert.False(t, r.Contains(Point{-5, 0}))
}
