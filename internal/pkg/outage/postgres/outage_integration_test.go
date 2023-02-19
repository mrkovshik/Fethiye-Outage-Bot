package postgres

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"go.uber.org/zap"
	"gotest.tools/assert"
)

func TestOutageStore_Save(t *testing.T) {
	// Open the configuration file
	configFile, err := os.Open("../../../../logger_config.json")
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
	// use local database, TODO mock
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		log.Fatal("Failed init configuration")
	}
	cfg := config.GetConfigInstance()
	db := database.ConnectDB(cfg, logger)
	defer db.Close()
	// TODO initialise out of test scope
	testStore := NewOutageStore(db)

	t.Run("Regular outage", func(t *testing.T) {
		err = testStore.Save([]outage.Outage{{
			Resource:  "water", // TODO enum
			City:      "test sity",
			District:  "test district",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(1 * time.Hour),
			SourceURL: "test entry",
		}})
		assert.NilError(t, err)
	})

	t.Run("Broken dates outage", func(t *testing.T) {
		err = testStore.Save([]outage.Outage{{
			Resource:  "water", // TODO enum
			City:      "test sity",
			District:  "test district",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(-1 * time.Hour),
			SourceURL: "",
		}})
		assert.Error(t, err, "Start date is after End date")
	})
}

func TestOutageStore_Get(t *testing.T) {
	// Open the configuration file
	configFile, err := os.Open("../../../../logger_config.json")
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
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		log.Fatal("Failed init configuration")
	}
	cfg := config.GetConfigInstance()
	db := database.ConnectDB(cfg, logger)
	defer db.Close()
	testStore := NewOutageStore(db)
	var getTests = []struct {
		name       string
		city       string
		district   string
		wantedQnty int
	}{
		{"just City", "Limpopo", "", 2},
		{"City and district", "Limpopo", "Ugadagada", 1},
		{"non existing city", "sadfsdfasd", "Ugadagada", 0},
		{"non existing district", "Limpopo", "sadfsdfasd", 0},
	}

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			outages, _ := testStore.GetActiveOutagesByCityDistrict(tt.district, tt.city)
			if len(outages) != tt.wantedQnty {
				t.Errorf("want %v, get %v", tt.wantedQnty, len(outages))
			}

		})
	}

}

func TestOutageStore_Validate(t *testing.T) {
	configFile, err := os.Open("../../../../logger_config.json")
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
	var StartDate = time.Now()
	var tests = []struct {
		name    string
		outages []outage.Outage
		want    []outage.Outage
	}{
		{"everything valid", []outage.Outage{{
			Resource:  "water",
			City:      "Fethiye",
			District:  "Babataşı",
			StartDate: StartDate,
			EndDate:   StartDate.Add(1 * time.Hour),
			SourceURL: "test entry",
		},
			{
				Resource:  "water",
				City:      "Ortaca",
				District:  "Dalyan",
				StartDate: time.Now(),
				EndDate:   time.Now().Add(1 * time.Hour),
				SourceURL: "test entry",
			}}, []outage.Outage{}},
		{"something invalid", []outage.Outage{{
			Resource:  "water",
			City:      "Fethiye",
			District:  "Babataşı",
			StartDate: StartDate,
			EndDate:   StartDate.Add(1 * time.Hour),
			SourceURL: "test entry",
		},
			{
				Resource:  "water",
				City:      "Ankara",
				District:  "Ankara",
				StartDate: StartDate,
				EndDate:   StartDate.Add(1 * time.Hour),
				SourceURL: "test entry",
			}},
			[]outage.Outage{{
				Resource:  "water",
				City:      "Ankara",
				District:  "Ankara",
				StartDate: StartDate,
				EndDate:   StartDate.Add(1 * time.Hour),
				SourceURL: "test entry",
			}},
		}}
	if err = config.ReadConfigYML("../../../../config.yml"); err != nil {
		log.Fatal("Failed init configuration")
	}
	cfg := config.GetConfigInstance()
	db := database.ConnectDB(cfg, logger)
	defer db.Close()
	testStore := NewOutageStore(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outages, err := testStore.ValidateDistricts(tt.outages)
			if err != nil {
				log.Fatal(err)
			}
			if len(outages) != len(tt.want) {
				t.Errorf("want %v, get %v", tt.want, outages)
			}
			if len(outages) > 0 {
				for i := range outages {
					if !outages[i].Equal(tt.want[i]) {
						t.Errorf("want %v, get %v", tt.want[i], outages[i])
					}
				}
			}
		})
	}
}
