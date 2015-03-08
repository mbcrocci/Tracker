package main

import (
	"log"
	"net/http"
	"strconv"

	"labix.org/v2/mgo/bson"
)

var animeList []Anime

func AnimeIndexHandler(w http.ResponseWriter, r *http.Request) {
	templ := LoadTempl()

	if err := colReturn(1).Find(nil).All(&animeList); err != nil {
		log.Println("Can't find any animes")
	}

	templ.Execute(w, animeList)
}

func AnimeAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ep, err := strconv.Atoi(r.Form["episode"][0])
	if err != nil {
		log.Println("Can't convert Form[\"episode\"] to int")
		http.Redirect(w, r, "/anime", http.StatusTemporaryRedirect)
	}

	log.Println("Adding anime:", r.Form)

	err = colReturn(1).Insert(Anime{
		Id:      bson.NewObjectId(),
		Title:   r.Form["title"][0],
		Episode: ep,
	})
	if err != nil {
		log.Println("Can't insert anime")
	}

	http.Redirect(w, r, "/anime", http.StatusTemporaryRedirect)
}

func AnimeIncrementHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	anime, err := SearchAnime(r.Form["title"][0], animeList)
	if err != nil {
		log.Println(err)
	}
	anime.Increment()
	err = colReturn(1).Update(
		bson.M{"title": anime.Title},
		bson.M{"$set": bson.M{"episode": anime.Episode}})
	if err != nil {
		log.Println("Can't update anime int database")
	}

	http.Redirect(w, r, "/anime", http.StatusTemporaryRedirect)
}

func AnimeRemoveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	anime, err := SearchAnime(r.Form["title"][0], animeList)
	if err != nil {
		log.Println(err)
	}

	err = colReturn(1).RemoveId(anime.Id)
	if err != nil {
		log.Println("Can't remove anime from database")
	}
}
