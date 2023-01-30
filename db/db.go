package db

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/mrkovshik/Fethiye-Outage-Bot/parsing"
	_ "github.com/lib/pq"
)


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
	dbName:   "outageDB",
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
func AddToDB (insertRows []parsing.Outage) error {
	db,err:=connectDB()		
	if err != nil {
		return err
	}
		defer db.Close()
	
	stmt, err := db.Prepare("INSERT INTO outages (city, district, resource, start_date ,  duration, end_date, source_url) VALUES($1, $2,$3,$4,$5,$6,$7)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _,i:=range insertRows{
	_, err = stmt.Exec( i.City,i.District,i.Resource, i.StartDate, i.Duration, i.EndDate, i.SourceURL)
	if err != nil {
		return err
	}
	
}
err=removeDup(db)
if err != nil {
	return err
}
	return err
}

func removeDup (db *sql.DB) error{
	stmt, err := db.Prepare("DELETE FROM outages WHERE id NOT IN (SELECT MIN(id) FROM outages GROUP BY district, start_date, resource);")
	if err != nil {
		return err
	}
		defer stmt.Close()
		_, err = stmt.Exec()
	return err
}