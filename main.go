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
	}

	return db.C("series")

}

// SearchAnime  for the anime on a lis
func SearchAnime(title string, list []Anime) (Anime, error) {
	for _, a := range list {
		if a.Title == title {
			return a, nil
		}
	}
	err := errors.New("Can't fin anime: " + title)
	return Anime{bson.NewObjectId(), "err", 0}, err
}

// Anime holds information on a certain anime
type Anime struct {
	ID      bson.ObjectId `bson:"id"`
	Title   string        `bson:"title"`
	Episode int           `bson:"episode"`
}

// Increment increases a anime episode by one
func (a *Anime) Increment() {
	a.Episode++
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

func test(heelo int) bool {
	return true
}

// Increment increase a series current episode by one.
// It checks to see if the current episode is the one of the season.
// If so, checks if the current season is the last one, in which case, it
// ca't increment.
// Other wise it increases the season and sets the current episode to 1.
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
