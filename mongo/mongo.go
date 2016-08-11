package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/marcusolsson/goddd"
)

type cargoRepository struct {
	db      string
	session *mgo.Session
}

func (r *cargoRepository) Store(cargo *goddd.Cargo) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("cargo")

	_, err := c.Upsert(bson.M{"trackingid": cargo.TrackingID}, bson.M{"$set": cargo})

	return err
}

func (r *cargoRepository) Find(trackingID goddd.TrackingID) (*goddd.Cargo, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("cargo")

	var result goddd.Cargo
	if err := c.Find(bson.M{"trackingid": trackingID}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, goddd.ErrUnknownCargo
		}
		return nil, err
	}

	return &result, nil
}

func (r *cargoRepository) FindAll() []*goddd.Cargo {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("cargo")

	var result []*goddd.Cargo
	if err := c.Find(bson.M{}).All(&result); err != nil {
		return []*goddd.Cargo{}
	}

	return result
}

// NewCargoRepository returns a new instance of a MongoDB cargo repository.
func NewCargoRepository(db string, session *mgo.Session) (goddd.CargoRepository, error) {
	r := &cargoRepository{
		db:      db,
		session: session,
	}

	index := mgo.Index{
		Key:        []string{"trackingid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("cargo")

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	return r, nil
}

type locationRepository struct {
	db      string
	session *mgo.Session
}

func (r *locationRepository) Find(locode goddd.UNLocode) (goddd.Location, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	var result goddd.Location
	if err := c.Find(bson.M{"unlocode": locode}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return goddd.Location{}, goddd.ErrUnknownLocation
		}
		return goddd.Location{}, err
	}

	return result, nil
}

func (r *locationRepository) FindAll() []goddd.Location {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	var result []goddd.Location
	if err := c.Find(bson.M{}).All(&result); err != nil {
		return []goddd.Location{}
	}

	return result
}

func (r *locationRepository) store(l goddd.Location) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	_, err := c.Upsert(bson.M{"unlocode": l.UNLocode}, bson.M{"$set": l})

	return err
}

// NewLocationRepository returns a new instance of a MongoDB location repository.
func NewLocationRepository(db string, session *mgo.Session) (goddd.LocationRepository, error) {
	r := &locationRepository{
		db:      db,
		session: session,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	index := mgo.Index{
		Key:        []string{"unlocode"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	initial := []goddd.Location{
		goddd.Stockholm,
		goddd.Melbourne,
		goddd.Hongkong,
		goddd.Tokyo,
		goddd.Rotterdam,
		goddd.Hamburg,
	}

	for _, l := range initial {
		r.store(l)
	}

	return r, nil
}

type voyageRepository struct {
	db      string
	session *mgo.Session
}

func (r *voyageRepository) Find(voyageNumber goddd.VoyageNumber) (*goddd.Voyage, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("voyage")

	var result goddd.Voyage
	if err := c.Find(bson.M{"number": voyageNumber}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, goddd.ErrUnknownVoyage
		}
		return nil, err
	}

	return &result, nil
}

func (r *voyageRepository) store(v *goddd.Voyage) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("voyage")

	_, err := c.Upsert(bson.M{"number": v.Number}, bson.M{"$set": v})

	return err
}

// NewVoyageRepository returns a new instance of a MongoDB voyage repository.
func NewVoyageRepository(db string, session *mgo.Session) (goddd.VoyageRepository, error) {
	r := &voyageRepository{
		db:      db,
		session: session,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("voyage")

	index := mgo.Index{
		Key:        []string{"number"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	initial := []*goddd.Voyage{
		goddd.V100,
		goddd.V300,
		goddd.V400,
		goddd.V0100S,
		goddd.V0200T,
		goddd.V0300A,
		goddd.V0301S,
		goddd.V0400S,
	}

	for _, v := range initial {
		r.store(v)
	}

	return r, nil
}

type handlingEventRepository struct {
	db      string
	session *mgo.Session
}

func (r *handlingEventRepository) Store(e goddd.HandlingEvent) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("handling_event")

	_ = c.Insert(e)
}

func (r *handlingEventRepository) QueryHandlingHistory(trackingID goddd.TrackingID) goddd.HandlingHistory {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("handling_event")

	var result []goddd.HandlingEvent
	_ = c.Find(bson.M{"trackingid": trackingID}).All(&result)

	return goddd.HandlingHistory{HandlingEvents: result}
}

// NewHandlingEventRepository returns a new instance of a MongoDB handling event repository.
func NewHandlingEventRepository(db string, session *mgo.Session) goddd.HandlingEventRepository {
	return &handlingEventRepository{
		db:      db,
		session: session,
	}
}
