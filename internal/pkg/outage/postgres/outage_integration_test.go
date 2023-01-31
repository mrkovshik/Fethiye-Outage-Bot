package postgres

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/database"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"gotest.tools/assert"
)


func TestOutageStore_Save(t *testing.T) {
	// use local database, TODO mock
	if err := config.ReadConfigYML("../../../../config.yml"); err != nil {
		log.Fatal("Failed init configuration")
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
	log.Printf(dsn)

	db, err := database.NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		log.Fatal("Failed init postgres")
	}
	defer db.Close()
	// TODO initialise out of test scope
	testStore := NewOutageStore(db)

	t.Run("Regular outage", func(t *testing.T) {
		err = testStore.Save(outage.Outage{
			Resource:  "water", // TODO enum
			City:      "Fethiye",
			District:  "Babatashi",
			StartDate: time.Now(),
			EndDate:   time.Now(),
			SourceURL: "",
		})
		assert.NilError(t, err)
	})

	t.Run("Broken dates outage", func(t *testing.T) {
		err = testStore.Save(outage.Outage{
			Resource:  "water", // TODO enum
			City:      "Fethiye",
			District:  "Babatashi",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(-1 * time.Hour),
			SourceURL: "",
		})
		assert.Error(t, err, "Start date is after End date")
	})

}
