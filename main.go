package main

import (
	"fmt"
	"log"
	"github.com/mrkovshik/Fethiye-Outage-Bot/DB"
	"github.com/mrkovshik/Fethiye-Outage-Bot/parsing"

)



func main() {
	fmt.Println("Here we go")
	muskiOutages:=parsing.ParceFromMuski()
	err:= DB.AddToDB(muskiOutages)
	if err != nil {
		log.Fatal(err)
	}

}
