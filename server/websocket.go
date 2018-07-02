package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type WebSocket struct {
	conn          *websocket.Conn
	outgoing      chan interface{}
	readLoopDone  chan struct{}
	writeLoopDone chan struct{}
	logger        logrus.FieldLogger
}

func NewWebSocket(logger logrus.FieldLogger, conn *websocket.Conn) *WebSocket {
	ret := &WebSocket{
		conn:          conn,
		outgoing:      make(chan interface{}, 1),
		readLoopDone:  make(chan struct{}),
		writeLoopDone: make(chan struct{}),
		logger:        logger,
	}
	go ret.writeLoop()
	go ret.readLoop()
	return ret
}

func (ws *WebSocket) Send(msg interface{}) {
	ws.outgoing <- msg
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				ws.logger.Error(errors.Wrap(err, "websocket read error"))
			}
			return
		}

		ws.logger.Info("received message: %v", msg)
	}
}
