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
	log.Println("Starting server on http://localhost:3000")

	path := "/Users/mbcrocci/Projects/gocode/src/github.com/mbcrocci/Tracker/"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"static"))))

	r := mux.NewRouter()

	a := r.PathPrefix("/anime").Subrouter()

	a.HandleFunc("/", AnimeIndexHandler)

	// Add new anime handler
	a.HandleFunc("/add", AnimeAddHandler)

	// Increment handler
	a.HandleFunc("/increment", AnimeIncrementHandler)

	// Remove handler
	a.HandleFunc("/remove", AnimeRemoveHandler)

	s := r.PathPrefix("/series").Subrouter()

	s.HandleFunc("/", SeriesIndexHandler)

	// Add new serie handler
	s.HandleFunc("/add", SeriesAddHandler)

	// Increment handler
	s.HandleFunc("/increment", SeriesIncrementHandler)

	// Remove handler
	s.HandleFunc("/remove", SeriesRemoveHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}

func LoadTempl() *template.Template {
	// Load html file
	path := "/Users/mbcrocci/Projects/gocode/src/github.com/mbcrocci/Tracker/"
	index, err := ioutil.ReadFile(path + "templates/index.html")
	if err != nil {
		log.Println("Can't read index.html")
		os.Exit(2)
	}
	// Generate template
	return template.Must(template.New("index").Parse(string(index[:])))

}
