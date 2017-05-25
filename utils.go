package edc

import (
	"log"
	"strings"
	"time"
)

func setStrMonth(d string) string {
	var result string
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		return result
	}
	str := t.Format("02.01.2006")
	spl := strings.Split(str, ".")
	month := map[string]string{
		"01": "января",
		"02": "февраля",
		"03": "марта",
		"04": "апреля",
		"05": "мая",
		"06": "июня",
		"07": "июля",
		"08": "августа",
		"09": "сентября",
		"10": "октября",
		"11": "ноября",
		"12": "декабря",
	}
	result = spl[0] + " " + month[spl[1]] + " " + spl[2] + " года"
	return result
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}
