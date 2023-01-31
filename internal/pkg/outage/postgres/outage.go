package postgres

import "time"

type outageRow struct {
	Resource  string	`db:"resource"`
	City      string	`db:"city"`
	District  string	`db:"district"`
	StartDate time.Time	`db:"start_date"`
	Duration  time.Duration		`db:"duration"`
	EndDate   time.Time `db:"end_date"`
	SourceURL string	`db:"duration"`
	// UpdatedAt  sql.NullTime `db:"updated_at"` 
	// TODO if changed?

	Resource         uuid.UUID `db:"id"`
	CustomerID string    `db:"customer_id"`
	Status     string    `db:"status"`
	Type       string    `db:"type"`

	Currency   money.Currency  `db:"currency"`
	Amount     decimal.Decimal `db:"amount"`

	CreatedAt  time.Time    `db:"created_at"`
}

func (or outageRow) unmarshal() postgresdb.outageRow {
	return postgresdb.outageRow{
		Resource:  o.Resource,
		City:      o.City,
		District:  o.District,
		StartDate: o.StartDate,
		Duration:  o.Duration,
		EndDate:   o.EndDate,
		SourceURL: o.SourceURL,
	}
}