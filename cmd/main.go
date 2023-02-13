package main

import (
	"flag"


	"log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"

	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/crawling"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"

	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/telegram"
	"github.com/pressly/goose/v3"
)

var err error





func main() {
	cfg := config.GetConfig()
	db := database.ConnectDB(cfg)
	defer db.Close()
	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()
	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Fatalf("Migration failed, %v", err)
			return
		}
	}

	store := postgres.NewOutageStore(db)
	ds := district.NewDistrictStore(db)
	go store.FetchOutages(cfg)
	telegram.BotRunner(ds, store)

}