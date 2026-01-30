package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(connString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}
