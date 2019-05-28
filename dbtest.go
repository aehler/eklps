package main

import (
	"./eklpsDb"
	"log"
	"fmt"
)

func main() {

	fmt.Println("Connectiong...")

	db, err := eklpsDb.NewDb()
	if(err == nil){
	} else {
		log.Fatal(err)
	}

	db.Connect()
	races := db.GetRaces()

	fmt.Println(races)
}
