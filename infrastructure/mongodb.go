package infrastructure

import (
	"log"
	"os"

	"github.com/marcusolsson/goddd/domain/location"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type locationRepositoryMongoDB struct {
}

func (r *locationRepositoryMongoDB) Store(l location.Location) {
	session, err := mgo.Dial(os.Getenv("MONGOHQ_URL"))

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("app30695645").C("locations")
	err = c.Insert(&l)

	if err != nil {
		log.Fatal(err)
	}
}

func (r *locationRepositoryMongoDB) Find(locode location.UNLocode) (location.Location, error) {
	return location.Location{}, nil
}

func (r *locationRepositoryMongoDB) FindAll() []location.Location {
	session, err := mgo.Dial(os.Getenv("MONGOHQ_URL"))

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("app30695645").C("locations")

	var result []location.Location
	err = c.Find(bson.M{}).All(&result)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func NewLocationRepositoryMongoDB() location.LocationRepository {
	session, err := mgo.Dial(os.Getenv("MONGOHQ_URL"))

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("app30695645").C("locations")
	err = c.EnsureIndexKey("unlocode")

	if err != nil {
		panic(err)
	}

	r := &locationRepositoryMongoDB{}

	r.Store(location.Stockholm)
	r.Store(location.Hamburg)
	r.Store(location.Chicago)

	return r
}