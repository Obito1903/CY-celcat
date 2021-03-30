package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"golang.org/x/net/publicsuffix"
)

type CalEvent struct {
	start    time.Time
	end      time.Time
	module   string
	category string
	prof     string
	location string
	id       string
}

type SessionData struct {
	token       string
	Antiforgery *http.Cookie
}

type Config struct {
	userId       string
	userPassword string
}

type Period struct {
	startDate time.Time
	endDate   time.Time
}

// Check if an error occured
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func getConfig(path string) Config {
	var result map[string]interface{}
	var config Config

	configFile, err := os.Open(path)
	checkErr(err)
	configFileStream, err := ioutil.ReadAll(configFile)
	checkErr(err)

	json.Unmarshal(configFileStream, &result)
	config.userId, config.userPassword = fmt.Sprintf("%s", result["userId"]), fmt.Sprintf("%s", result["userPassword"])
	return config
}

// Get RequestToken and cookie from services-web.u-cergy.fr wich are necessary to send a login request to celcat
func getRequestLoginToken(client *http.Client) SessionData {
	var token SessionData
	//send a get request to the calendar service to obtain the requestToken
	resp, err := client.Get("https://services-web.u-cergy.fr/calendar/LdapLogin")
	checkErr(err)
	//Store the request cookie
	for _, cookie := range resp.Cookies() {
		token.Antiforgery = cookie
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
func logon(client *http.Client, session SessionData, config Config) SessionData {
	// Init the header to be sent
	headerData := url.Values{
		"Name":                       {config.userId},
		"Password":                   {config.userPassword},
		"__RequestVerificationToken": {session.token},
	}
	// Setup a login request
	req, err := http.NewRequest("POST", "https://services-web.u-cergy.fr/calendar/LdapLogin/Logon", strings.NewReader(headerData.Encode()))
	checkErr(err)
	req.AddCookie(session.Antiforgery)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Send the login request
	resp, err := client.Do(req)
	checkErr(err)
	fmt.Println(resp.StatusCode)
	return session
}

// Query the calendar and return the interface containing all calendar info
func getCalData(client *http.Client, session SessionData, queryPeriod Period) []interface{} {
	var calendarInterface []interface{}

	fmt.Println("Requesting calendar from", queryPeriod.startDate.Format("2006-01-02"), "to", queryPeriod.endDate.Format("2006-01-02"))
	// Setup request header
	headerData := url.Values{
		"start":           {queryPeriod.startDate.Format("2006-01-02")},
		"end":             {queryPeriod.endDate.Format("2006-01-02")},
		"resType":         {"104"},
		"calView":         {"agendaWeek"},
		"federationIds[]": {"22014815"},
	}
	// Send request with cookies from the cookieJar
	resp, err := client.PostForm("https://services-web.u-cergy.fr/calendar/Home/GetCalendarData", headerData)
	checkErr(err)

	// Parse the response and put all data into an interface
	body, err := io.ReadAll(resp.Body)
	checkErr(err)
	fmt.Println(resp.StatusCode)
	err = json.Unmarshal(body, &calendarInterface)
	checkErr(err)
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
		for i := 2; i < len(parsed); i++ {
			location = location + strings.Trim(parsed[i], "\r\n")
		}
		prof = strings.Trim(parsed[len(parsed)-1], "\r\n")
	}

	return module, location, prof
}

// Create a callEvent from the json data of one celcat event
func parseCalEvent(baseEvent interface{}) CalEvent {
	var parsedEvent CalEvent
	loc, _ := time.LoadLocation("Europe/Paris")

	mappedCelcatEvent := baseEvent.(map[string]interface{})

	parsedEvent.category = fmt.Sprintf("%s", mappedCelcatEvent["eventCategory"])

	parsedEvent.module, parsedEvent.location, parsedEvent.prof = parseDescription(fmt.Sprintf("%s", mappedCelcatEvent["description"]), parsedEvent.category)
	parsedEvent.id = fmt.Sprintf("%s", mappedCelcatEvent["id"])

	startDate := fmt.Sprintf("%s", mappedCelcatEvent["start"])
	parsedEvent.start, _ = time.ParseInLocation("2006-01-02T15:04:05", startDate, loc)

	endDate := fmt.Sprintf("%s", mappedCelcatEvent["start"])
	parsedEvent.end, _ = time.ParseInLocation("2006-01-02T15:04:05", endDate, loc)

	return parsedEvent
}

// Create the ICS serialized string from the list of json celcat event
func creatICS(calendarInterface []interface{}) string {

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for _, celcatEvent := range calendarInterface {
		mappedCelcatEvent := celcatEvent.(map[string]interface{})
		calEvent := parseCalEvent(mappedCelcatEvent)
		event := cal.AddEvent(calEvent.id)
		event.SetStartAt(calEvent.start)
		event.SetEndAt(calEvent.end)
		event.SetLocation(calEvent.location)
		event.SetOrganizer(calEvent.prof)
		event.SetSummary(calEvent.module)
	}
	return cal.Serialize()
}

// Save the ICS serialized string to a file
func saveICS(icsCal string) {
	f, err := os.Create("data.ics")
	checkErr(err)
	defer f.Close()
	_, err = f.WriteString(icsCal)
	checkErr(err)
}

func main() {

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	checkErr(err)
	client := &http.Client{
		Jar: jar,
	}

	configPath := "config.json"
	queryPeriod := Period{
		startDate: time.Now(),
		endDate:   time.Now().Add(time.Hour * 24 * 30),
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
			queryPeriod.startDate, err = time.ParseInLocation("2006-01-02", os.Args[i+1], loc)
			checkErr(err)
			queryPeriod.endDate, err = time.ParseInLocation("2006-01-02", os.Args[i+2], loc)
			checkErr(err)
		}

	}

	config := getConfig(configPath)

	session := getRequestLoginToken(client)
	session = logon(client, session, config)
	saveICS(creatICS(getCalData(client, session, queryPeriod)))

}
