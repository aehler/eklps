package horse

import (
	"../race"
	"strings"
	"strconv"
	"regexp"
	"bytes"
	"fmt"
)

type Season uint

type Horse struct {
	ID uint `db:"id"`
	Name string `db:"name"`
	Sex string `db:"sex"`
	Born Season `db:"born"`
	Height map[Season]uint
	Owner string `db:"owner"`
	Breeder string `db:"breeder"`
	WinsTotal float64 `db:"wins_total"`
	WinsBySeason map[Season]float64
	RaceCareerBySeason map[Season]*Career
	RaceCareer *Career
	RaceHistory []RaceResult
}

type Career struct {
	TotalRaces uint
	RacesWon uint
	RacesSecond uint
	RacesThird uint
	RacesFourth uint
}

type RaceResult struct {
	R *race.Race
	Position uint
	Result string
}

func (h *Horse) UnmarshalOld(b []byte) error {

	re := regexp.MustCompile(`<table id='infoblocktbl'>((?s).+?)</div>`)
	retab1 := regexp.MustCompile(`<td id='textleft2'>(.*?)</td>`)

	data := re.FindAllStringSubmatch(string(b), -1)

	sp := strings.Split(data[0][1], "\n")

	var curSeason Season = 2015

	for _, s := range sp {

		switch {
		case strings.Contains(s, "<font color='#0c2e7d' size=4>"):

			re := regexp.MustCompile(`<strong>([0-9]+) (.+)</strong>`)
			hn := re.FindAllStringSubmatch(s, -1)

			h.Name = hn[0][2]

			id, err := strconv.ParseUint(hn[0][1], 10, 32)
			if err != nil {
				return err
			}

			h.ID = uint(id)

		case strings.Contains(s, "Сезон рождения:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			seas, err := strconv.ParseUint(tx[0][1], 10, 32)
			if err != nil {
				return err
			}

			h.Born = Season(uint(seas))

		case strings.Contains(s, "Пол:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			h.Sex = tx[0][1]

		case strings.Contains(s, "Высота в холке, см:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			seas, err := strconv.ParseUint(tx[0][1], 10, 32)
			if err != nil {
				return err
			}

			h.Height = map[Season]uint{curSeason : uint(seas)}

		case strings.Contains(s, "Заводчик:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			h.Breeder = tx[0][1]

		case strings.Contains(s, "Владелец:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			h.Owner = tx[0][1]

			retab2 := regexp.MustCompile(`<td id='textright2'>Скаковая карьера:</td><td id='textleft2'>(.*?)</td>`)
			tx = retab2.FindAllStringSubmatch(s, -1)

			h.RaceCareer = &Career{}
			h.RaceCareer.Fill(tx[0][1])

		case strings.Contains(s, "Выигрыш:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			money := strings.Replace(tx[0][1], "$ ", "", -1)
			money = strings.Replace(money, " ", "", -1)
			money = strings.Replace(money, ",", "", -1)

			w, err := strconv.ParseFloat(money, 64)
			if err != nil {
				fmt.Println("Couldn't parse wins")
			}
			h.WinsTotal = w

			retab2 := regexp.MustCompile(`<td id='textright2'>За сезон:</td><td id='textleft2'>(.*?)</td>`)
			tx = retab2.FindAllStringSubmatch(s, -1)

			h.RaceCareerBySeason = map[Season]*Career{curSeason : &Career{}}
			h.RaceCareerBySeason[curSeason].Fill(tx[0][1])

		case strings.Contains(s, "Выигрыш за сезон:"):

			tx := retab1.FindAllStringSubmatch(s, -1)
			money := strings.Replace(tx[0][1], "$ ", "", -1)
			money = strings.Replace(money, " ", "", -1)
			money = strings.Replace(money, ",", "", -1)

			w, err := strconv.ParseFloat(money, 64)
			if err != nil {
				fmt.Println("Couldn't parse wins")
			}

			h.WinsBySeason = map[Season]float64{curSeason : w}

		}

		if strings.Contains(s, "<table id='infoblocktbl2'><tr id='infoblockheader'><td colspan=4>Скаковая карьера") {

			re3 := regexp.MustCompile("<table id='infoblocktbl2'><tr id='infoblockheader'><td colspan=4>Скаковая карьера</td></tr>(.*)$")

			rs := re3.FindAllStringSubmatch(s, -1)

			rer := regexp.MustCompile(`<tr id='charter' ><td><a href='raceinfo4.php\?race=([0-9]+)' class='menulink'>([0-9]+)</a></td><td align=left>(.*?)</td><td>([0-9]+)</td><td>(.*?)</td></tr>`)

			rss := rer.FindAllStringSubmatch(rs[0][0], -1)

			for _, racesrc := range rss {

				race := &race.Race{
					RaceID : racesrc[1],
				}

				race.Unmarshal(racesrc[3])

				pos, err := strconv.ParseUint(racesrc[4], 10, 32)
				if err != nil {
					pos = 0
				}

				bb := bytes.NewBufferString(racesrc[1])

				year := []byte{}

				for i:=0; i<4; i++ {
					b, _ := bb.ReadByte()
					year = append(year, b)
				}

				raceyear, err := strconv.ParseUint(string(year), 10, 64)
				if err != nil {
					raceyear = 2012
				}

				race.Season = uint(raceyear)

				raceResult := RaceResult{
					R : race,
					Position : uint(pos),
					Result : racesrc[5],
				}

				h.RaceHistory = append(h.RaceHistory, raceResult)
			}


		}

	}

	return nil
}

func (c *Career) Fill (src string) {

	t := strings.Split(src, " ")

	seas, err := strconv.ParseUint(t[0], 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	c.TotalRaces = uint(seas)

	rw := strings.Split(t[1], "-")

	f := []uint{}

	for _, p := range rw {

		pp, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			fmt.Println(err)
			pp = 0
		}

		f = append(f, uint(pp))
	}

	c.RacesWon = f[0]
	c.RacesSecond = f[1]
	c.RacesThird = f[2]
	c.RacesFourth = f[3]

}
