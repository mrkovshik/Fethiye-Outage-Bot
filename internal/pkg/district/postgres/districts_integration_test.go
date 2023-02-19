//go:build integration

package district

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	"go.uber.org/zap"
)

func TestDistrictStore_StrictMatch(t *testing.T) {
	// use local database, TODO mock
	// Open the configuration file
	configFile, err := os.Open("logger_config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	// Decode the configuration file into a zap.Config struct
	var logCfg zap.Config
	if err := json.NewDecoder(configFile).Decode(&logCfg); err != nil {
		log.Fatal(err)
	}

	// Create a logger from the configuration
	logger, err := logCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	//nolint:errcheck
	defer logger.Sync()

	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		logger.Fatal("Failed init configuration",
			zap.Error(err),
		)
	}
	cfg := config.GetConfigInstance()
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	logger.Info("dsn: ",
		zap.String("", dsn),
	)
	db, err := database.NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		logger.Fatal("Failed init postgres",
			zap.Error(err),
		)
	}
	defer db.Close()
	// TODO initialise out of test scope
	testStore := NewDistrictStore(db)

	var tests = []struct {
		name     string
		city     string
		district string
		want     bool
	}{
		{"normal query", "Fethiye", "Karaçulha", true},
		{"caps", "FETHIYE", "Karaçulha", true},
		{"non existing city", "sadfsdfasd", "Karaçulha", false},
		{"non existing district", "Fethiye", "sadfsdfasd", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := testStore.CheckStrictMatch(tt.city, tt.district)
			if res != tt.want {
				t.Errorf("want %v, get %v", tt.want, res)
			}

		})
	}

}
func TestDistrictStore_FuzzyMatch(t *testing.T) {
	// use local database, TODO mock
	// Open the configuration file
	configFile, err := os.Open("logger_config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	// Decode the configuration file into a zap.Config struct
	var logCfg zap.Config
	if err := json.NewDecoder(configFile).Decode(&logCfg); err != nil {
		log.Fatal(err)
	}

	// Create a logger from the configuration
	logger, err := logCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	//nolint:errcheck
	defer logger.Sync()

	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		logger.Fatal("Failed init configuration",
			zap.Error(err),
		)
	}
	cfg := config.GetConfigInstance()
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)
	logger.Info("dsn: ",
		zap.String("", dsn),
	)
	db, err := database.NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		logger.Fatal("Failed init postgres",
			zap.Error(err),
		)
	}
	defer db.Close()
	// TODO initialise out of test scope
	testStore := NewDistrictStore(db)

	var tests = []struct {
		name      string
		input     string
		wantCity  string
		wantDistr string
	}{
		{"normal query", "Fethiye Karaçulha", "Fethiye", "Karaçulha"},
		{"caps", "FETHIYE Karaçulha", "Fethiye", "Karaçulha"},
		{"wrong spelling", "Fetie menteseolu", "Fethiye", "Menteşeoğlu"},
		{"non existing city", "sadfsdfasd Karaçulha", "Fethiye", "Karaçulha"},
		{"only district", "Karaçulha", "Fethiye", "Karaçulha"},
		{"total nonsense", "lsdfhjk iorewjg", "no matches", "no matches"},
		{"no space", "FethiyeKaraçulha", "Fethiye", "Karaçulha"},
		{"vice versa", "Karaçulha Fethiye", "Fethiye", "Karaçulha"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, _ := testStore.GetFuzzyMatch(tt.input)
			if res.Name != tt.wantDistr || res.City != tt.wantCity {
				t.Errorf("want %v and %v, get %v and %v", tt.wantDistr, tt.wantCity, res.Name, res.City)
			}

		})
	}

}
