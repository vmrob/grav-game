package server

import "github.com/vmrob/grav-game/game"

type WebSocketBody struct {
	Position game.Point
	Mass     float64
	Radius   float64
	NetForce game.Vector
}

func NewWebSocketBody(body *game.Body) *WebSocketBody {
	return &WebSocketBody{
		Position: body.Position,
		Mass:     body.Mass,
		Radius:   body.Radius,
		NetForce: body.NetForce,
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
