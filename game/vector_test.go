package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert.Equal(t, Vector{0, 1}.Add(Vector{1, 2}), Vector{1, 3})
	assert.Equal(t, Vector{-1, -10}.Add(Vector{1, -2}), Vector{0, -12})
}

func TestSub(t *testing.T) {
	assert.Equal(t, Vector{0, 1}.Sub(Vector{1, 2}), Vector{-1, -1})
	assert.Equal(t, Vector{15, -10}.Sub(Vector{1, 2}), Vector{14, -12})
}

func TestScale(t *testing.T) {
	assert.Equal(t, Vector{1, 10}.Scale(2), Vector{2, 20})
	assert.Equal(t, Vector{10, 10}.Scale(2), Vector{20, 20})
}

func TestMagnitude(t *testing.T) {
	assert.Equal(t, Vector{0, 10}.Magnitude(), float64(10))
	assert.Equal(t, Vector{10, 0}.Magnitude(), float64(10))
	assert.Equal(t, Vector{3, 4}.Magnitude(), float64(5))
}

func TestMagnitudeSquared(t *testing.T) {
	assert.Equal(t, Vector{0, 10}.MagnitudeSquared(), float64(100))
	assert.Equal(t, Vector{10, 0}.MagnitudeSquared(), float64(100))
	assert.Equal(t, Vector{3, 4}.MagnitudeSquared(), float64(25))
}

func TestWithMagnitude(t *testing.T) {
	assert.Equal(t, Vector{0, 10}.WithMagnitude(5), Vector{0, 5})
	assert.Equal(t, Vector{10, 0}.WithMagnitude(5), Vector{5, 0})
	assert.Equal(t, Vector{3, 4}.WithMagnitude(10), Vector{6, 8})
}
