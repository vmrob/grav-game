package game

const PlayerStartMass = 10000

const decayPerStep = 0.0001
const outOfBoundsDecayPerStep = 0.001
const minDecayMass = PlayerStartMass
const minDecayMassForced = 500
const gravitationalConstant = 100
const thrustBaseMagnitude = 5000000

var (
	North = Vector{0, 1}
	South = Vector{0, -1}
	East  = Vector{1, 0}
	West  = Vector{-1, 0}

	NorthEast = Vector{1, 1}
	NorthWest = Vector{-1, 1}
	SouthEast = Vector{1, -1}
	SouthWest = Vector{-1, -1}
)
