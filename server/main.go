package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/vmrob/grav-game/game"
)

func main() {
	logger := logrus.StandardLogger()

	universe := game.NewUniverse(game.Rect{X: 0, Y: 0, W: 100, H: 100})

	server := NewServer(logger, universe)
	defer server.Close()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
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
