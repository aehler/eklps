package eklpsDb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"bytes"
)

type Eklps struct {
	Dsn string
	Conn *sql.DB
	FilterMap map[string]string
	ClassFilterMap map[string]string
}

type Race struct {
	RaceId int
	RaceTitle string
	RaceDate string
	RaceClass string
}

type Dist struct {
	Distance string
	Fl string
	Sc string
}

type Params struct {
	Age int
	Sex string
	Spec string
	Distances []Dist
	Filters []string
}

type QueryWriter interface {
	WriteString(string) (int, error)
}

func (eklps *Eklps) Connect() error {

	db, err := sql.Open("mysql", eklps.Dsn)
	if err == nil {
	} else {
		return err
	}

	eklps.Conn = db

	return nil
}


func (eklps *Eklps) GetSeason() (int, error) {

	var season int

	rows, err := eklps.Conn.Query("select season from season limit 1")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&season); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return season, nil
}


/**
runs a query like

select r.id as RaceId, r.fullname as RaceName, date(substr(r.id, 1,8)) as RaceDate
from races r
left join race_conditions rc1 on r.id = rc1.id_races and r.distance = 1000 and
        case rc1.class_eq
                when "=" then rc1.class = "E"
                when "<=" then ASCII(rc1.class) <= ASCII("E")
                when ">=" then ASCII(rc1.class) >= ASCII("E")
                else 1
        end
left join race_conditions rc2 on r.id = rc2.id_races and r.distance = 1900 and
        case rc2.class_eq
                when "=" then rc1.class = "K"
                when "<=" then ASCII(rc2.class) <= ASCII("K")
                when ">=" then ASCII(rc2.class) >= ASCII("K")
                else 1
        end
left join race_conditions rc3 on r.id = rc3.id_races and r.distance = 1100 and
        case rc3.class_eq
                when "=" then rc3.class = "C"
                when "<=" then ASCII(rc3.class) <= ASCII("C")
                when ">=" then ASCII(rc3.class) >= ASCII("C")
                else 1
        end
left join race_conditions rc_void on r.id = rc_void.id_races and rc_void.class = ""
where
(r.age_conditions like "%2%"
or (substr(r.age_conditions,1,1) <= "2" && substr(r.age_conditions,2,1) = "+"))
and r.sex in ('f', 'all')
and sc = 1
and case
        when rc1.class is not null then 1
        when rc2.class is not null then 1
        when rc3.class is not null then 1
        when rc_void.class is not null then 1
        else 0
end
*/
func (eklps *Eklps) BuildSql(p Params) (bytes.Buffer, error) {

	var query bytes.Buffer

	season, err := eklps.GetSeason()
	if err != nil {
		return query, err
	}

	age := season - p.Age
	sex := "f"
	if p.Sex == "ж" {
		sex = "m"
	}

	var subfields []string
	var letter byte
	var sc int
	var classFilter bytes.Buffer

	query.WriteString("select r.id as RaceId, r.fullname as RaceTitle, date(substr(r.id, 1,8)) as RaceDate,")
	query.WriteString(" concat(rcw.class, replace(replace(replace(rcw.class_eq,'=',''), '<', '-'),'>','+')) as RaceClass")
	query.WriteString(" from races r ")

	for i, d := range p.Distances {

		if d.Fl == "" {
			sc = 1
			letter = []byte(d.Sc)[0]
		} else {
			sc = 0
			letter = []byte(d.Fl)[0]
		}

		query.WriteString(fmt.Sprintf(" left join race_conditions rc%d on r.id = rc%d.id_races and r.distance = '%s'", i, i, d.Distance))
		query.WriteString(fmt.Sprintf(" and case rc%d.class_eq ", i))
		query.WriteString(fmt.Sprintf(" when '=' then rc%d.class = '%s'", i, string(letter)))
		query.WriteString(fmt.Sprintf(" when '<=' then ASCII(rc%d.class) <= ASCII('%s')", i, string(letter)))
		query.WriteString(fmt.Sprintf(" when '>=' then ASCII(rc%d.class) >= ASCII('%s')", i, string(letter)))
		query.WriteString(" else 1 end ")
		query.WriteString(fmt.Sprintf(" and r.sc = %d ", sc))

		subfields = append(subfields, fmt.Sprintf("rc%d", i))
	}

	query.WriteString("left join race_conditions rc_void on r.id = rc_void.id_races and rc_void.class = '' ")
	query.WriteString("left join race_conditions rcw on r.id = rcw.id_races ")

	query.WriteString("where")
	query.WriteString("(r.age_conditions like '%")
	query.WriteString(fmt.Sprintf("%d", age))
	query.WriteString("%' or (substr(r.age_conditions,1,1) <= '")
	query.WriteString(fmt.Sprintf("%d", age))
	query.WriteString("' && substr(r.age_conditions,2,1) = '+'))")
	query.WriteString(fmt.Sprintf("and r.sex in ('%s', 'all')", sex))
	query.WriteString(" and case ")

	for _, r := range subfields {
		query.WriteString(fmt.Sprintf(" when %s.class is not null then 1", r))
	}

	query.WriteString(" when rc_void.class is not null then 1 else 0 end ")

	for _, f := range p.Filters {
		eklps.AddOrFilters(f, &classFilter)
		eklps.AddFilters(f, &query)
	}

	if classFilter.Len() > 0 {
		query.WriteString(fmt.Sprintf(" and (0%s)", classFilter.String()))
	}

	return query, nil
}

func (eklps *Eklps) AddOrFilters(f string, q QueryWriter) {
	if eklps.ClassFilterMap[f] != "" {
		q.WriteString(fmt.Sprintf(" or %s",eklps.ClassFilterMap[f]))
	}
}

func (eklps *Eklps) AddFilters(f string, q QueryWriter) {
	if eklps.FilterMap[f] != "" {
		q.WriteString(fmt.Sprintf(" and %s ", eklps.FilterMap[f]))
	}
}


func (eklps *Eklps) GetRaces(p Params) ([]Race, error) {

	query, err := eklps.BuildSql(p)
	if err != nil {
		return []Race{}, err
	}

	fmt.Println(query.String())

	rows, err := eklps.Conn.Query(query.String())
	if err != nil {
		return []Race{}, err
	}
	defer rows.Close()

	var race Race
	var races []Race

	for rows.Next() {

		if err := rows.Scan(&race.RaceId, &race.RaceTitle, &race.RaceDate, &race.RaceClass); err != nil {
			return []Race{}, err
		}

		races = append(races, race)
	}

	if err := rows.Err(); err != nil {
		return []Race{}, err
	}

	return races, nil
}
