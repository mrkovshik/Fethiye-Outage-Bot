package main

import (
	"flag"
	"fmt"

	"log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"

	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/crawling"
	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"

	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/telegram"
	"github.com/pressly/goose/v3"
)

var err error





func main() {
	cfg := config.GetConfig()
	db := database.ConnectDB(cfg)
	defer db.Close()
	migration := flag.Bool("migration", false, "Defines the migration start option")
	flag.Parse()
	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Fatalf("Migration failed, %v", err)
			return
		}
	}

	// muskiStore := postgres.NewOutageStore(db)
	// ds := district.NewDistrictStore(db)
	// go muskiStore.FetchOutages(cfg)
	// telegram.BotRunner(ds, muskiStore)
	aydemURL:=cfg.CrawlersURL.Aydem
	var Aydem = crawling.OutageAydem {
		Url:aydemURL,
		Resource: "power",
	}
	r:=crawling.CrawlOutages(Aydem)
	for _,i:= range r {
		fmt.Printf("%+v\n\n", i)
	}

}