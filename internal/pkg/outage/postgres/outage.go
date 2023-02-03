package postgres

import (
	"fmt"

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
		EndDate: from.EndDate,
		SourceURL: from.SourceURL, 
	}
	return nil
}

func (or *outageRow) Columns() []string {
	return []string{ 
		"id", "resource", "city", "district", "start_date", "source_url",
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

func (os *OutageStore) Save (o [] outage.Outage) error {
	// TODO from Columns
	query := `INSERT INTO outages (resource, city, district, start_date, end_date, source_url)
  			VALUES (:resource, :city, :district, :start_date, :end_date, :source_url);`
	
	var orow outageRow
	for _,i:=range o{
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

func (os *OutageStore) Read (query string) ([] outage.Outage, error) {
	rows, err := os.db.Query(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return [] outage.Outage{}, err
	}
	defer rows.Close()

	qryRes := make([]outage.Outage, 0)
	for rows.Next() {
		var o = outageRow{}
		if err := rows.Scan(&o.Resource, &o.City, &o.District, &o.StartDate, &o.EndDate, &o.SourceURL); err != nil {
			fmt.Println("Failed to scan row:", err)
			return[] outage.Outage{}, err
		}
		qryRes = append(qryRes, o.marshal())
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through rows:", err)
		return [] outage.Outage{}, err
	}

	return qryRes,err
}

func (os *OutageStore) GetActiveOutagesByCityDistrict (distr string, city string) ([] outage.Outage, error){
	qrtime :=time.Now().UTC().String()[:19]
	query := fmt.Sprintf("SELECT resource, city, district, start_date, end_date, source_url	FROM outages WHERE district ILIKE '%v' AND city ILIKE '%v' AND \"end_date\" > '%v';", ("%"+distr+"%"),("%"+city+"%"),qrtime)
	return os.Read(query)
}

func (os *OutageStore) GetOutagesByEndTime (t time.Time) ([] outage.Outage, error){
qrtime:= (t.String()[:19])
	query := fmt.Sprintf("SELECT resource, city, district, start_date, end_date, source_url FROM outages WHERE \"end_date\" > '%v';",qrtime)
	fmt.Println("query = ", query)
	return os.Read(query)
}

func (os *OutageStore) FindNew (crawled [] outage.Outage) ([] outage.Outage,error) {
result:=make([] outage.Outage,0)
readed,err:=os.GetActiveOutagesByCityDistrict("","")
if err != nil {
	fmt.Println("Failed to query database:", err)
	return [] outage.Outage{},err
}
fmt.Println("readed:")
for _,i:= range readed{
	fmt.Printf("\n%+v\n",i)
}
if len(readed) == 0{
	return crawled,err
}
for _,i:= range crawled {
for n,j:= range readed{
if i.Equal(j){
	break
}
if n==len(readed)-1{
result = append(result, i)	
}
}

} 
return result, err
}