package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/telegram"
	"github.com/pressly/goose/v3"
)

var err error

func getConfig() config.Config {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatalf("Failed init configuration %v", err)
	}
	return config.GetConfigInstance()
}

func connectDB(cfg config.Config) *sqlx.DB {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	log.Printf(dsn)
	db, err := database.NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		log.Fatal("Failed init postgres")
	}
	return db
}

func main() {
	cfg := getConfig()
	db := connectDB(cfg)
	defer db.Close()
	migration := flag.Bool("migration", false, "Defines the migration start option")
	flag.Parse()
	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Fatalf("Migration failed, %v", err)
			return
		}
	}

	muskiStore := postgres.NewOutageStore(db)
	ds := district.NewDistrictStore(db)
	go muskiStore.FetchOutages(cfg)
	telegram.BotRunner(ds, muskiStore)

}
