package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spankie/tw-interview/blockparser"
)

func gracefulShutdown(ctx context.Context, apiServer *http.Server) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("Server force to shutdown: %v", err))
	}

	slog.Info("Server exiting")
}

func configureLogger(logLevel string) error {
	var level slog.Level

	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return fmt.Errorf("error setting log level: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		Level: level,
	}))
	slog.SetDefault(logger)

	return nil
}

func run(ctx context.Context, apiServer *http.Server) {
	go func() {
		slog.Info("server listening on port " + apiServer.Addr)

		err := apiServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("HTTP server error: %v", err))
			os.Exit(1)
		}
	}()

	gracefulShutdown(ctx, apiServer)
}

func main() {
	err := configureLogger(os.Getenv("TW_LOG_LEVEL"))
	if err != nil {
		log.Fatalf("could not configure logger: %v", err)
	}

	blockParser := blockparser.NewBlockParser()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	blockParser.StartBlockScanning(ctx)

	run(ctx, newServer(blockParser))
}
