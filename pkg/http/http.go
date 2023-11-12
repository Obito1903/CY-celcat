package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	config "github.com/Obito1903/CY-celcat/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	config config.Config
}

func (serv Server) icsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varName := vars["groupe"]
	http.ServeFile(w, r, serv.config.ICSPath+varName)
}

func (serv Server) pngHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varName := strings.TrimSuffix(vars["groupe"], ".png")
	if serv.config.HA {
		http.Redirect(w, r, fmt.Sprintf("%s/screenshot?url=%s/%s&width=%d&height=%d", serv.config.PuppeterUrl, serv.config.Url, varName, serv.config.PNGWidth, serv.config.PNGHeigh), http.StatusFound)
	}
	http.ServeFile(w, r, serv.config.PNGPath+varName)
}

func (serv Server) htmlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varName := vars["groupe"]
	http.ServeFile(w, r, serv.config.HTMLPath+varName+".html")
}

func (serv Server) nextAlarmHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varName := vars["groupe"]
	http.ServeFile(w, r, serv.config.NextAlarmPath+varName+".json")
}

func (serv Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, serv.config.HTMLPath+"index.html")
}

func StartServer(config config.Config) {
	serv := Server{config: config}
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{groupe:[[:alnum:]]+.ics}", serv.icsHandler)
	rtr.HandleFunc("/{groupe:[[:alnum:]]+.png}", serv.pngHandler)
	rtr.HandleFunc("/{groupe:[[:alnum:]]+}", serv.htmlHandler)
	rtr.HandleFunc("/", serv.indexHandler)
	rtr.HandleFunc("/{groupe:[[:alnum:]]+}/nextAlarm", serv.nextAlarmHandler)
	http.Handle("/", rtr)

	var handler http.Handler
	if config.AllowCORS {
		handler = cors.Default().Handler(rtr)
	}

	if err := http.ListenAndServe(":"+config.WebPort, handler); err != nil {
		log.Fatal(err)
	}
}
