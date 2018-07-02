package server

import (
	"fmt"
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

	ticker := time.NewTicker(tickDuration)
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

type GameStateMessage struct {
	Universe struct {
		Bounds game.Rect
		Bodies map[string]*game.Body
	}
}

func (s *Server) tick() {
	s.universe.Step(tickDuration)
	var gameState GameStateMessage
	gameState.Universe.Bounds = s.universe.Bounds()
	gameState.Universe.Bodies = make(map[string]*game.Body)
	for id, body := range s.universe.Bodies() {
		gameState.Universe.Bodies[fmt.Sprintf("%v", id)] = body
	}

	s.webSocketsMutex.Lock()
	defer s.webSocketsMutex.Unlock()

	for ws := range s.webSockets {
		if !ws.IsAlive() {
			delete(s.webSockets, ws)
			continue
		}

		ws.Send(&gameState)
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

	s.webSocketsMutex.Lock()
	defer s.webSocketsMutex.Unlock()
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
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>Gravity Game</title>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
	</head>
	<body>
		<center>
			<canvas id="gameCanvas" width="1200" height="700" style="border:1px solid #000000;"></canvas>
		</center>
		<span id="message"></span>
		<script>
		const GRID_LINE_INTERVAL = 250;

		class Rect {
			constructor(x, y, width, height) {
				this.x = x;
				this.y = y;
				this.width = width;
				this.height = height;
			}
		}

		class Vector {
			constructor(x, y) {
				this.x = x;
				this.y = y;
			}

			withMagnitude(m) {
				var ret = new Vector(this.x, this.y);
				var scale = m / this.magnitude();
				ret.x *= scale;
				ret.y *= scale;
				return ret;
			}

			magnitude() {
				return Math.sqrt(this.x * this.x + this.y * this.y);
			}
		}

		class Universe {
			constructor() {
				this.state = null;
			}

			draw(context) {
				var min = new Vector(0, 0);
				var max = new Vector(context.canvas.width, context.canvas.height);

				var padding = 100;
				for (const [id, body] of Object.entries(this.state["Bodies"])) {
					if (body["Static"]) {
						continue;
					}
					var pos = new Vector(body["Position"]["X"], body["Position"]["Y"]);
					var r = body["Radius"];
					if (pos.x - r - padding < min.x) {
							min.x = pos.x - r - padding;
					}
					if (pos.y - r - padding < min.y) {
							min.y = pos.y - r - padding;
					}
					if (pos.x + r + padding > max.x) {
							max.x = pos.x + r + padding;
					}
					if (pos.y + r + padding > max.y) {
							max.y = pos.y + r + padding;
					}
				}

				context.clearRect(0, 0, context.canvas.width, context.canvas.height);
				var scaleX = context.canvas.width / (max.x - min.x);
				var scaleY = context.canvas.height / (max.y - min.y);
				var scale = scaleX > scaleY ? scaleY : scaleX;
				context.scale(scale, scale);
				context.translate(-min.x, -min.y);

				var bounds = new Rect(this.state["Bounds"]["X"], this.state["Bounds"]["Y"], this.state["Bounds"]["W"], this.state["Bounds"]["H"])

				for (var x = bounds.x; x < bounds.x + bounds.width; x += GRID_LINE_INTERVAL) {
					context.beginPath();
					context.strokeStyle = '#000000';
					context.moveTo(x, bounds.y);
					context.lineTo(x, bounds.y + bounds.height);
					context.stroke();
				}
				for (var y = bounds.y; y < bounds.y + bounds.height; y += GRID_LINE_INTERVAL) {
					context.beginPath();
					context.strokeStyle = '#000000';
					context.moveTo(bounds.x, y);
					context.lineTo(bounds.x + bounds.width, y);
					context.stroke();
				}

				this.drawBodies(context);
				this.drawBounds(context);

				context.translate(min.x, min.y);
				context.scale(1.0 / scale, 1.0 / scale);
			}

			drawBodies(context) {
				for (const [id, body] of Object.entries(this.state["Bodies"])) {
					var r = body["Radius"];
					var f = new Vector(body["NetForce"]["X"], body["NetForce"]["Y"]);
					var pos = new Vector(body["Position"]["X"], body["Position"]["Y"])
					var mass = body["Mass"]

					var fMag = f.magnitude();
					var fNorm = new Vector(f.x / fMag, f.y / fMag);

					context.beginPath();
					context.arc(pos.x, pos.y, r, 0, 2 * Math.PI);
					context.fillStyle = this.color;
					context.fill();
					context.lineWidth = 5;
					context.strokeStyle = '#003300';
					context.stroke();

					context.lineWidth = 2;
					context.strokeStyle = '#FF00FF';
					context.globalAlpha = 0.7;
					context.setLineDash([20, 15]);
					context.beginPath();
					var lStart = new Vector(pos.x + r * fNorm.x, pos.y + r * fNorm.y);
					context.moveTo(lStart.x, lStart.y);
					context.lineTo(lStart.x + f.x / mass, lStart.y + f.y / mass);
					context.stroke();
					context.setLineDash([]);
					context.globalAlpha = 1.0;
				}
			}

			drawBounds(context) {
				context.rect(this.state["Bounds"]["X"], this.state["Bounds"]["Y"], this.state["Bounds"]["W"], this.state["Bounds"]["H"]);
				context.stroke();
			}
		}

		var canvas = document.getElementById('gameCanvas');
		var context = canvas.getContext("2d");
		var universe = new Universe()

		function update(state) {
			universe.state = state;
			universe.draw(context);
		}

		var ws = new WebSocket('ws://127.0.0.1:8080/game');
		ws.onmessage = function(e) {
			document.getElementById('message').innerText = e.data;
			console.log(e)
			update(JSON.parse(e.data)["Universe"])
		};
		ws.onerror = function(e) {
			document.getElementById('message').innerText = 'unable to connect';
		};
		</script>
	</body>
</html>
`))
