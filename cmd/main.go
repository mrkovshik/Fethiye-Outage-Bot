package main

import (
	"flag"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/telegram"
	"github.com/pressly/goose/v3"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

func main() {
	//reading config file
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	// Create a logger from the configuration
	logger, err := cfg.LoggerConfig.Build()
	if err != nil {
		panic(err)
	}
	//nolint:errcheck
	defer logger.Sync()

	logger.Info("App started",
		zap.Any("Version", cfg.Project.Version))
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
	if err != nil {
		logger.Fatal("Sceduler error",
			zap.Error(err),
		)
	}
	c.Start()
	store.FetchOutages(cfg, logger)
	telegram.BotRunner(ds, store, logger)
}
