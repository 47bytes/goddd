package inspection

import (
	"testing"

	"github.com/marcusolsson/goddd"
)

type stubEventHandler struct {
	events []interface{}
}

func (h *stubEventHandler) CargoWasMisdirected(c *goddd.Cargo) {
	h.events = append(h.events, c)
}

func (h *stubEventHandler) CargoHasArrived(c *goddd.Cargo) {
	h.events = append(h.events, c)
}

func TestInspectMisdirectedCargo(t *testing.T) {
	var cargos mockCargoRepository

	events := mockHandlingEventRepository{
		events: make(map[goddd.TrackingID][]goddd.HandlingEvent),
	}

	handler := stubEventHandler{make([]interface{}, 0)}

	s := NewService(&cargos, &events, &handler)

	id := goddd.TrackingID("ABC123")
	c := goddd.NewCargo(id, goddd.RouteSpecification{
		Origin:      goddd.SESTO,
		Destination: goddd.CNHKG,
	})

	voyage := goddd.VoyageNumber("001A")

	c.AssignToRoute(goddd.Itinerary{Legs: []goddd.Leg{
		{VoyageNumber: voyage, LoadLocation: goddd.SESTO, UnloadLocation: goddd.AUMEL},
		{VoyageNumber: voyage, LoadLocation: goddd.AUMEL, UnloadLocation: goddd.CNHKG},
	}})

	if err := cargos.Store(c); err != nil {
		t.Fatal(err)
	}

	storeEvent(&events, id, voyage, goddd.Receive, goddd.SESTO)
	storeEvent(&events, id, voyage, goddd.Load, goddd.SESTO)
	storeEvent(&events, id, voyage, goddd.Unload, goddd.USNYC)

	if len(handler.events) != 0 {
		t.Errorf("no events should be handled")
	}

	s.InspectCargo(id)

	if len(handler.events) != 1 {
		t.Errorf("1 event should be handled")
	}

	s.InspectCargo("no_such_id")

	// no events was published
	if len(handler.events) != 1 {
		t.Errorf("len(handler.events) = %d; want = %d", len(handler.events), 1)
	}
}

func TestInspectUnloadedCargo(t *testing.T) {
	var cargos mockCargoRepository

	events := mockHandlingEventRepository{
		events: make(map[goddd.TrackingID][]goddd.HandlingEvent),
	}

	handler := stubEventHandler{make([]interface{}, 0)}

	s := &service{
		cargoRepository:         &cargos,
		handlingEventRepository: &events,
		cargoEventHandler:       &handler,
	}

	id := goddd.TrackingID("ABC123")
	unloadedCargo := goddd.NewCargo(id, goddd.RouteSpecification{
		Origin:      goddd.SESTO,
		Destination: goddd.CNHKG,
	})

	var voyage goddd.VoyageNumber = "001A"

	unloadedCargo.AssignToRoute(goddd.Itinerary{Legs: []goddd.Leg{
		{VoyageNumber: voyage, LoadLocation: goddd.SESTO, UnloadLocation: goddd.AUMEL},
		{VoyageNumber: voyage, LoadLocation: goddd.AUMEL, UnloadLocation: goddd.CNHKG},
	}})

	cargos.Store(unloadedCargo)

	storeEvent(&events, id, voyage, goddd.Receive, goddd.SESTO)
	storeEvent(&events, id, voyage, goddd.Load, goddd.SESTO)
	storeEvent(&events, id, voyage, goddd.Unload, goddd.AUMEL)
	storeEvent(&events, id, voyage, goddd.Load, goddd.AUMEL)
	storeEvent(&events, id, voyage, goddd.Unload, goddd.CNHKG)

	if len(handler.events) != 0 {
		t.Errorf("len(handler.events) = %d; want = %d", len(handler.events), 0)
	}

	s.InspectCargo(id)

	if len(handler.events) != 1 {
		t.Errorf("len(handler.events) = %d; want = %d", len(handler.events), 1)
	}
}

func storeEvent(r goddd.HandlingEventRepository, id goddd.TrackingID, voyage goddd.VoyageNumber, typ goddd.HandlingEventType, loc goddd.UNLocode) {
	e := goddd.HandlingEvent{
		TrackingID: id,
		Activity: goddd.HandlingActivity{
			VoyageNumber: voyage,
			Type:         typ,
			Location:     loc,
		},
	}

	r.Store(e)
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

type mockHandlingEventRepository struct {
	events map[goddd.TrackingID][]goddd.HandlingEvent
}

func (r *mockHandlingEventRepository) Store(e goddd.HandlingEvent) {
	if _, ok := r.events[e.TrackingID]; !ok {
		r.events[e.TrackingID] = make([]goddd.HandlingEvent, 0)
	}
	r.events[e.TrackingID] = append(r.events[e.TrackingID], e)
}

func (r *mockHandlingEventRepository) QueryHandlingHistory(id goddd.TrackingID) goddd.HandlingHistory {
	return goddd.HandlingHistory{HandlingEvents: r.events[id]}
}
