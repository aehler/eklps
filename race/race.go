package race

import (
	"time"
	"strings"
)

type Conditions struct {
	Class string
	ClassEq string
	Timing string
	TimingEq string
}

type Race struct {
	Id int
	RaceID string
	RaceDate int
	Dt time.Time
	Class string
	Distance int64
	Sc bool
	AgeConditions string
	Sex string
	Name string
	Conditions Conditions
	Season uint
}

func (r *Race) Unmarshal(src string) {

	rd := strings.Split(src, ",")

	r.Name = src
	r.Class = strings.Trim(rd[0], " ")

}
