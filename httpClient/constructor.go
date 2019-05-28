package httpClient

import (
	"log"
	"net/http"
	"crypto/tls"
)

func NewClient() (*client) {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &client {
		Url: 		"http://eklps.com/stables/",
		NewUrl:		"http://eklps.com/eclipse/inc/",
		Port:		"80",
		TimeLayout: "02.01.2006 3:04:05.000000 pm (MST)",
		Client: 	http.Client{Transport: tr},
		Login:		"", //site login
		Password:	"", //site password
	}

	if err := c.eklpsLogin(); err != nil {

		log.Fatal(err)

	}

	log.Println("logged in, uid: ", c.Uid)

	return c
}
