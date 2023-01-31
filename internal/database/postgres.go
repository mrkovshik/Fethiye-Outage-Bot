package database

import (
	"log"
	"github.com/jmoiron/sqlx"
)

// NewPostgres returns DB
func NewPostgres(dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Fatalf("failed to create database connection, %v", err)

		return nil, err
	}

	return db, nil
}