package server

import (
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/vmrob/grav-game/game"
)

const tickDuration = time.Second / 30

type Server struct {
	logger          logrus.FieldLogger
	universe        *game.Universe
	router          *mux.Router
	webSockets      map[*WebSocket]struct{}
	webSocketsMutex sync.Mutex
	stop            chan struct{}
	stopped         chan struct{}
}

func DefaultUniverse() *game.Universe {
	universe := game.NewUniverse(game.Rect{X: -10000, Y: -10000, W: 20000, H: 20000})
	for i := 0; i < 500; i++ {
		universe.AddBody(&game.Body{
			Position: game.Point{X: rand.Float64()*20000 - 10000, Y: rand.Float64()*20000 - 10000},
			Mass:     rand.Float64() * 1000000,
			Velocity: game.Vector{X: rand.Float64()*1000 - 500, Y: rand.Float64()*1000 - 500},
		})
	}
	return universe
}

func NewServer(logger logrus.FieldLogger) *Server {
	ret := &Server{
		logger:     logger,
		universe:   DefaultUniverse(),
		router:     mux.NewRouter(),
		webSockets: make(map[*WebSocket]struct{}),
		stop:       make(chan struct{}),
		stopped:    make(chan struct{}),
	}
	ret.router.HandleFunc("/", ret.indexHandler)
	ret.router.HandleFunc("/game", ret.gameHandler)
	ret.router.NotFoundHandler = http.FileServer(http.Dir("dist"))
	go ret.run()
	return ret
}

func (s *Server) run() {
	defer close(s.stopped)

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.tick()
			// TODO: win condition?
			if len(s.universe.Bodies()) < 2 {
				s.universe = DefaultUniverse()
			}
		}
	}
}

func (s *Server) tick() {
	s.universe.Step(tickDuration)

	var gameState WebSocketGameState
	gameState.Universe.Bounds = s.universe.Bounds()
	gameState.Universe.Bodies = make(map[string]*game.Body)
	for id, body := range s.universe.Bodies() {
		gameState.Universe.Bodies[id.String()] = body
	}

	s.webSocketsMutex.Lock()
	defer s.webSocketsMutex.Unlock()

	for ws := range s.webSockets {
		if !ws.IsAlive() {
			delete(s.webSockets, ws)
			continue
		}

		ws.Send(&WebSocketOutput{
			GameState: &gameState,
		})
	}
}

// Close closes any hijacked connections.
func (s *Server) Close() error {
	close(s.stop)
	<-s.stopped

	var closers []io.Closer

	s.webSocketsMutex.Lock()
	for ws := range s.webSockets {
		closers = append(closers, ws)
	}
	s.webSocketsMutex.Unlock()

	for _, closer := range closers {
		closer.Close()
	}
	return nil
}

var upgrader = websocket.Upgrader{}

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Warn(err)
		return
	}

	logger := s.logger.WithField("connection_id", uuid.NewV4())
	logger.Info("accepted websocket connection")

	ws := NewWebSocket(logger, conn, s.universe)

	s.webSocketsMutex.Lock()
	defer s.webSocketsMutex.Unlock()
	s.webSockets[ws] = struct{}{}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "dist/index.html")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
