package main

import (
	"weasel/app"
	"log"
	"io/ioutil"
	"./horse"
	"fmt"
)

//var Db *eklpsDb.Eklps

func main() {

	a := app.New("conf.d")

	fmt.Println(a)

	horse := getRacerOldCard(uint(5845))

	fmt.Println(horse)

}

func getRacerOldCard(id uint) *horse.Horse{

	b, err := ioutil.ReadFile(fmt.Sprintf("horsedata/%d.horse", id))
	if err != nil {
		log.Fatal(err)
	}

	h := &horse.Horse{}

	if err := h.UnmarshalOld(b); err != nil {

		log.Println(err)

	}

	return h

}
