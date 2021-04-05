package common

import (
	"log"
	"os"
	"time"
)

type CalEvent struct {
	Start    time.Time
	End      time.Time
	Module   string
	Category string
	Prof     string
	Location string
	Id       string
}

type Calendar []CalEvent

type Config struct {
	UserId       string
	UserPassword string
}

type Period struct {
	StartDate time.Time
	EndDate   time.Time
}

// Check if an error occured
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
