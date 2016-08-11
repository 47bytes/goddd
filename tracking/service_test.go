package tracking

import (
	"testing"

	"github.com/marcusolsson/goddd"
	"github.com/marcusolsson/goddd/mock"
)

func TestTrack(t *testing.T) {
	var cargos mock.CargoRepository
	cargos.FindFn = func(id goddd.TrackingID) (*goddd.Cargo, error) {
		return goddd.NewCargo("FTL456", goddd.RouteSpecification{
			Origin:      goddd.AUMEL,
			Destination: goddd.SESTO,
		}), nil
	}

	var events mock.HandlingEventRepository
	events.QueryHandlingHistoryFn = func(id goddd.TrackingID) goddd.HandlingHistory {
		return goddd.HandlingHistory{}
	}

	s := NewService(&cargos, &events)

	c, err := s.Track("FTL456")
	if err != nil {
		t.Fatal(err)
	}

	if c.TrackingID != "FTL456" {
		t.Errorf("c.TrackingID = %v; want = %v", c.TrackingID, "FTL456")
	}
	if c.Origin != "AUMEL" {
		t.Errorf("c.Origin = %v; want = %v", c.Destination, "AUMEL")
	}
	if c.Destination != "SESTO" {
		t.Errorf("c.Destination = %v; want = %v", c.Destination, "SESTO")
	}
	if c.StatusText != goddd.NotReceived.String() {
		t.Errorf("c.StatusText = %v; want = %v", c.StatusText, goddd.NotReceived.String())
	}
}
