package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/crawling"
	district "github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district/postgres"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type outageRow struct {
	ID                 string    `db:"id"`
	Resource           string    `db:"resource"`
	City               string    `db:"city"`
	District           string    `db:"district"`
	StartDate          time.Time `db:"start_date"`
	EndDate            time.Time `db:"end_date"`
	SourceURL          string    `db:"source_url"`
	Notes              string    `db:"notes"`
	DateAdded          time.Time `db:"date_added"`
	CityNormalized     string    `db:"city_normalized"`
	DistrictNormalized string    `db:"district_normalized"`
	Alerted            bool      `db:"alerted"`
	// UpdatedAt  sql.NullTime `db:"updated_at"`
	// TODO if changed?
	// Add enabled or processed
}

func (from *outageRow) marshal() outage.Outage {
	return outage.Outage{
		ID:                 from.ID,
		Resource:           from.Resource,
		City:               from.City,
		District:           from.District,
		StartDate:          from.StartDate,
		EndDate:            from.EndDate,
		SourceURL:          from.SourceURL,
		Notes:              from.Notes,
		CityNormalized:     from.CityNormalized,
		DistrictNormalized: from.DistrictNormalized,
		Alerted:            from.Alerted,
	}
}

func (or *outageRow) unmarshal(from outage.Outage) error {
	if from.EndDate.Before(from.StartDate) {
		return errors.New("Start date is after End date")
	}
	*or = outageRow{
		Resource:           from.Resource,
		City:               from.City,
		District:           from.District,
		StartDate:          from.StartDate,
		EndDate:            from.EndDate,
		SourceURL:          from.SourceURL,
		Notes:              from.Notes,
		DateAdded:          time.Now().UTC(),
		CityNormalized:     from.CityNormalized,
		DistrictNormalized: from.DistrictNormalized,
		Alerted:            from.Alerted,
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

func (os *OutageStore) SetAlerted(o outage.Outage) error {
	query := fmt.Sprintf(`UPDATE outages set alerted=true where id='%v';`, o.ID)
	var orow outageRow
		err := orow.unmarshal(o)
		if err != nil {
			return err
		}
		if _, err := os.db.Query(query); err != nil {
			return errors.Wrap(err, "Error updating table")
		}
	return nil

}

func (os *OutageStore) Save(o []outage.Outage) error {
	// TODO from Columns
	query := `INSERT INTO outages (resource, city, city_normalized, district, district_normalized, start_date, end_date, source_url, notes, date_added, alerted)
  			VALUES (:resource, :city, :city_normalized, :district, :district_normalized, :start_date, :end_date, :source_url, :notes, :date_added, :alerted);`

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
		err = errors.Wrap(err, "Failed to query database:")
		return []outage.Outage{}, err
	}
	defer rows.Close()

	qryRes := make([]outage.Outage, 0)
	for rows.Next() {
		var o = outageRow{}
		if err := rows.Scan(&o.ID, &o.Resource, &o.City, &o.District, &o.CityNormalized, &o.DistrictNormalized, &o.StartDate, &o.EndDate, &o.SourceURL, &o.Notes); err != nil {
			err = errors.Wrap(err, "Failed to scan row:")
			return []outage.Outage{}, err
		}
		qryRes = append(qryRes, o.marshal())
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error iterating through rows:")
		return []outage.Outage{}, err
	}
	return qryRes, err
}

func (os *OutageStore) GetActiveOutagesByCityDistrict(distr string, city string) ([]outage.Outage, error) {
	qrtime := time.Now().UTC().String()[:19]

	query := fmt.Sprintf("SELECT id, resource, city, district, city_normalized, district_normalized, start_date, end_date, source_url, notes	FROM outages WHERE district_normalized='%v' AND city_normalized='%v' AND end_date > '%v';", distr, city, qrtime)
	return os.Read(query)
}

func (os *OutageStore) GetAllActiveOutages() ([]outage.Outage, error) {
	qrtime := time.Now().UTC().String()[:19]
	query := fmt.Sprintf("SELECT id, resource, city, district, city_normalized, district_normalized, start_date, end_date, source_url, notes	FROM outages WHERE end_date > '%v';", qrtime)
	return os.Read(query)
}

func (os *OutageStore) GetOutageByID(id string) (outage.Outage, error) {
	query := fmt.Sprintf("SELECT id, resource, city, district, city_normalized, district_normalized, start_date, end_date, source_url, notes	FROM outages WHERE id = '%v';", id)
	o, err:=os.Read(query)
	
	return o[0],err
}

func (os *OutageStore) GetAllActiveUnalertedOutages() ([]outage.Outage, error) {
	qrtime := time.Now().UTC().String()[:19]
	query := fmt.Sprintf("SELECT id, resource, city, district, city_normalized, district_normalized, start_date, end_date, source_url, notes	FROM outages WHERE end_date > '%v' AND alerted=false;", qrtime)
	return os.Read(query)
}

func (os *OutageStore) ValidateDistricts(crawled []outage.Outage) ([]district.District, error) {
	var err error
	unValidated := make([]district.District, 0)
	ds := district.NewDistrictStore(os.db)
	for _, i := range crawled {
		ok, err := ds.CheckNormMatch(i.CityNormalized, i.DistrictNormalized)
		if err != nil {
			err = errors.Wrap(err, "Failed to validate:")
			return []district.District{}, err
		}
		if !ok {
			unValidated = append(unValidated, district.District{Name: i.District, City: i.City})
		}
	}
	return unValidated, err
}

func (os *OutageStore) FindNew(crawled []outage.Outage) ([]outage.Outage, error) {
	result := make([]outage.Outage, 0)
	readed, err := os.GetAllActiveOutages()
	if err != nil {
		err = errors.Wrap(err, "Failed to query database:")
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

func (o OutageStore) FetchOutages(cfg config.Config, logger *zap.Logger) {
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
	logger.Sugar().Infoln("Crawling started")
	crawled := make([]outage.Outage, 0)
	for _, crw := range crawlers {
		res, err := crawling.CrawlOutages(crw)
		if err != nil {
			logger.Sugar().Warn(err)
		}
		crawled = append(crawled, res...)
	}
	invalidDistr, err := o.ValidateDistricts(crawled)
	if err != nil {
		logger.Sugar().Warn(err)
	}
	if invalidDistr != nil {
		logger.Warn("Attempt to add folowing invalid Districts to DB",
			zap.Any("Invalid districts", invalidDistr),
		)
	}
	f, err := o.FindNew(crawled)
	if err != nil {
		logger.Sugar().Warn(err)
	}
	logger.Sugar().Infoln("Crawling finished")
	err = o.Save(f)
	if err != nil {
		logger.Sugar().Warn(err)
	}

}
