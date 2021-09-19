package main

import (
	"log"
	"os"
	"time"
	"fmt"

	"github.com/Obito1903/CY-celcat/celcat"
	"github.com/Obito1903/CY-celcat/celcat/common"
	"github.com/Obito1903/CY-celcat/celcat/fetch"
)

func firstDayOfISOWeek(timezone *time.Location) time.Time {
	year, week := time.Now().ISOWeek()
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		fmt.Println("bingo!")
		week = week + 1
	}
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}

func main() {
	var err error
	genSVG := false
	svgName := "cal.svg"
	configPath := "example.config.json"
	firstDayofWeek := firstDayOfISOWeek(time.UTC)
	queryPeriod := common.Period{
		StartDate: firstDayofWeek,
		EndDate:   firstDayofWeek.Add(time.Hour * 24 * 30),
	}

	for i, arg := range os.Args {
		switch arg {
		case "-c":
			configPath = os.Args[i+1]
		case "-d":
			if len(os.Args) < i+2 {
				log.Fatal("Arguments manquant")
				os.Exit(1)

			}
			loc, _ := time.LoadLocation("Europe/Paris")
			queryPeriod.StartDate, err = time.ParseInLocation("2006-01-02", os.Args[i+1], loc)
			common.CheckErr(err)
			queryPeriod.EndDate, err = time.ParseInLocation("2006-01-02", os.Args[i+2], loc)
			common.CheckErr(err)
		case "-svg":
			genSVG = true
			queryPeriod = common.Period{
				StartDate: firstDayofWeek,
				EndDate:   firstDayofWeek.Add(time.Hour * 24 * 6),
			}
			if len(os.Args) < i+2 {

			} else {
				svgName = os.Args[i+1]
			}
		}

	}

	config := celcat.GetConfig(configPath)
	calendar := fetch.GetCalendar(config, queryPeriod)
	celcat.ToICS(calendar)
	if genSVG {
		celcat.ToSVG(calendar, svgName)
	}

}
