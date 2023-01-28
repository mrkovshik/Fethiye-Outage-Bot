package parsing

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Outage struct {
	ID        int           `db:"id"`
	City      string        `db:"city"`
	District  string        `db:"district"`
	StartDate time.Time     `db:"start_date"`
	Duration  int		    `db:"duration"`
	EndDate   time.Time     `db:"end_date"` 
}

type OutageStore struct {
	db           *sql.DB
	// queryBuilder sq.StatementBuilderType
}

func NewOutageStore(db *sql.DB) *OutageStore {
	return &OutageStore {
		db:           db,
		// queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

const oldTimeFormat = "02.01.2006 15:04"
type dbCred struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}
var cred  = dbCred {
	host: "localhost",
	port:     5432,
	user:     "postgres",
	password: "17pasHres19!",
	dbName:   "vacancies",
	}

func connectDB() (*sql.DB, error) {
	// Connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cred.host, cred.port, cred.user, cred.password, cred.dbName)
		db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
func addToDB (insertRows []Outage) string {
	db,err:=connectDB()		
	if err != nil {
		errMsg:=fmt.Sprintf("Error connecting to DB: %v", err)
		return errMsg
	}
		defer db.Close()
	
	stmt, err := db.Prepare("INSERT INTO outages (city, district, start_date ,  duration) VALUES($1, $2,$3,$4)")
	if err != nil {
		errMsg:=fmt.Sprintf("Error inserting registries to DB: %v", err)
	
		return errMsg
	}	
	for _,i:=range insertRows{
	_, err = stmt.Exec(i.City,i.District, i.StartDate, i.Duration)
	if err != nil {
		errMsg:=fmt.Sprintf("Error inserting registries to DB: %v", err)
	
		return errMsg
	}
	stmt.Close()
}
	return "Запись добавлена"
}
func ParceFromMuski() {
	rowSlice := make([]Outage, 0)
	doc, err := goquery.NewDocument("https://www.muski.gov.tr")
	if err != nil {
		log.Fatal(err)
	}
	table := doc.Find("table#plansiz")
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		fmt.Println("row ", i)
		if i > 2 {
			rowSlice = append(rowSlice, Outage{})
			k := i - 3
			row.Find("td").Each(func(j int, cell *goquery.Selection) {

				fmt.Println("cell ", j, cell.Text())
				switch {
				case j == 2:
					rowSlice[k].City = cell.Text()
				case j == 3:
					rowSlice[k].District = cell.Text()
				case j == 4:
					parsedDur, err := strconv.ParseInt(strings.Trim(cell.Text(), " Saat"), 0, 64)
					if err != nil {
						log.Fatal(err)
					}
					rowSlice[k].Duration = int(parsedDur)
					case j == 5:
					parsedTime := strings.Trim(cell.Text(), " ")
					rowSlice[k].StartDate, err = time.Parse(oldTimeFormat, parsedTime)
					if err != nil {
						log.Fatal(err)
					}
					rowSlice[k].EndDate=rowSlice[k].StartDate.Add(time.Duration(rowSlice[k].Duration)*time.Hour)
				}

			})
		}
		for _,i:=range rowSlice {
			fmt.Printf("%+v\n", i)
		}
	
	})

}
