package server

import (
	"math"

	"github.com/vmrob/grav-game/game"
)

func WebSocketFloat(f float64) float32 {
	return float32(math.Round(f*100) / 100)
}

type WebSocketPoint struct {
	X float32
	Y float32
}

type WebSocketVector struct {
	X float32
	Y float32
}

type WebSocketBody struct {
	MinorName string `json:",omitempty"`
	MajorName string `json:",omitempty"`
	Position  WebSocketPoint
	Mass      float32
	Radius    float32
	NetForce  WebSocketVector
}

func NewWebSocketBody(body *game.Body) *WebSocketBody {
	return &WebSocketBody{
		MinorName: body.MinorName,
		MajorName: body.MajorName,
		Position: WebSocketPoint{
			X: WebSocketFloat(body.Position.X),
			Y: WebSocketFloat(body.Position.Y),
		},
		Mass:   WebSocketFloat(body.Mass),
		Radius: WebSocketFloat(body.Radius),
		NetForce: WebSocketVector{
			X: WebSocketFloat(body.NetForce.X),
			Y: WebSocketFloat(body.NetForce.Y),
		},
	}
}

type WebSocketGameState struct {
	Universe struct {
		Bounds game.Rect
		Bodies map[string]*WebSocketBody
	}
}

type WebSocketOutput struct {
	GameState      *WebSocketGameState `json:",omitempty"`
	AssignedBodyId string              `json:",omitempty"`
}

type WebSocketInput struct {
	Thrust *game.Vector `json:",omitempty"`
}
