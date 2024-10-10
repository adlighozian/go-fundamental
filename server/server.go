package server

import (
	"context"
	"database/sql"
	"fmt"
	"go-axiata/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	cfg    *config.Config
	logger zerolog.Logger
	db     *sql.DB
	server *http.Server
}

func NewServer(cfg *config.Config, logger zerolog.Logger, db *sql.DB) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	if err := s.MapHandlers(mux); err != nil {
		s.logger.Error().Msg(err.Error())
		return err
	}

	portMsg := fmt.Sprintf("Server is listening on port: %s", s.cfg.GetEnv.Port)
	s.logger.Info().Msg(portMsg)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Warn().Msg("Error starting server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error().Msg("Server forced to shutdown: " + err.Error())
	}

	s.logger.Warn().Msg("Server has been shut down")

	return nil
}
