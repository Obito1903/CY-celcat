package cyCelcat

import (
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
		log.Fatal("Could not parse Celcat resoinse")
		os.Exit(1)
	}
	celcat.Login(client, *url, config.UserName, config.UserPassword)

	log.Print("Query period : ", period)
	for _, groupe := range config.Groupes {
		log.Print("Processing groupe : ", groupe.Name, " | ", groupe.Id)
		period.Start = calendar.FirstDayOfISOWeek(period.Start)
		celcatCalendar := celcat.GetCalendar(client, *url, groupe.Id, period.Start, period.End)
		calendar := calendar.FromCelcat(celcatCalendar, groupe.Name)
		if config.ICS {
			ics.IcsToFile(ics.CalendarToICS(calendar), config.ICSPath+calendar.Name+".ics")
		}

		if config.HTML {
			htmlCal := html.CalToHtmlCal(calendar, period.Start)
			htmlCal.ToFile(config.HtmlTemplate, config.HTMLPath+calendar.Name+".html")
			if config.PNG {
				html.ToPng(config, config.HTMLPath+calendar.Name+".html", config.PNGPath+calendar.Name+".png")
			}
		}
	}

}
