//go:build integration

package district

import (
	"fmt"
	"log"
	"testing"

	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	"go.uber.org/zap"
)

func TestDistrictStore_StrictMatch(t *testing.T) {
	// use local database, TODO mock
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		err = errors.Wrap(err, "Failed init configuration")
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
			res, _ := testStore.CheckNormMatch(tt.city, tt.district)
			if res != tt.want {
				t.Errorf("want %v, get %v", tt.want, res)
			}

		})
	}

}

