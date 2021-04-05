package fetch

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/Obito1903/CY-celcat/celcat/common"
	"golang.org/x/net/publicsuffix"
)

type SessionData struct {
	token         string
	antiforgery   *http.Cookie
	federationIds string
}

// Get RequestToken and cookie from services-web.u-cergy.fr wich are necessary to send a login request to celcat
func getRequestLoginToken(client *http.Client) SessionData {
	var token SessionData
	//send a get request to the calendar service to obtain the requestToken
	resp, err := client.Get("https://services-web.u-cergy.fr/calendar/LdapLogin")
	common.CheckErr(err)
	//Store the request cookie
	for _, cookie := range resp.Cookies() {
		token.antiforgery = cookie
	}
	// Search for the request token in the request response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "__RequestVerificationToken") {
			requestTokenTag := scanner.Text()
			posStartValue := strings.Index(requestTokenTag, "value=")
			posEndValue := strings.Index(requestTokenTag, "\" />")
			token.token = requestTokenTag[posStartValue+7 : posEndValue]
			break
		}
	}
	return token
}

// Log onto celcat using provided login and password
func logon(client *http.Client, session SessionData, config common.Config) SessionData {
	// Init the header to be sent
	headerData := url.Values{
		"Name":                       {config.UserId},
		"Password":                   {config.UserPassword},
		"__RequestVerificationToken": {session.token},
	}
	// Setup a login request
	req, err := http.NewRequest("POST", "https://services-web.u-cergy.fr/calendar/LdapLogin/Logon", strings.NewReader(headerData.Encode()))
	common.CheckErr(err)
	req.AddCookie(session.antiforgery)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Send the login request
	resp, err := client.Do(req)
	common.CheckErr(err)
	if resp.StatusCode != 200 {
		common.CheckErr(errors.New(fmt.Sprint(resp.StatusCode) + ": Request Error"))
	}
	responseUrl, err := resp.Request.Response.Location()
	common.CheckErr(err)
	session.federationIds = responseUrl.Query().Get("FederationIds")
	return session
}

// Query the calendar and return the interface containing all calendar info
func getCalData(client *http.Client, session SessionData, queryPeriod common.Period) []interface{} {
	var calendarInterface []interface{}

	fmt.Println("Requesting calendar from", queryPeriod.StartDate.Format("2006-01-02"), "to", queryPeriod.EndDate.Format("2006-01-02"))
	// Setup request header
	headerData := url.Values{
		"start":           {queryPeriod.StartDate.Format("2006-01-02")},
		"end":             {queryPeriod.EndDate.Format("2006-01-02")},
		"resType":         {"104"},
		"calView":         {"agendaWeek"},
		"federationIds[]": {session.federationIds},
	}
	// Send request with cookies from the cookieJar
	resp, err := client.PostForm("https://services-web.u-cergy.fr/calendar/Home/GetCalendarData", headerData)
	common.CheckErr(err)

	// Parse the response and put all data into an interface
	body, err := io.ReadAll(resp.Body)
	common.CheckErr(err)
	fmt.Println(resp.StatusCode)
	err = json.Unmarshal(body, &calendarInterface)
	common.CheckErr(err)
	return calendarInterface
}

// Parse the json description to extract the subject, location and organizer of the event
func parseDescription(desc string, category string) (string, string, string) {
	parsed := strings.Split(desc, "<br />")
	var module, location, prof string

	if category == "Indisponibilité" {
		module, location, prof = "férié", "", ""
	} else {
		module = html.UnescapeString(strings.Trim(parsed[1], "\r\n"))
		for i := 2; i < len(parsed)-1; i++ {
			location = location + " " + strings.Trim(parsed[i], "\r\n")
		}
		location = html.UnescapeString(location)
		prof = strings.Trim(parsed[len(parsed)-1], "\r\n")
	}

	return module, location, prof
}

// Create a callEvent from the json data of one celcat event
func parseCalEvent(baseEvent interface{}) common.CalEvent {
	var parsedEvent common.CalEvent
	loc, _ := time.LoadLocation("Europe/Paris")

	mappedCelcatEvent := baseEvent.(map[string]interface{})

	parsedEvent.Category = fmt.Sprintf("%s", mappedCelcatEvent["eventCategory"])

	parsedEvent.Module, parsedEvent.Location, parsedEvent.Prof = parseDescription(fmt.Sprintf("%s", mappedCelcatEvent["description"]), parsedEvent.Category)
	parsedEvent.Id = fmt.Sprintf("%s", mappedCelcatEvent["id"])

	startDate := fmt.Sprintf("%s", mappedCelcatEvent["start"])
	parsedEvent.Start, _ = time.ParseInLocation("2006-01-02T15:04:05", startDate, loc)

	endDate := fmt.Sprintf("%s", mappedCelcatEvent["end"])
	parsedEvent.End, _ = time.ParseInLocation("2006-01-02T15:04:05", endDate, loc)

	return parsedEvent
}

func GetCalendar(config common.Config, queryPeriod common.Period) common.Calendar {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client := &http.Client{
		Jar: jar,
	}
	common.CheckErr(err)

	session := getRequestLoginToken(client)
	session = logon(client, session, config)

	var calendar common.Calendar
	for _, celcatEvent := range getCalData(client, session, queryPeriod) {
		mappedCelcatEvent := celcatEvent.(map[string]interface{})

		calendar = append(calendar, parseCalEvent(mappedCelcatEvent))
	}

	return calendar
}
