package main

import (
	"./eklpsDb"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var Db *eklpsDb.Eklps


func RaceServer(w http.ResponseWriter, req *http.Request) {

	t := time.Now()
	fmt.Printf("[%s] %s", t, "Selecting races... ");

	params, err := unPack(req.RequestURI)
	if err != nil {
		io.WriteString(w, "Wrong params!\n")
		fmt.Println("ListenAndServe: wrong params", err)
		return
	}

	races, err := selectRaces(params)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("[%s] An error occured, please copy this message and report it on the forum!\n", t))
		fmt.Println("ListenAndServe: ", err)
		return
	}

	jrace, err := json.Marshal(races)
	if err != nil {
		fmt.Println("error:", err)
		io.WriteString(w, "Wrong params!\n")
		return
	}

	io.WriteString(w, string(jrace))
}

func selectRaces(p eklpsDb.Params) ([]eklpsDb.Race, error) {

	races, err := Db.GetRaces(p)
	if err != nil {
		return []eklpsDb.Race{}, err
	}

	return races, nil
}

func unPack(req string) (eklpsDb.Params, error) {

	req = strings.Replace(req, "/getraces?data=", ``, -1)

	var p eklpsDb.Params

	req, err := url.QueryUnescape(req)

	//fmt.Println(req)

	if err != nil {
		fmt.Println("error:", err)
		return eklpsDb.Params{}, err
	}

	err = json.Unmarshal([]byte(req), &p)
	if err != nil {
		fmt.Println("error:", err)
		return eklpsDb.Params{}, err
	}

	fmt.Println(p)

	return p, nil
}

func main() {

	db, err := eklpsDb.NewDb()
	if(err == nil){
	} else {
		panic(err)
	}

	err = db.Connect()
	if err != nil {
		panic(err)
	}

	Db = db

	s, err := Db.GetSeason()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current season",s)

	http.HandleFunc("/getraces", RaceServer)
	err = http.ListenAndServe(":8011", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
