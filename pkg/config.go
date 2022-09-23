package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	// Username for celcat.
	UserName string `json:"userName"`
	// Password for celcat.
	UserPassword string `json:"userPassword"`
	// The host of the celcat instance.
	CelcatHost string `json:"celcatHost"`
	// Run in continuous mode. Will query the calendar periodicly according to the period defined in the config.
	Continuous bool `json:"continuous"`
	// Time in seconds between each query in daemon mode. Default : 1800
	QueryDelay int `json:"queryDelay"`
	// Number of weeks to query. Default : 4
	Weeks int `json:"weeks"`
	// Path to the chrome executable. Default : "/usr/bin/chromium"
	ChromePath string `json:"chromePath"`
	// Enable PNG output (Require Chromium on your computer). Default : false
	PNG bool `json:"png"`
	// Output directory for the PNG output. Default : "out/calendar/png/"
	PNGPath string `json:"pngPath"`
	// Width of the PNG output. Default : 1920
	PNGWidth int `json:"pngWidth"`
	// Height of the PNG output. Default : 1080
	PNGHeigh int `json:"pngHeigh"`
	// Enable HTML output. Default : false
	HTML bool `json:"html"`
	// The template used to render the html page.
	// Default : "web/templates/calendar.go.html"
	HtmlTemplate string `json:"htmlTemplate"`
	// Output directory for the HTML output. Default : "out/calendar/html/"
	HTMLPath string `json:"htmlPath"`
	// Enable ICS output. Default : true
	ICS bool `json:"ics"`
	// Output directory for the ICS output. Default : "out/calendar/ics/"
	ICSPath string `json:"icsPath"`
	// Enable the web server.
	Web bool `json:"web"`
	//Web listen Port. Default : 8080
	WebPort string `json:"webPort"`
	// List of groups to query
	Groupes []Groupe `json:"groupes"`
}

type Groupe struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func ReadConfig(path string) Config {
	configJson, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal("Couldn't open config file", err)
	}
	// Wait for the file to be fully read
	defer configJson.Close()
	// Convert file int a byte field
	configByte, _ := ioutil.ReadAll(configJson)

	config := Config{
		Continuous:   false,
		QueryDelay:   1800,
		Weeks:        4,
		ChromePath:   "/usr/bin/chromium",
		PNG:          false,
		PNGPath:      "out/calendar/png/",
		PNGWidth:     1920,
		PNGHeigh:     1080,
		HTML:         false,
		HtmlTemplate: "web/templates/calendar.go.html",
		HTMLPath:     "out/calendar/html/",
		ICS:          false,
		ICSPath:      "out/calendar/ics/",
		Web:          false,
		WebPort:      "8080",
	}

	json.Unmarshal(configByte, &config)
	return config
}

func Configure() Config {
	config := ReadConfig("./config.json")

	flag.StringVar(&config.UserName, "user", config.UserName, "Username for celcat.")
	flag.StringVar(&config.UserPassword, "pass", config.UserPassword, "Password for celcat.")
	flag.StringVar(&config.CelcatHost, "host", config.CelcatHost, "Host of the celcat instance to scrap.")

	flag.BoolVar(&config.Web, "web", config.Web, "Enable the web server.")
	flag.StringVar(&config.WebPort, "port", config.WebPort, "Port to listen to for http request.")

	// Daemon config
	flag.BoolVar(&config.Continuous, "loop", config.Continuous, "Run in continuous mode. Will query the calendar periodicly according to the period defined in the config.")
	flag.IntVar(&config.QueryDelay, "delay", config.QueryDelay, "time in seconds between each query in daemon mode.")
	flag.IntVar(&config.Weeks, "weeks", config.Weeks, "Number of weeks to query.")

	// HTML config
	flag.StringVar(&config.HTMLPath, "htmlOut", config.HTMLPath, "Output directory for the HTML output")
	flag.BoolVar(&config.HTML, "html", config.HTML, "Enable HTML output.")
	flag.StringVar(&config.HtmlTemplate, "template", config.HtmlTemplate, "Path to the HTML template for the HTML Calendar.")

	// ICS config
	flag.StringVar(&config.ICSPath, "icsOut", config.ICSPath, "Output directory for the ICS output.")
	flag.BoolVar(&config.ICS, "ics", config.ICS, "Enable ICS output.")

	// PNG config
	flag.StringVar(&config.ChromePath, "chromePath", config.ChromePath, "Path to the chrome executable to render PNG.")
	flag.StringVar(&config.PNGPath, "pngOut", config.PNGPath, "Output directory for the PNG output.")
	flag.BoolVar(&config.PNG, "png", config.PNG, "Enable PNG output (Require HTML output enable and Chromium on your computer).")
	flag.IntVar(&config.PNGHeigh, "height", config.PNGHeigh, "Height of the PNG output.")
	flag.IntVar(&config.PNGWidth, "width", config.PNGWidth, "Width of the PNG output.")

	flag.Parse()

	os.MkdirAll(config.HTMLPath, os.ModePerm)
	os.MkdirAll(config.ICSPath, os.ModePerm)
	os.MkdirAll(config.PNGPath, os.ModePerm)

	return config
}
