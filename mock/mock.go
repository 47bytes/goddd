package mock

import "github.com/marcusolsson/goddd"

// CargoRepository is a mock cargo repository.
type CargoRepository struct {
	StoreFn      func(c *goddd.Cargo) error
	StoreInvoked bool

	FindFn      func(id goddd.TrackingID) (*goddd.Cargo, error)
	FindInvoked bool

	FindAllFn      func() []*goddd.Cargo
	FindAllInvoked bool
}

// Store calls the StoreFn.
func (r *CargoRepository) Store(c *goddd.Cargo) error {
	r.StoreInvoked = true
	return r.StoreFn(c)
}

// Find calls the FindFn.
func (r *CargoRepository) Find(id goddd.TrackingID) (*goddd.Cargo, error) {
	r.FindInvoked = true
	return r.FindFn(id)
}

// FindAll calls the FindAllFn.
func (r *CargoRepository) FindAll() []*goddd.Cargo {
	r.FindAllInvoked = true
	return r.FindAllFn()
}

// LocationRepository is a mock location repository.
type LocationRepository struct {
	FindFn      func(goddd.UNLocode) (goddd.Location, error)
	FindInvoked bool

	FindAllFn      func() []goddd.Location
	FindAllInvoked bool
}

// Find calls the FindFn.
func (r *LocationRepository) Find(locode goddd.UNLocode) (goddd.Location, error) {
	r.FindInvoked = true
	return r.FindFn(locode)
}

// FindAll calls the FindAllFn.
func (r *LocationRepository) FindAll() []goddd.Location {
	r.FindAllInvoked = true
	return r.FindAllFn()
}

// VoyageRepository is a mock voyage repository.
type VoyageRepository struct {
	FindFn      func(goddd.VoyageNumber) (*goddd.Voyage, error)
	FindInvoked bool
}

// Find calls the FindFn.
func (r *VoyageRepository) Find(number goddd.VoyageNumber) (*goddd.Voyage, error) {
	r.FindInvoked = true
	return r.FindFn(number)
}

// HandlingEventRepository is a mock handling events repository.
type HandlingEventRepository struct {
	StoreFn      func(goddd.HandlingEvent)
	StoreInvoked bool

	QueryHandlingHistoryFn      func(goddd.TrackingID) goddd.HandlingHistory
	QueryHandlingHistoryInvoked bool
}

// Store calls the StoreFn.
func (r *HandlingEventRepository) Store(e goddd.HandlingEvent) {
	r.StoreInvoked = true
	r.StoreFn(e)
}

// QueryHandlingHistory calls the QueryHandlingHistoryFn.
func (r *HandlingEventRepository) QueryHandlingHistory(id goddd.TrackingID) goddd.HandlingHistory {
	r.QueryHandlingHistoryInvoked = true
	return r.QueryHandlingHistoryFn(id)
}

// RoutingService provides a mock routing service.
type RoutingService struct {
	FetchRoutesFn      func(goddd.RouteSpecification) []goddd.Itinerary
	FetchRoutesInvoked bool
}

// FetchRoutesForSpecification calls the FetchRoutesFn.
func (s *RoutingService) FetchRoutesForSpecification(rs goddd.RouteSpecification) []goddd.Itinerary {
	s.FetchRoutesInvoked = true
	return s.FetchRoutesFn(rs)
}
