package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
)

func ConnectDB(cfg config.Config) *sqlx.DB {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	log.Printf(dsn)
	db, err := NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		log.Fatal("Failed init postgres")
	}
	return db
}
// NewPostgres returns DB
func NewPostgres(dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Fatalf("failed to create database connection, %v", err)

		return nil, err
	}

	return db, nil
}