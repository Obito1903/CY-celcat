package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Obito1903/CY-celcat/internal/cyCelcat"
	"github.com/Obito1903/CY-celcat/pkg/http"

	config "github.com/Obito1903/CY-celcat/pkg"
)

func initDirs(config config.Config) {
	err := os.MkdirAll(config.HTMLPath, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create HTML Path.", err)
	}
	err = os.MkdirAll(config.ICSPath, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create ICS Path.", err)
	}
	err = os.MkdirAll(config.PNGPath, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create PNG Path.", err)
	}
	os.MkdirAll(config.NextAlarmPath, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create NextAlarm Path.", err)
	}
}

func main() {
	config := config.Configure()
	initDirs(config)
	fmt.Println(config)
	if config.Web {
		go http.StartServer(config)
	}
	if config.Continuous {
		for {
			today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
			cyCelcat.Query(config, cyCelcat.Period{Start: today, End: today.Add(time.Hour * 24 * 7 * time.Duration(config.Weeks))})
			time.Sleep(time.Duration(config.QueryDelay) * time.Second)
		}
	} else {
		today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		cyCelcat.Query(config, cyCelcat.Period{Start: today, End: today.Add(time.Hour * 24 * 7 * time.Duration(config.Weeks))})

	}

}
