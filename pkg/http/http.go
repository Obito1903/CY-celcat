package http

import (
	"log"
	"net/http"

	config "github.com/Obito1903/CY-celcat/pkg"
	"github.com/gorilla/mux"
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
	varName := vars["groupe"]
	http.ServeFile(w, r, serv.config.PNGPath+varName)
}
func (serv Server) htmlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varName := vars["groupe"]
	http.ServeFile(w, r, serv.config.HTMLPath+varName+".html")
}

func StartServer(config config.Config) {
	serv := Server{config: config}
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{groupe:[[:alnum:]]+.ics}", serv.icsHandler)
	rtr.HandleFunc("/{groupe:[[:alnum:]]+.png}", serv.pngHandler)
	rtr.HandleFunc("/{groupe:[[:alnum:]]+}", serv.htmlHandler)
	http.Handle("/", rtr)
	if err := http.ListenAndServe(":"+config.WebPort, nil); err != nil {
		log.Fatal(err)
	}
}
