package cyCelcat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	config "github.com/Obito1903/CY-celcat/pkg"
	"github.com/Obito1903/CY-celcat/pkg/calendar"
	"github.com/Obito1903/CY-celcat/pkg/calendar/html"
	"github.com/Obito1903/CY-celcat/pkg/calendar/ics"
	"github.com/Obito1903/CY-celcat/pkg/celcat"
)

type Period struct {
	Start time.Time
	End   time.Time
}

func Query(config config.Config, period Period) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		os.Exit(1)
	}

	client := &http.Client{
		Jar: jar,
	}
	url, err := url.Parse(config.CelcatHost)
	if err != nil {
		log.Fatal("Could not parse Celcat response")
		os.Exit(1)
	}
	celcat.Login(client, *url, config.UserName, config.UserPassword)

	log.Print("Query period : ", period)

	calendarList := []calendar.Calendar{}

	for _, groupe := range config.Groupes {
		log.Print("Processing groupe : ", groupe.Name, " | ", groupe.Id)
		period.Start = calendar.FirstDayOfISOWeek(period.Start)
		period.End = calendar.FirstDayOfISOWeek(period.End.Add(time.Hour * 24 * 7))
		celcatCalendar := celcat.GetCalendar(client, *url, groupe.Id, period.Start, period.End)
		calendarList = append(calendarList, calendar.FromCelcat(celcatCalendar, groupe.Name))
	}

	for _, calendar := range calendarList {
		if config.ICS {
			ics.IcsToFile(ics.CalendarToICS(calendar), config.ICSPath+calendar.Name+".ics")
		}

		fmt.Printf("Calendar %s next event : %+v\n", calendar.Name, calendar.TomorrowFirstEvent())

		if config.NextAlarm {
			nextDayEvent := calendar.NextEventToJson()
			_ = ioutil.WriteFile(config.NextAlarmPath+calendar.Name+".json", []byte(nextDayEvent), 0644)
		}

		if config.HTML {
			for week := 0; week < config.Weeks; week++ {
				htmlCal := html.CalToHtmlCal(calendar, period.Start.Add(time.Hour*24*7*time.Duration(week)))
				name := calendar.Name
				if week != 0 {
					name = name + "+" + fmt.Sprint(week)
				}

				htmlCal.ToFile(config.HtmlTemplate, config.HTMLPath+name+".html")

				if config.PNG && week == 0 {
					html.ToPng(config, config.HTMLPath+calendar.Name+".html", config.PNGPath+calendar.Name+".png")
				}

			}

		}
	}

}
