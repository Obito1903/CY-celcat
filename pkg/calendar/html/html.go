package html

import (
	"log"
	"os"

	"github.com/Obito1903/CY-celcat/pkg/calendar"
	"golang.org/x/net/html"
)

func CalendarToHtml(cal calendar.Calendar, pathToTemplate string) {
	file, _ := os.Open(pathToTemplate)
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal("Could not open the template.", err)
	}
}
