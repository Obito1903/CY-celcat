package html

import (
	"log"
	"os"
	"time"

	"html/template"

	"github.com/Obito1903/CY-celcat/pkg/calendar"
)

type htmlEvent struct {
	Event    calendar.Event
	Top      float32
	Height   float32
	TimeSpan string
}

type htmlCalendar struct {
	Horaires []string
	Days     map[time.Weekday]htmlEvent
	MaxStart time.Time
	MaxEnd   time.Time
}

func calcHoraires(cal calendar.Calendar, week time.Time, htmlCal *htmlCalendar) {
	htmlCal.MaxStart = week.Add(time.Hour * 24)
	htmlCal.MaxEnd = week
	for _, event := range cal.Events {
		if event.End.Before(week.Add(7*24*time.Hour)) && event.Start.Before(week) {
			if event.End.After(htmlCal.MaxEnd) {
				htmlCal.MaxEnd = event.End
			}
			if event.Start.Before(htmlCal.MaxStart) {

			}
		}
	}
}

// func eventToHtmlEvent(htmlCal htmlCalendar, event calendar.Event) htmlEvent {
// 	var htmlEv htmlEvent
// }

func calToHtmlCal(cal calendar.Calendar, week time.Time) htmlCalendar {
	var htmlCal htmlCalendar
	for _, event := range cal.Events {
		if event.End.Before(week.Add(7*24*3600*1000*1000*1000)) && event.Start.Before(week) {
			htmlCal.Days[event.Start.Weekday()] = htmlEvent{
				Event:    event,
				TimeSpan: event.Start.Format("15h04") + "-" + event.End.Format("15h04"),
			}

		}
	}

}

func CalendarToHtml(cal calendar.Calendar, templatePath string, week time.Time) string {
	t := template.New("calendar")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("Could not load template.", err)
	}

}

func HtmlToFile(htmlCal string, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("Could not save HTML.", err)
	}
	defer f.Close()
	_, err = f.WriteString(htmlCal)
	if err != nil {
		log.Fatal("Could not save HTML.", err)
	}
}
