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
	federationId string
	location     url.URL
}

// type celcatCalEventSmall struct {
// 	Id     string `json:"id"`
// 	Start  string `json:"start"`
// 	End    string `json:"end"`
// 	AllDay bool   `json:"allDay"`
// }

type celcatCalEvent struct {
	Id       string               `json:"id"`
	Start    string               `json:"start"`
	End      string               `json:"end"`
	AllDay   bool                 `json:"allDay"`
	Elements []celcatEventElement `json:"elements"`
}

type celcatEventElement struct {
	Label             string `json:"label"`
	Content           string `json:"content"`
	EntityType        int    `json:"entityType"`
	IsStudentSpecific bool   `json:"isStudentSpecific"`
}

// Query the event list from celcat
func getEventList(client *http.Client, celcatUrl url.URL, groupeId string, start time.Time, end time.Time) []celcatCalEvent {
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
	var celcatEventList []celcatCalEvent
	err = json.Unmarshal(body, &celcatEventList)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Could not parse calendar data : ", celcatUrl.String(), err)
	}
	return celcatEventList
}

func getEventDetails(client *http.Client, celcatUrl url.URL, event *celcatCalEvent) {
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

func GetCalendar(client *http.Client, celcatUrl url.URL, groupeId string, start time.Time, end time.Time) []celcatCalEvent {
	events := getEventList(client, celcatUrl, groupeId, start, end)
	for _, event := range events {
		getEventDetails(client, celcatUrl, &event)
	}
	return events
}
