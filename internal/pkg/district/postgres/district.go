package district

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/config"
	"github.com/mrkovshik/Fethiye-Outage-Bot/internal/pkg/district"
)
func getConfig() config.Config {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatalf("Failed init configuration %v", err)
	}
	return config.GetConfigInstance()
}

type districtRow struct {
	City      string	`db:"city"`
	Name  string	`db:"district"`
}

func (d *districtRow) marshal() district.District {
	return district.District{
		City: d.City,
		Name: d.Name,
	}
}

func (d *districtRow) unmarshal(from district.District) error {
	*d = districtRow{
		City: from.City,
		Name: from.Name,	
	}
	return nil
}

func (su *districtRow) Columns() []string {
	return []string{ 
		"id", "city", "district",
	}
}

type DistrictStore struct {
	db           *sqlx.DB
}

func NewDistrictStore(db *sqlx.DB) *DistrictStore {
	return &DistrictStore {
		db: db,
	}
}




func (sstr *DistrictStore) Read (query string) ([]district.District, error) {
	rows, err := sstr.db.Query(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return [] district.District{}, err
	}
	defer rows.Close()

	qryRes := make([]district.District, 0)
	for rows.Next() {
		var s districtRow
		if err := rows.Scan(&s.City, &s.Name); err != nil {
			fmt.Println("Failed to scan row:", err)
			return[] district.District{}, err
		}
		qryRes = append(qryRes, s.marshal())
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through rows:", err)
		return [] district.District{}, err
	}

	return qryRes,err
}

func (d *DistrictStore) CheckStrictMatch (cit string, dis string) (bool, error) {
	query := fmt.Sprintf("SELECT city, district FROM districts WHERE district ILIKE '%v' AND city ILIKE '%v';", dis, cit)
	found,err:=d.Read(query)
	if err != nil {
		fmt.Println("Failed to query database:", err)
		return false, err
	}
	if len(found)<1 {
		return false,err
	}
	
return true,err
}
func (d *DistrictStore) fuzzyQuery (city string, dist string ) ([] district.District, error) {
	cfg:=getConfig()
	levRatio:= cfg.SearchConfig.LevRatio //Levenstein searching ratio from config
	simRatio:= cfg.SearchConfig.SimRatio //Similarity searching ratio from config
	query := fmt.Sprintf("SELECT city, district FROM districts where LEVENSHTEIN(district, '%v')<%v or Similarity(district, '%v')>%v ORDER BY LEVENSHTEIN(district, '%v') asc, LEVENSHTEIN(city, '%v') asc, Similarity(district, '%v') desc, Similarity(city, '%v') desc LIMIT 1;",dist,levRatio, dist, simRatio, dist, city,dist, city, )
		found,err:=d.Read(query)
		if err != nil {
		fmt.Println("Failed to query database:", err)
		
	}
	result:=found
		return result, err
}

func (d *DistrictStore) GetFuzzyMatch (input string) (district.District, error) {

	var city, dist string
	var err error
	s:=strings.Split(input," ")
	if len(s)==0{
		return district.District{},err
	}
	if len(s)==1{
		city=""
		dist=s[0]
	}
	if len(s)==2{
		city=s[0]
		dist=s[1]
	}
	found,err:=d.fuzzyQuery(city,dist)
	if err != nil {
		return district.District{},err		
	}

	if len(found)<1 {
		found,err:=d.fuzzyQuery(dist,city)
	if err != nil {
		return district.District{},err		
	}
	if len(found)<1 {		
	 		return district.District{
			Name: "no matches",
			City: "no matches",
		},err
	} else {
		return found[0],err
		}
	}

return found[0],err
}