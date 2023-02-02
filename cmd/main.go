package main

import (
	"flag"
	"fmt"
	"log"
	// "time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"

	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/crawling"
	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"

	// "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"

	"github.com/pressly/goose/v3"
)


func main() {	
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatalf("Failed init configuration %v", err)
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
			log.Fatalf("Migration failed, %v", err)
			return
		}
	}


	var Muski = crawling.OutageMuski {
		Url:"https://www.muski.gov.tr/",
		Resource: "water",
	}
r:= Muski.Crawl() 
fmt.Println("Crawled from muski:")
for _,i:= range r{
	fmt.Printf("\n%+v\n",i)
}

muskiStore:=postgres.NewOutageStore(db)
f,_:=muskiStore.FindNew(r)
fmt.Println("Crawled remains:")
for _,i:= range f{
	fmt.Printf("\n%+v\n",i)
}

muskiStore.Save(f)


// k,err:=muskiStore.GetOutagesByDistrict("")
// if err != nil {
// 	log.Fatalf("reading failed, %v", err)
// 	}
// for _,i:= range k{
// 	fmt.Printf("%+v\n",i)
// }
// }
}
