package main

import (
	"fmt"
	"log"
	"github.com/mrkovshik/Fethiye-Outage-Bot/db"
	"github.com/mrkovshik/Fethiye-Outage-Bot/parsing"

)

type crawler interface {
	crawl() [] parsing.Outage
}



func main() {
	var muskiOutages parsing.WaterOutage
	fmt.Println("Here we go")
	muskiOutages.Crawl
	err:= db.AddToDB(muskiOutages)
	if err != nil {
		log.Fatal(err)
	}

}
