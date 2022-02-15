package calendar

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	config "github.com/Obito1903/CY-celcat/pkg"
	"github.com/Obito1903/CY-celcat/pkg/celcat"
)

func TestEventDetails(t *testing.T) {
	config := config.ReadConfig("../../example.config.json")
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

	celcatCalendar := celcat.GetCalendar(client, *url, data.FederationId, time.Date(2022, 01, 24, 0, 0, 0, 0, time.Local), time.Date(2022, 01, 25, 0, 0, 0, 0, time.Local))
	for _, celcatEvent := range celcatCalendar {
		event := eventFromCelcat(celcatEvent)
		t.Log("Id: ", event.Id, "Category: ", event.Category, " | Subjects: ", event.Subjects, " | Location: ", event.Location)

	}
}
