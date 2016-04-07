package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"labix.org/v2/mgo/bson"
)

var animeList []Anime

func AnimeIndexHandler(w http.ResponseWriter, r *http.Request) {
	// Load html file
	path := os.Getenv("GOPATH") + "/src/github.com/mbcrocci/Tracker/"
	index, err := ioutil.ReadFile(path + "templates/aindex.html")
	if err != nil {
		log.Println("Can't read aindex.html")
		os.Exit(2)
	}
	// Generate template
	templ := template.Must(template.New("aindex").Parse(string(index[:])))

	err = colReturn(1).Find(nil).Sort("-watching", "completed", "title").All(&animeList)
	if err != nil {
		log.Println("Can't find any animes")
	}

	templ.Execute(w, animeList)
}

func AnimeAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ep, err := strconv.Atoi(r.Form["episode"][0])
	if err != nil {
		log.Println("Can't convert Form[\"episode\"] to int")
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	err = colReturn(1).Insert(Anime{
		ID:      bson.NewObjectId(),
		Title:   r.Form["title"][0],
		Episode: ep,
	})
	if err != nil {
		log.Println("Can't insert anime")
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	log.Println("Added: ", r.Form)
	http.Redirect(w, r, "/anime/", http.StatusFound)
}

func AnimeIncrementHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	anime, err := SearchAnime(r.Form["Title"][0], animeList)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	err = anime.Increment()
	if err != nil {
		log.Println("Anime is completed cant increment")
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	err = colReturn(1).Update(
		bson.M{"title": anime.Title},
		bson.M{"$set": bson.M{"episode": anime.Episode}},
	)
	if err != nil {
		log.Println("Can't update anime int database")
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	log.Println("Incremented: ", r.Form)
	http.Redirect(w, r, "/anime/", http.StatusFound)
}

func AnimeCompleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	anime, err := SearchAnime(r.Form["Title"][0], animeList)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	anime.Complete()

	err = colReturn(1).Update(
		bson.M{"title": anime.Title},
		bson.M{"$set": bson.M{"completed": anime.Completed}},
	)
	if err != nil {
		log.Println("Can't update anime")
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	log.Println("Completed: ", r.Form)
	http.Redirect(w, r, "/anime/", http.StatusFound)
}

func AnimeWatchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	anime, err := SearchAnime(r.Form["Title"][0], animeList)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}
	anime.Show()
	anime.Watch()
	anime.Show()

	err = colReturn(1).Update(
		bson.M{"title": anime.Title},
		bson.M{"$set": bson.M{"watching": anime.Watching}},
	)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}
	log.Println("Watching: ", r.Form)
	http.Redirect(w, r, "/anime/", http.StatusFound)
}

func AnimeRemoveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	err := colReturn(1).Remove(bson.M{"title": r.Form["Title"][0]})
	if err != nil {
		log.Println("Can't remove anime from database: ")
		log.Println(err)
		http.Redirect(w, r, "/anime/", http.StatusFound)
		return
	}

	log.Println("Removed: ", r.Form)
	http.Redirect(w, r, "/anime/", http.StatusFound)
}
