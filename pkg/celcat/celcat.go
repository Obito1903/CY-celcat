package celcat

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type SessionData struct {
	token        string
	FederationId string
	location     url.URL
}

// type celcatCalEventSmall struct {
// 	Id     string `json:"id"`
// 	Start  string `json:"start"`
// 	End    string `json:"end"`
// 	AllDay bool   `json:"allDay"`
// }

type CelcatCalEvent struct {
	Id       string               `json:"id"`
	Start    string               `json:"start"`
	End      string               `json:"end"`
	AllDay   bool                 `json:"allDay"`
	Elements []CelcatEventElement `json:"elements"`
}

type CelcatEventElement struct {
	Label             string `json:"label"`
	Content           string `json:"content"`
	EntityType        int    `json:"entityType"`
	IsStudentSpecific bool   `json:"isStudentSpecific"`
}

// Query the event list from celcat
func getEventList(client *http.Client, celcatUrl url.URL, groupeId string, start time.Time, end time.Time) []CelcatCalEvent {
	headerData := url.Values{
		"start":           {start.Format("2006-01-02")},
		"end":             {end.Format("2006-01-02")},
		"resType":         {"104"},
		"calView":         {"agendaWeek"},
		"federationIds[]": {groupeId},
	}
	resp, err := client.PostForm(celcatUrl.String()+"/Home/GetCalendarData", headerData)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Could not querry calendar data : ", celcatUrl.String(), err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err)
	}
	var celcatEventList []CelcatCalEvent
	err = json.Unmarshal(body, &celcatEventList)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Could not parse calendar data : ", celcatUrl.String(), err)
	}
	return celcatEventList
}

func getEventDetails(client *http.Client, celcatUrl url.URL, event *CelcatCalEvent) {
	headerData := url.Values{
		"eventId": {event.Id},
	}
	resp, err := client.PostForm(celcatUrl.String()+"/Home/GetSideBarEvent", headerData)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Could not querry event data : ", celcatUrl.String(), err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &event)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Could not parse calendar data : ", celcatUrl.String(), err)
	}
}

func GetCalendar(client *http.Client, celcatUrl url.URL, groupeId string, start time.Time, end time.Time) []CelcatCalEvent {
	events := getEventList(client, celcatUrl, groupeId, start, end)
	for idx := range events {
		getEventDetails(client, celcatUrl, &events[idx])
	}
	return events
}
