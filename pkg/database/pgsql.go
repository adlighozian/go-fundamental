package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func NewPostgresDB(conn string, logger zerolog.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, errors.New("could not connect to the database: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		return nil, errors.New("could not ping the database: " + err.Error())
	}

	return db, nil
}
