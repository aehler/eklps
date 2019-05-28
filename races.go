package main

import (
	"fmt"
	"regexp"
	"encoding/xml"
	"strings"
	"strconv"
	"./race"
	//"./httpClient"
//	"database/sql"
//	"github.com/go-sql-driver/mysql"
)

type Horse struct {
	Id int
	Name string
	NextRace int
	Sc bool
	Rested int
	Sex string
}

func main() {


//	hc := httpClient.NewClient()
//	races, err := hc.GetRaces()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(len(races))


	s := getSchedule()

	fmt.Println(len(s))

	re := regexp.MustCompile(`<tr id="charter">(.+?)</tr>`)

	races := re.FindAllStringSubmatch(s, -1)

	r := race.Race{}
	for i, _ := range races {
		r = xmlUnmarshalSchedule(races[i][0])
		r.Conditions = prepareConditions(r.Class)
		fmt.Println(r)
	}

//
//	hr := getHorses()
//
//	reh := regexp.MustCompile("<tr id='charter'><td align='right'><a id='(.*)")
//
//	horses := reh.FindAllStringSubmatch(hr, -1)
//
//	for i, _ := range horses {
//		h := xmlHorseUnmarshal(horses[i][0])
//		h = h;
//	}

//	db, err := sql.Open("mysql", "user:password@/dbname")

}

func xmlHorseUnmarshal (r string) Horse {

	h := Horse{}

	r = r + "</tr>"

	r = strings.Replace(r, "'", `"`, -1)
	r = strings.Replace(r, "<br><font size=1 color=#666>", `<br /><font size="1" color="#666">`, -1)

	fmt.Println(r)

	type Result struct {
		XMLName xml.Name `xml:"tr"`
		Id string `xml:"id,attr"`
		Td []string `xml:"td"`
	}

	v := Result{}
	err := xml.Unmarshal([]byte(r), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	fmt.Println(v)

	return h
}

func prepareConditions(s string) race.Conditions {

	c := race.Conditions{}

	switch {
	case regexp.MustCompile("время на данной дистанции не резвее").MatchString(s) :
		re := regexp.MustCompile("[0-9:]+")
		c.Timing = re.FindString(s)
		c.TimingEq = ">="
	case regexp.MustCompile("время на данной дистанции не медленнее").MatchString(s) :
		re := regexp.MustCompile("[0-9:]+")
		c.Timing = re.FindString(s)
		c.TimingEq = "<="
	case regexp.MustCompile("минимальный").MatchString(s) :
		re := regexp.MustCompile("[A-Z]+")
		c.Class = re.FindString(s)
		c.ClassEq = ">="
	case regexp.MustCompile("максимальный").MatchString(s) :
		re := regexp.MustCompile("[A-Z]+")
		c.Class = re.FindString(s)
		c.ClassEq = "<="
	case len(s) == 3 :
		re := regexp.MustCompile("[A-Z]+")
		c.Class = re.FindString(s)
		c.ClassEq = "="
	}

	return c
}

func xmlUnmarshal (r string) race.Race {

	r = strings.Replace(r, "&", ",", -1)

	re := regexp.MustCompile("</strong>.*<td>")
	class := strings.Replace(re.FindAllString(r, 1)[0], "<td>", "", -1)
	class = strings.Replace(class, "</strong>", "", -1)
	class = strings.Replace(class, "</td>", "", -1)
	class = strings.Replace(class, "(", "", -1)
	class = strings.Replace(class, ")", "", -1)
	class = strings.Replace(class, "класс резвости", "", -1)

	rz := regexp.MustCompile("Заявлено.*")
	class = rz.ReplaceAllString(class, "")

	type Result struct {
		XMLName xml.Name `xml:"tr"`
		Id    int   `xml:"td>a"`
		Conditions string `xml:"td>strong"`
	}

	v := Result{}
	err := xml.Unmarshal([]byte(r), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	rsc := false
	dist_ := ""
	age := ""
	name := ""

	if regexp.MustCompile("ст-з|торф").MatchString(v.Conditions) {
		sc := regexp.MustCompile(", ").Split(v.Conditions, 4)
		name = sc[0]
		rsc = true
		dist_ = strings.Replace(sc[2], " м", "", -1)
		age = sc[3]
	} else {
		sc := regexp.MustCompile(", ").Split(v.Conditions, 3)
		name = sc[0]
		dist_ = strings.Replace(sc[1], " м", "", -1)
		age = sc[2]
	}

	age = strings.Replace(age, "yo", "", -1)

	dist, err := strconv.ParseInt(strings.Replace(dist_, "ст-з", "", -1), 10, 32)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	sex := ""

	switch {
		case regexp.MustCompile("К").MatchString(v.Conditions) : sex = "f"
		case regexp.MustCompile("Ж").MatchString(v.Conditions) : sex = "m"
		default : sex = "all"
	}

	date := v.Id / 100

	return race.Race{v.Id, date, class, dist ,rsc, age, sex, name, race.Conditions{}}
}


func xmlUnmarshalSchedule (r string) race.Race {

	type Result struct {
		XMLName string
		Id    int
		Conditions string
	}

	v := Result{}

	r = strings.Replace(r, "&", ",", -1)

	re := regexp.MustCompile(`<font color="#666" size="1">(.*)</font>`)
	class := re.FindAllStringSubmatch(r, -1)[0][0]
	class = strings.Replace(class, `<font color="#666" size="1">`, "", -1)
	class = strings.Replace(class, "</font>", "", -1)
	class = strings.Replace(class, "(", "", -1)
	class = strings.Replace(class, ")", "", -1)
	class = strings.Replace(class, "класс резвости", "", -1)

	reid := regexp.MustCompile(`<td width="100">([0-9]+)?</td>`)
	id, err := strconv.ParseInt(reid.FindStringSubmatch(r)[1], 10, 0)
	if err != nil {
		fmt.Println(err)
	}

	v.Id = int(id)

	ren := regexp.MustCompile(`<td align="left">(.+?)<font`)
	v.XMLName = ren.FindStringSubmatch(r)[1]

	v.Conditions = v.XMLName

	rsc := false
	dist_ := ""
	age := ""
	name := ""

	if regexp.MustCompile("ст-з|торф").MatchString(v.Conditions) {
		sc := regexp.MustCompile(", ").Split(v.Conditions, 4)
		name = sc[0]
		rsc = true
		dist_ = strings.Replace(sc[2], " м", "", -1)
		age = sc[3]
	} else {
		sc := regexp.MustCompile(", ").Split(v.Conditions, 3)
		name = sc[0]
		dist_ = strings.Replace(sc[1], " м", "", -1)
		age = sc[2]
	}

	age = strings.Replace(age, "yo", "", -1)

	dist, err := strconv.ParseInt(strings.Replace(dist_, "ст-з", "", -1), 10, 32)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	sex := ""

	switch {
	case regexp.MustCompile("К").MatchString(v.Conditions) : sex = "f"
	case regexp.MustCompile("Ж").MatchString(v.Conditions) : sex = "m"
	default : sex = "all"
	}

	date := v.Id / 100

	return race.Race{v.Id, date, class, dist ,rsc, age, sex, name, race.Conditions{}}

}


func getSchedule() string {

	s := `<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">25 августа 2014</td></tr>
<tr id="charter"><td width="100">2014082501</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082502</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082503</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082504</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082505</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082511</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082512</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082513</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082514</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082515</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">26 августа 2014</td></tr>
<tr id="charter"><td width="100">2014082601</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082602</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082603</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082604</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082605</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082609</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082610</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082611</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082612</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082613</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">27 августа 2014</td></tr>
<tr id="charter"><td width="100">2014082701</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082702</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082703</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082704</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082705</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082706</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082707</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082708</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082709</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082710</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">28 августа 2014</td></tr>
<tr id="charter"><td width="100">2014082801</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082802</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082803</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082804</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082805</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082806</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082807</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082808</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082809</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082810</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">29 августа 2014</td></tr>
<tr id="charter"><td width="100">2014082901</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082902</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082903</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082904</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082905</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082906</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082907</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082908</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082909</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014082910</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">30 августа 2014</td></tr>
<tr id="charter"><td width="100">2014083001</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083002</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083003</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083004</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083005</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083006</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083007</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083008</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083009</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083010</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">31 августа 2014</td></tr>
<tr id="charter"><td width="100">2014083101</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083102</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083103</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083104</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083105</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083106</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083107</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083108</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083109</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014083110</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">1 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090101</td><td align="left">Гр.III Приз Адама, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014090102</td><td align="left">Гр.III Венец славы, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014090103</td><td align="left">Гр.II Amethyst Plate, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014090104</td><td align="left">Гр.I Свита Артемиды, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014090105</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090106</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090107</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090108</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090109</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090110</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090111</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090112</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090113</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090114</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">2 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090201</td><td align="left">Гр.III Love Hearts Cup, ст-з, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014090202</td><td align="left">Гр.II Rozarium Sprint, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014090203</td><td align="left">Гр.II Non Stop Flight Stakes, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014090204</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090205</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090206</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090207</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090208</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090209</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090210</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090211</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090213</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">3 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090301</td><td align="left">Гр.II Приз Бирюзы, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014090302</td><td align="left">Гр.II Starfish Stakes, ст-з, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014090303</td><td align="left">Гр.II Russian Bazzar, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014090304</td><td align="left">Гр.I Future Stars Plate, 1700 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014090306</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090308</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090309</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090311</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">4 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090401</td><td align="left">Гр.III Success Stakes, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014090402</td><td align="left">Гр.II Joncol Plate, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014090403</td><td align="left">Гр.I Доспехи Гора, ст-з, 1800 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014090405</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090406</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090407</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090408</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090409</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">5 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090501</td><td align="left">Гр.III 6 Furlong Battle, 1200 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014090502</td><td align="left">Гр.III Ночные фиалки, 1700 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014090503</td><td align="left">Гр.III Parad Cup, ст-з, 1600 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014090504</td><td align="left">Гр.III Триумф Цезаря, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014090505</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090506</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090507</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090510</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090511</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090512</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090513</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090514</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090515</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090516</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">6 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090601</td><td align="left">Гр.II Seren Sea Stakes, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014090602</td><td align="left">Гр.II Flamingo Bay Plate, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014090603</td><td align="left">Гр.II Sandals Stakes, ст-з, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014090604</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090605</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090606</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090607</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090608</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090609</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090612</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090614</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090615</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">7 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090701</td><td align="left">AA El Grande, 1400 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014090702</td><td align="left">AA Golden Ring Plate, 1900 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014090703</td><td align="left">AA Lexington Handicap, ст-з, 1000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014090704</td><td align="left">AA Saratoga Springs Handicap, ст-з, 1800 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014090705</td><td align="left">AA Приз Пражский, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014090706</td><td align="left">AA Oaklawn Park Gold, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014090707</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090708</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090709</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090710</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090711</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090712</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090714</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090715</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090716</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090717</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">8 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090802</td><td align="left">Гр.III Альпийские Луга, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014090803</td><td align="left">Гр.III Moon N Sea Chase, ст-з, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014090804</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090805</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090806</td><td align="left">Гр.I Brilliants Cup, 1300 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">9 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014090901</td><td align="left">Гр.III Гавайский Приз, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014090902</td><td align="left">Гр.II Приз Ниагары, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014090903</td><td align="left">Гр.II Morning Chase, ст-з, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014090904</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090905</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090906</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090907</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090908</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090910</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090911</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090912</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090913</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014090914</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">10 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091001</td><td align="left">Гр.II Canada Stakes, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014091002</td><td align="left">Гр.II Mirage Grade Run, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014091003</td><td align="left">Гр.I Long Way Home Cup, 2600 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091004</td><td align="left">Гр.I Rose Velvet Stakes, ст-з, 1700 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091005</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091006</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091008</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091009</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091010</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091011</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091012</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091013</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091016</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091017</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">11 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091101</td><td align="left">Гр.III Green Fairy, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091102</td><td align="left">Гр.III Приз Ришелье, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091103</td><td align="left">Гр.II Rockefeller Stakes, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014091104</td><td align="left">Гр.I Roses Sprint Chase, ст-з, 1000 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091105</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091106</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091107</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091108</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">12 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091201</td><td align="left">Гр.III Sonata Handicap, ст-з, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091202</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091203</td><td align="left">Гр.II Львиное Сердце, 1800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014091204</td><td align="left">Гр.II Montclair Stakes, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091205</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091206</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091207</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091208</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091209</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091210</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091211</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091215</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091216</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">13 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091301</td><td align="left">Гр.III Приз Романтический, 1300 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014091302</td><td align="left">Гр.III Saratoga Stakes, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091303</td><td align="left">Гр.III Whitewater Chase, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014091304</td><td align="left">Гр.II Tarquinius Chase, ст-з, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014091305</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091306</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091307</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091308</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091309</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091310</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091311</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">14 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091401</td><td align="left">AA Big Red Springs Handicap, 1100 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014091402</td><td align="left">AA Indigo Chase, 1800 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014091403</td><td align="left">AA Infinity Cup, 2400 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014091404</td><td align="left">AA Emerald Isle, 4400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014091405</td><td align="left">AA Oman Oasis Plate, торф, 1100 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014091406</td><td align="left">AA Saltan Plate, торф, 1700 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014091407</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091408</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091411</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091412</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091415</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091416</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091417</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">15 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091501</td><td align="left">Гр.III Bluez Plate, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091502</td><td align="left">Гр.III Приз Сергиев Посад, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091503</td><td align="left">Гр.I Fires Sprint, 1000 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091504</td><td align="left">Гр.I Приз Пражский, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014091505</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091508</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091509</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091510</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091511</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091512</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091513</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">16 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091601</td><td align="left">Гр.I Приз Мистраль, 2400 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091602</td><td align="left">Гр.III Movement Stakes, 2600 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014091603</td><td align="left">Гр.I Pink Panter Run, 3400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091604</td><td align="left">Гр.III Greenlands Stakes, ст-з, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091607</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091609</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091610</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091611</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091615</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091616</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091617</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091618</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">17 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091701</td><td align="left">Гр.III Olay Plate, ст-з, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091702</td><td align="left">Гр.II Intercoila Stakes, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014091703</td><td align="left">Гр.I Кубок Колумба, 1700 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091707</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091710</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091711</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091712</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091713</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091714</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091715</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091716</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">18 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091801</td><td align="left">Гр.III Two Miles Plate, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014091802</td><td align="left">Гр.II Гольфстрим, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014091803</td><td align="left">Гр.I Woodbine Dash, ст-з, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014091804</td><td align="left">Гр.I Dreams And Hopes Stakes, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014091808</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091809</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091810</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">19 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014091901</td><td align="left">Гр.I Fashion Sprint, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014091902</td><td align="left">Гр.I Luck Of The Draw Endurance, 2800 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014091903</td><td align="left">Гр.II Beautiful Stranger Stakes, ст-з, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014091904</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091907</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091908</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091909</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091910</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091911</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091912</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091913</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091914</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014091915</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">20 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092001</td><td align="left">Гр.II Irish Futurity, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014092002</td><td align="left">Гр.II National Day Stakes, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014092003</td><td align="left">Гр.II Три Богатыря, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014092006</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092007</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092009</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092010</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092011</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092012</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092013</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092014</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092015</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">21 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092101</td><td align="left">AA Blue Diamond Filly Prelude, 1000 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092102</td><td align="left">AA Keeneland Challenge, 2800 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092103</td><td align="left">AA Park Plaza Stakes, ст-з, 1200 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092104</td><td align="left">AA Черная Жемчужина, торф, 1900 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092105</td><td align="left">AA Legend Cup, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092106</td><td align="left">AA Ocean Star Chase, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092110</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092112</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092113</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092114</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092115</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092116</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">22 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092201</td><td align="left">Гр.III Desert Dreams Stakes, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014092202</td><td align="left">Гр.II Приз Классика, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014092203</td><td align="left">Гр.I Golden Range Cup, ст-з, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092207</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092208</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092210</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092211</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092212</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092213</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092214</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092215</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">23 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092301</td><td align="left">Гр.II Two Power Stakes, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014092302</td><td align="left">Гр.II My Fair Lady, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014092303</td><td align="left">Гр.I Impresario Handicap, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092304</td><td align="left">Гр.I Park Plaza Stakes, ст-з, 1100 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014092313</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092314</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092315</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092316</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092317</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092318</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092320</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092321</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092322</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">24 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092401</td><td align="left">Гр.II Ice King Stakes, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014092402</td><td align="left">Гр.I Hear The Ghost Plate, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092403</td><td align="left">Гр.I Донна Анна, ст-з, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092404</td><td align="left">Гр.I Lilies Stakes, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014092408</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092409</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092410</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092411</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092412</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092413</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092414</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">25 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092501</td><td align="left">Гр.I Green Ocianic Plate, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092502</td><td align="left">Гр.I Plan B Stakes, 2000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014092503</td><td align="left">Гр.I Долгая Дорога, 4000 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014092504</td><td align="left">Гр.I Приз Лазури, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014092505</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092507</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092510</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092512</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092513</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">26 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092601</td><td align="left">Гр.III Приз Яшмы, 1800 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014092602</td><td align="left">Гр.III Наше Будущее, ст-з, 1200 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014092603</td><td align="left">Гр.II Rodeo Valley Stakes, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014092604</td><td align="left">Гр.I Wind Music Stakes, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092608</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092609</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092610</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092611</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092612</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092614</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092615</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">27 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092701</td><td align="left">Гр.III Monrovia Cup, ст-з, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014092702</td><td align="left">Гр.III San Fernando Stakes, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014092703</td><td align="left">Гр.II Capricon Cup, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014092704</td><td align="left">Гр.I Бегущая по Волнам, 1100 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014092708</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092709</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092710</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">28 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092801</td><td align="left">AA Crown Horse Plate, 1300 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092802</td><td align="left">AA Freedom Cup, 1600 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092803</td><td align="left">AA Grand Prix Du Mans, 2200 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092804</td><td align="left">AA Brooklyn Stakes, 3600 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092805</td><td align="left">AA National Day Stakes, торф, 1000 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092806</td><td align="left">AA Prima Plate, торф, 1300 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092807</td><td align="left">AA Приз Аквамарин, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014092811</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092813</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092814</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092815</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092816</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092817</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092818</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092819</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">29 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014092901</td><td align="left">Гр.III Gliss Stakes, 2000 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014092902</td><td align="left">Гр.III Stop the Clock Run, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014092903</td><td align="left">Гр.I Woodbine Classic, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014092906</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092907</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014092909</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">30 сентября 2014</td></tr>
<tr id="charter"><td width="100">2014093001</td><td align="left">Гр.II Zanter Plate, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014093002</td><td align="left">Гр.I East Lake Stakes, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014093003</td><td align="left">Гр.I Never Ever Chase, ст-з, 1800 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014093006</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093007</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093008</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093009</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093010</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093012</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093013</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093014</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093015</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014093017</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">1 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100101</td><td align="left">Гр.III Long Run Uphill Cup, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014100102</td><td align="left">Гр.II Sharp Novices Hurdle, ст-з, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014100103</td><td align="left">Гр.I Glamour Handicap, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014100104</td><td align="left">Гр.I Подвески Королевы, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014100107</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100108</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100110</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100111</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">2 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100201</td><td align="left">Гр.III Strawberry Fields, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014100202</td><td align="left">Гр.III Приз Рябины, 3000 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100203</td><td align="left">Гр.III Sunset Distaff Chase, ст-з, 1700 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100207</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100212</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100213</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">3 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100301</td><td align="left">Гр.III Kanmar Plate, ст-з, 1400 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100302</td><td align="left">Гр.III Randir Handicap, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100303</td><td align="left">Гр.III Трилогия, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100335</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100336</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100337</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100338</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100339</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100340</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100341</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100342</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100344</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100345</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">4 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100401</td><td align="left">Гр.III Carlsberg Cup, ст-з, 1700 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100402</td><td align="left">Гр.III Lynns Handicap, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100403</td><td align="left">Гр.I Dunes Stakes, 1300 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014100404</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100410</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">5 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100501</td><td align="left">AA Black Crystal Cup, 3200 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014100502</td><td align="left">AA Millennium Challenge, ст-з, 1100 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014100503</td><td align="left">AA Presidents Handicap, ст-з, 1400 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014100504</td><td align="left">AA Gold Rush Stakes, 1600 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014100505</td><td align="left">AA United Nations Invitational, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014100506</td><td align="left">AA Ocean Pearl Stakes, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014100509</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100510</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100511</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100512</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100513</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100514</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100515</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100516</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100517</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">6 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100601</td><td align="left">Гр.III NYRA Sprint Cup, ст-з, 1000 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014100602</td><td align="left">Гр.II Clear Finish Stakes, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014100603</td><td align="left">Гр.I The Black Mesa, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014100607</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100608</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100610</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100611</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">7 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100701</td><td align="left">Гр.II White Lake Run, ст-з, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014100702</td><td align="left">Гр.II No Return Hurdle, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014100703</td><td align="left">Гр.I My Perfect Lady Stakes, 1000 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014100704</td><td align="left">Гр.I Честь Короны, 3400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014100706</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100707</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100708</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100709</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100711</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100712</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100713</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100714</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100715</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100716</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">8 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100801</td><td align="left">Гр.II Orchids Stakes, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014100802</td><td align="left">Гр.II Приз Белого Лотоса, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014100803</td><td align="left">Гр.II Возвращение Одиссея, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014100804</td><td align="left">Гр.I Bogota Stakes, ст-з, 2600, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014100805</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100806</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100809</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">9 октября 2014</td></tr>
<tr id="charter"><td width="100">2014100901</td><td align="left">Гр.III Бег Времени, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014100902</td><td align="left">Гр.II Magic Art Stakes, ст-з, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014100903</td><td align="left">Гр.I Limpopo Dash, ст-з, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014100904</td><td align="left">Гр.I Battle Front Spring Prix, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014100907</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100908</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100909</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100910</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100911</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100912</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100913</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100914</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014100916</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">10 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101001</td><td align="left">Гр.II Strawberry Time, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014101002</td><td align="left">Гр.II Country Ruffian Cup, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014101003</td><td align="left">Гр.I Note Bianko Stakes, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014101006</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101007</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101008</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101009</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101010</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101011</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101012</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101013</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101014</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">11 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101101</td><td align="left">Гр.III Tasotti Stakes, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014101102</td><td align="left">Гр.II Цветочные скачки, ст-з, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014101103</td><td align="left">Гр.II Приз Луары, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014101105</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101106</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101109</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101110</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101111</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101112</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101113</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101114</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101115</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">12 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101201</td><td align="left">AA Astor Plate, 1000 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101202</td><td align="left">AA Grand Prix de Saint-Cloud, 2000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101203</td><td align="left">AA New Ideal Handicap, торф, 1100 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101204</td><td align="left">AA Queen Cup Stakes, ст-з, 1600 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101205</td><td align="left">AA Triple Vivat Stakes, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101206</td><td align="left">AA Queen Victoria Stakes, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101207</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101211</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101212</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101213</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101215</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101216</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">13 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101301</td><td align="left">Гр.III Delta X Stakes, 1900 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014101302</td><td align="left">Гр.II Полотно Пенелопы, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014101303</td><td align="left">Гр.II The Hamilton Classic, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014101307</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101308</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101309</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101310</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101311</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101312</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">14 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101401</td><td align="left">Гр.III Emax Handicap, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014101402</td><td align="left">Гр.II Discovery Gardens Handicap, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014101403</td><td align="left">Гр.I Desert Stakes, ст-з, 1600 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014101406</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101407</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">15 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101501</td><td align="left">Гр.III Восточная Сказка, 1100 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014101502</td><td align="left">Гр.III Randevu Stakes, 4000 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014101503</td><td align="left">Гр.II Приз Матроны, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014101504</td><td align="left">Гр.I Triple Vivat Stakes, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014101507</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101509</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101510</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101511</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101512</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101513</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101514</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">16 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101601</td><td align="left">Гр.III Осенний круиз, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014101602</td><td align="left">Гр.II Vannil Flowers Sprint, ст-з, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014101603</td><td align="left">Гр.I Latin America, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014101606</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">17 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101701</td><td align="left">Гр.III Sweet Cacao Stakes, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014101702</td><td align="left">Гр.II Огонь Прометея, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014101703</td><td align="left">Гр.II Ladies First, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014101704</td><td align="left">Гр.I Приз Санкт-Петербурга, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014101707</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101708</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101709</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101710</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101711</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101712</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101713</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101714</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101715</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">18 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101801</td><td align="left">Гр.III Кубок Арианны, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014101802</td><td align="left">Гр.III Приз Пивовара, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014101803</td><td align="left">Гр.III Луна Памира, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014101806</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101807</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101809</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101810</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101811</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101812</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">19 октября 2014</td></tr>
<tr id="charter"><td width="100">2014101901</td><td align="left">AA Impresario Handicap, 1800 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101902</td><td align="left">AA Concord Stakes, 4000 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101903</td><td align="left">AA Maktub Run, ст-з, 1000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101904</td><td align="left">AA Queen Mary Stakes, торф, 1400 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101905</td><td align="left">AA Saratoga Stakes, ст-з, 1900 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014101908</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101910</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101911</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101912</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014101913</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">20 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102001</td><td align="left">Гр.III Великая Армада, 3200 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102002</td><td align="left">Гр.II Retry Sprint Stakes, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014102003</td><td align="left">Гр.I Bingo Ringo, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014102006</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102007</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102008</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102009</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102012</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102013</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102014</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102015</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102016</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102017</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">21 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102101</td><td align="left">Гр.III Faradei Stakes, 1300 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102102</td><td align="left">Гр.II No Country For Old Man Chase, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014102103</td><td align="left">Гр.I Приз Портоса, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014102104</td><td align="left">Гр.I Saami Cup, ст-з, 1700 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014102107</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102108</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102109</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102110</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102111</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102112</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102113</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102114</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102115</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102116</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">22 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102201</td><td align="left">Гр.III Young Talent Chase, ст-з, 1400 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102202</td><td align="left">Гр.I Prix de Diane, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014102203</td><td align="left">Гр.I Великолепный Век, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014102207</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102208</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102209</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102210</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102211</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102212</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102213</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102216</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102217</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102218</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102229</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">23 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102301</td><td align="left">Гр.II Lunoref Chase, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014102302</td><td align="left">Гр.I The Green Wasabi, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014102303</td><td align="left">Гр.I The Alley Cat Dash, ст-з, 1100 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014102304</td><td align="left">Гр.I Грозовой Перевал, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014102307</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102310</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102311</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">24 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102401</td><td align="left">Гр.III Old Forest Plate, 3600 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102402</td><td align="left">Гр.I Oxbow Cup, ст-з, 1300 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014102403</td><td align="left">Гр.I Creat Vor Dash, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014102406</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102407</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102408</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102410</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102411</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">25 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102501</td><td align="left">Гр.III Альпийские Луга, 2200 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102502</td><td align="left">Гр.III Xorital Stakes, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014102503</td><td align="left">Гр.III Алтайский Край, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102507</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102508</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102511</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102517</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102519</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">26 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102601</td><td align="left">AA Black Orchid Plate, 1200 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014102602</td><td align="left">AA Grand Reef Stakes, 1600 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014102603</td><td align="left">AA Givenchy Cup, 1700 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014102604</td><td align="left">AA Ladies First, 3000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014102605</td><td align="left">AA Золотое Кольцо, торф, 1800 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014102606</td><td align="left">AA Lis-de-Fleur Stakes, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014102609</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102610</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102611</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102612</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">27 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102701</td><td align="left">Гр.III Hurdle Stakes, ст-з, 1900 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102702</td><td align="left">Гр.III Quarl Ego Stakes, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014102703</td><td align="left">Гр.II Blueflower Stakes, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014102704</td><td align="left">Гр.II Фрекен Бок Приз, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014102707</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102708</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102712</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102713</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102714</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102715</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102716</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102717</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102718</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">28 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102801</td><td align="left">Гр.III Emerald Handicap, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102802</td><td align="left">Гр.II Prima Plate, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014102803</td><td align="left">Гр.II Приз Столичный, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014102807</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102808</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102809</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">29 октября 2014</td></tr>
<tr id="charter"><td width="100">2014102901</td><td align="left">Гр.III Twist N Shout Stakes, 4400 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014102902</td><td align="left">Гр.II Amanda Stakes, торф, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014102903</td><td align="left">Гр.I Adios Amigo Stakes, 2200 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014102904</td><td align="left">Гр.I Patrizia Stakes, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014102905</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102908</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014102909</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">30 октября 2014</td></tr>
<tr id="charter"><td width="100">2014103001</td><td align="left">Гр.II Осенний Марафон, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014103002</td><td align="left">Гр.II Holiday Cheer Plate, торф, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014103003</td><td align="left">Гр.II Вавилонские Воины, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014103004</td><td align="left">Гр.I Примадона Стипльчеза, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014103007</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103008</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103009</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103015</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103016</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">31 октября 2014</td></tr>
<tr id="charter"><td width="100">2014103101</td><td align="left">Гр.II Snow Blossom Stakes, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014103102</td><td align="left">Гр.II Savanna Stakes, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014103103</td><td align="left">Гр.II Морской Круиз, ст-з, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014103108</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103109</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103110</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014103111</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">1 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110101</td><td align="left">Гр.III Бег Спартанца, 3000 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014110102</td><td align="left">Гр.III Четыре Тысячи, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014110103</td><td align="left">Гр.III Golden Ring Plate, ст-з, 1300 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014110104</td><td align="left">Гр.I Note Bianko Stakes, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014110107</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110108</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110109</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">2 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110201</td><td align="left">TC Hopeful Stakes, 1200 м, 2yo Ж (2yoTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110202</td><td align="left">TC Spinaway Stakes, 1100 м, 2yo К (2yoFTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110203</td><td align="left">TC Lightning Stakes, 1000 м, 4+yo (SprTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110204</td><td align="left">TC Fullmoon Stakes, 1600 м, 4+yo (ClassTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110205</td><td align="left">TC Churchill Downs Oaks, 1800 м, 3yo К (FTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110206</td><td align="left">TC Churchill Downs Derby, 2200 м, 3yo Ж (TC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110207</td><td align="left">TC Fortune Stakes, 2200 м, 4yo (TrTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110208</td><td align="left">TC Kimidar Stakes, 4400 м, 4+yo (EndTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110209</td><td align="left">TC Green Grass Stakes, торф, 1200 м, 2yo Ж (2yoTTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110210</td><td align="left">TC Sierra Stakes, торф, 1100 м, 2yo К (2yoFTTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110211</td><td align="left">TC Salinger Chase, ст-з, 1200 м, 4+yo (MelbSCTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110212</td><td align="left">TC Filly Hurdle Challenge, ст-з, 1400 м, 3yo К (FSCTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110213</td><td align="left">TC Hurdle Challenge, ст-з, 1600 м, 3yo Ж (SCTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110214</td><td align="left">TC Metropol SC Handicap, ст-з, 2000 м, 4+yo (SCHTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110215</td><td align="left">TC Excellence Chase Stakes, ст-з, 2200 м, 4yo (TrSCTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110216</td><td align="left">TC Roland Chase, ст-з, 4800 м, 4+yo (SCEndTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110219</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110220</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110221</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110226</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110227</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110228</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110229</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110230</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110232</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">3 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110301</td><td align="left">Гр.III Rubin Stakes, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014110302</td><td align="left">Гр.II Premium Stakes, 1800 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014110303</td><td align="left">Гр.I Saboro Stakes, 4400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014110308</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110314</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">4 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110401</td><td align="left">Гр.II Brigitte Stakes, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014110402</td><td align="left">Гр.II Muerdo Stakes, ст-з, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014110403</td><td align="left">Гр.I The Sleeping Beauty, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014110409</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110410</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110411</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110413</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">5 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110501</td><td align="left">Гр.III Приз Ноябрьский, 1400 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014110502</td><td align="left">Гр.III Knickerbocker Handicap, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014110503</td><td align="left">Гр.II Old Ladies Handicap, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014110504</td><td align="left">Гр.I Set In Her Ways Distaff, 1900 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014110507</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110508</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110510</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110511</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110512</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110518</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110525</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">6 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110601</td><td align="left">Гр.III Сокровище Клана, 2600 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014110602</td><td align="left">Гр.II Приз Созвездие Льва, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014110603</td><td align="left">Гр.I Кубок Королевской Конницы, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014110606</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110621</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110622</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">7 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110701</td><td align="left">Гр.III Shine At Me Stakes, 1000 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014110702</td><td align="left">Гр.II Дорога к славе, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014110703</td><td align="left">Гр.I Geishas Song, 4000 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014110721</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110722</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110723</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110724</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110725</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110728</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110730</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">8 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110801</td><td align="left">Гр.III Blue Rose Plate, 1800 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014110802</td><td align="left">Гр.III Красота Сибири, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014110803</td><td align="left">Гр.II Sophia Stakes, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014110814</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110817</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110818</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110822</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110823</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110824</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110825</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014110828</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">9 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014110901</td><td align="left">AA Festival Stakes, 1600 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110902</td><td align="left">AA Great Storm, 1700 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110903</td><td align="left">AA Brilliants Cup, 3400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110904</td><td align="left">AA Pearl Sands, торф, 1200 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110905</td><td align="left">AA Золотая Лихорадка, торф, 1800 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014110906</td><td align="left">AA Old Gold Plate, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">10 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111001</td><td align="left">Гр.III Королевская Гавань, 1000 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014111002</td><td align="left">Гр.II Королева препятствий, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014111003</td><td align="left">Гр.I Северное Сияние, 2800 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111006</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">11 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111101</td><td align="left">Гр.III Over And Over Chase, ст-з, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014111102</td><td align="left">Гр.II Ocean Pearl Stakes, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014111103</td><td align="left">Гр.I Приз Мимозы, ст-з, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111104</td><td align="left">Гр.I Providencia Stakes, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111107</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111108</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111109</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111111</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111113</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111114</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">12 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111201</td><td align="left">Гр.III Arange Stakes, 4800 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014111202</td><td align="left">Гр.III Sarafan Run, ст-з, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014111203</td><td align="left">Гр.I Under World Stakes, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111204</td><td align="left">Гр.I Приз Русские Просторы, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111207</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111210</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">13 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111301</td><td align="left">Гр.II Oman Oasis Plate, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014111302</td><td align="left">Гр.I Морской Бриз, 1200 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111303</td><td align="left">Гр.I Каменный Гость, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111310</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111311</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111313</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111314</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">14 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111401</td><td align="left">Гр.III Приз Красноярска, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014111402</td><td align="left">Гр.II Подвиг Геракла, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014111403</td><td align="left">Гр.II Sessill Cup, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014111406</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111409</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">15 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111501</td><td align="left">Гр.II Machinehead Stakes, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014111502</td><td align="left">Гр.II Парижские Огни, ст-з, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014111503</td><td align="left">Гр.II Neon Night Run, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014111504</td><td align="left">Гр.I The Golden Flames Sprint, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111507</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111510</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111514</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111516</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">16 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111601</td><td align="left">AA Karoline Trials, 2600 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014111602</td><td align="left">AA Oxbow Cup, ст-з, 1200 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014111603</td><td align="left">AA Royal Mile Stakes, торф, 1600 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014111604</td><td align="left">AA Queens Stakes, ст-з, 1700 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014111605</td><td align="left">AA Grandness of Waterfall Cup, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014111609</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111612</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111615</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">17 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111701</td><td align="left">Гр.III Saroque Stakes, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014111702</td><td align="left">Гр.III Young Ladies Cup, ст-з, 1900 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014111703</td><td align="left">Гр.I Вечерняя Москва, ст-з, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111706</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111710</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111711</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111712</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">18 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111801</td><td align="left">Гр.I Ramada Stakes, 1600 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111802</td><td align="left">Гр.I Space Stakes, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014111803</td><td align="left">Гр.I Приз Императрицы, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111804</td><td align="left">Гр.I Scandinavia Pride, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111808</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111809</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111810</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111811</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">19 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014111901</td><td align="left">Гр.II Turbulent Chase, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014111902</td><td align="left">Гр.I Yung Yong Plate, 1800 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111903</td><td align="left">Гр.I Luna Plate, 3200 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014111912</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014111914</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">20 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112001</td><td align="left">Гр.II Земляничная Поляна, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014112002</td><td align="left">Гр.I Золотая Лихорадка, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014112003</td><td align="left">Гр.I Ladies First Stakes, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014112007</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112008</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112012</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112013</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112014</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112016</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112017</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112018</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112019</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112020</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112021</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">21 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112101</td><td align="left">Гр.III Loose Me Not Run, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014112102</td><td align="left">Гр.II Приз Двойная Звезда, ст-з, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014112103</td><td align="left">Гр.I In The Next Life, 2400 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014112106</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112107</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">22 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112201</td><td align="left">Гр.III Vanity Fair, 1700 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014112202</td><td align="left">Гр.III Kingab Mile, ст-з, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112203</td><td align="left">Гр.III Veramonda Plate, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112206</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112208</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112211</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">23 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112301</td><td align="left">AA Blue Rose Plate, 1200 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014112302</td><td align="left">AA Delta X Stakes, 1300 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014112303</td><td align="left">AA Grand Prix de Paris, 2000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014112304</td><td align="left">AA Granat Handicap, 4800 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014112305</td><td align="left">AA Prix Royal Oak, ст-з, 1400 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014112306</td><td align="left">AA Приз Российский, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014112315</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">24 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112401</td><td align="left">Гр.III Полуночная Сага, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112402</td><td align="left">Гр.III Восточная  Сладость, ст-з, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112403</td><td align="left">Гр.III Проделки Посейдона, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112404</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112407</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112408</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112410</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112411</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112412</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112413</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112414</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">25 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112501</td><td align="left">Гр.II Go-Baby-Go Stakes, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014112502</td><td align="left">Гр.II Nagoue Chase, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014112503</td><td align="left">Гр.I Great Shuttle Run, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014112506</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112508</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112509</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112510</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112511</td><td align="left">Тестовый класс, ст-з, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">26 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112601</td><td align="left">Гр.II Russian Song Cup, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014112602</td><td align="left">Гр.II Приз Августа, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014112603</td><td align="left">Гр.II Большой Круиз, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014112604</td><td align="left">Гр.I Aurora Sprint, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">27 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112701</td><td align="left">Гр.III Драконий Остров, ст-з, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112702</td><td align="left">Гр.I Collibri Cup, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014112703</td><td align="left">Гр.I Приз Фрезии, 3000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014112706</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112707</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">28 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112801</td><td align="left">Гр.III Греческая Мифология, 2400 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014112802</td><td align="left">Гр.III Приз Грация, ст-з, 1100 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014112803</td><td align="left">Гр.II Battle Time, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014112806</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">29 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014112901</td><td align="left">Гр.III December Nights Stakes, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014112902</td><td align="left">Гр.III Our Pride Stakes, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014112903</td><td align="left">Гр.I Crystal Plate, 1600 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014112904</td><td align="left">Гр.I Wolatur Cup, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014112908</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112909</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112910</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112911</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014112912</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">30 ноября 2014</td></tr>
<tr id="charter"><td width="100">2014113001</td><td align="left">AA Breeders Plate, 1200 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014113002</td><td align="left">AA Golden Globe Handicap, 1800 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014113003</td><td align="left">AA Black Hills Stakes, 3400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014113004</td><td align="left">AA Royal Whip Stakes, торф, 1600 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014113005</td><td align="left">AA The King Of the Pacific, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014113006</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014113014</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014113015</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">1 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120101</td><td align="left">Гр.III Malibu Stakes, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014120102</td><td align="left">Гр.II Retry Stakes, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014120103</td><td align="left">Гр.I Приз Сердце Льва, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">2 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120201</td><td align="left">Гр.II Post Factum Stakes, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014120202</td><td align="left">Гр.II Superstud Chase, ст-з, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014120203</td><td align="left">Гр.I Dagmars Cup, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014120204</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120205</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120207</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120209</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120210</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120212</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120213</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">3 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120301</td><td align="left">Гр.II Moon Lyric Stakes, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014120302</td><td align="left">Гр.II Mawingo Run, ст-з, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014120303</td><td align="left">Гр.II Poitvar Plate, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014120304</td><td align="left">Гр.I Quiet Harbor, 3600 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">4 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120401</td><td align="left">Гр.III Приз Баллисты, 3400 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014120402</td><td align="left">Гр.II Sunset Distaff, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014120403</td><td align="left">Гр.I Восточный Экспресс, ст-з, 1900 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014120406</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120407</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120410</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120411</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120412</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">5 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120501</td><td align="left">Гр.III Приз Рождественский, 1600 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014120502</td><td align="left">Гр.III Infinity Cup, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014120503</td><td align="left">Гр.II Buwarad Plate, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014120504</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120505</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120508</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120509</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120510</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120511</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120512</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120513</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120514</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">6 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120601</td><td align="left">Гр.III Forgive Me Stakes, 1200 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014120602</td><td align="left">Гр.III Лесной Кросс, ст-з, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014120603</td><td align="left">Гр.II San Antonio Stakes, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014120606</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120607</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120608</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120611</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">7 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120701</td><td align="left">AA Dream Supreme Plate, 1400 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014120702</td><td align="left">AA Goldencents Handicap, 1900 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014120703</td><td align="left">AA Mantoux SC Cup, ст-з, 1100 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014120704</td><td align="left">AA Rubin Stakes, ст-з, 1700 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014120705</td><td align="left">AA Триумф Цезаря, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014120708</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">8 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120801</td><td align="left">Гр.III Trifollet Cup, 1100 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014120802</td><td align="left">Гр.III Зимний Фестиваль, 3400 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014120803</td><td align="left">Гр.I Diadem Stakes, ст-з, 1900 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014120805</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120811</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120812</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014120813</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">9 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014120901</td><td align="left">Гр.III Приз Белорусский Зимний, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014120902</td><td align="left">Гр.II Приз Самоцвет, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014120903</td><td align="left">Гр.I Shadow Chaser, ст-з, 1200 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">10 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121001</td><td align="left">Гр.III Millennium Challenge, 2400 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014121002</td><td align="left">Гр.I Ares Cup, ст-з, 1300 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014121003</td><td align="left">Гр.I Bit O Honey Endurance, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014121010</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121011</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121012</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121013</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">11 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121101</td><td align="left">Гр.III Sundance Stakes, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014121102</td><td align="left">Гр.III Grandiatta Plate, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014121103</td><td align="left">Гр.I Trio Stakes, ст-з, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121104</td><td align="left">Гр.I Bursledon Run, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121105</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121111</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121112</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121113</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121114</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">12 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121201</td><td align="left">Гр.I Irish 2000 Guineas, 1600 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121202</td><td align="left">Гр.I Приз Тангоры, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014121203</td><td align="left">Гр.I Frentur Stakes, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121206</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121211</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121212</td><td align="left">Тестовый класс, ст-з, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">13 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121301</td><td align="left">Гр.III Великий Путь, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014121302</td><td align="left">Гр.III Sentinel Stakes, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014121303</td><td align="left">Гр.I Шаль Маргариты, ст-з, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121312</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">14 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121401</td><td align="left">AA Alpen Silver Cup, 1000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014121402</td><td align="left">AA Kingab Mile, 2800 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014121403</td><td align="left">AA Gold Cup, 4800 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014121404</td><td align="left">AA Plaza Stakes, торф, 1200 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014121405</td><td align="left">AA Приз Монако, торф, 1900 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014121408</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121409</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121412</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121413</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121415</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">15 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121501</td><td align="left">Гр.III Loraine Cup, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014121502</td><td align="left">Гр.III Приз Карнавал, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014121503</td><td align="left">Гр.III Дорога Гладиатора, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014121506</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121507</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">16 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121601</td><td align="left">Гр.II 5th Avenue Stakes, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014121602</td><td align="left">Гр.II Caramel Wizard, ст-з, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014121603</td><td align="left">Гр.I Yamarta Plate, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121611</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121612</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121613</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">17 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121701</td><td align="left">Гр.II Приз Акации, ст-з, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014121702</td><td align="left">Гр.II Приз Кремлевские Звезды, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014121703</td><td align="left">Гр.I Приз Амазонии, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121704</td><td align="left">Гр.I Приз Воли к Победе, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014121707</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121709</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121710</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121711</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121715</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121716</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">18 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121801</td><td align="left">Гр.II Queen Victoria Stakes, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014121802</td><td align="left">Гр.II Long Island Chase, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014121803</td><td align="left">Гр.I Karavai Chase, ст-з, 1200 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014121809</td><td align="left">Тестовый класс, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121811</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121812</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121813</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">19 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014121901</td><td align="left">Гр.II Black Orchid Plate, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014121902</td><td align="left">Гр.II Ruby Stone Stakes, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014121903</td><td align="left">Гр.II Грозовой Фронт, ст-з, 1800 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014121906</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121907</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014121911</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">20 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122001</td><td align="left">Гр.III Огни Москвы, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014122002</td><td align="left">Гр.III Weekend Prize, 4400 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122003</td><td align="left">Гр.III Silver Ring Plate, ст-з, 1300 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122007</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122009</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">21 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122101</td><td align="left">TC Oakleigh Plate, 1100 м, 4+yo (SprTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122102</td><td align="left">TC Matron Stakes, 1300 м, 2yo К (2yoFTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122103</td><td align="left">TC Futurity Stakes, 1600 м, 2yo Ж (2yoTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122104</td><td align="left">TC Pimlico Oaks, 1600 м, 3yo К (FTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122105</td><td align="left">TC Azur Plate, 1800 м, 4+yo (ClassTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122106</td><td align="left">TC Pimlico Derby, 1900 м, 3yo Ж (TC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122107</td><td align="left">TC Nakheel Cup, 2000 м, 3yo (TrTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122108</td><td align="left">TC Albatros Plate, 4000 м, 4+yo (EndTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122109</td><td align="left">TC Manikato Chase, ст-з, 1000 м, 4+yo (MelbSCTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122110</td><td align="left">TC Pyramisa Stakes, торф, 1300 м, 2yo К (2yoFTTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122111</td><td align="left">TC Ritz Stakes, торф, 1600 м, 2yo Ж (2yoTTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122112</td><td align="left">TC Filly Hurdle Challenge, ст-з, 1800 м, 3yo К (FSCTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122113</td><td align="left">TC Hurdle Challenge, ст-з, 1900 м, 3yo Ж (SCTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122114</td><td align="left">TC Grand Chase Cup, ст-з, 2000 м, 3yo (TrSCTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122115</td><td align="left">TC Brooklyn SC Handicap, ст-з, 2400 м, 4+yo (SCHTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122116</td><td align="left">TC Ambassador Chase, ст-з, 4000 м, 4+yo (SCEndTC-2) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122120</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122121</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122122</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122123</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122124</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">22 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122201</td><td align="left">Гр.III Duke of Edinburgh Stakes, ст-з, 1200 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122202</td><td align="left">Гр.III Легенды Старого Леса, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122203</td><td align="left">Гр.II Северная Звезда, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014122208</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122209</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">23 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122301</td><td align="left">Гр.III Самый Стойкий, 1700 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122302</td><td align="left">Гр.II Mirovik Stakes, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014122303</td><td align="left">Гр.I Grand Prix Du Mans, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014122310</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">24 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122401</td><td align="left">Гр.II Приз Орбиты, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014122402</td><td align="left">Гр.I Givenchy Cup, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014122403</td><td align="left">Гр.I Sands Still Cup, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014122408</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">25 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122501</td><td align="left">Гр.III Sunday Stakes, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014122502</td><td align="left">Гр.I Pine Lane Stakes, 3200 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014122503</td><td align="left">Гр.I Black Onyx Stakes, ст-з, 1600, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014122504</td><td align="left">Гр.I Great Taxis Cup, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014122509</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122510</td><td align="left">Тестовый класс, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122511</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">26 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122601</td><td align="left">Гр.III Lamodar Stakes, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014122602</td><td align="left">Гр.II Краса Страны, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014122603</td><td align="left">Гр.I Lis-de-Fleur Stakes, ст-з, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014122610</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122612</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">27 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122701</td><td align="left">Гр.III Приз Констанции, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014122702</td><td align="left">Гр.III Ледяное Сердце, ст-з, 1800 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122703</td><td align="left">Гр.III Каникулы в Париже, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014122711</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122713</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122715</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122716</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122719</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">28 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122801</td><td align="left">AA Buckinghams Cup, 1100 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122802</td><td align="left">AA Kingarvie Stakes, 1900 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122803</td><td align="left">AA Parlament Stakes, ст-з, 1300 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122804</td><td align="left">AA Super Ninety Nine Run, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2014122809</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122811</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122813</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">29 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014122901</td><td align="left">Гр.III Liam Stakes, 3200 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014122902</td><td align="left">Гр.II Первоцвет, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2014122903</td><td align="left">Гр.II Angellic Chase, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2014122904</td><td align="left">Гр.I Венок Дриады, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014122907</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2014122910</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">30 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014123001</td><td align="left">Гр.III Долгая Дорога в Дюнах, 2800 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014123002</td><td align="left">Гр.III Camelia Plate, ст-з, 1000 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2014123003</td><td align="left">Гр.I Приз Тройки, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2014123010</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">31 декабря 2014</td></tr>
<tr id="charter"><td width="100">2014123101</td><td align="left">Гр.III Приз Вальса, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014123102</td><td align="left">Гр.III Приз Фемиды, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2014123103</td><td align="left">Гр.I White Warrior, 1100 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2014123106</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">1 января 2015</td></tr>
<tr id="charter"><td width="100">2015010101</td><td align="left">Гр.III Хозяин Гор, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015010102</td><td align="left">Гр.II Приз Авроры, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015010103</td><td align="left">Гр.I Be Fashionable Stakes, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015010104</td><td align="left">Гр.I Надежда сезона, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015010109</td><td align="left">Тестовый класс, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">2 января 2015</td></tr>
<tr id="charter"><td width="100">2015010201</td><td align="left">Гр.I Морская Пена, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015010202</td><td align="left">Гр.I Seven Seas, 4800 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015010203</td><td align="left">Гр.I Кубок Аванпост, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015010206</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010207</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010210</td><td align="left">Тестовый класс, ст-з, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010211</td><td align="left">Тестовый класс, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010215</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">3 января 2015</td></tr>
<tr id="charter"><td width="100">2015010301</td><td align="left">Гр.III Стрела Аполлона, 4000 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015010302</td><td align="left">Гр.III Танец Наяды, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015010303</td><td align="left">Гр.II Destiny Flame Stakes, 3000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015010306</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">4 января 2015</td></tr>
<tr id="charter"><td width="100">2015010401</td><td align="left">AA Celebration Stakes, 1300 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015010402</td><td align="left">AA Legacy Cup, 1900 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015010403</td><td align="left">AA Classic Trinity Plate, 3600 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015010404</td><td align="left">AA Luna Plate, торф, 1000 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015010405</td><td align="left">AA Secret Word Stakes, ст-з, 1900 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015010406</td><td align="left">AA The Melbourne Cup, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015010410</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010412</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">5 января 2015</td></tr>
<tr id="charter"><td width="100">2015010501</td><td align="left">Гр.III Rose on the Rain Stakes, 2000 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015010502</td><td align="left">Гр.III Sandy Beach Chase, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015010503</td><td align="left">Гр.II Agava Cup, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015010504</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010505</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010506</td><td align="left">Тестовый класс, 3200 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">6 января 2015</td></tr>
<tr id="charter"><td width="100">2015010601</td><td align="left">Гр.II Приз Аквамарин, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015010602</td><td align="left">Гр.II Горная Песня, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015010603</td><td align="left">Гр.I Brightwood Stakes, 1100 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015010605</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">7 января 2015</td></tr>
<tr id="charter"><td width="100">2015010701</td><td align="left">Гр.III Xeronix Stakes, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015010702</td><td align="left">Гр.II Abraj Stakes, ст-з, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015010703</td><td align="left">Гр.II Приз Мануфактуры, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015010704</td><td align="left">Гр.I France Plate, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015010707</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010708</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">8 января 2015</td></tr>
<tr id="charter"><td width="100">2015010801</td><td align="left">Гр.III Лунное Сокровище, 2600 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015010802</td><td align="left">Гр.III Unique Stakes, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015010803</td><td align="left">Гр.II Almas Sprint, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015010805</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010809</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010810</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">9 января 2015</td></tr>
<tr id="charter"><td width="100">2015010901</td><td align="left">Гр.III Приз Фестивальный, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015010922</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010923</td><td align="left">Тестовый класс, ст-з, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010924</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015010925</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">10 января 2015</td></tr>
<tr id="charter"><td width="100">2015011001</td><td align="left">Гр.III Queen Cup Stakes, 1600 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015011002</td><td align="left">Гр.III Кубок Атлантис, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015011003</td><td align="left">Гр.II Jureart Stakes, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015011004</td><td align="left">Тестовый класс, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">11 января 2015</td></tr>
<tr id="charter"><td width="100">2015011101</td><td align="left">AA Energy Star Preview, 1400 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011102</td><td align="left">AA Hollywood Handicap, 2200 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011103</td><td align="left">AA Ice King Stakes, 2400 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011104</td><td align="left">AA PR Premier Stakes, ст-з, 1300 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011105</td><td align="left">AA Great Jump Chase, ст-з, 3200 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011109</td><td align="left">Тестовый класс, 3400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011115</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011116</td><td align="left">Тестовый класс, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">12 января 2015</td></tr>
<tr id="charter"><td width="100">2015011201</td><td align="left">Гр.III Богемское Стекло, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015011202</td><td align="left">Гр.II Королева Бала, ст-з, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015011203</td><td align="left">Гр.I Apronto Cup, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015011204</td><td align="left">Гр.I Mariinsky Theatre, ст-з, 4000 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015011210</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011211</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">13 января 2015</td></tr>
<tr id="charter"><td width="100">2015011301</td><td align="left">Гр.II Omega Handicap, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015011302</td><td align="left">Гр.II Great Storm, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015011303</td><td align="left">Гр.I Milan Distaff, ст-з, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015011308</td><td align="left">Тестовый класс, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011310</td><td align="left">Тестовый класс, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">14 января 2015</td></tr>
<tr id="charter"><td width="100">2015011401</td><td align="left">Гр.III Улыбка Джоконды, ст-з, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015011402</td><td align="left">Гр.II Lander Prize, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015011403</td><td align="left">Гр.I Blue Frost Prize, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015011405</td><td align="left">Тестовый класс, ст-з, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011411</td><td align="left">Тестовый класс, ст-з, 1300 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011412</td><td align="left">Тестовый класс, 1200 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011413</td><td align="left">Тестовый класс, ст-з, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">15 января 2015</td></tr>
<tr id="charter"><td width="100">2015011501</td><td align="left">Гр.III Рассвет в Сахаре, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015011502</td><td align="left">Гр.II Opera House Classic, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015011503</td><td align="left">Гр.I Arktika Stakes, 1400 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015011504</td><td align="left">Гр.I Oxford Endurance, 4800 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015011505</td><td align="left">Тестовый класс, ст-з, 3000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011507</td><td align="left">Тестовый класс, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011508</td><td align="left">Тестовый класс, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011509</td><td align="left">Тестовый класс, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011510</td><td align="left">Тестовый класс, ст-з, 1000 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">16 января 2015</td></tr>
<tr id="charter"><td width="100">2015011601</td><td align="left">Гр.III Maple Leaf Stakes, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015011602</td><td align="left">Гр.III Пражский Сувенир, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015011603</td><td align="left">Гр.II Oriental Bazzar, ст-з, 1800 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015011604</td><td align="left">Тестовый класс, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011605</td><td align="left">Тестовый класс, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011606</td><td align="left">Тестовый класс, 2800 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011607</td><td align="left">Тестовый класс, ст-з, 2000 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011608</td><td align="left">Тестовый класс, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011609</td><td align="left">Тестовый класс, 1100 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011615</td><td align="left">Тестовый класс, 3600 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011619</td><td align="left">Тестовый класс, ст-з, 4400 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011620</td><td align="left">Тестовый класс, ст-з, 2400 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011621</td><td align="left">Тестовый класс, ст-з, 1700 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011622</td><td align="left">Тестовый класс, ст-з, 1400 м, 2+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011623</td><td align="left">Тестовый класс, ст-з, 1800 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">17 января 2015</td></tr>
<tr id="charter"><td width="100">2015011701</td><td align="left">Гр.III Black And White Lilly Run, 2800 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015011702</td><td align="left">Гр.II Far-Far-Away Stakes, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015011703</td><td align="left">Гр.I Aegean Paradise, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015011711</td><td align="left">Тестовый класс, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">18 января 2015</td></tr>
<tr id="charter"><td width="100">2015011801</td><td align="left">AA Carlsberg Cup, 1200 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011802</td><td align="left">AA International Stakes, 2600 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011803</td><td align="left">AA Crystal Plate, 4000 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011804</td><td align="left">AA Prix De l Atrium, торф, 1400 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011805</td><td align="left">AA The Golden Flames Sprint, торф, 1700 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015011809</td><td align="left">Тестовый класс, ст-з, 2600 м, 3+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011810</td><td align="left">Тестовый класс, ст-з, 2200 м, 3+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">19 января 2015</td></tr>
<tr id="charter"><td width="100">2015011901</td><td align="left">Гр.I Magia Cup, 2000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015011902</td><td align="left">Гр.I Rochester Stakes, ст-з, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015011903</td><td align="left">Гр.I Phoenix Stakes, ст-з, 3000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015011904</td><td align="left">Тестовый класс, ст-з, 4000 м, 4+yo <font color="#666" size="1"></font></td></tr><tr id="charter"><td width="100">2015011905</td><td align="left">Тестовый класс, ст-з, 1900 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">20 января 2015</td></tr>
<tr id="charter"><td width="100">2015012001</td><td align="left">Гр.III Ballyogan Stakes, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015012002</td><td align="left">Гр.II Asfara Stakes, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015012003</td><td align="left">Гр.I Приз Итаки, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015012004</td><td align="left">Гр.I Battle Front Cup, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015012006</td><td align="left">Тестовый класс, ст-з, 1600 м, 2+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">21 января 2015</td></tr>
<tr id="charter"><td width="100">2015012101</td><td align="left">Гр.II Quel Esprit Stakes, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015012102</td><td align="left">Гр.II Daffodil Garden Stakes, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015012103</td><td align="left">Гр.I Fiery Chase Dash, ст-з, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">22 января 2015</td></tr>
<tr id="charter"><td width="100">2015012201</td><td align="left">Гр.III Florida Distaff Chase, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015012202</td><td align="left">Гр.II Emma Plate, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015012203</td><td align="left">Гр.II Tango Mahogany Chase, ст-з, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015012204</td><td align="left">Гр.I Белорусская Весна, 2600 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015012205</td><td align="left">Тестовый класс, ст-з, 4800 м, 4+yo <font color="#666" size="1"></font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">23 января 2015</td></tr>
<tr id="charter"><td width="100">2015012301</td><td align="left">Гр.III Afrodita Stakes, 3600 м, 4+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015012302</td><td align="left">Гр.III Zakat Chase, ст-з, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015012303</td><td align="left">Гр.III Песня Коралловых Рифов, ст-з, 1800 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">24 января 2015</td></tr>
<tr id="charter"><td width="100">2015012401</td><td align="left">Гр.III Кубок Шангрила, 2000 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015012402</td><td align="left">Гр.III Rainbow Stakes, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015012403</td><td align="left">Гр.II Wellfarm Stakes, ст-з, 1300 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015012404</td><td align="left">Гр.II Приз Субботы, ст-з, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">25 января 2015</td></tr>
<tr id="charter"><td width="100">2015012501</td><td align="left">AA Arabian Cup, 1000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015012502</td><td align="left">AA Caesar Stakes, 1300 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015012503</td><td align="left">AA Gold Medal Stakes, 1800 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015012504</td><td align="left">AA The Backdraft Special, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015012505</td><td align="left">AA Ideal Finish Chase, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">26 января 2015</td></tr>
<tr id="charter"><td width="100">2015012601</td><td align="left">Гр.III Ellington Stakes, 1400 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015012602</td><td align="left">Гр.III The Great Wall Stakes, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015012603</td><td align="left">Гр.II Love Affaire Chase, ст-з, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">27 января 2015</td></tr>
<tr id="charter"><td width="100">2015012701</td><td align="left">Гр.II Приз Ритмики, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015012702</td><td align="left">Гр.II Приз Монако, ст-з, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015012703</td><td align="left">Гр.I Hot Springs Dash, 1200 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015012704</td><td align="left">Гр.I Grandness of Waterfall Cup, ст-з, 4000 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">28 января 2015</td></tr>
<tr id="charter"><td width="100">2015012801</td><td align="left">Гр.III The Cherry Blossoms, ст-з, 1600 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015012802</td><td align="left">Гр.I Greek Sweets, 1800 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015012803</td><td align="left">Гр.I Moment In The Sun Plate, 3200 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">29 января 2015</td></tr>
<tr id="charter"><td width="100">2015012901</td><td align="left">Гр.III Приз Закат, 4800 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015012902</td><td align="left">Гр.II Sherwood Stakes, ст-з, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015012903</td><td align="left">Гр.II Terra Cotta Chase, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015012904</td><td align="left">Гр.I Veren Stakes, 2200 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">30 января 2015</td></tr>
<tr id="charter"><td width="100">2015013001</td><td align="left">Гр.II Berliner Luft Stakes, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015013002</td><td align="left">Гр.II La Scala Stakes, ст-з, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015013003</td><td align="left">Гр.I Sonata Chase, ст-з, 2000 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">31 января 2015</td></tr>
<tr id="charter"><td width="100">2015013101</td><td align="left">Гр.III Advantage Stakes, 1900 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015013102</td><td align="left">Гр.II Приз Цирцеи, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015013103</td><td align="left">Гр.I Приз Жасмина, ст-з, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">1 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020101</td><td align="left">Гр.I Шамаханская Царица, 1200 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015020102</td><td align="left">AA Blue Diamond Prelude, 1100 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020103</td><td align="left">AA Alexander The Great Stakes, 3200 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020104</td><td align="left">AA Red Rum Crown, ст-з, 1600 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020105</td><td align="left">AA Violla Handicap, ст-з, 2600 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">2 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020201</td><td align="left">Гр.III Ranz Hurdle, ст-з, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015020202</td><td align="left">Гр.II Приз Кирова, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015020203</td><td align="left">Гр.II Кладрубский стипльчез, ст-з, 3600 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">3 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020301</td><td align="left">Гр.III Castania Stakes, 4800 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015020302</td><td align="left">Гр.I Snow Flowers Cup, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015020303</td><td align="left">Гр.I Oleandrs Dreams, 2400 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015020304</td><td align="left">Гр.I Приз Азалии, ст-з, 2200 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">4 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020401</td><td align="left">Гр.II Грандфортский Порт, ст-з, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015020402</td><td align="left">Гр.I Karakao Plate, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015020403</td><td align="left">Гр.I Santa Margarita Stakes, 3600 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">5 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020501</td><td align="left">Гр.III Приз Гарема, ст-з, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015020502</td><td align="left">Гр.II Adagio Cup, 2600 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015020503</td><td align="left">Гр.I Pinkwater Chase, ст-з, 1400 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015020504</td><td align="left">Гр.I Stars Final, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">6 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020601</td><td align="left">Гр.III Beauty Fashion Stakes, ст-з, 1800 м, 2+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015020602</td><td align="left">Гр.III Заморские Дары, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015020603</td><td align="left">Гр.II Песня Кружевницы, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015020604</td><td align="left">Гр.II Saint Are Cup, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">7 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020701</td><td align="left">Гр.III Хозяин Морей, 1800 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015020702</td><td align="left">Гр.II Три Русалки, ст-з, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015020703</td><td align="left">Гр.II Приз Миланды, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">8 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020801</td><td align="left">TC Newmarket Handicap, 1200 м, 4+yo (SprTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020802</td><td align="left">TC Frizette Stakes, 1700 м, 2yo К (2yoFTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020803</td><td align="left">TC Champagne Stakes, 1800 м, 2yo Ж (2yoTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020804</td><td align="left">TC Damas Plate, 1800 м, 2yo (TrTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020805</td><td align="left">TC Belmont Oaks, 2000 м, 3yo К (FTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020806</td><td align="left">TC Black Horse Handicap, 2000 м, 4+yo (ClassTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020807</td><td align="left">TC Belmont Derby, 2600 м, 3yo Ж (TC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020808</td><td align="left">TC Garranah Handicap, 4800 м, 4+yo (EndTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020809</td><td align="left">TC Schillaci Chase, ст-з, 1100 м, 4+yo (MelbSCTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020810</td><td align="left">TC Jasmine Stakes, торф, 1700 м, 2yo К (2yoFTTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020811</td><td align="left">TC Radisson Stakes, торф, 1800 м, 2yo Ж (2yoTTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020812</td><td align="left">TC Power Turf Plate, торф, 1800 м, 2yo (TrSCTC-1) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020813</td><td align="left">TC Suburban SC Handicap, ст-з, 1800 м, 4+yo (SCHTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020814</td><td align="left">TC Filly Hurdle Challenge, ст-з, 2200 м, 3yo К (FSCTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020815</td><td align="left">TC Hurdle Challenge, ст-з, 2400 м, 3yo Ж (SCTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015020816</td><td align="left">TC Cascades Chase, ст-з, 4400 м, 4+yo (SCEndTC-3) <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">9 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015020901</td><td align="left">Гр.II Приз Черная Жемчужина, 1000 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015020902</td><td align="left">Гр.II Indigo Chase, ст-з, 1700 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015020903</td><td align="left">Гр.I Аленький Цветочек, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">10 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021001</td><td align="left">Гр.III Млечный Путь, 3600 м, 4+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015021002</td><td align="left">Гр.II Rasta Pasta Run, ст-з, 2400 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015021003</td><td align="left">Гр.I Приз Каприз, ст-з, 1300 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">11 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021101</td><td align="left">Гр.II Приз Темзы, 1100 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015021102</td><td align="left">Гр.II Кубок Авангард, ст-з, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015021103</td><td align="left">Гр.I Glassico Stakes, 1900 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015021104</td><td align="left">Гр.I Cottage Acre Run, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">12 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021201</td><td align="left">Гр.III Букет Белых Лилий, 2200 м, 3+yo К <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015021202</td><td align="left">Гр.I The Metropolitan, 1000 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015021203</td><td align="left">Гр.I Золотое Кольцо России, 3400 м, 4+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">13 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021301</td><td align="left">Гр.III Сокровища Могола, 2800 м, 3+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015021302</td><td align="left">Гр.II Кубок Нельсона, 1400 м, 2+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015021303</td><td align="left">Гр.I Снежная Королева, ст-з, 1800 м, 2+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">14 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021401</td><td align="left">Гр.III Aventura Stakes, ст-з, 1100 м, 2+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015021402</td><td align="left">Гр.III The Mister Flyer, ст-з, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015021403</td><td align="left">Гр.III Whittington Plate, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015021404</td><td align="left">Гр.II Royal Whip Stakes, 2200 м, 3+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">15 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021501</td><td align="left">AA Delight Stakes, 1100 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015021502</td><td align="left">AA Glamour Handicap, 1700 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015021503</td><td align="left">AA Dubizzle Stakes, 4400 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015021504</td><td align="left">AA Prix de l Arc de Triomphe, торф, 1300 м, 2yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015021505</td><td align="left">AA Seven Kings Cup, ст-з, 2000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015021506</td><td align="left">AA Vatican Treasures Handicap, ст-з, 2600 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">16 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021601</td><td align="left">Гр.III Алые Паруса, 2000 м, 3+yo Ж <font color="#666" size="1"> (максимальный класс резвости F и ниже)</font></td></tr><tr id="charter"><td width="100">2015021602</td><td align="left">Гр.I Sir Villington Stakes, 1400 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015021603</td><td align="left">Гр.I End Of Season Trophy, 3400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">17 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021701</td><td align="left">Гр.II Мечта Поэта, ст-з, 1400 м, 2+yo К <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015021702</td><td align="left">Гр.I Alamir Plate, 1800 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr><tr id="charter"><td width="100">2015021703</td><td align="left">Гр.I Приз Тайги, ст-з, 2800 м, 3+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">18 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021801</td><td align="left">Гр.II Приз Калина, 3200 м, 4+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015021802</td><td align="left">Гр.II The Birth Of Venus, ст-з, 1900 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015021803</td><td align="left">Гр.II Романтическое Увлечение, ст-з, 4800 м, 4+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015021804</td><td align="left">Гр.I Приз Тюльпана, ст-з, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">19 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015021901</td><td align="left">Гр.II Приз Руны, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015021902</td><td align="left">Гр.I Дорогой Длинною, 4400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015021903</td><td align="left">Гр.I Сладкий Вкус Победы, ст-з, 1400 м, 2+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015021904</td><td align="left">Гр.I Новый Век, ст-з, 3200 м, 4+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">20 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015022001</td><td align="left">Гр.II Серебряная Подкова, 1000 м, 2+yo Ж <font color="#666" size="1"> (класс резвости C)</font></td></tr><tr id="charter"><td width="100">2015022002</td><td align="left">Гр.I Rensai Cup, 3000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости A и выше)</font></td></tr><tr id="charter"><td width="100">2015022003</td><td align="left">Гр.I Утренний Бриз, ст-з, 1900 м, 2+yo Ж <font color="#666" size="1"> (класс резвости B)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">21 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015022101</td><td align="left">Гр.III Приз Каролины, 1600 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015022102</td><td align="left">Гр.III Lotus Aroma Cup, ст-з, 1700 м, 2+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015022103</td><td align="left">Гр.III Villa Verde Stakes, ст-з, 3000 м, 3+yo К <font color="#666" size="1"> (класс резвости E)</font></td></tr><tr id="charter"><td width="100">2015022104</td><td align="left">Гр.II Приз Дружбы Народов, ст-з, 4400 м, 4+yo Ж <font color="#666" size="1"> (класс резвости D)</font></td></tr></tbody></table>
<table id="infoblocktbl"><tbody><tr id="infoblockheader"><td colspan="2">22 февраля 2015</td></tr>
<tr id="charter"><td width="100">2015022201</td><td align="left">Гр.II Queen Mary Stakes, 1200 м, 2+yo К <font color="#666" size="1"> (класс резвости D)</font></td></tr><tr id="charter"><td width="100">2015022202</td><td align="left">AA Diadem Stakes, 1400 м, 3+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015022203</td><td align="left">AA Great Shuttle Run, 1700 м, 2yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015022204</td><td align="left">AA Lander Prize, 3000 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015022205</td><td align="left">AA Russian Derby, ст-з, 1800 м, 3+yo Ж <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr><tr id="charter"><td width="100">2015022206</td><td align="left">AA Prix de Diane, ст-з, 4400 м, 4+yo К <font color="#666" size="1"> (минимальный класс резвости D и выше)</font></td></tr></tbody></table>`

	return s
}


func getHorses() string {
	horses := `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>"ЭклипС" :: Виртуальные скачки</title>

<!-- CSS styles -->
<link href="pagelayout.css" rel="stylesheet" type="text/css">

<!-- Vertical Main Menu Bar -->
<link href="SpryAssets/SpryMenuBarVertical.css" rel="stylesheet" type="text/css">
<script src="SpryAssets/SpryMenuBar.js" type="text/javascript"></script>

<!-- JavaScripts -->
<!-- JQuery -->
<script type="text/javascript" src="js/jquery-1.10.2.min.js"></script>
<script type="text/javascript" src="js/jquery-migrate-1.2.1.min.js"></script>


<script type="text/javascript">
function GetXmlHttpObject()
{
if (window.XMLHttpRequest)
  {
  // code for IE7+, Firefox, Chrome, Opera, Safari
  return new XMLHttpRequest();
  }
if (window.ActiveXObject)
  {
  // code for IE6, IE5
  return new ActiveXObject("Microsoft.XMLHTTP");
  }
return null;
}

/////////////////////////////////
var xmlhttp10;

function changerest(hrrid,ropt)
{
xmlhttp10=GetXmlHttpObject();
if (xmlhttp10==null)
  {
  alert ("Browser does not support HTTP Request");
  return;
  }
var url="inc/change_restoption.php";
url=url+"?q="+hrrid+"&newopt="+ropt;
url=url+"&sid="+Math.random();
xmlhttp10.onreadystatechange=stateChanged10;
xmlhttp10.open("GET",url,true);
xmlhttp10.send(null);
}

function stateChanged10()
{
if (xmlhttp10.readyState==4)
{
document.getElementById("txtHint"+hrrid).innerHTML=xmlhttp10.responseText;
}
}

///////////
var xmlhttp1;

function moveHorseTo(bldgid,horseid,depID)
{
xmlhttp1=GetXmlHttpObject();
if (xmlhttp1==null)
  {
  alert ("Browser does not support HTTP Request");
  return;
  }
divID="txtHint"+horseid;
var url="inc/stbldgchange.php";
url=url+"?q="+bldgid+"&h="+horseid+"&d="+depID;
url=url+"&sid="+Math.random();
xmlhttp1.onreadystatechange=stateChanged1;
xmlhttp1.open("GET",url,true);
xmlhttp1.send(null);
}

function stateChanged1()
{
if (xmlhttp1.readyState==4)
{
document.getElementById(divID).innerHTML=xmlhttp1.responseText;
}
}
</script>

</head><body>
<div id="pageblock">
  <div id="header"><img src="img/logo.png" width="168" height="100">
</div>
  <div id="loginboard">
  <div style="display: inline; margin: 0 400px 0 0; text-aling: left"><b>Конюшня: "Penumbra" [Рег.Номер: 194] </b> &#8226; <a href="logout.php" class="menulink">Выход</a></div> <div style="display: inline; margin-right: 0; text-align: right">
 Онлайн: 4  &#8226; <span id="servertime"></span>
  <script type="text/javascript">
var currenttime = 'October 06, 2014 08:24:33' //PHP method of getting server date

///////////Stop editting here/////////////////////////////////

var montharray=new Array("января","февраля","марта","апреля","мая","июня","июля","августа","сентября","октября","ноября","декабря")
var serverdate=new Date(currenttime)

function padlength(what){
var output=(what.toString().length==1)? "0"+what : what
return output
}

function displaytime(){
serverdate.setSeconds(serverdate.getSeconds()+1)
var datestring=padlength(serverdate.getDate())+" "+montharray[serverdate.getMonth()]+" "+serverdate.getFullYear()
var timestring="&#8226; "+padlength(serverdate.getHours())+":"+padlength(serverdate.getMinutes())+":"+padlength(serverdate.getSeconds())
document.getElementById("servertime").innerHTML=datestring+"  "+timestring
}

window.onload=function(){
setInterval("displaytime()", 1000)
}

</script>
</div>
</div>
<div id="page">  <div id="menublock">
  <!-- Main menu -->
  <ul id="MenuBar1" class="MenuBarVertical">
  <li><a href="news.php">Новости</a>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Конюшня</a>
    <ul>
      <li><a href="stableinfo.php">&#9830; Общая информация</a></li>
      <li><a href="stableinfo_studs.php">&#9830; Жеребцы-производители</a></li>
      <li><a href="stableinfo_broods.php">&#9830; Племенные кобылы</a></li>
      <li><a href="stableinfo_racers.php">&#9830; Скаковое отделение</a></li>
      <li><a href="stableinfo_foals.php">&#9830; Племенной молодняк</a></li>
      <li><a href="stableinfo_retired.php">&#9830; Рабочий состав</a></li>
      <li><a href="stableinfo_rent.php">&#9830; Аренда</a></li>
      <li><a href="stableinfo_rip.php">&#9830; Мемориал</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Ипподром</a>
      <ul>
      <li><a href="track_info.php">&#9830; Ипподромы</a></li>
      <li><a href="races_schedule.php">&#9830; Расписание скачек</a></li>
      <li><a href="races_all_order.php">&#9830; Заказ скачек</a></li>
      <li><a href="races_entries.php">&#9830; Прием заявок</a></li>
      <li><a href="removefromrace.php">&#9830; Снять заявку</a></li>
      <li><a href="races_preracecards.php">&#9830; Предварительная карточка скачек</a></li>
      <li><a href="races_results.php">&#9830; Результаты скачек</a></li>
      <li><a href="raceresults_search.php">&#9830; Поиск результатов скачек</a></li>
      <li><a href="races_order.php">&#9830; Заказ субботних скачек</a></li>
      <li><a class="MenuBarItemSubmenu" href="#">Спецмитинги</a>
		<ul>
      <li><a class="MenuBarItemSubmenu" href="#">Кубки "Пегаса"</a>
		<ul>
      <li><a href="race_pegascupsinfo.php">&#9830; Общая информация</a></li>
      <li><a href="pg_races_entries.php">&#9830; Регистрация</a></li>
      <li><a href="pg_races_results.php">&#9830; Результаты</a></li>
		</ul></li>
		</ul>
	  </li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Букмекер</a>
      <ul>
      <li><a href="toto_bets.php">&#9830; Ставки</a></li>
      <li><a href="toto_results.php">&#9830; Результаты</a></li>
      <li><a href="toto_statistics.php">&#9830; Статистика</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Тренинг</a>
      <ul>
      <li><a href="trening_normal.php">&#9830; Тренинг</a></li>
      <li><a href="trening_group.php">&#9830; Результаты</a></li>
      <li><a href="trening.php">&#9830; Тренинг (тестовый)</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Банк</a>
      <ul>
      <li><a href="bank.php">&#9830; Общая информация</a></li>
      <li><a href="bank_transfer.php">&#9830; Переводы</a></li>
      <li><a href="bankreport.php">&#9830; Выписка со счета</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Офис</a>
      <ul>
      <li><a class="MenuBarItemSubmenu" href="#">Смена владельца</a>
		<ul>
		<li><a href="horsesales.php">&#9830; Продажа лошади</a></li>
		<li><a href="horserents.php">&#9830; Аренда лошади</a></li>
		</ul>
	  </li>
	  <li><a class="MenuBarItemSubmenu" href="#">Райский Сад</a>
		<ul>
		<li><a href="horsesales_rs.php">&#9830; Передача лошади</a></li>
		</ul>
	  </li>
      <li><a class="MenuBarItemSubmenu" href="#">Регистрат</a>
		<ul>
		<li><a href="changestatus_racer.php">&#9830; Регистрация скаковой лошади</a></li>
		<li><a href="changestatus_stud.php">&#9830; Регистрация жеребца-производителя</a></li>
		<li><a href="changestatus_private_stud.php">&#9830; Регистрация частного жеребца-производителя</a></li>
		<li><a href="changestatus_brood.php">&#9830; Регистрация племкобылы</a></li>
       <li><a href="changestatus_retired.php">&#9830; Перевод в рабочий состав</a></li>
	   <li><a href="changestatus_horsename.php">&#9830; Изменение клички лошади</a></li>
	  </ul>
	  </li>
	  <li><a href="horsegeld.php">&#9830; Кастрация жеребца</a></li>
	  <li><a href="changestatus_racersctofl.php">&#9830; Перевод из стипльчеза в гладкие</a></li>
	  <li><a class="MenuBarItemSubmenu" href="#">Строительная контора</a>
		<ul>
		<li><a href="st_build.php">&#9830; Прораб</a></li>
		<li><a href="st_manage.php">&#9830; Управляющий</a></li>
		</ul>
	  </li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Центр заводчиков</a>
      <ul>
      <li><a href="studs_reiting.php">&#9830; Рейтинги ЖП</a></li>
      <li><a href="studs_catalog.php">&#9830; Каталог ЖП</a></li>
      <li><a href="embrio.php">&#9830; ЦИО "Эмбрио"</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Аукционы</a>
      <ul>
		<li><a class="MenuBarItemSubmenu" href="#">Сезонные</a>
			<ul>
				<li><a href="auction.php">&#9830; Каталог аукциона</a></li>
				<li><a href="auctionlotsubmit.php">&#9830; Регистрация лотов</a></li>
				<li><a href="auction_archive.php">&#9830; Архивы</a></li>
			</ul>
		</li>
        <li><a class="MenuBarItemSubmenu" href="#">Частные</a>
			<ul>
				<li><a href="privauction.php">&#9830; Каталог аукциона</a></li>
				<li><a href="privauctionlotsubmit.php">&#9830; Регистрация лотов</a></li>
				<li><a href="privauction_archive.php">&#9830; Архивы</a></li>
			</ul>
		</li>
		<li><a href="zkauction.php">&#9830; Распродажа</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Справочная</a>
      <ul>
      <li><a href="visit_stable.php">&#9830; Визитка конюшни</a></li>
      <li><a href="info_horsecard.php">&#9830; Карточка лошади</a></li>
      <li><a href="info_horsename.php">&#9830; Проверка клички лошади</a></li>
      <li><a href="race_best10.php">&#9830; Лучшие лошади</a></li>
      <li><a href="info_speedclass.php">&#9830; Класс резвости</a></li>
      <li><a href="race_records.php">&#9830; Рекорды</a></li>
      <li><a href="triplecrowns.php">&#9830; "Тройные Короны"</a></li>
    </ul>
  </li>
    <li><a href="horseimport.php">Импорт лошади</a>
    </li>
  <li><a class="MenuBarItemSubmenu" href="#">Чемпионат</a>
      <ul>
      <li><a href="championship_points.php">&#9830; Баллы</a></li>
      <li><a href="championship_entries.php">&#9830; Регистрация</a></li>
    </ul>
  </li>
      <li><a href="edit_profil.php"><br>Профиль<br><br></a>
    </li>

</ul>
<script type="text/javascript">
<!--
var MenuBar1 = new Spry.Widget.MenuBar("MenuBar1", {imgRight:"SpryAssets/SpryMenuBarRightHover.gif"});
//-->
</script>
 </div>
</div><!-- Начало основного блока -->
<div id="mainblock">
<div id="infoblock">
<span id="infoblocktitle">Отделение Скаковых Лошадей</span>

<table id='infoblocktbl'>
<tr><td width=100 rowspan='3'><img src='img/racestall.jpg'></td>
<td id='textright'>Вместимость:</td><td id='textleft'>50</td>
<td width=200 rowspan='3'><center><font size=3 color='C00E2D'>Важно!</font><br>Каждая лошадь должна находиться в своем деннике. Лошади в леваде не могут быть проданы, записаны на скачки или случку к ЖП</center></td></tr>
<tr id='charter'><td id='textright'>Занято:</td><td id='textleft'>21</td></tr><tr id='charter'><td id='textright'>Свободных денников:</td><td id='textleft'>29</td></tr></table><table id='infoblocktbl'><tr align='center'><td><br>
<form id='chooseBuildNo' name='chooseBuildNo' method='post' action='stableinfo_racers.php'>
Отделение: <select id='stdepid' name='stdepid'>
<option value=0 style='color: #ccc'> ... выбираем из списка ... </option><option value='590'>Отделение 1</option><option value='674'>Отделение 2</option>
</select>
<input name="submitStID" type="submit" value=" показать отделение ">
</form>
<br><br>
</td></tr>
</table>

<table id='infoblocktbl'>
<tr id='infoblockheader'><td colspan=8>№ 590 Отделение 1</td></tr>
<tr id='infoblockheader'><td>№ и Кличка</td><td>Пол</td><td>Возраст</td><td>Спец.</td><td>Заявлена в скачку</td><td>Состояние</td><td>ДПС</td><td>Отд</td></tr><tr id='charter'><td align='right'><a id='2564'> </a><a href='horsecard.php?horseid=2564' class='menulink'>2564 Карча</a><br><font size=1 color=#666>(105 Чинар - 1135 Картахена)</font></td><td>Ж</td><td>7</td><td>fl</td><td>2014100914</td><td>100</td><td>--</td>
<td>
<span id='rest2564'>
<input type=CHECKBOX name='restlist2564' value='1'   onchange='changerest(2564,this.value)'>
</span>
<div id='txtHint2564'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='2977'> </a><a href='horsecard.php?horseid=2977' class='menulink'>2977 Tiger</a><br><font size=1 color=#666>(19 Гранат - 1223 Точка)</font></td><td>Ж</td><td>7</td><td>fl</td><td>--</td><td>100</td><td>0</td>
<td>
<span id='rest2977'>
<input type=CHECKBOX name='restlist2977' value='1'   onchange='changerest(2977,this.value)'>
</span>
<div id='txtHint2977'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='3036'> </a><a href='horsecard.php?horseid=3036' class='menulink'>3036 Варяг</a><br><font size=1 color=#666>(19 Гранат - 1219 Вариация)</font></td><td>Ж</td><td>7</td><td>fl</td><td>--</td><td>100</td><td>0</td>
<td>
<span id='rest3036'>
<input type=CHECKBOX name='restlist3036' value='1'   onchange='changerest(3036,this.value)'>
</span>
<div id='txtHint3036'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='2979'> </a><a href='horsecard.php?horseid=2979' class='menulink'>2979 Амстердам</a><br><font size=1 color=#666>(67 Минор - 1070 Алыча)</font></td><td>Ж</td><td>7</td><td>sc</td><td>2014100806</td><td>100</td><td>115</td>
<td>
<span id='rest2979'>
<input type=CHECKBOX name='restlist2979' value='1'   onchange='changerest(2979,this.value)'>
</span>
<div id='txtHint2979'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='4157'> </a><a href='horsecard.php?horseid=4157' class='menulink'>4157 Кампари</a><br><font size=1 color=#666>(885 Метрополь - 712 Керамика)</font></td><td>Ж</td><td>5</td><td>sc</td><td>--</td><td>100</td><td>170</td>
<td>
<span id='rest4157'>
<input type=CHECKBOX name='restlist4157' value='1'   onchange='changerest(4157,this.value)'>
</span>
<div id='txtHint4157'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='4583'> </a><a href='horsecard.php?horseid=4583' class='menulink'>4583 Boon Of Battle</a><br><font size=1 color=#666>(792 Battle Front - 271 Блокада)</font></td><td>Ж</td><td>5</td><td>sc</td><td>--</td><td>96</td><td>3</td>
<td>
<span id='rest4583'>
<input type=CHECKBOX name='restlist4583' value='1'   onchange='changerest(4583,this.value)'>
</span>
<div id='txtHint4583'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='4611'> </a><a href='horsecard.php?horseid=4611' class='menulink'>4611 Gemini Trend</a><br><font size=1 color=#666>(607 Gladiator (Ire) - 264 Трактовка)</font></td><td>Ж</td><td>5</td><td>sc</td><td>--</td><td>100</td><td>128</td>
<td>
<span id='rest4611'>
<input type=CHECKBOX name='restlist4611' value='1'   onchange='changerest(4611,this.value)'>
</span>
<div id='txtHint4611'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='5345'> </a><a href='horsecard.php?horseid=5345' class='menulink'>5345 Параллакс</a><br><font size=1 color=#666>(587 Лидер II - 1281 Пара-Пит)</font></td><td>Ж</td><td>4</td><td>fl</td><td>--</td><td>100</td><td>118</td>
<td>
<span id='rest5345'>
<input type=CHECKBOX name='restlist5345' value='1'   onchange='changerest(5345,this.value)'>
</span>
<div id='txtHint5345'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='5467'> </a><a href='horsecard.php?horseid=5467' class='menulink'>5467 Фамилия</a><br><font size=1 color=#666>(2489 Fatal Treasure - 2067 Lady Perfect)</font></td><td>К</td><td>4</td><td>fl</td><td>2014100710</td><td>100</td><td>121</td>
<td>
<span id='rest5467'>
<input type=CHECKBOX name='restlist5467' value='1'   onchange='changerest(5467,this.value)'>
</span>
<div id='txtHint5467'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='5625'> </a><a href='horsecard.php?horseid=5625' class='menulink'>5625 Кадавр</a><br><font size=1 color=#666>(927 Diplomat - 956 Celtic Queen)</font></td><td>Ж</td><td>4</td><td>fl</td><td>2014100713</td><td>100</td><td>134</td>
<td>
<span id='rest5625'>
<input type=CHECKBOX name='restlist5625' value='1'   onchange='changerest(5625,this.value)'>
</span>
<div id='txtHint5625'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='5655'> </a><a href='horsecard.php?horseid=5655' class='menulink'>5655 Джигси</a><br><font size=1 color=#666>(484 Грануш - 438 Дачница)</font></td><td>К</td><td>4</td><td>fl</td><td>--</td><td>100</td><td>127</td>
<td>
<span id='rest5655'>
<input type=CHECKBOX name='restlist5655' value='1'   onchange='changerest(5655,this.value)'>
</span>
<div id='txtHint5655'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='4932'> </a><a href='horsecard.php?horseid=4932' class='menulink'>4932 Дайкири</a><br><font size=1 color=#666>(582 Доминик - 1362 Germiona)</font></td><td>К</td><td>4</td><td>sc</td><td>--</td><td>100</td><td>118</td>
<td>
<span id='rest4932'>
<input type=CHECKBOX name='restlist4932' value='1'   onchange='changerest(4932,this.value)'>
</span>
<div id='txtHint4932'></div>
</td>
</tr><tr id='charter'><td align='right'><a id='5314'> </a><a href='horsecard.php?horseid=5314' class='menulink'>5314 Интуиция</a><br><font size=1 color=#666>(34 Тюльпан - 700 Ирония)</font></td><td>К</td><td>4</td><td>sc</td><td>2014100808</td><td>100</td><td>113</td>
<td>
<span id='rest5314'>
<input type=CHECKBOX name='restlist5314' value='1'   onchange='changerest(5314,this.value)'>
</span>
<div id='txtHint5314'></div>
</td>
</tr></table>

</div>
<!-- End of основного блока -->
</div>
<div id="footer">
<br />
<a href="http://www.eklps.com" class="menulink">Игровой проект "ЭклипС"</a> &copy; 2003-2014<br />
Вся информация на данном сайте является частной собственностью и защищена законом <br />
<div style="padding: 40px 0px 20px 0px;">
<span>
		<!--Akavita counter start-->
<script type="text/javascript">var AC_ID=55830;var AC_TR=false;
(function(){var l='http://adlik.akavita.com/acode.js'; var t='text/javascript';
try {var h=document.getElementsByTagName('head')[0];
var s=document.createElement('script'); s.src=l;s.type=t;h.appendChild(s);}catch(e){
document.write(unescape('%3Cscript src="'+l+'" type="'+t+'"%3E%3C/script%3E'));}})();
</script><span id="AC_Image"></span>
<noscript><a target='_blank' href='http://www.akavita.by/'>
<img src='http://adlik.akavita.com/bin/lik?id=55830&it=1'
border='0' height='1' width='1' alt='Akavita'/>
</a></noscript>
<!--Akavita counter end-->


&nbsp;&nbsp;&nbsp;
<!-- HotLog -->
<script type="text/javascript">
var hotlog_counter_id = 2319785;
var hotlog_hit = 25;
var hotlog_counter_type = 565;
</script>
<script src="http://js.hotlog.ru/counter.js" type="text/javascript"></script>
<noscript>
<a href="http://click.hotlog.ru/?2319785" target="_blank">
<img src="http://hit25.hotlog.ru/cgi-bin/hotlog/count?s=2319785&im=565" border="0"
title="HotLog" alt="HotLog"></a>
</noscript>
<!-- /HotLog -->
&nbsp;&nbsp;&nbsp;


<!--Openstat-->
<span id="openstat586556"></span>
<script type="text/javascript">
var openstat = { counter: 586556, image: 87, color: "ff9822", next: openstat, track_links: "all" };
(function(d, t, p) {
var j = d.createElement(t); j.async = true; j.type = "text/javascript";
j.src = ("https:" == p ? "https:" : "http:") + "//openstat.net/cnt.js";
var s = d.getElementsByTagName(t)[0]; s.parentNode.insertBefore(j, s);
})(document, "script", document.location.protocol);
</script>
<!--/Openstat-->

<!--Google-->
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-43673244-1', 'eklps.com');
  ga('send', 'pageview');

</script>
<!--//Google-->

</span>
		</div>
</div>
</div>
</body>
</html>`

	return horses
}


func getRaces() string {

	races := `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>"ЭклипС" :: Виртуальные скачки</title>

<!-- CSS styles -->
<link href="pagelayout.css" rel="stylesheet" type="text/css">

<!-- Vertical Main Menu Bar -->
<link href="SpryAssets/SpryMenuBarVertical.css" rel="stylesheet" type="text/css">
<script src="SpryAssets/SpryMenuBar.js" type="text/javascript"></script>

<!-- Validate the Input -->
<script src="SpryAssets/SpryValidationTextField.js" type="text/javascript"></script>
<link href="SpryAssets/SpryValidationTextField.css" rel="stylesheet" type="text/css">

<!-- JavaScripts -->

<script type="text/javascript">
function GetXmlHttpObject()
{
if (window.XMLHttpRequest)
  {
  // code for IE7+, Firefox, Chrome, Opera, Safari
  return new XMLHttpRequest();
  }
if (window.ActiveXObject)
  {
  // code for IE6, IE5
  return new ActiveXObject("Microsoft.XMLHTTP");
  }
return null;
}

/////////////////////////////////
var xmlhttp7;

function raceentrydelete(racecodeid,horseid,stableid)
{
xmlhttp7=GetXmlHttpObject();
if (xmlhttp7==null)
  {
  alert ("Browser does not support HTTP Request");
  return;
  }
var url="inc/deleteraceentry.php";
url=url+"?q="+racecodeid+"&&q2="+horseid+"&&q3="+stableid;
url=url+"&sid="+Math.random();
xmlhttp7.onreadystatechange=stateChanged7;
xmlhttp7.open("GET",url,true);
xmlhttp7.send(null);
}
function stateChanged7()
{
if (xmlhttp7.readyState==4)
{
document.getElementById("txtHint"+racecodeid+horseid).innerHTML=xmlhttp7.responseText;
}
}
</script>

</head><body>
<div id="pageblock">
  <div id="header"><img src="img/logo.png" width="168" height="100">
</div>
  <div id="loginboard">
  <div style="display: inline; margin: 0 400px 0 0; text-aling: left"><b>Конюшня: "Penumbra" [Рег.Номер: 194] </b> &#8226; <a href="logout.php" class="menulink">Выход</a></div> <div style="display: inline; margin-right: 0; text-align: right">
 Онлайн: 4  &#8226; <span id="servertime"></span>
  <script type="text/javascript">
var currenttime = 'October 06, 2014 01:45:28' //PHP method of getting server date

///////////Stop editting here/////////////////////////////////

var montharray=new Array("января","февраля","марта","апреля","мая","июня","июля","августа","сентября","октября","ноября","декабря")
var serverdate=new Date(currenttime)

function padlength(what){
var output=(what.toString().length==1)? "0"+what : what
return output
}

function displaytime(){
serverdate.setSeconds(serverdate.getSeconds()+1)
var datestring=padlength(serverdate.getDate())+" "+montharray[serverdate.getMonth()]+" "+serverdate.getFullYear()
var timestring="&#8226; "+padlength(serverdate.getHours())+":"+padlength(serverdate.getMinutes())+":"+padlength(serverdate.getSeconds())
document.getElementById("servertime").innerHTML=datestring+"  "+timestring
}

window.onload=function(){
setInterval("displaytime()", 1000)
}

</script>
</div>
</div>
<div id="page">  <div id="menublock">
  <!-- Main menu -->
  <ul id="MenuBar1" class="MenuBarVertical">
  <li><a href="news.php">Новости</a>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Конюшня</a>
    <ul>
      <li><a href="stableinfo.php">&#9830; Общая информация</a></li>
      <li><a href="stableinfo_studs.php">&#9830; Жеребцы-производители</a></li>
      <li><a href="stableinfo_broods.php">&#9830; Племенные кобылы</a></li>
      <li><a href="stableinfo_racers.php">&#9830; Скаковое отделение</a></li>
      <li><a href="stableinfo_foals.php">&#9830; Племенной молодняк</a></li>
      <li><a href="stableinfo_retired.php">&#9830; Рабочий состав</a></li>
      <li><a href="stableinfo_rent.php">&#9830; Аренда</a></li>
      <li><a href="stableinfo_rip.php">&#9830; Мемориал</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Ипподром</a>
      <ul>
      <li><a href="track_info.php">&#9830; Ипподромы</a></li>
      <li><a href="races_schedule.php">&#9830; Расписание скачек</a></li>
      <li><a href="races_all_order.php">&#9830; Заказ скачек</a></li>
      <li><a href="races_entries.php">&#9830; Прием заявок</a></li>
      <li><a href="removefromrace.php">&#9830; Снять заявку</a></li>
      <li><a href="races_preracecards.php">&#9830; Предварительная карточка скачек</a></li>
      <li><a href="races_results.php">&#9830; Результаты скачек</a></li>
      <li><a href="raceresults_search.php">&#9830; Поиск результатов скачек</a></li>
      <li><a href="races_order.php">&#9830; Заказ субботних скачек</a></li>
      <li><a class="MenuBarItemSubmenu" href="#">Спецмитинги</a>
		<ul>
      <li><a class="MenuBarItemSubmenu" href="#">Кубки "Пегаса"</a>
		<ul>
      <li><a href="race_pegascupsinfo.php">&#9830; Общая информация</a></li>
      <li><a href="pg_races_entries.php">&#9830; Регистрация</a></li>
      <li><a href="pg_races_results.php">&#9830; Результаты</a></li>
		</ul></li>
		</ul>
	  </li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Букмекер</a>
      <ul>
      <li><a href="toto_bets.php">&#9830; Ставки</a></li>
      <li><a href="toto_results.php">&#9830; Результаты</a></li>
      <li><a href="toto_statistics.php">&#9830; Статистика</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Тренинг</a>
      <ul>
      <li><a href="trening_normal.php">&#9830; Тренинг</a></li>
      <li><a href="trening_group.php">&#9830; Результаты</a></li>
      <li><a href="trening.php">&#9830; Тренинг (тестовый)</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Банк</a>
      <ul>
      <li><a href="bank.php">&#9830; Общая информация</a></li>
      <li><a href="bank_transfer.php">&#9830; Переводы</a></li>
      <li><a href="bankreport.php">&#9830; Выписка со счета</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Офис</a>
      <ul>
      <li><a class="MenuBarItemSubmenu" href="#">Смена владельца</a>
		<ul>
		<li><a href="horsesales.php">&#9830; Продажа лошади</a></li>
		<li><a href="horserents.php">&#9830; Аренда лошади</a></li>
		</ul>
	  </li>
	  <li><a class="MenuBarItemSubmenu" href="#">Райский Сад</a>
		<ul>
		<li><a href="horsesales_rs.php">&#9830; Передача лошади</a></li>
		</ul>
	  </li>
      <li><a class="MenuBarItemSubmenu" href="#">Регистрат</a>
		<ul>
		<li><a href="changestatus_racer.php">&#9830; Регистрация скаковой лошади</a></li>
		<li><a href="changestatus_stud.php">&#9830; Регистрация жеребца-производителя</a></li>
		<li><a href="changestatus_private_stud.php">&#9830; Регистрация частного жеребца-производителя</a></li>
		<li><a href="changestatus_brood.php">&#9830; Регистрация племкобылы</a></li>
       <li><a href="changestatus_retired.php">&#9830; Перевод в рабочий состав</a></li>
	   <li><a href="changestatus_horsename.php">&#9830; Изменение клички лошади</a></li>
	  </ul>
	  </li>
	  <li><a href="horsegeld.php">&#9830; Кастрация жеребца</a></li>
	  <li><a href="changestatus_racersctofl.php">&#9830; Перевод из стипльчеза в гладкие</a></li>
	  <li><a class="MenuBarItemSubmenu" href="#">Строительная контора</a>
		<ul>
		<li><a href="st_build.php">&#9830; Прораб</a></li>
		<li><a href="st_manage.php">&#9830; Управляющий</a></li>
		</ul>
	  </li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Центр заводчиков</a>
      <ul>
      <li><a href="studs_reiting.php">&#9830; Рейтинги ЖП</a></li>
      <li><a href="studs_catalog.php">&#9830; Каталог ЖП</a></li>
      <li><a href="embrio.php">&#9830; ЦИО "Эмбрио"</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Аукционы</a>
      <ul>
		<li><a class="MenuBarItemSubmenu" href="#">Сезонные</a>
			<ul>
				<li><a href="auction.php">&#9830; Каталог аукциона</a></li>
				<li><a href="auctionlotsubmit.php">&#9830; Регистрация лотов</a></li>
				<li><a href="auction_archive.php">&#9830; Архивы</a></li>
			</ul>
		</li>
        <li><a class="MenuBarItemSubmenu" href="#">Частные</a>
			<ul>
				<li><a href="privauction.php">&#9830; Каталог аукциона</a></li>
				<li><a href="privauctionlotsubmit.php">&#9830; Регистрация лотов</a></li>
				<li><a href="privauction_archive.php">&#9830; Архивы</a></li>
			</ul>
		</li>
		<li><a href="zkauction.php">&#9830; Распродажа</a></li>
    </ul>
  </li>
  <li><a class="MenuBarItemSubmenu" href="#">Справочная</a>
      <ul>
      <li><a href="visit_stable.php">&#9830; Визитка конюшни</a></li>
      <li><a href="info_horsecard.php">&#9830; Карточка лошади</a></li>
      <li><a href="info_horsename.php">&#9830; Проверка клички лошади</a></li>
      <li><a href="race_best10.php">&#9830; Лучшие лошади</a></li>
      <li><a href="info_speedclass.php">&#9830; Класс резвости</a></li>
      <li><a href="race_records.php">&#9830; Рекорды</a></li>
      <li><a href="triplecrowns.php">&#9830; "Тройные Короны"</a></li>
    </ul>
  </li>
    <li><a href="horseimport.php">Импорт лошади</a>
    </li>
  <li><a class="MenuBarItemSubmenu" href="#">Чемпионат</a>
      <ul>
      <li><a href="championship_points.php">&#9830; Баллы</a></li>
      <li><a href="championship_entries.php">&#9830; Регистрация</a></li>
    </ul>
  </li>
      <li><a href="edit_profil.php"><br>Профиль<br><br></a>
    </li>

</ul>
<script type="text/javascript">
<!--
var MenuBar1 = new Spry.Widget.MenuBar("MenuBar1", {imgRight:"SpryAssets/SpryMenuBarRightHover.gif"});
//-->
</script>
 </div>
</div><!-- Начало основного блока -->
<div id="mainblock">
<div id="infoblock">
<span id="infoblocktitle">Прием заявок</span>


<table id='infoblocktbl'>
<tr id='charter'><td><br>
<font color='red'>Внимание!</font> Регистрация на скачки заканчивается за ДВА ДНЯ до даты проведения.
<br><br>
</td></tr>
</table>
<table id='infoblocktbl'><tr align='center'><td><br>

Первым делом необходимо выбрать здание, в котором стоят скаковые лошади, доступные для регистрации на скачки:

<br><br>
<form id='chooseBuildNo' name='chooseBuildNo' method='post' action='races_entries.php'>
Отделение: <select id='stdepid' name='stdepid'>
<option value=0 style='color: #ccc'> ... выбираем из списка ... </option><option value='590'>Отделение 1</option><option value='674'>Отделение 2</option>
</select>
<input name="submitStID" type="submit" value=" выбрать отделение ">
</form>
<br><br>
</td></tr>
</table>

<table id='infoblocktbl'>
<tr id='infoblockheader'><td>Поиск скачек</td></tr>
<tr id='charter'><td><br>
<form id='racesearch' name='racesearch' method='post' action='races_entries.php'>
Класс: <select name='racetype'><option value=0>&nbsp;</option>
<option value=1>Тестовый класс</option><option value=2>Медный класс</option><option value=3>Бронзовый класс</option><option value=4>Серебряный класс</option><option value=5>Золотой класс</option><option value=6>FNX Гр.III</option><option value=7>FNX Гр.II</option><option value=8>FNX Гр.I</option><option value=10>Гр.III</option><option value=11>Гр.II</option><option value=12>Гр.I</option><option value=200>AA</option><option value=300>DWC</option><option value=500>BRC</option><option value=700>TC</option><option value=900>CH</option>
</select><br><br>
Дистанция:  <select name='distance2'><option value=0>&nbsp;</option>
<option value='40'>1000 м</option><option value='44'>1100 м</option><option value='48'>1200 м</option><option value='52'>1300 м</option><option value='56'>1400 м</option><option value='64'>1600 м</option><option value='68'>1700 м</option><option value='72'>1800 м</option><option value='76'>1900 м</option><option value='80'>2000 м</option><option value='88'>2200 м</option><option value='96'>2400 м</option><option value='104'>2600 м</option><option value='112'>2800 м</option><option value='120'>3000 м</option><option value='128'>3200 м</option><option value='136'>3400 м</option><option value='144'>3600 м</option><option value='160'>4000 м</option><option value='176'>4400 м</option><option value='192'>4800 м</option>
</select>
или  от <select name='distance3'><option value=0>&nbsp;</option>
<option value='40'>1000 м</option><option value='44'>1100 м</option><option value='48'>1200 м</option><option value='52'>1300 м</option><option value='56'>1400 м</option><option value='64'>1600 м</option><option value='68'>1700 м</option><option value='72'>1800 м</option><option value='76'>1900 м</option><option value='80'>2000 м</option><option value='88'>2200 м</option><option value='96'>2400 м</option><option value='104'>2600 м</option><option value='112'>2800 м</option><option value='120'>3000 м</option><option value='128'>3200 м</option><option value='136'>3400 м</option><option value='144'>3600 м</option><option value='160'>4000 м</option><option value='176'>4400 м</option><option value='192'>4800 м</option>
</select>
до <select name='distance4'><option value=0>&nbsp;</option>
<option value='40'>1000 м</option><option value='44'>1100 м</option><option value='48'>1200 м</option><option value='52'>1300 м</option><option value='56'>1400 м</option><option value='64'>1600 м</option><option value='68'>1700 м</option><option value='72'>1800 м</option><option value='76'>1900 м</option><option value='80'>2000 м</option><option value='88'>2200 м</option><option value='96'>2400 м</option><option value='104'>2600 м</option><option value='112'>2800 м</option><option value='120'>3000 м</option><option value='128'>3200 м</option><option value='136'>3400 м</option><option value='144'>3600 м</option><option value='160'>4000 м</option><option value='176'>4400 м</option><option value='192'>4800 м</option>
</select><br><br>

Специализация:<br>
<input type='radio' id='track_pref_0' value='0' name='sc'> гладкие скачки<br>
<input type='radio' id='track_pref_1' value='1' name='sc'> стипльчез / торф<br><br>

Возраст: <select name='agecriteria'>
<option value='0'></option>
<option value='2'> 2х лет </option>
<option value='2up'> 2х лет и старше </option>
<option value='3'> 3х лет </option>
<option value='3up'> 3х лет и старше </option>
<option value='4'> 4х лет </option>
<option value='4up'> 4х лет и старше </option>
</select><br><br>

<p><input type='radio' id='track_pref_0' value='3' name='sex' /> общая<br />
<input type='radio' id='track_pref_1' value='2' name='sex' /> только кобылы<br />
<input type='radio' id='track_pref_2' value='1' name='sex' /> только жеребцы</p>

<p><input name='submit' type='submit' value='поиск'></p>
</form>
</td></tr>
</table>
<table id='infoblocktbl'>
<tr id='infoblockheader'><td>9 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100913'>2014100913</a></td><td align='left'><strong>Тестовый класс, 1800 м, 2+yo</strong></td><td>Заявлено: 24</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100913>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100909'>2014100909</a></td><td align='left'><strong>Тестовый класс, 1900 м, 2+yo</strong></td><td>Заявлено: 29</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100909>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100912'>2014100912</a></td><td align='left'><strong>Тестовый класс, 2800 м, 3+yo</strong></td><td>Заявлено: 8</td></tr><div id='txtHint20141009126757'><tr id='charter'><td colspan='2' align='left'>6757 <strong>Нельсон</strong><font size=1>, Ж 3,  (402 Ночник - 435 Незабудка), Penumbra, 0 0-0-0 $0</font></td><td><br><form id='cancellentry' name='cancellentry' method='post' action='races_entries.php#race2014100912'><input type='image' src='img/cancel_but.jpg' alt='Отменить заявку' onclick='raceentrydelete(2014100912,6757,194)'></form></td></tr></div><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100912>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100914'>2014100914</a></td><td align='left'><strong>Тестовый класс, 3200 м, 4+yo</strong></td><td>Заявлено: 9</td></tr><div id='txtHint20141009142564'><tr id='charter'><td colspan='2' align='left'>2564 <strong>Карча</strong><font size=1>, Ж 7,  (105 Чинар - 1135 Картахена), Penumbra, 23 1-3-3 $54,275</font></td><td><br><form id='cancellentry' name='cancellentry' method='post' action='races_entries.php#race2014100914'><input type='image' src='img/cancel_but.jpg' alt='Отменить заявку' onclick='raceentrydelete(2014100914,2564,194)'></form></td></tr></div><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100914>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100907'>2014100907</a></td><td align='left'><strong>Тестовый класс, ст-з, 1000 м, 2+yo</strong></td><td>Заявлено: 16</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100907>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100911'>2014100911</a></td><td align='left'><strong>Тестовый класс, ст-з, 1100 м, 2+yo</strong></td><td>Заявлено: 13</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100911>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100916'>2014100916</a></td><td align='left'><strong>Тестовый класс, ст-з, 1800 м, 2+yo</strong></td><td>Заявлено: 6</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100916>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100908'>2014100908</a></td><td align='left'><strong>Тестовый класс, ст-з, 2600 м, 3+yo</strong></td><td>Заявлено: 5</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100908>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100910'>2014100910</a></td><td align='left'><strong>Тестовый класс, ст-з, 3200 м, 4+yo</strong></td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100910>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer_rs.php' method='post'>
<input type='hidden' name='racecode' value=2014100910>
<input type='submit' value='из РайСада'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100901'>2014100901</a></td><td align='left'><strong>Гр.III Бег Времени, 3200 м, 4+yo Ж</strong> (класс резвости E)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100901>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer_rs.php' method='post'>
<input type='hidden' name='racecode' value=2014100901>
<input type='submit' value='из РайСада'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100902'>2014100902</a></td><td align='left'><strong>Гр.II Magic Art Stakes, ст-з, 1700 м, 2+yo К</strong> (класс резвости D)</td><td>Заявлено: 12</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100902>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100903'>2014100903</a></td><td align='left'><strong>Гр.I Limpopo Dash, ст-з, 1100 м, 2+yo К</strong> (класс резвости B)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100903>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer_rs.php' method='post'>
<input type='hidden' name='racecode' value=2014100903>
<input type='submit' value='из РайСада'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014100904'>2014100904</a></td><td align='left'><strong>Гр.I Battle Front Spring Prix, ст-з, 4400 м, 4+yo Ж</strong> (класс резвости B)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100904>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer_rs.php' method='post'>
<input type='hidden' name='racecode' value=2014100904>
<input type='submit' value='из РайСада'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>9 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014100915'>2014100915</a></td><td align='left'><strong>Бронзовый класс, 1400 м, 2yo</strong> (максимальный класс резвости F и ниже)</td><td>Заявлено: 5</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100915>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014100905'>2014100905</a></td><td align='left'><strong>Золотой класс, 1600 м, 3+yo</strong> (максимальный класс резвости C и ниже)</td><td>Заявлено: 19</td></tr><div id='txtHint20141009056568'><tr id='charter'><td colspan='2' align='left'>6568 <strong>Приправа</strong><font size=1>, К 3,  (1329 Проспект - 1211 Прогулка), Penumbra, 1 0-0-0 $0</font></td><td><br><form id='cancellentry' name='cancellentry' method='post' action='races_entries.php#race2014100905'><input type='image' src='img/cancel_but.jpg' alt='Отменить заявку' onclick='raceentrydelete(2014100905,6568,194)'></form></td></tr></div><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014100905>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014100906'>2014100906</a></td><td align='left'><strong>Золотой класс, 1800 м, 3yo</strong> (максимальный класс резвости E и ниже)</td><td>Заявлено: 9</td></tr><div id='txtHint20141009066302'><tr id='charter'><td colspan='2' align='left'>6302 <strong>Момент</strong><font size=1>, Ж 3,  (2490 Marvellous - 1207 Марсельеза), Penumbra, 2 0-1-0 $1,630</font></td><td><br><form id='cancellentry' name='cancellentry' method='post' action='races_entries.php#race2014100906'><input type='image' src='img/cancel_but.jpg' alt='Отменить заявку' onclick='raceentrydelete(2014100906,6302,194)'></form></td></tr></div><div id='txtHint20141009066738'><tr id='charter'><td colspan='2' align='left'>6738 <strong>El Alza</strong><font size=1>, Ж 3,  (2509 Attar Of Roses - 2103 El Lazize), Penumbra, 0 0-0-0 $0</font></td><td><br><form id='cancellentry' name='cancellentry' method='post' action='races_entries.php#race2014100906'><input type='image' src='img/cancel_but.jpg' alt='Отменить заявку' onclick='raceentrydelete(2014100906,6738,194)'></form></td></tr></div>
<tr id='infoblockheader'><td>10 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101010'>2014101010</a></td><td align='left'><strong>Тестовый класс, 1200 м, 2+yo</strong></td><td>Заявлено: 15</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101010>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101013'>2014101013</a></td><td align='left'><strong>Тестовый класс, 1300 м, 2+yo</strong></td><td>Заявлено: 9</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101013>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101009'>2014101009</a></td><td align='left'><strong>Тестовый класс, 2200 м, 3+yo</strong></td><td>Заявлено: 19</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101009>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101012'>2014101012</a></td><td align='left'><strong>Тестовый класс, 4000 м, 4+yo</strong></td><td>Заявлено: 6</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101012>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101011'>2014101011</a></td><td align='left'><strong>Тестовый класс, 4800 м, 4+yo</strong></td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101011>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101014'>2014101014</a></td><td align='left'><strong>Тестовый класс, ст-з, 1300 м, 2+yo</strong></td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101014>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101008'>2014101008</a></td><td align='left'><strong>Тестовый класс, ст-з, 1300 м, 2+yo</strong></td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101008>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101007'>2014101007</a></td><td align='left'><strong>Тестовый класс, ст-з, 1300 м, 2+yo</strong></td><td>Заявлено: 4</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101007>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101006'>2014101006</a></td><td align='left'><strong>Тестовый класс, ст-з, 4000 м, 4+yo</strong></td><td>Заявлено: 5</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101006>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101001'>2014101001</a></td><td align='left'><strong>Гр.II Strawberry Time, 2200 м, 3+yo К</strong> (класс резвости C)</td><td>Заявлено: 12</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101001>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101002'>2014101002</a></td><td align='left'><strong>Гр.II Country Ruffian Cup, ст-з, 2800 м, 3+yo Ж</strong> (класс резвости D)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101002>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101003'>2014101003</a></td><td align='left'><strong>Гр.I Note Bianko Stakes, ст-з, 1600 м, 2+yo К</strong> (класс резвости B)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101003>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>10 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#E4FFAE'><td width='100'><a name='race2014101005'>2014101005</a></td><td align='left'><strong>Золотой класс, 1900 м, 2+yo</strong> (максимальный класс резвости A и ниже)</td><td>Заявлено: 10</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101005>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#E4FFAE'><td width='100'><a name='race2014101004'>2014101004</a></td><td align='left'><strong>Золотой класс, 2800 м, 3+yo Ж</strong> (максимальный класс резвости B и ниже)</td><td>Заявлено: 7</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101004>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>11 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101115'>2014101115</a></td><td align='left'><strong>Тестовый класс, 1400 м, 2+yo</strong></td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101115>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101111'>2014101111</a></td><td align='left'><strong>Тестовый класс, 1400 м, 2+yo</strong></td><td>Заявлено: 7</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101111>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101110'>2014101110</a></td><td align='left'><strong>Тестовый класс, 1600 м, 2+yo</strong></td><td>Заявлено: 10</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101110>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101109'>2014101109</a></td><td align='left'><strong>Тестовый класс, 1900 м, 2+yo</strong></td><td>Заявлено: 7</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101109>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101113'>2014101113</a></td><td align='left'><strong>Тестовый класс, 4400 м, 4+yo</strong></td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101113>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101114'>2014101114</a></td><td align='left'><strong>Тестовый класс, 4400 м, 4+yo</strong></td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101114>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101106'>2014101106</a></td><td align='left'><strong>Тестовый класс, 4800 м, 4+yo</strong></td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101106>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101105'>2014101105</a></td><td align='left'><strong>Тестовый класс, ст-з, 1400 м, 2+yo</strong></td><td>Заявлено: 6</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101105>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101112'>2014101112</a></td><td align='left'><strong>Тестовый класс, ст-з, 2200 м, 3+yo</strong></td><td>Заявлено: 7</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101112>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101101'>2014101101</a></td><td align='left'><strong>Гр.III Tasotti Stakes, 2800 м, 3+yo К</strong> (класс резвости E)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101101>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101102'>2014101102</a></td><td align='left'><strong>Гр.II Цветочные скачки, ст-з, 1300 м, 2+yo К</strong> (класс резвости C)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101102>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101103'>2014101103</a></td><td align='left'><strong>Гр.II Приз Луары, ст-з, 3200 м, 4+yo К</strong> (класс резвости D)</td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101103>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>11 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#CCCC66'><td width='100'><a name='race2014101108'>2014101108</a></td><td align='left'><strong>Золотой класс, 1400 м, 3+yo К</strong> (максимальный класс резвости C и ниже)</td><td>Заявлено: 6</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101108>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#CCCC66'><td width='100'><a name='race2014101107'>2014101107</a></td><td align='left'><strong>Золотой класс, 2000 м, 4+yo К</strong> (максимальный класс резвости B и ниже)</td><td>Заявлено: 5</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101107>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>11 октября 2014</td><td colspan='2'>Ипподром: Феникс</td></tr>
<tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014101104'>2014101104</a></td><td align='left'><strong>FNX Гр.I  Stars Wins, 1700 м, 2yo</strong> (время на данной дистанции не резвее 1:50)</td><td>Заявлено: 4</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101104>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>12 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101211'>2014101211</a></td><td align='left'><strong>Тестовый класс, 1700 м, 2+yo</strong></td><td>Заявлено: 5</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101211>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101215'>2014101215</a></td><td align='left'><strong>Тестовый класс, 2000 м, 3+yo</strong></td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101215>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101207'>2014101207</a></td><td align='left'><strong>Тестовый класс, ст-з, 1200 м, 2+yo</strong></td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101207>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101213'>2014101213</a></td><td align='left'><strong>Тестовый класс, ст-з, 1900 м, 2+yo</strong></td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101213>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101212'>2014101212</a></td><td align='left'><strong>Тестовый класс, ст-з, 1900 м, 2+yo</strong></td><td>Заявлено: 6</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101212>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101216'>2014101216</a></td><td align='left'><strong>Тестовый класс, ст-з, 2000 м, 3+yo</strong></td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101216>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101201'>2014101201</a></td><td align='left'><strong>AA Astor Plate, 1000 м, 2yo Ж</strong> (минимальный класс резвости D и выше)</td><td>Заявлено: 4</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101201>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101202'>2014101202</a></td><td align='left'><strong>AA Grand Prix de Saint-Cloud, 2000 м, 3+yo К</strong> (минимальный класс резвости D и выше)</td><td>Заявлено: 3</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101202>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101203'>2014101203</a></td><td align='left'><strong>AA New Ideal Handicap, торф, 1100 м, 2yo Ж</strong> (минимальный класс резвости D и выше)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101203>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101204'>2014101204</a></td><td align='left'><strong>AA Queen Cup Stakes, ст-з, 1600 м, 3+yo К</strong> (минимальный класс резвости D и выше)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101204>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101205'>2014101205</a></td><td align='left'><strong>AA Triple Vivat Stakes, ст-з, 2400 м, 3+yo Ж</strong> (минимальный класс резвости D и выше)</td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101205>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101206'>2014101206</a></td><td align='left'><strong>AA Queen Victoria Stakes, ст-з, 4800 м, 4+yo К</strong> (минимальный класс резвости D и выше)</td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101206>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>12 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101214'>2014101214</a></td><td align='left'><strong>Серебряный класс, 1000 м, 4+yo</strong> (максимальный класс резвости F и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101214>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101210'>2014101210</a></td><td align='left'><strong>Серебряный класс, 1400 м, 3+yo</strong> (максимальный класс резвости D и ниже)</td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101210>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101217'>2014101217</a></td><td align='left'><strong>Серебряный класс, 2400 м, 3yo</strong> (максимальный класс резвости C и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101217>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101208'>2014101208</a></td><td align='left'><strong>Золотой класс, торф, 1100 м, 2&3yo</strong> (максимальный класс резвости D и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101208>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101209'>2014101209</a></td><td align='left'><strong>Золотой класс, торф, 1800 м, 2&3yo</strong> (максимальный класс резвости E и ниже)</td><td>Заявлено: 3</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101209>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>13 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#E4FFAE'><td width='100'><a name='race2014101301'>2014101301</a></td><td align='left'><strong>Гр.III Delta X Stakes, 1900 м, 2+yo Ж</strong> (максимальный класс резвости F и ниже)</td><td>Заявлено: 6</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101301>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#E4FFAE'><td width='100'><a name='race2014101302'>2014101302</a></td><td align='left'><strong>Гр.II Полотно Пенелопы, 4000 м, 4+yo К</strong> (класс резвости D)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101302>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#E4FFAE'><td width='100'><a name='race2014101303'>2014101303</a></td><td align='left'><strong>Гр.II The Hamilton Classic, ст-з, 2000 м, 3+yo К</strong> (класс резвости D)</td><td>Заявлено: 2</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101303>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>13 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101306'>2014101306</a></td><td align='left'><strong>Серебряный класс, 1200 м, 3yo</strong> (максимальный класс резвости C и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101306>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101305'>2014101305</a></td><td align='left'><strong>Золотой класс, 4000 м, 4+yo</strong> (максимальный класс резвости C и ниже)</td><td>Заявлено: 3</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101305>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#D8FEFE'><td width='100'><a name='race2014101304'>2014101304</a></td><td align='left'><strong>Золотой класс, 4800 м, 4+yo</strong> (максимальный класс резвости D и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101304>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>14 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#CCCC66'><td width='100'><a name='race2014101401'>2014101401</a></td><td align='left'><strong>Гр.III Emax Handicap, ст-з, 2800 м, 3+yo К</strong> (класс резвости E)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101401>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#CCCC66'><td width='100'><a name='race2014101402'>2014101402</a></td><td align='left'><strong>Гр.II Discovery Gardens Handicap, 3200 м, 4+yo К</strong> (класс резвости C)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101402>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#CCCC66'><td width='100'><a name='race2014101403'>2014101403</a></td><td align='left'><strong>Гр.I Desert Stakes, ст-з, 1600 м, 2+yo Ж</strong> (минимальный класс резвости A и выше)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101403>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>14 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014101404'>2014101404</a></td><td align='left'><strong>Золотой класс, 1400 м, 2+yo</strong> (максимальный класс резвости E и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101404>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#C1F999'><td width='100'><a name='race2014101405'>2014101405</a></td><td align='left'><strong>Золотой класс, 2400 м, 6yo Ж</strong> (максимальный класс резвости B и ниже)</td><td>Заявлено: 1</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101405>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>15 октября 2014</td><td colspan='2'>Ипподром: Центральный</td></tr>
<tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101501'>2014101501</a></td><td align='left'><strong>Гр.III Восточная Сказка, 1100 м, 2+yo К</strong> (максимальный класс резвости F и ниже)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101501>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101502'>2014101502</a></td><td align='left'><strong>Гр.III Randevu Stakes, 4000 м, 4+yo К</strong> (максимальный класс резвости F и ниже)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101502>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101503'>2014101503</a></td><td align='left'><strong>Гр.II Приз Матроны, ст-з, 3200 м, 4+yo К</strong> (класс резвости C)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101503>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FED6FE'><td width='100'><a name='race2014101504'>2014101504</a></td><td align='left'><strong>Гр.I Triple Vivat Stakes, ст-з, 2000 м, 3+yo Ж</strong> (минимальный класс резвости A и выше)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101504>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr>
<tr id='infoblockheader'><td>15 октября 2014</td><td colspan='2'>Ипподром: Зеленый Луг</td></tr>
<tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101506'>2014101506</a></td><td align='left'><strong>Золотой класс, 2400 м, 3yo</strong> (максимальный класс резвости E и ниже)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101506>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr><tr id='charter' bgcolor='#FBF9B5'><td width='100'><a name='race2014101505'>2014101505</a></td><td align='left'><strong>Золотой класс, торф, 1900 м, 2+yo</strong> (максимальный класс резвости D и ниже)</td><td>Заявлено: 0</td></tr><tr id='charter'><td colspan='3' align='right'><br><form action='enter_racer.php' method='post'>
<input type='hidden' name='racecode' value=2014101505>
<input type='submit' value='Регистрация'>
</form><br><br></td></tr></table>


</div>
<!-- End of основного блока -->
</div>
<div id="footer">
<br />
<a href="http://www.eklps.com" class="menulink">Игровой проект "ЭклипС"</a> &copy; 2003-2014<br />
Вся информация на данном сайте является частной собственностью и защищена законом <br />
<div style="padding: 40px 0px 20px 0px;">
<span>
		<!--Akavita counter start-->
<script type="text/javascript">var AC_ID=55830;var AC_TR=false;
(function(){var l='http://adlik.akavita.com/acode.js'; var t='text/javascript';
try {var h=document.getElementsByTagName('head')[0];
var s=document.createElement('script'); s.src=l;s.type=t;h.appendChild(s);}catch(e){
document.write(unescape('%3Cscript src="'+l+'" type="'+t+'"%3E%3C/script%3E'));}})();
</script><span id="AC_Image"></span>
<noscript><a target='_blank' href='http://www.akavita.by/'>
<img src='http://adlik.akavita.com/bin/lik?id=55830&it=1'
border='0' height='1' width='1' alt='Akavita'/>
</a></noscript>
<!--Akavita counter end-->


&nbsp;&nbsp;&nbsp;
<!-- HotLog -->
<script type="text/javascript">
var hotlog_counter_id = 2319785;
var hotlog_hit = 25;
var hotlog_counter_type = 565;
</script>
<script src="http://js.hotlog.ru/counter.js" type="text/javascript"></script>
<noscript>
<a href="http://click.hotlog.ru/?2319785" target="_blank">
<img src="http://hit25.hotlog.ru/cgi-bin/hotlog/count?s=2319785&im=565" border="0"
title="HotLog" alt="HotLog"></a>
</noscript>
<!-- /HotLog -->
&nbsp;&nbsp;&nbsp;


<!--Openstat-->
<span id="openstat586556"></span>
<script type="text/javascript">
var openstat = { counter: 586556, image: 87, color: "ff9822", next: openstat, track_links: "all" };
(function(d, t, p) {
var j = d.createElement(t); j.async = true; j.type = "text/javascript";
j.src = ("https:" == p ? "https:" : "http:") + "//openstat.net/cnt.js";
var s = d.getElementsByTagName(t)[0]; s.parentNode.insertBefore(j, s);
})(document, "script", document.location.protocol);
</script>
<!--/Openstat-->

<!--Google-->
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-43673244-1', 'eklps.com');
  ga('send', 'pageview');

</script>
<!--//Google-->

</span>
		</div>
</div>
</div>
</body>
</html>
`

	return races
}
