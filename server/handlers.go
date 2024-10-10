package server

import (
	"go-axiata/https"
	"go-axiata/project/repository"
	"go-axiata/project/usecase"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func (s *Server) MapHandlers(mux *http.ServeMux) error {
	repository := repository.NewRepository(s.db, s.cfg, s.logger)
	usecase := usecase.NewUsecase(s.cfg, s.logger, repository)
	handlers := https.NewHandlers(s.cfg, usecase, s.logger)

	MapRouters(mux, handlers)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	s.server = &http.Server{
		Handler:      corsHandler,
		Addr:         ":" + s.cfg.GetEnv.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return nil

}
