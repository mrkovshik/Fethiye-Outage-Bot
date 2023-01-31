package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"

	// "github.com/mrkovshik/Fethiye-Outage-Bot/parsing"
	// "github.com/mrkovshik/Fethiye-Outage-Bot/postgresdb"

	"github.com/pressly/goose/v3"
)

// type crawler interface {
// 	crawl() []outage.Outage
// }

func main() {	
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

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
	defer db.Close()

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Fatal("Migration failed")

			return
		}
	}

	// TODO
	// var muskiOutages []parsing.WaterOutage
	// fmt.Println("Here we go")
	// muskiOutages.Crawl()
	// err = postgresdb.AddToDB(muskiOutages)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
