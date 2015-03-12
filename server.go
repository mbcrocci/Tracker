package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

func RunServer() {
	log.Println("Starting server on http://localhost:4000")

	path := os.Getenv("GOPATH") + "/src/github.com/mbcrocci/Tracker/"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"templates"))))

	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)

	a := r.PathPrefix("/anime").Subrouter()
	s := r.PathPrefix("/series").Subrouter()

	a.HandleFunc("/", AnimeIndexHandler)

	a.HandleFunc("/add", AnimeAddHandler)

	a.HandleFunc("/increment", AnimeIncrementHandler)

	a.HandleFunc("/remove", AnimeRemoveHandler)

	s.HandleFunc("/", SeriesIndexHandler)

	s.HandleFunc("/new", SeriesNewHandler)

	s.HandleFunc("/add", SeriesAddHandler)

	s.HandleFunc("/increment", SeriesIncrementHandler)

	s.HandleFunc("/remove", SeriesRemoveHandler)

	http.Handle("/", r)
	http.ListenAndServe(":4000", nil)
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	path := os.Getenv("GOPATH") + "/src/github.com/mbcrocci/Tracker/"
	index, err := ioutil.ReadFile(path + "templates/index.html")
	if err != nil {
		log.Println("Can't read index.html")
		os.Exit(2)
	}
	// Generate template
	templ := template.Must(template.New("index").Parse(string(index[:])))

	templ.Execute(rw, nil)
}
