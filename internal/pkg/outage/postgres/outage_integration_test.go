package postgres

import (
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/pkg/errors"
	"gotest.tools/assert"
)

func TestOutageStore_Save(t *testing.T) {
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		err=errors.Wrap(err, "Failed init configuration")
		log.Fatal(err)
	}
	cfg := config.GetConfigInstance()
	// Create a logger from the configuration
	logger, err := cfg.LoggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	//nolint:errcheck
	defer logger.Sync()
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
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		err=errors.Wrap(err, "Failed init configuration")
		log.Fatal(err)
	}
	cfg := config.GetConfigInstance()
	// Create a logger from the configuration
	logger, err := cfg.LoggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	//nolint:errcheck
	defer logger.Sync()
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
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		err=errors.Wrap(err, "Failed init configuration")
		log.Fatal(err)
	}
	cfg := config.GetConfigInstance()
	// Create a logger from the configuration
	logger, err := cfg.LoggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	//nolint:errcheck
	defer logger.Sync()
	var StartDate = time.Now()
	var tests = []struct {
		name    string
		outages []outage.Outage
		want    []district.District
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
			}}, []district.District{}},
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
			[]district.District{{
				City: "Ankara",
				Name: "Ankara",
			}},
		}}

	db := database.ConnectDB(cfg, logger)
	defer db.Close()
	testStore := NewOutageStore(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invDistr, err := testStore.ValidateDistricts(tt.outages)
			if err != nil {
				log.Fatal(err)
			}
			if len(invDistr) != len(tt.want) {
				t.Errorf("want %v, get %v", tt.want, invDistr)
			}
			if len(invDistr) > 0 {
				for i := range invDistr {
					if invDistr[i].Name != tt.want[i].Name || invDistr[i].City != tt.want[i].City {
						t.Errorf("want %v, get %v", tt.want[i], invDistr[i])
					}
				}
			}
		})
	}
}
