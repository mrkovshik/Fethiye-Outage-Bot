package district

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type District struct {
	Name           string
	City           string
	NameNormalized string
	CityNormalized string
}

type districtRow struct {
	City           string `db:"city"`
	Name           string `db:"district"`
	NameNormalized string `db:"city_normalized"`
	CityNormalized string `db:"district_normalized"`
}

func (d *districtRow) marshal() District {
	return District{
		City:           d.City,
		Name:           d.Name,
		NameNormalized: d.NameNormalized,
		CityNormalized: d.CityNormalized,
	}
}

func (su *districtRow) Columns() []string {
	return []string{
		"id", "city", "district",
	}
}

type DistrictStore struct {
	db *sqlx.DB
}

func NewDistrictStore(db *sqlx.DB) *DistrictStore {
	return &DistrictStore{
		db: db,
	}
}

func (sstr *DistrictStore) Read(query string) ([]District, error) {
	rows, err := sstr.db.Query(query)
	if err != nil {
		err = errors.Wrap(err, "Error query to database:")
		return []District{}, err
	}
	defer rows.Close()
	qryRes := make([]District, 0)
	for rows.Next() {
		var s districtRow
		if err := rows.Scan(&s.City, &s.Name, &s.CityNormalized, &s.NameNormalized); err != nil {
			err = errors.Wrap(err, "Failed to scan DB row:")
			return []District{}, err
		}
		qryRes = append(qryRes, s.marshal())
	}
	if err := rows.Err(); err != nil {
		err = errors.Wrap(err, "Error iterating through DB rows:")
		return []District{}, err
	}
	return qryRes, err
}


func (d *DistrictStore) CheckNormMatch(cit string, dis string) (bool, error) {
	query := fmt.Sprintf("SELECT city, district, city_normalized, district_normalized FROM districts WHERE district_normalized='%v' AND city_normalized='%v';", dis, cit)
	found, err := d.Read(query)
	if err != nil {
		err = errors.Wrap(err, "Failed to read from database: ")
		return false, err
	}
	if len(found) < 1 {
		return false, err
	}
	return true, err
}

func (d *DistrictStore) GetNormFromDB (cit string, dis string) (District, error) {
	query := fmt.Sprintf("SELECT city, district, city_normalized, district_normalized FROM districts WHERE district='%v' AND city='%v';", dis, cit)
	found, err := d.Read(query)
	if err != nil {
		err = errors.Wrap(err, "Failed to read from database: ")
		return District{}, err
	}
	if len(found) < 1 {
		return District{}, err
	}

	return found[0], err
}


