package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestGame *Game
var BenchmarkResult interface{}

func init() {
	TestGame = New()
	for i := 1; i <= 1000; i++ {
		TestGame.AddBody(&Body{
			Position: Vector{
				X: float64(i * 100), Y: float64(i * 200),
			},
			Mass: float64(i * 1000),
		})
	}
}

func TestGravitationalForce(t *testing.T) {
	f := TestGame.GravitationalForce(Vector{
		X: 0, Y: 0,
	}, 2000)
	assert.Equal(t, Vector{
		X: 0.0002669632, Y: 0.0005339264,
	}, f)
}

func BenchmarkGravitationalForce(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BenchmarkResult = TestGame.GravitationalForce(Vector{
			X: float64(n), Y: float64(n),
		}, 2000)
	}
}
