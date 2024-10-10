package main

import (
	"go-axiata/config"
	"go-axiata/pkg/database"
	"go-axiata/pkg/logger"
	"go-axiata/server"
	"log"
)

func main() {
	logger := logger.CreateLogger()

	logger.Info().Msg("Starting api server")

	cfgFile, err := config.LoadConfig("config")
	if err != nil {
		logger.Error().Msg("LoadConfig: " + err.Error())
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		logger.Error().Msg("ParseConfig: " + err.Error())
	}

	dbPgSql, err := database.NewPostgresDB(cfg.GetEnv.DbPg, logger)
	if err != nil {
		logger.Warn().Msg(err.Error())
	} else {
		logger.Info().Msg("Connected to PostgreSQL database!")
	}
	defer dbPgSql.Close()

	s := server.NewServer(cfg, logger, dbPgSql)
	if err := s.Run(); err != nil {
		log.Fatalf("error while run server, %v", err)
	}

}
