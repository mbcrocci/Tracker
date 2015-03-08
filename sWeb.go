package main

import (
	"log"
	"net/http"
	"strconv"

	"labix.org/v2/mgo/bson"
)

var seriesList []Serie

func SeriesIndexHandler(w http.ResponseWriter, r *http.Request) {
	templ := LoadTempl()

	if err := colReturn(2).Find(nil).All(&seriesList); err != nil {
		log.Println("Can't find any series")
	}

	templ.Execute(w, seriesList)
}

func SeriesaNewHandler(w http.ResponseWriter, r *http.Request) {
	//(TODO): new template
}

func SeriesAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	n_seasons, err := strconv.Atoi(r.Form["n_seasons"][0])
	if err != nil {
		log.Println("Can't convert n_seasons to int")
		http.Redirect(w, r, "/serues", http.StatusTemporaryRedirect)
	}

	curr_season, err := strconv.Atoi(r.Form["curr_season"][0])
	if err != nil {
		log.Println("Can't convert curr_seasons to int")
		http.Redirect(w, r, "/serues", http.StatusTemporaryRedirect)
	}

	curr_ep, err := strconv.Atoi(r.Form["curr_ep"][0])
	if err != nil {
		log.Println("Can't convert curr_seasons to int")
		http.Redirect(w, r, "/serues", http.StatusTemporaryRedirect)
	}

	err = colReturn(2).Insert(Serie{
		Id:            bson.NewObjectId(),
		Title:         r.Form["title"][0],
		N_seasons:     n_seasons,
		Curr_season:   curr_season,
		Ep_per_season: []int{curr_ep},
	})
	if err != nil {
		log.Println("Can't insert serie")
	}

	http.Redirect(w, r, "/serie", http.StatusTemporaryRedirect)

}
func SeriesIncrementHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	serie, err := SearchSerie(r.Form["title"][0], seriesList)
	if err != nil {
		log.Println(err)
	}

	serie.Increment(serie.Curr_season)

	err = colReturn(2).Update(
		bson.M{"title": serie.Title},
		bson.M{"$set": bson.M{
			"episode": serie.Ep_per_season[serie.Curr_season]}})
	if err != nil {
		log.Println("Can't update serie into database")
	}

	http.Redirect(w, r, "/serie", http.StatusTemporaryRedirect)
}

func SeriesRemoveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	serie, err := SearchSerie(r.Form["title"][0], seriesList)
	if err != nil {
		log.Println(err)
	}

	err = colReturn(2).RemoveId(serie.Id)
	if err != nil {
		log.Println("Can't remove anime from database")
	}

	http.Redirect(w, r, "/serie", http.StatusTemporaryRedirect)
}
