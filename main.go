package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/vmrob/grav-game/game"
	"github.com/vmrob/grav-game/server"
)

func main() {
	logger := logrus.StandardLogger()

	universe := game.NewUniverse(game.Rect{X: -10000, Y: -10000, W: 20000, H: 20000})
	for i := 0; i < 10; i++ {
		universe.AddBody(&game.Body{
			Position: game.Point{rand.Float64()*20000 - 10000, rand.Float64()*20000 - 10000},
			Mass:     rand.Float64() * 1000000,
			Velocity: game.Vector{rand.Float64()*1000 - 500, rand.Float64()*1000 - 500},
		})
	}

	s := server.NewServer(logger, universe)
	defer s.Close()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	done := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		logger.Info("signal caught. shutting down...")

		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Error(err)
		}
		close(done)
	}()

	logger.Info("listening at http://127.0.0.1:8080")
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error(err)
	}
	<-done
}
