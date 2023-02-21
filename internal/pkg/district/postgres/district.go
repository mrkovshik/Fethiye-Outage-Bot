package district

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/pkg/errors"
)

type District struct {
	Name string
	City string
}

type districtRow struct {
	City string `db:"city"`
	Name string `db:"district"`
}

func (d *districtRow) marshal() District {
	return District{
		City: d.City,
		Name: d.Name,
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
		fmt.Println("Failed to query database:", err)
		return []District{}, err
	}
	defer rows.Close()

	qryRes := make([]District, 0)
	for rows.Next() {
		var s districtRow
		if err := rows.Scan(&s.City, &s.Name); err != nil {
			fmt.Println("Failed to scan row:", err)
			return []District{}, err
		}
		qryRes = append(qryRes, s.marshal())
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through rows:", err)
		return []District{}, err
	}

	return qryRes, err
}

func countRatio (s string, levRatio int) int {
	res := len(s) - levRatio - 1
	switch {
	case res < 0:
		res = 0
	case res > levRatio:
		res = levRatio
	}
	return res
}

func (d *DistrictStore) CheckStrictMatch(cit string, dis string) (bool, error) {
	query := fmt.Sprintf("SELECT city, district FROM districts WHERE district ILIKE '%v' AND city ILIKE '%v';", dis, cit)
	found, err := d.Read(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return false, err
	}
	if len(found) < 1 {
		return false, err
	}

	return true, err
}
func (d *DistrictStore) fuzzyQuery(city string, dist string) ([]District, error) {
	var query string
	cfg := config.GetConfig()

	levRatio,err := strconv.Atoi(cfg.SearchConfig.LevRatio) //Levenstein searching ratio from config
	if err != nil {
		errors.Wrap(err,"Error converting levRatio")
	}
	cityLevRatio := countRatio(city, levRatio)
	districtLevRatio := countRatio(dist, levRatio)
	simRatio := cfg.SearchConfig.SimRatio //Similarity searching ratio from config
	if city == "" {
		query = fmt.Sprintf("SELECT city, district FROM districts where LEVENSHTEIN(district, '%v')<%v or Similarity(district, '%v')>%v ORDER BY LEVENSHTEIN(district, '%v') asc, Similarity(district, '%v') desc LIMIT 1;", dist, districtLevRatio, dist, simRatio, dist, dist)
	} else {
		query = fmt.Sprintf("SELECT city, district FROM districts where (LEVENSHTEIN(district, '%v')<%v or Similarity(district, '%v')>%v) and (LEVENSHTEIN(city, '%v')<%v or Similarity(city, '%v') >%v) ORDER BY LEVENSHTEIN(district, '%v') asc, LEVENSHTEIN(city, '%v') asc, Similarity(district, '%v') desc, Similarity(city, '%v') desc LIMIT 1;", dist, districtLevRatio, dist, simRatio, city, cityLevRatio, city, simRatio, dist, city, dist, city)
	}
	found, err := d.Read(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)

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
		if len(found)<1{
			found, err = d.fuzzyQuery("",city+dist)
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
