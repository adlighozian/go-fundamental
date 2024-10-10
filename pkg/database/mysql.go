package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

func NewMysqlDB(conn string, logger zerolog.Logger) (*sql.DB, error) {

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Connect to DB!")

	return db, nil
}
