package server

import (
	"io"
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
const threatSpawnInterval = time.Second * 5
const foodSpawnInterval = time.Millisecond * 100

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
	return game.NewUniverse(game.Rect{X: -5000, Y: -5000, W: 10000, H: 10000})
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

	tickTicker := time.NewTicker(tickDuration)
	threatTicker := time.NewTicker(threatSpawnInterval)
	foodTicker := time.NewTicker(foodSpawnInterval)
	defer foodTicker.Stop()
	defer threatTicker.Stop()
	defer tickTicker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-threatTicker.C:
			s.universe.AddEvent(game.ThreatSpawnEvent(s.universe))
		case <-foodTicker.C:
			s.universe.AddEvent(game.FoodSpawnEvent(s.universe))
		case <-tickTicker.C:
			s.tick()
		}
	}
}

func (s *Server) tick() {
	s.universe.Step(tickDuration)

	var gameState WebSocketGameState
	gameState.Universe.Bounds = s.universe.Bounds()
	gameState.Universe.Bodies = make(map[string]*WebSocketBody)
	for id, body := range s.universe.Bodies() {
		gameState.Universe.Bodies[id.String()] = NewWebSocketBody(body)
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Warn(err)
		return
	}
	conn.EnableWriteCompression(true)

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
