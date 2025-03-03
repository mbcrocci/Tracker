package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"labix.org/v2/mgo/bson"
)

var seriesList []Serie

func SeriesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// Load html file
	path := os.Getenv("GOPATH") + "/src/github.com/mbcrocci/Tracker/"
	index, err := ioutil.ReadFile(path + "templates/sindex.html")
	if err != nil {
		log.Println("Can't read sindex.html")
		os.Exit(2)
	}
	// Generate template
	templ, err := template.New("series").Parse(string(index[:]))
	if err != nil {
		log.Println("Can't generate template beacause:", err)
		os.Exit(1)
	}
	if err := colReturn(2).Find(nil).Sort("title").All(&seriesList); err != nil {
		log.Println("Can't find any series")
	}

	templ.Execute(w, seriesList)
}

func SeriesNewHandler(w http.ResponseWriter, r *http.Request) {
	// Load html file
	path := os.Getenv("GOPATH") + "/src/github.com/mbcrocci/Tracker/"
	index, err := ioutil.ReadFile(path + "templates/new.html")
	if err != nil {
		log.Println("Can't read new.html")
		os.Exit(2)
	}
	// Generate template
	templ := template.Must(template.New("new").Parse(string(index[:])))
	templ.Execute(w, nil)
}

func SeriesAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	nSeasons, err := strconv.Atoi(r.Form["n_seasons"][0])
	if err != nil {
		log.Println("Can't convert n_seasons to int")
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	currSeason, err := strconv.Atoi(r.Form["curr_season"][0])
	if err != nil {
		log.Println("Can't convert curr_seasons to int")
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	currEp, err := strconv.Atoi(r.Form["curr_ep"][0])
	if err != nil {
		log.Println("Can't convert curr_seasons to int")
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	var epPerSeason []int
	for i := 0; i < nSeasons; i++ {
		sName := "s" + strconv.Itoa(i)
		ep, err := strconv.Atoi(r.Form[sName][0])
		if err != nil {
			log.Println("Can't conver ep to int")
			http.Redirect(w, r, "/series/", http.StatusFound)
			return
		}
		epPerSeason = append(epPerSeason, ep)
	}

	err = colReturn(2).Insert(Serie{
		ID:          bson.NewObjectId(),
		Title:       r.Form["title"][0],
		NSeasons:    nSeasons,
		CurrSeason:  currSeason,
		CurrEp:      currEp,
		EpPerSeason: epPerSeason,
	})
	if err != nil {
		log.Println("Can't insert serie")
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	log.Println("Added: ", r.Form)
	http.Redirect(w, r, "/series/", http.StatusFound)

}
func SeriesIncrementHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	serie, err := SearchSerie(r.Form["Title"][0], seriesList)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	if err = serie.Increment(); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	err = colReturn(2).Update(
		bson.M{"title": serie.Title},
		bson.M{"$set": bson.M{
			"currSeason": serie.CurrSeason,
			"currEp":     serie.CurrEp,
		}},
	)
	if err != nil {
		log.Println("Can't update serie into database")
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	log.Println("Incrementing: ", r.Form)
	http.Redirect(w, r, "/series/", http.StatusFound)
}

func SeriesRemoveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	err := colReturn(2).Remove(bson.M{"title": r.Form["Title"][0]})
	if err != nil {
		log.Println("Can't remove anime from database")
		http.Redirect(w, r, "/series/", http.StatusFound)
		return
	}

	log.Println("Removed: ", r.Form)
	http.Redirect(w, r, "/series/", http.StatusFound)
}
