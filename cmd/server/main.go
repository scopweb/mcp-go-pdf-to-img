package main

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/scopweb/mcp-go-pdf-to-img/internal/config"
	"github.com/scopweb/mcp-go-pdf-to-img/internal/logging"
	"github.com/scopweb/mcp-go-pdf-to-img/pkg/converter"
)

func main() {
	// Load configuration from environment
	cfg := config.NewServerConfig()

	// Initialize logger
	var logger logging.Logger
	if cfg.LogFormat == "json" {
		logger = logging.NewJSON(cfg.LogLevel)
	} else {
		logger = logging.New(cfg.LogLevel)
	}

	// Initialize PDF converter
	conv, err := converter.New()
	if err != nil {
		logger.Error("failed to initialize PDF converter", err)
		return
	}
	defer conv.Close()

	// Initialize handlers
	handlers := NewHandlers(conv, logger, cfg)

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/api/v1/pdf/convert", handlers.Convert)
	mux.HandleFunc("/api/v1/pdf/info", handlers.Info)

	// Create HTTP server with configuration
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server
	logger.Info("starting HTTP server",
		slog.String("addr", addr),
		slog.String("log_level", cfg.LogLevel))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", err)
	}
}