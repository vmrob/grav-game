package main

import (
	"html/template"
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

type Server struct {
	logger          logrus.FieldLogger
	universe        *game.Universe
	router          *mux.Router
	webSockets      map[*WebSocket]struct{}
	webSocketsMutex sync.Mutex
	stop            chan struct{}
	stopped         chan struct{}
}

func NewServer(logger logrus.FieldLogger, universe *game.Universe) *Server {
	ret := &Server{
		logger:     logger,
		universe:   universe,
		router:     mux.NewRouter(),
		webSockets: make(map[*WebSocket]struct{}),
		stop:       make(chan struct{}),
		stopped:    make(chan struct{}),
	}
	ret.router.HandleFunc("/", ret.indexHandler)
	ret.router.HandleFunc("/game", ret.gameHandler)
	go ret.run()
	return ret
}

func (s *Server) run() {
	defer close(s.stopped)

	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

func (s *Server) tick() {
	s.webSocketsMutex.Lock()
	defer s.webSocketsMutex.Unlock()

	type GameState struct{}

	for ws := range s.webSockets {
		// TODO: actually send something
		ws.Send(GameState{})
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
	s.webSockets[NewWebSocket(logger, conn)] = struct{}{}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexTemplate.Execute(w, nil)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

var indexTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<title>grav-game</title>
	</head>
	<body>
		<p>Nothing to see here.</p>
	</body>
</html>
`))
