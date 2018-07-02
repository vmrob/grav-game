package server

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmrob/grav-game/game"
)

func newWebsocketConnection(server *Server) (*websocket.Conn, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	defer l.Close()

	httpServer := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer l.Close()
			server.ServeHTTP(w, r)
		}),
	}
	go httpServer.Serve(l)

	var client *websocket.Conn
	var firstError error
	for attempts := 0; true; attempts++ {
		client, _, err = websocket.DefaultDialer.Dial("ws://"+l.Addr().String()+"/game", http.Header{})
		if err == nil {
			return client, nil
		} else if firstError == nil {
			firstError = err
		}
		if attempts > 100 {
			return nil, firstError
		}
		time.Sleep(time.Millisecond * 10)
	}
	return nil, nil
}

func TestServer(t *testing.T) {
	universe := game.NewUniverse(game.Rect{X: 0, Y: 0, W: 100, H: 100})

	s := NewServer(logrus.StandardLogger(), universe)
	defer s.Close()

	client, err := newWebsocketConnection(s)
	require.NoError(t, err)
	defer client.Close()

	var msg GameStateMessage
	assert.NoError(t, client.ReadJSON(&msg))
	assert.NoError(t, client.ReadJSON(&msg))
	assert.NoError(t, client.ReadJSON(&msg))
}
