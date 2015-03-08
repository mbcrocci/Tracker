package main

import (
	"errors"
	"log"
	"os"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var db *mgo.Database

func main() {
	session, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Println("Can't connect to mongo")
		time.Sleep(time.Second * 10)
		session, err = mgo.Dial(os.Getenv("Mongo_URL"))
	}
	defer session.Close()
	session.SetSafe(&mgo.Safe{})

	db = session.DB("tracker")

	RunServer()
}

func colReturn(op int) *mgo.Collection {
	if op == 1 {
		return db.C("anime")
	} else {

		return db.C("series")
	}
}

func SearchAnime(title string, list []Anime) (Anime, error) {
	for _, a := range list {
		if a.Title == title {
			return a, nil
		}
	}
	err := errors.New("Can't fin anime: " + title)
	return Anime{bson.NewObjectId(), "err", 0}, err
}

type Anime struct {
	Id      bson.ObjectId `bson:"id"`
	Title   string        `bson:"title"`
	Episode int           `bson:"episode"`
}

func (a *Anime) Increment() {
	a.Episode += 1
}

type Serie struct {
	Id            bson.ObjectId `bson:"id"`
	Title         string        `bson:"title"`
	N_seasons     int           `bson:"n_seasons"`
	Curr_season   int           `bson:"curr_season"`
	Ep_per_season []int         `bson:"ep_per_season"`
}

func (s *Serie) Increment(pos int) {
	s.Ep_per_season[pos] += 1
}
func SearchSerie(title string, list []Serie) (Serie, error) {
	for _, s := range list {
		if s.Title == title {
			return s, nil
		}
	}
	err := errors.New("Can't fin anime: " + title)
	return Serie{bson.NewObjectId(), "err", 0, 0, []int{}}, err
}
