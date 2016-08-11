package handling

import (
	"testing"
	"time"

	"github.com/marcusolsson/goddd"
	"github.com/marcusolsson/goddd/mock"
)

type stubEventHandler struct {
	events []interface{}
}

func (h *stubEventHandler) CargoWasHandled(e goddd.HandlingEvent) {
	h.events = append(h.events, e)
}

func TestRegisterHandlingEvent(t *testing.T) {
	var cargos mock.CargoRepository
	cargos.StoreFn = func(c *goddd.Cargo) error {
		return nil
	}
	cargos.FindFn = func(id goddd.TrackingID) (*goddd.Cargo, error) {
		if id == "no_such_id" {
			return nil, goddd.ErrUnknownCargo
		}
		return new(goddd.Cargo), nil
	}

	var voyages mock.VoyageRepository
	voyages.FindFn = func(n goddd.VoyageNumber) (*goddd.Voyage, error) {
		return new(goddd.Voyage), nil
	}

	var locations mock.LocationRepository
	locations.FindFn = func(l goddd.UNLocode) (goddd.Location, error) {
		return goddd.Location{}, nil
	}

	var events mock.HandlingEventRepository
	events.StoreFn = func(e goddd.HandlingEvent) {}

	eh := &stubEventHandler{events: make([]interface{}, 0)}
	ef := goddd.HandlingEventFactory{
		CargoRepository:    &cargos,
		VoyageRepository:   &voyages,
		LocationRepository: &locations,
	}

	s := NewService(&events, ef, eh)

	var (
		completed = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
		id        = goddd.TrackingID("ABC123")
		voyage    = goddd.VoyageNumber("V100")
	)

	var err error

	err = cargos.Store(goddd.NewCargo(id, goddd.RouteSpecification{}))
	if err != nil {
		t.Fatal(err)
	}

	err = s.RegisterHandlingEvent(completed, id, voyage, goddd.SESTO, goddd.Load)
	if err != nil {
		t.Fatal(err)
	}

	err = s.RegisterHandlingEvent(completed, "no_such_id", voyage, goddd.SESTO, goddd.Load)
	if err != goddd.ErrUnknownCargo {
		t.Errorf("err = %s; want = %s", err, goddd.ErrUnknownCargo)
	}

	if len(eh.events) != 1 {
		t.Errorf("len(eh.events) = %d; want = %d", len(eh.events), 1)
	}
}
