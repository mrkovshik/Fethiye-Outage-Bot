package postgres

import (
	"fmt"
	"log"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/crawling"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/pkg/errors"
)

type outageRow struct {
	Resource  string    `db:"resource"`
	City      string    `db:"city"`
	District  string    `db:"district"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	SourceURL string    `db:"source_url"`
	Notes     string    `db:"notes"`
	DateAdded time.Time `db:"date_added"`
	// UpdatedAt  sql.NullTime `db:"updated_at"`
	// TODO if changed?
	// Add enabled or processed
}

func (from *outageRow) marshal() outage.Outage {
	return outage.Outage{
		Resource:  from.Resource,
		City:      from.City,
		District:  from.District,
		StartDate: from.StartDate,
		EndDate:   from.EndDate,
		SourceURL: from.SourceURL,
		Notes:     from.Notes,
	}
}

func (or *outageRow) unmarshal(from outage.Outage) error {
	if from.EndDate.Before(from.StartDate) {
		return errors.New("Start date is after End date")
	}
	*or = outageRow{
		Resource:  from.Resource,
		City:      from.City,
		District:  from.District,
		StartDate: from.StartDate,
		EndDate:   from.EndDate,
		SourceURL: from.SourceURL,
		Notes:     from.Notes,
		DateAdded: time.Now().UTC(),
	}
	return nil
}

func (or *outageRow) Columns() []string {
	return []string{
		"id", "resource", "city", "district", "start_date", "source_url",
	}
}

type OutageStore struct {
	db *sqlx.DB
}

func NewOutageStore(db *sqlx.DB) *OutageStore {
	return &OutageStore{
		db: db,
	}
}

func (os *OutageStore) Save(o []outage.Outage) error {
	// TODO from Columns
	query := `INSERT INTO outages (resource, city, district, start_date, end_date, source_url, notes, date_added)
  			VALUES (:resource, :city, :district, :start_date, :end_date, :source_url, :notes, :date_added);`

	var orow outageRow
	for _, i := range o {
		err := orow.unmarshal(i)
		if err != nil {
			return err
		}
		if _, err := sqlx.NamedExec(os.db, query, orow); err != nil {
			return errors.Wrap(err, "could not save New synthetic order")
		}
	}
	return nil
}

func (os *OutageStore) Read(query string) ([]outage.Outage, error) {
	rows, err := os.db.Query(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return []outage.Outage{}, err
	}
	defer rows.Close()

	qryRes := make([]outage.Outage, 0)
	for rows.Next() {
		var o = outageRow{}
		if err := rows.Scan(&o.Resource, &o.City, &o.District, &o.StartDate, &o.EndDate, &o.SourceURL, &o.Notes); err != nil {
			fmt.Println("Failed to scan row:", err)
			return []outage.Outage{}, err
		}
		qryRes = append(qryRes, o.marshal())
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through rows:", err)
		return []outage.Outage{}, err
	}

	return qryRes, err
}

func (os *OutageStore) GetActiveOutagesByCityDistrict(distr string, city string) ([]outage.Outage, error) {
	qrtime := time.Now().UTC().String()[:19]
	query := fmt.Sprintf("SELECT resource, city, district, start_date, end_date, source_url, notes	FROM outages WHERE district ILIKE '%v' AND city ILIKE '%v' AND end_date > '%v';", ("%" + distr + "%"), ("%" + city + "%"), qrtime)
	return os.Read(query)
}

func (os *OutageStore) GetOutagesByEndTime(t time.Time) ([]outage.Outage, error) {
	qrtime := (t.String()[:19])
	query := fmt.Sprintf("SELECT resource, city, district, start_date, end_date, source_url, notes FROM outages WHERE end_date > '%v';", qrtime)
	return os.Read(query)
}

func (os *OutageStore) ValidateDistricts(crawled []outage.Outage) ([]outage.Outage, error) {
	var err error
	unValidated := make([]outage.Outage, 0)
	ds := district.NewDistrictStore(os.db)
	for _, i := range crawled {
		ok, err := ds.CheckStrictMatch(i.City, i.District)
		if err != nil {
			fmt.Println("Failed to validate:", err)
			return []outage.Outage{}, err
		}
		if !ok {
			unValidated = append(unValidated, i)
		}
	}
	return unValidated, err
}

func (os *OutageStore) FindNew(crawled []outage.Outage) ([]outage.Outage, error) {
	result := make([]outage.Outage, 0)
	readed, err := os.GetActiveOutagesByCityDistrict("", "")
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return []outage.Outage{}, err
	}

	if len(readed) == 0 {
		return crawled, err
	}
	for _, i := range crawled {
		for n, j := range readed {
			if i.Equal(j) {
				break
			}
			if n == len(readed)-1 {
				result = append(result, i)
			}
		}

	}
	return result, err
}

func (os OutageStore) FetchOutages(cfg config.Config) {
	muskiURL := cfg.CrawlersURL.Muski
	aydemURL := cfg.CrawlersURL.Aydem
	var Aydem = crawling.OutageAydem{
		Url:      aydemURL,
		Resource: "power",
	}
	var Muski = crawling.OutageMuski{
		Url:      muskiURL,
		Resource: "water",
	}

	var crawlers = []crawling.Crawler{
		Aydem,
		Muski,
	}

	fmt.Println("Crawling started")
	crawled := make([]outage.Outage, 0)
	for _, crw := range crawlers {
		crawled = append(crawled, crawling.CrawlOutages(crw)...)
	}
	f, err := os.FindNew(crawled)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Crawling completed")
	err = os.Save(f)
	if err != nil {
		log.Fatal(err)
	}

}
