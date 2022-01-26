package celcat

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	config "github.com/Obito1903/CY-celcat/pkg"
)

func TestRequestToken(t *testing.T) {
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
		os.Exit(1)
	}
	t.Log("Token: " + getRequestVerificationToken(client, *url))
}

func TestLogin(t *testing.T) {
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
	data := Login(client, *url, config.UserName, config.UserPassword)
	t.Log("Token : " + data.token + " | Id : " + data.federationId)
}

func TestCalendar(t *testing.T) {
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
	data := Login(client, *url, config.UserName, config.UserPassword)

	events := getEventList(client, *url, data.federationId, time.Date(2022, 01, 24, 0, 0, 0, 0, time.Local), time.Date(2022, 01, 25, 0, 0, 0, 0, time.Local))
	t.Log("Loged as : ", data.federationId)
	for _, event := range events {
		t.Log(event.Id, "| Start : ", event.Start, ", End : ", event.End)
	}
}

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
	data := Login(client, *url, config.UserName, config.UserPassword)

	events := getEventList(client, *url, data.federationId, time.Date(2022, 01, 24, 0, 0, 0, 0, time.Local), time.Date(2022, 01, 25, 0, 0, 0, 0, time.Local))
	t.Log("Loged as : ", data.federationId)
	for _, event := range events {
		getEventDetails(client, *url, &event)
		t.Log(event.Id, "| Start : ", event.Start, ", End : ", event.End)
		for _, element := range event.Elements {
			t.Log(element.Label)
		}

	}
}
