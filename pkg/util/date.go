package util

import (
	ptime "github.com/yaa110/go-persian-calendar"
	"log"
	"time"
)

func ConvertShamsiToGregorian(year, month, day int) time.Time {
	var pt ptime.Time = ptime.Date(1394, ptime.Mehr, 2, 12, 0, 0, 0, ptime.Iran())

	// Get a new instance of time.Time
	t := pt.Time()

	return t

}

func ConvertGregorianToShamsi(t time.Time) {
	//g
	date, err := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatalln("err:", err)
	}
	pt := ptime.New(date)
	pt.Time()

}
