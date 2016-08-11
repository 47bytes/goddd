// Package inmem provides in-memory implementations of all the domain repositories.
package inmem

import "github.com/marcusolsson/goddd"

import "sync"

type cargoRepository struct {
	mtx    sync.RWMutex
	cargos map[goddd.TrackingID]*goddd.Cargo
}

func (r *cargoRepository) Store(c *goddd.Cargo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.cargos[c.TrackingID] = c
	return nil
}

func (r *cargoRepository) Find(trackingID goddd.TrackingID) (*goddd.Cargo, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.cargos[trackingID]; ok {
		return val, nil
	}
	return nil, goddd.ErrUnknownCargo
}

func (r *cargoRepository) FindAll() []*goddd.Cargo {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	c := make([]*goddd.Cargo, 0, len(r.cargos))
	for _, val := range r.cargos {
		c = append(c, val)
	}
	return c
}

// NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewCargoRepository() goddd.CargoRepository {
	return &cargoRepository{
		cargos: make(map[goddd.TrackingID]*goddd.Cargo),
	}
}

type locationRepository struct {
	locations map[goddd.UNLocode]goddd.Location
}

func (r *locationRepository) Find(locode goddd.UNLocode) (goddd.Location, error) {
	if l, ok := r.locations[locode]; ok {
		return l, nil
	}
	return goddd.Location{}, goddd.ErrUnknownLocation
}

func (r *locationRepository) FindAll() []goddd.Location {
	l := make([]goddd.Location, 0, len(r.locations))
	for _, val := range r.locations {
		l = append(l, val)
	}
	return l
}

// NewLocationRepository returns a new instance of a in-memory location repository.
func NewLocationRepository() goddd.LocationRepository {
	r := &locationRepository{
		locations: make(map[goddd.UNLocode]goddd.Location),
	}

	r.locations[goddd.SESTO] = goddd.Stockholm
	r.locations[goddd.AUMEL] = goddd.Melbourne
	r.locations[goddd.CNHKG] = goddd.Hongkong
	r.locations[goddd.JNTKO] = goddd.Tokyo
	r.locations[goddd.NLRTM] = goddd.Rotterdam
	r.locations[goddd.DEHAM] = goddd.Hamburg

	return r
}

type voyageRepository struct {
	voyages map[goddd.VoyageNumber]*goddd.Voyage
}

func (r *voyageRepository) Find(voyageNumber goddd.VoyageNumber) (*goddd.Voyage, error) {
	if v, ok := r.voyages[voyageNumber]; ok {
		return v, nil
	}

	return nil, goddd.ErrUnknownVoyage
}

// NewVoyageRepository returns a new instance of a in-memory voyage repository.
func NewVoyageRepository() goddd.VoyageRepository {
	r := &voyageRepository{
		voyages: make(map[goddd.VoyageNumber]*goddd.Voyage),
	}

	r.voyages[goddd.V100.Number] = goddd.V100
	r.voyages[goddd.V300.Number] = goddd.V300
	r.voyages[goddd.V400.Number] = goddd.V400

	r.voyages[goddd.V0100S.Number] = goddd.V0100S
	r.voyages[goddd.V0200T.Number] = goddd.V0200T
	r.voyages[goddd.V0300A.Number] = goddd.V0300A
	r.voyages[goddd.V0301S.Number] = goddd.V0301S
	r.voyages[goddd.V0400S.Number] = goddd.V0400S

	return r
}

type handlingEventRepository struct {
	mtx    sync.RWMutex
	events map[goddd.TrackingID][]goddd.HandlingEvent
}

func (r *handlingEventRepository) Store(e goddd.HandlingEvent) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	// Make array if it's the first event with this tracking ID.
	if _, ok := r.events[e.TrackingID]; !ok {
		r.events[e.TrackingID] = make([]goddd.HandlingEvent, 0)
	}
	r.events[e.TrackingID] = append(r.events[e.TrackingID], e)
}

func (r *handlingEventRepository) QueryHandlingHistory(trackingID goddd.TrackingID) goddd.HandlingHistory {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return goddd.HandlingHistory{HandlingEvents: r.events[trackingID]}
}

// NewHandlingEventRepository returns a new instance of a in-memory handling event repository.
func NewHandlingEventRepository() goddd.HandlingEventRepository {
	return &handlingEventRepository{
		events: make(map[goddd.TrackingID][]goddd.HandlingEvent),
	}
}
