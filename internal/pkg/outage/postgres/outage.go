package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/outage"
	"github.com/pkg/errors"
)

type outageRow struct {
	Resource  string	`db:"resource"`
	City      string	`db:"city"`
	District  string	`db:"district"`
	StartDate time.Time	`db:"start_date"`
	Duration  time.Duration
	EndDate   time.Time `db:"end_date"`
	SourceURL string	`db:"source_url"`
	// UpdatedAt  sql.NullTime `db:"updated_at"` 
	// TODO if changed?
	// Add enabled or processed
}

func (or *outageRow) marshal() outage.Outage {
	return outage.Outage{
		Resource: or.Resource,
		City: or.City,
		District: or.District,
		StartDate: or.StartDate,
		Duration: or.EndDate.Sub(or.StartDate),
		EndDate: or.EndDate,
		SourceURL: or.SourceURL, 
	}
}

func (or *outageRow) unmarshal(from outage.Outage) error {
	if from.EndDate.Before(from.StartDate) {
		return errors.New("Start date is after End date")
	}
	*or = outageRow{
		Resource: from.Resource,
		City: from.City,
		District: from.District,
		StartDate: from.StartDate,
		Duration: from.EndDate.Sub(from.StartDate),
		EndDate: from.EndDate,
		SourceURL: from.SourceURL, 
	}
	return nil
}

func (or *outageRow) Columns() []string {
	return []string{ 
		"id", "resource", "city", "district", "start_date", "duration", "source_url",
	}
}

type OutageStore struct {
	db           *sqlx.DB
}

func NewOutageStore(db *sqlx.DB) *OutageStore {
	return &OutageStore {
		db: db,
	}
}

func (os *OutageStore) Save (o outage.Outage) error {
	// TODO from Columns
	query := `INSERT INTO outages (resource, city, district, start_date, end_date, source_url)
  			VALUES (:resource, :city, :district, :start_date, :end_date, :source_url);`
	
	var orow outageRow
	err := orow.unmarshal(o)
	if err != nil {
		return err
	}
	
	if _, err := sqlx.NamedExec(os.db, query, orow); err != nil {
		return errors.Wrap(err, "could not save New synthetic order")
	}

	return nil
}
