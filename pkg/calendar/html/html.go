package html

import (
	"fmt"
	"log"
	"os"
	"time"

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
		if event.End.Before(week.Add(7*24*time.Hour)) && event.Start.After(week) {
			start := time.Date(week.Year(), week.Month(), week.Day(), event.Start.Hour(), event.Start.Minute(), 0, 0, time.Local)
			end := time.Date(week.Year(), week.Month(), week.Day(), event.End.Hour(), event.End.Minute(), 0, 0, time.Local)
			fmt.Println(start, end)
			if end.After(htmlCal.MaxEnd) {
				htmlCal.MaxEnd = end
			}
			if start.Before(htmlCal.MaxStart) {
				htmlCal.MaxStart = start
			}
		}
	}
	htmlCal.MaxStart = time.Date(htmlCal.MaxStart.Year(), htmlCal.MaxStart.Month(), htmlCal.MaxStart.Day(), htmlCal.MaxStart.Hour(), (htmlCal.MaxStart.Minute()/10)*10, 0, 0, time.Local)
	for h := htmlCal.MaxStart; h.Unix() < htmlCal.MaxEnd.Unix(); h = h.Add(time.Minute * 30) {
		htmlCal.Horaires = append(htmlCal.Horaires, h.Format("15h04"))
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
	return htmlCal
}

// func CalendarToHtml(cal calendar.Calendar, templatePath string, week time.Time) string {
// 	t := template.New("calendar")
// 	t, err := template.ParseFiles(templatePath)
// 	if err != nil {
// 		log.Fatal("Could not load template.", err)
// 	}

// }

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
