package booking

import (
	"testing"
	"time"

	"github.com/marcusolsson/goddd"
	"github.com/marcusolsson/goddd/mock"
)

func TestBookNewCargo(t *testing.T) {
	var (
		origin      = goddd.SESTO
		destination = goddd.AUMEL
		deadline    = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
	)

	var cargos mockCargoRepository

	s := NewService(&cargos, nil, nil, nil)

	id, err := s.BookNewCargo(origin, destination, deadline)
	if err != nil {
		t.Fatal(err)
	}

	c, err := cargos.Find(id)
	if err != nil {
		t.Fatal(err)
	}

	if c.TrackingID != id {
		t.Errorf("c.TrackingID = %s; want = %s", c.TrackingID, id)
	}
	if c.Origin != origin {
		t.Errorf("c.Origin = %s; want = %s", c.Origin, origin)
	}
	if c.RouteSpecification.Destination != destination {
		t.Errorf("c.RouteSpecification.Destination = %s; want = %s",
			c.RouteSpecification.Destination, destination)
	}
	if c.RouteSpecification.ArrivalDeadline != deadline {
		t.Errorf("c.RouteSpecification.ArrivalDeadline = %s; want = %s",
			c.RouteSpecification.ArrivalDeadline, deadline)
	}
}

type stubRoutingService struct{}

func (s *stubRoutingService) FetchRoutesForSpecification(rs goddd.RouteSpecification) []goddd.Itinerary {
	legs := []goddd.Leg{
		{LoadLocation: rs.Origin, UnloadLocation: rs.Destination},
	}

	return []goddd.Itinerary{
		{Legs: legs},
	}
}

func TestRequestPossibleRoutesForCargo(t *testing.T) {
	var (
		origin      = goddd.SESTO
		destination = goddd.AUMEL
		deadline    = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
	)

	var cargos mockCargoRepository

	var rs stubRoutingService

	s := NewService(&cargos, nil, nil, &rs)

	r := s.RequestPossibleRoutesForCargo("no_such_id")

	if len(r) != 0 {
		t.Errorf("len(r) = %d; want = %d", len(r), 0)
	}

	id, err := s.BookNewCargo(origin, destination, deadline)
	if err != nil {
		t.Fatal(err)
	}

	i := s.RequestPossibleRoutesForCargo(id)

	if len(i) != 1 {
		t.Errorf("len(i) = %d; want = %d", len(i), 1)
	}
}

func TestAssignCargoToRoute(t *testing.T) {
	var cargos mockCargoRepository

	var rs stubRoutingService

	s := NewService(&cargos, nil, nil, &rs)

	var (
		origin      = goddd.SESTO
		destination = goddd.AUMEL
		deadline    = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
	)

	id, err := s.BookNewCargo(origin, destination, deadline)
	if err != nil {
		t.Fatal(err)
	}

	i := s.RequestPossibleRoutesForCargo(id)

	if len(i) != 1 {
		t.Errorf("len(i) = %d; want = %d", len(i), 1)
	}

	if err := s.AssignCargoToRoute(id, i[0]); err != nil {
		t.Fatal(err)
	}

	if err := s.AssignCargoToRoute("no_such_id", goddd.Itinerary{}); err != ErrInvalidArgument {
		t.Errorf("err = %s; want = %s", err, ErrInvalidArgument)
	}
}

func TestChangeCargoDestination(t *testing.T) {
	var cargos mockCargoRepository
	var locations mock.LocationRepository

	locations.FindFn = func(loc goddd.UNLocode) (goddd.Location, error) {
		if loc != goddd.AUMEL {
			return goddd.Location{}, goddd.ErrUnknownLocation
		}
		return goddd.Melbourne, nil
	}

	var rs stubRoutingService

	s := NewService(&cargos, &locations, nil, &rs)

	c := goddd.NewCargo("ABC", goddd.RouteSpecification{
		Origin:          goddd.SESTO,
		Destination:     goddd.CNHKG,
		ArrivalDeadline: time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC),
	})

	if err := s.ChangeDestination("no_such_id", goddd.SESTO); err != goddd.ErrUnknownCargo {
		t.Errorf("err = %s; want = %s", err, goddd.ErrUnknownCargo)
	}

	if err := cargos.Store(c); err != nil {
		t.Fatal(err)
	}

	if err := s.ChangeDestination(c.TrackingID, "no_such_unlocode"); err != goddd.ErrUnknownLocation {
		t.Errorf("err = %s; want = %s", err, goddd.ErrUnknownLocation)
	}

	if c.RouteSpecification.Destination != goddd.CNHKG {
		t.Errorf("c.RouteSpecification.Destination = %s; want = %s",
			c.RouteSpecification.Destination, goddd.CNHKG)
	}

	if err := s.ChangeDestination(c.TrackingID, goddd.AUMEL); err != nil {
		t.Fatal(err)
	}

	uc, err := cargos.Find(c.TrackingID)
	if err != nil {
		t.Fatal(err)
	}

	if uc.RouteSpecification.Destination != goddd.AUMEL {
		t.Errorf("uc.RouteSpecification.Destination = %s; want = %s",
			uc.RouteSpecification.Destination, goddd.AUMEL)
	}
}

type mockCargoRepository struct {
	cargo *goddd.Cargo
}

func (r *mockCargoRepository) Store(c *goddd.Cargo) error {
	r.cargo = c
	return nil
}

func (r *mockCargoRepository) Find(id goddd.TrackingID) (*goddd.Cargo, error) {
	if r.cargo != nil {
		return r.cargo, nil
	}
	return nil, goddd.ErrUnknownCargo
}

func (r *mockCargoRepository) FindAll() []*goddd.Cargo {
	return []*goddd.Cargo{r.cargo}
}
