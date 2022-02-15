package html

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	config "github.com/Obito1903/CY-celcat/pkg"
	"github.com/Obito1903/CY-celcat/pkg/calendar"
	"github.com/Obito1903/CY-celcat/pkg/celcat"
)

func TestCalHtmlStatic(t *testing.T) {
	cal := calendar.Calendar{
		Name: "GIG1",
		Events: []calendar.Event{
			{
				Id:         "Salut",
				Start:      time.Now(),
				End:        time.Now().Add(time.Hour),
				AllDay:     false,
				Category:   "CM",
				Subjects:   []string{"Algo"},
				Location:   []string{"102"},
				Professors: []string{"Ines"},
				Notes:      "nothing",
			},
			{
				Id:         "Salut2",
				Start:      time.Now().Add(time.Hour),
				End:        time.Now().Add(time.Hour * 2),
				AllDay:     false,
				Category:   "CM",
				Subjects:   []string{"Anglais"},
				Location:   []string{"102"},
				Professors: []string{"Magda"},
				Notes:      "nothing",
			},
		},
	}

	htmlCal := CalToHtmlCal(cal, time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local))
	t.Log(htmlCal)
	htmlCal.ToFile("../../../web/templates/calendar.go.html", "../../../web/static/"+cal.Name+".html")
}

func TestCalHtml(t *testing.T) {
	config := config.ReadConfig("../../../example.config.json")
	jar, err := cookiejar.New(nil)
	if err != nil {
		os.Exit(1)
	}

	client := &http.Client{
		Jar: jar,
	}
	url, err := url.Parse(config.CelcatHost)
	if err != nil {
		t.Log("Hello")
		os.Exit(1)
	}
	data := celcat.Login(client, *url, config.UserName, config.UserPassword)

	celcatCalendar := celcat.GetCalendar(client, *url, data.FederationId, time.Date(2022, 01, 24, 0, 0, 0, 0, time.Local), time.Date(2022, 01, 27, 0, 0, 0, 0, time.Local))
	calendar := calendar.FromCelcat(celcatCalendar, "GIG1")

	htmlCal := CalToHtmlCal(calendar, time.Date(2022, 01, 24, 0, 0, 0, 0, time.Local))
	// t.Log(htmlCal)
	htmlCal.ToFile("../../../web/templates/calendar.go.html", "../../../web/static/calendars/"+calendar.Name+".html")
}
