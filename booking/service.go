// Package booking provides the use-case of booking a cargo. Used by views
// facing an administrator.
package booking

import (
	"errors"
	"time"

	"github.com/marcusolsson/goddd"
	"github.com/marcusolsson/goddd/routing"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// Service is the interface that provides booking methods.
type Service interface {
	// BookNewCargo registers a new cargo in the tracking system, not yet
	// routed.
	BookNewCargo(origin goddd.UNLocode, destination goddd.UNLocode, arrivalDeadline time.Time) (goddd.TrackingID, error)

	// LoadCargo returns a read model of a cargo.
	LoadCargo(trackingID goddd.TrackingID) (Cargo, error)

	// RequestPossibleRoutesForCargo requests a list of itineraries describing
	// possible routes for this cargo.
	RequestPossibleRoutesForCargo(trackingID goddd.TrackingID) []goddd.Itinerary

	// AssignCargoToRoute assigns a cargo to the route specified by the
	// itinerary.
	AssignCargoToRoute(trackingID goddd.TrackingID, itinerary goddd.Itinerary) error

	// ChangeDestination changes the destination of a cargo.
	ChangeDestination(trackingID goddd.TrackingID, unLocode goddd.UNLocode) error

	// Cargos returns a list of all cargos that have been booked.
	Cargos() []Cargo

	// Locations returns a list of registered locations.
	Locations() []Location
}

type service struct {
	cargos         goddd.CargoRepository
	locations      goddd.LocationRepository
	handlingEvents goddd.HandlingEventRepository
	routingService routing.Service
}

func (s *service) AssignCargoToRoute(id goddd.TrackingID, itinerary goddd.Itinerary) error {
	if id == "" || len(itinerary.Legs) == 0 {
		return ErrInvalidArgument
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return err
	}

	c.AssignToRoute(itinerary)

	return s.cargos.Store(c)
}

func (s *service) BookNewCargo(origin, destination goddd.UNLocode, arrivalDeadline time.Time) (goddd.TrackingID, error) {
	if origin == "" || destination == "" || arrivalDeadline.IsZero() {
		return "", ErrInvalidArgument
	}

	id := goddd.NextTrackingID()
	rs := goddd.RouteSpecification{
		Origin:          origin,
		Destination:     destination,
		ArrivalDeadline: arrivalDeadline,
	}

	c := goddd.NewCargo(id, rs)

	if err := s.cargos.Store(c); err != nil {
		return "", err
	}

	return c.TrackingID, nil
}

func (s *service) LoadCargo(trackingID goddd.TrackingID) (Cargo, error) {
	if trackingID == "" {
		return Cargo{}, ErrInvalidArgument
	}

	c, err := s.cargos.Find(trackingID)
	if err != nil {
		return Cargo{}, err
	}

	return assemble(c, s.handlingEvents), nil
}

func (s *service) ChangeDestination(id goddd.TrackingID, destination goddd.UNLocode) error {
	if id == "" || destination == "" {
		return ErrInvalidArgument
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return err
	}

	l, err := s.locations.Find(destination)
	if err != nil {
		return err
	}

	c.SpecifyNewRoute(goddd.RouteSpecification{
		Origin:          c.Origin,
		Destination:     l.UNLocode,
		ArrivalDeadline: c.RouteSpecification.ArrivalDeadline,
	})

	if err := s.cargos.Store(c); err != nil {
		return err
	}

	return nil
}

func (s *service) RequestPossibleRoutesForCargo(id goddd.TrackingID) []goddd.Itinerary {
	if id == "" {
		return nil
	}

	c, err := s.cargos.Find(id)
	if err != nil {
		return []goddd.Itinerary{}
	}

	return s.routingService.FetchRoutesForSpecification(c.RouteSpecification)
}

func (s *service) Cargos() []Cargo {
	var result []Cargo
	for _, c := range s.cargos.FindAll() {
		result = append(result, assemble(c, s.handlingEvents))
	}
	return result
}

func (s *service) Locations() []Location {
	var result []Location
	for _, v := range s.locations.FindAll() {
		result = append(result, Location{
			UNLocode: string(v.UNLocode),
			Name:     v.Name,
		})
	}
	return result
}

// NewService creates a booking service with necessary dependencies.
func NewService(cr goddd.CargoRepository, lr goddd.LocationRepository, her goddd.HandlingEventRepository, rs routing.Service) Service {
	return &service{
		cargos:         cr,
		locations:      lr,
		handlingEvents: her,
		routingService: rs,
	}
}

// Location is a read model for booking views.
type Location struct {
	UNLocode string `json:"locode"`
	Name     string `json:"name"`
}

// Cargo is a read model for booking views.
type Cargo struct {
	ArrivalDeadline time.Time   `json:"arrival_deadline"`
	Destination     string      `json:"destination"`
	Legs            []goddd.Leg `json:"legs,omitempty"`
	Misrouted       bool        `json:"misrouted"`
	Origin          string      `json:"origin"`
	Routed          bool        `json:"routed"`
	TrackingID      string      `json:"tracking_id"`
}

func assemble(c *goddd.Cargo, her goddd.HandlingEventRepository) Cargo {
	return Cargo{
		TrackingID:      string(c.TrackingID),
		Origin:          string(c.Origin),
		Destination:     string(c.RouteSpecification.Destination),
		Misrouted:       c.Delivery.RoutingStatus == goddd.Misrouted,
		Routed:          !c.Itinerary.IsEmpty(),
		ArrivalDeadline: c.RouteSpecification.ArrivalDeadline,
		Legs:            c.Itinerary.Legs,
	}
}
