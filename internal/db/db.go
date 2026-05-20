package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Open(dsn string) (*sqlx.DB, error) {
	database, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(15)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(30 * time.Minute)
	return database, database.Ping()
}
