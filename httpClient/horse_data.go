package httpClient

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"net/url"
	"log"
)

func (client *client) GetOldCard(id uint) (string, error) {

	log.Printf("getting racer %d data from %s\n", id, client.Url)

	form := url.Values{}
	form.Add("racer", fmt.Sprintf("%d", id))
	form.Add("submit", "поиск")

	b := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", client.Url+"info_horsecard.php", b)
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

	log.Printf("got data, total %d bytes\n", len(body))

	return string(body), nil
}
