package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres DB pool
var Db *pgxpool.Pool

// Initialize PostgreSQL connection pool
func InitDB() error {
	connStr := "postgres://postgres:postgres@localhost:5432/docdb?sslmode=disable"
	var err error
	Db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}
	return Db.Ping(context.Background())
}
