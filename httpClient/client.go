package httpClient

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"io/ioutil"
	"bytes"
)

type client struct {
	Url			string
	NewUrl      string
	Port		string
	TimeLayout 	string
	Uid			string
	Sid			string
	Client		http.Client
	Login		string
	Password	string
}

func (client *client) GetRaces() (string, error) {
	log.Printf("getting races from %s\n", client.Url)
	fmt.Println("")
	if err := client.eklpsLogin(); err != nil {
		return client.Uid, err
	}

	req, err := http.NewRequest("GET", client.Url+"race_entries.php", nil)
	if err != nil {
		return "", err
	}

	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		"PHPSESSID",
		client.Uid,
		"/",
		"eklps.com",
		expire,
		expire.Format(time.UnixDate),
		86400,
		true,
		true,
		fmt.Sprintf("PHPSESSID=%s", client.Uid),
		[]string{fmt.Sprintf("PHPSESSID=%s", client.Uid)},
	}
	req.AddCookie(&cookie)

	resp, err := client.Client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("got race entries, total %d bytes\n", len(body))
	fmt.Println()
	return string(body), nil
}


func (client *client) eklpsLogin () error {

	resp, err := http.Get(fmt.Sprintf("%s/login.php", client.Url))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	client.Uid = strings.Replace(fmt.Sprintf("%s", resp.Header["Set-Cookie"]), "[PHPSESSID=", ``, -1)
	client.Uid = strings.Replace(client.Uid, "; path=/]", ``, -1)

	data := url.Values{}
	data.Set("username", client.Login)
	data.Add("userpass", client.Password)

	req, err := http.NewRequest("POST", fmt.Sprintf("%slogin.php", client.Url), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		"PHPSESSID",
		client.Uid,
		"/",
		"eklps.com",
		expire,
		expire.Format(time.UnixDate),
		86400,
		true,
		true,
		fmt.Sprintf("PHPSESSID=%s", client.Uid),
		[]string{fmt.Sprintf("PHPSESSID=%s", client.Uid)},
	}
	req.AddCookie(&cookie)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err = client.Client.Do(req)

	defer resp.Body.Close()

	return nil
}
