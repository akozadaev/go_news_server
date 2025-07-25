package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go_news_server/pkg/config"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App, logger *zap.SugaredLogger, cfg *config.Config) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Errorw("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	fiberConnURL := fmt.Sprintf(
		"%v:%v",
		cfg.ServeHost,
		cfg.ServerPort,
	)

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App, logger *zap.SugaredLogger, cfg *config.Config) {
	fiberConnURL := fmt.Sprintf(
		"%v:%v",
		cfg.ServeHost,
		cfg.ServerPort,
	)

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		logger.Errorw("Oops... Server is not running! Reason: %v", err)
	}
}
