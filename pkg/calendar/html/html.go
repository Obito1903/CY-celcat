package html

import (
	"log"
	"os"
	"time"

	"html/template"

	"github.com/Obito1903/CY-celcat/pkg/calendar"
)

type HtmlEvent struct {
	Event    calendar.Event
	Top      float32
	Height   float32
	TimeSpan string
}

type HtmlDay struct {
	Name   string
	Events []HtmlEvent
}

type HtmlCalendar struct {
	Horaires []string
	Days     []HtmlDay
	MaxSpan  time.Duration
}

func CalendarToHtml(cal calendar.Calendar, templatePath string) string {
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
