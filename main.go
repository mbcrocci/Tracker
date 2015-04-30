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

		// try to reconnect to database after 10 seconds
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
	}

	return db.C("series")

}

// SearchAnime  for the anime on a list
func SearchAnime(title string, list []Anime) (Anime, error) {
	for _, a := range list {
		if a.Title == title {
			return a, nil
		}
	}
	err := errors.New("Can't fin anime: " + title)
	return Anime{bson.NewObjectId(), "err", 0, false}, err
}

// Anime holds information on a certain anime
type Anime struct {
	ID        bson.ObjectId `bson:"id"`
	Title     string        `bson:"title"`
	Episode   int           `bson:"episode"`
	Completed bool          `bson:"completed"`
}

// Increment increases a anime episode by one unless it is completed
func (a *Anime) Increment() error {
	if !a.Completed {
		a.Episode++
		return nil
	}

	return errors.New("Anime completed")
}

// Complete changes the Completed field to true
func (a *Anime) Complete() {
	a.Completed = true
}

// Serie is a more complex tracker than anime
type Serie struct {
	ID          bson.ObjectId `bson:"id"`
	Title       string        `bson:"title"`
	NSeasons    int           `bson:"nseasons"`
	CurrSeason  int           `bson:"currSeason"`
	CurrEp      int           `bson:"currEp"`
	EpPerSeason []int         `bson:"epPerSeason"`
}

// Increment increases a series current episode by one.
// It checks to see if the current episode is the last one of the season.
// If so, checks if the current season is the last one, in which case, it can't increment.
// Other wise, it increases the season and sets the current episode to 1.
func (s *Serie) Increment() error {
	if s.CurrEp == s.EpPerSeason[s.CurrSeason-1] {
		if s.CurrSeason == s.NSeasons {
			return errors.New("Can't increment")
		}

		s.CurrSeason++
		s.CurrEp = 1
		return nil
	}
	s.CurrEp++
	return nil
}

// SearchSerie looks for a serie with a certain title inside a slice of Series
func SearchSerie(title string, list []Serie) (Serie, error) {
	for _, s := range list {
		if s.Title == title {
			return s, nil
		}
	}
	err := errors.New("Can't fin serie: " + title)
	return Serie{bson.NewObjectId(), "err", 0, 0, 0, []int{}}, err
}
