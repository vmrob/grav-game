package server

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmrob/grav-game/game"
)

type WebSocket struct {
	conn          *websocket.Conn
	outgoing      chan *WebSocketMessage
	readLoopDone  chan struct{}
	writeLoopDone chan struct{}
	logger        logrus.FieldLogger
	universe      *game.Universe
}

type WebSocketGameState struct {
	Universe struct {
		Bounds game.Rect
		Bodies map[string]*game.Body
	}
}

type WebSocketMessage struct {
	GameState      *WebSocketGameState `json:",omitempty"`
	AssignedBodyId string              `json:",omitempty"`
}

func NewWebSocket(logger logrus.FieldLogger, conn *websocket.Conn, universe *game.Universe) *WebSocket {
	ret := &WebSocket{
		conn:          conn,
		outgoing:      make(chan *WebSocketMessage, 10),
		readLoopDone:  make(chan struct{}),
		writeLoopDone: make(chan struct{}),
		logger:        logger,
		universe:      universe,
	}
	go ret.writeLoop()
	go ret.readLoop()
	return ret
}

func (ws *WebSocket) Send(msg *WebSocketMessage) {
	select {
	case ws.outgoing <- msg:
	default:
		ws.logger.Warn("dropping outgoing websocket message")
	}
}

func (ws *WebSocket) IsAlive() bool {
	select {
	case <-ws.writeLoopDone:
		return false
	default:
		return true
	}
}

func (ws *WebSocket) Close() error {
	close(ws.outgoing)
	<-ws.readLoopDone
	<-ws.writeLoopDone
	return nil
}

func (ws *WebSocket) writeLoop() {
	defer close(ws.writeLoopDone)

	defer ws.conn.Close()

	for {
		msg, ok := <-ws.outgoing
		if !ok {
			break
		}

		ws.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

		if err := ws.conn.WriteJSON(msg); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) && err != websocket.ErrCloseSent {
				ws.logger.Error(errors.Wrap(err, "websocket write error"))
			}
			break
		}
	}
}

func (ws *WebSocket) readLoop() {
	defer close(ws.readLoopDone)

	for {
		ws.conn.SetReadLimit(4 * 1024)

		var msg json.RawMessage
		err := ws.conn.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				ws.logger.Error(errors.Wrap(err, "websocket read error"))
			}
			return
		}

		ws.logger.Info("received message: %v", msg)
	}
}
