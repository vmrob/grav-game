package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorTo(t *testing.T) {
	assert.Equal(t, Point{0, 1}.VectorTo(Point{1, 2}), Vector{1, 1})
	assert.Equal(t, Point{-1, -10}.VectorTo(Point{1, -2}), Vector{2, 8})
}
