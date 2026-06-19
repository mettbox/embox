package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"embox/internal/api/routes"
	"embox/internal/config"

	"gorm.io/gorm"
)

// InitServer initializes the HTTP server with the provided database connection.
// It sets up the router and handles graceful shutdown on interrupt signals.
// The server runs in a goroutine and listens for incoming requests.
func InitServer(db *gorm.DB, apiConfig *config.ApiConfig) {
	router := routes.Init(db, apiConfig)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", apiConfig.Server.Host, apiConfig.Server.Port),
		Handler: router,
	}

	go func() {
		slog.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	// Signal catching for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}

	slog.Info("server stopped")
}
