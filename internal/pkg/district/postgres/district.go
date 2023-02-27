package district

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
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
		City: d.City,
		Name: d.Name,
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

func countRatio(s string, levRatio int) int {
	res := len(s) * levRatio / 10
	if res > levRatio {
		res = levRatio
	}
	return res
}

func (d *DistrictStore) CheckStrictMatch(cit string, dis string) (bool, error) {
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

func (d *DistrictStore) fuzzyQuery(city string, dist string) ([]District, error) {
	var query string
	cfg, err := config.GetConfig()
	if err != nil {
		err = errors.Wrap(err, "Error fuzzy query to database:")
		return []District{}, err
	}
	levRatio, err := strconv.Atoi(cfg.SearchConfig.LevRatio) //Levenstein searching ratio from config
	if err != nil {
		err = errors.Wrap(err, "Error converting levRatio")
		return []District{}, err
	}
	cityLevRatio := countRatio(city, levRatio)
	districtLevRatio := countRatio(dist, levRatio)
	simRatio := cfg.SearchConfig.SimRatio //Similarity searching ratio from config
	if city == "" {
		query = fmt.Sprintf("SELECT city, district, city_normalized, district_normalized FROM districts where LEVENSHTEIN(district_normalized, '%v')<%v or Similarity(district_normalized, '%v')>%v ORDER BY LEVENSHTEIN(district_normalized, '%v') asc, Similarity(district_normalized, '%v') desc LIMIT 1;", dist, districtLevRatio, dist, simRatio, dist, dist)
	} else {
		query = fmt.Sprintf("SELECT city, district, city_normalized, district_normalized FROM districts where (LEVENSHTEIN(district_normalized, '%v')<%v or Similarity(district_normalized, '%v')>%v) and (LEVENSHTEIN(city_normalized, '%v')<%v or Similarity(city_normalized, '%v') >%v) ORDER BY LEVENSHTEIN(district_normalized, '%v') asc, LEVENSHTEIN(city_normalized, '%v') asc, Similarity(district_normalized, '%v') desc, Similarity(city_normalized, '%v') desc LIMIT 1;", dist, districtLevRatio, dist, simRatio, city, cityLevRatio, city, simRatio, dist, city, dist, city)
	}
	found, err := d.Read(query)
	if err != nil {
		err = errors.Wrap(err, "Error reading from database:")
		return []District{}, err
	}
	result := found
	return result, err
}

func (d *DistrictStore) GetFuzzyMatch(s []string) (District, error) {

	var city, dist string
	var err error
	if len(s) == 0 {
		return District{}, err
	}
	if len(s) == 1 {
		city = ""
		dist = s[0]
	}
	if len(s) == 2 {
		city = s[0]
		dist = s[1]
	}
	if len(s) > 2 {
		city = s[0]
		dist = s[1] + " " + s[2]
	}
	found, err := d.fuzzyQuery(city, dist)
	if err != nil {
		return District{}, err
	}
	if len(found) < 1 {
		found, err = d.fuzzyQuery(dist, city)
		if err != nil {
			return District{}, err
		}
		if len(found) < 1 {
			found, err = d.fuzzyQuery("", city+dist)
			if err != nil {
				return District{}, err
			}
		}
		if len(found) < 1 {
			return District{
				Name: "no matches",
				City: "no matches",
			}, err
		} else {
			return found[0], err
		}
	}
	return found[0], err
}
