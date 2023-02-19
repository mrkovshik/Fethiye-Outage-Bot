package main

import (
	"encoding/json"
	"flag"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/robfig/cron"
	"go.uber.org/zap"

	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/telegram"
	"github.com/pressly/goose/v3"
)

func main() {
	// Open the configuration file
	configFile, err := os.Open("logger_config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	// Decode the configuration file into a zap.Config struct
	var logCfg zap.Config
	if err := json.NewDecoder(configFile).Decode(&logCfg); err != nil {
		panic(err)
	}

	// Create a logger from the configuration
	logger, err := logCfg.Build()
	if err != nil {
		panic(err)
	}

	//nolint:errcheck
	defer logger.Sync()

	//reading config file
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error("",
			zap.Error(err),
		)
	}
	db := database.ConnectDB(cfg, logger)
	defer db.Close()
	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()
	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			logger.Fatal("Migration failed",
				zap.Error(err),
			)
		}
	}
	c := cron.New()
	store := postgres.NewOutageStore(db)
	ds := district.NewDistrictStore(db)
	err = c.AddFunc(cfg.SchedulerConfig.FetchPeriod, func() { store.FetchOutages(cfg, logger) })
	c.Start()
	if err != nil {
		logger.Fatal("Sceduler error",
			zap.Error(err),
		)
	}
	telegram.BotRunner(ds, store, logger)
}
