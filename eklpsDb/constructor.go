package eklpsDb

import (

)

func NewDb() (*Eklps, error) {

	eklps := &Eklps {
		//Dsn : "dev:qwr123rrr@tcp(192.168.56.101:3306)/eklps",
		Dsn : "eklps:eklps@tcp(127.0.0.1:3306)/eklps",
		FilterMap : map[string]string{

			"Sc" : "r.sc=1",
			"Fl" : "r.sc=0",

			"SSex" : `r.sex = "m"`,
			"MSex" : `r.sex = "f"`,
			"Strict" : `1`,
			"StrictPlus" : `1`,
			"StrictPlusPlus" : `1`,
		},

		ClassFilterMap : map[string]string{
			"TClass" : `r.name like "Тестовый класс%"`,
			"MClass" : `r.name like "Медный класс%"`,
			"BClass" : `r.name like "Бронзовый класс%"`,
			"SClass" : `r.name like "Серебряный класс%"`,
			"GClass" : `r.name like "Золотой класс%"`,
			"PClass" : `r.name like "Платиновый класс%"`,
			"GR3" : `r.name like "Гр.III%"`,
			"GR2" : `r.name like "Гр.II%"`,
			"GR1" : `r.name like "Гр.I%"`,
			"AAClass" : `r.name like "AA %"`,
		},
	}

	return eklps, nil
}
