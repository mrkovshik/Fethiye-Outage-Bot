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

	var getTests = [] struct {
		name string
		city string
		district string
		wantedQnty int
		}{
{"just City","Limpopo","",2 },
{"City and district","Limpopo","Ugadagada",1 },
{"non existing city","sadfsdfasd","Ugadagada",0 },
{"non existing district","Limpopo","sadfsdfasd",0 },

		}

		for _,tt:=range getTests {
	t.Run(tt.name, func(t *testing.T) {
		outages:=make([] outage.Outage,0)
		outages,err = testStore.GetActiveOutagesByCityDistrict(tt.district,tt.city)
	   if len(outages) !=tt.wantedQnty {
		t.Errorf("want %v, get %v",tt.wantedQnty,len(outages))
	   } 
	
	})
}

}
