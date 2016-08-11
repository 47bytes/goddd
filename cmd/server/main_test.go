package main

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/marcusolsson/goddd"
	"github.com/marcusolsson/goddd/booking"
	"github.com/marcusolsson/goddd/handling"
	"github.com/marcusolsson/goddd/inmem"
	"github.com/marcusolsson/goddd/inspection"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestCargoFromHongkongToStockholm(chk *C) {
	var err error

	var (
		cargoRepository         = inmem.NewCargoRepository()
		locationRepository      = inmem.NewLocationRepository()
		voyageRepository        = inmem.NewVoyageRepository()
		handlingEventRepository = inmem.NewHandlingEventRepository()
	)

	handlingEventFactory := goddd.HandlingEventFactory{
		CargoRepository:    cargoRepository,
		VoyageRepository:   voyageRepository,
		LocationRepository: locationRepository,
	}

	routingService := &stubRoutingService{}

	cargoEventHandler := &stubCargoEventHandler{}
	cargoInspectionService := inspection.NewService(cargoRepository, handlingEventRepository, cargoEventHandler)
	handlingEventHandler := &stubHandlingEventHandler{cargoInspectionService}

	var (
		bookingService       = booking.NewService(cargoRepository, locationRepository, handlingEventRepository, routingService)
		handlingEventService = handling.NewService(handlingEventRepository, handlingEventFactory, handlingEventHandler)
	)

	var (
		origin          = goddd.CNHKG // Hongkong
		destination     = goddd.SESTO // Stockholm
		arrivalDeadline = time.Date(2009, time.March, 18, 12, 00, 00, 00, time.UTC)
	)

	//
	// Use case 1: booking
	//

	trackingID, err := bookingService.BookNewCargo(origin, destination, arrivalDeadline)

	chk.Assert(err, IsNil)

	c, err := cargoRepository.Find(trackingID)

	chk.Assert(err, IsNil)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.NotReceived)
	chk.Check(c.Delivery.RoutingStatus, Equals, goddd.NotRouted)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.Itinerary.IsEmpty(), Equals, true)
	chk.Check(c.Delivery.ETA, Equals, time.Time{})
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{})

	//
	// Use case 2: routing
	//

	itineraries := bookingService.RequestPossibleRoutesForCargo(trackingID)
	itinerary := selectPreferredItinerary(itineraries)

	c.AssignToRoute(itinerary)

	cargoRepository.Store(c)

	chk.Check(c.Delivery.TransportStatus, Equals, goddd.NotReceived)
	chk.Check(c.Delivery.RoutingStatus, Equals, goddd.Routed)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.Itinerary.IsEmpty(), Equals, false)
	chk.Check(c.Delivery.ETA, Not(Equals), time.Time{})
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{Type: goddd.Receive, Location: goddd.CNHKG})

	//
	// Use case 3: handling
	//

	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 1), trackingID, "", goddd.CNHKG, goddd.Receive)
	chk.Check(err, IsNil)

	// Ensure we're not working with stale cargo.
	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.TransportStatus, Equals, goddd.InPort)
	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.CNHKG)
	chk.Check(c.Delivery.Itinerary.IsEmpty(), Equals, false)

	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 3), trackingID, goddd.V100.Number, goddd.CNHKG, goddd.Load)
	chk.Check(err, IsNil)

	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.TransportStatus, Equals, goddd.OnboardCarrier)
	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.CNHKG)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.Itinerary.IsEmpty(), Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.V100.Number)
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{Type: goddd.Unload, Location: goddd.USNYC, VoyageNumber: goddd.V100.Number})

	//
	// Here's an attempt to register a handling event that's not valid
	// because there is no voyage with the specified voyage number,
	// and there's no location with the specified UN Locode either.
	//

	noSuchVoyageNumber := goddd.VoyageNumber("XX000")
	noSuchUNLocode := goddd.UNLocode("ZZZZZ")
	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 5), trackingID, noSuchVoyageNumber, noSuchUNLocode, goddd.Load)
	chk.Check(err, NotNil)

	//
	// Cargo is incorrectly unloaded in Tokyo
	//

	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 5), trackingID, goddd.V100.Number, goddd.JNTKO, goddd.Unload)
	chk.Check(err, IsNil)

	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.JNTKO)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.InPort)
	chk.Check(c.Delivery.Itinerary.IsEmpty(), Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.VoyageNumber(""))
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{})

	// Cargo is now misdirected
	chk.Check(c.Delivery.IsMisdirected, Equals, true)

	//
	// Cargo needs to be rerouted
	//

	rs := goddd.RouteSpecification{
		Origin:          goddd.JNTKO,
		Destination:     goddd.SESTO,
		ArrivalDeadline: arrivalDeadline,
	}

	// Specify a new route, this time from Tokyo (where it was incorrectly unloaded) to Stockholm
	c.SpecifyNewRoute(rs)

	cargoRepository.Store(c)

	chk.Check(c.Delivery.RoutingStatus, Equals, goddd.Misrouted)
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{})

	// Repeat procedure of selecting one out of a number of possible routes satisfying the route spec
	newItineraries := bookingService.RequestPossibleRoutesForCargo(trackingID)
	newItinerary := selectPreferredItinerary(newItineraries)

	c.AssignToRoute(newItinerary)

	cargoRepository.Store(c)

	chk.Check(c.Delivery.RoutingStatus, Equals, goddd.Routed)

	//
	// Cargo has been rerouted, shipping continues
	//

	// Load in Tokyo
	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 8), trackingID, goddd.V300.Number, goddd.JNTKO, goddd.Load)
	chk.Check(err, IsNil)

	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.JNTKO)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.OnboardCarrier)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.V300.Number)
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{Type: goddd.Unload, Location: goddd.DEHAM, VoyageNumber: goddd.V300.Number})

	// Unload in Hamburg
	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 12), trackingID, goddd.V300.Number, goddd.DEHAM, goddd.Unload)
	chk.Check(err, IsNil)

	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.DEHAM)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.InPort)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.VoyageNumber(""))
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{Type: goddd.Load, Location: goddd.DEHAM, VoyageNumber: goddd.V400.Number})

	// Load in Hamburg
	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 14), trackingID, goddd.V400.Number, goddd.DEHAM, goddd.Load)
	chk.Check(err, IsNil)

	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.DEHAM)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.OnboardCarrier)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.V400.Number)
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{Type: goddd.Unload, Location: goddd.SESTO, VoyageNumber: goddd.V400.Number})

	// Unload in Stockholm
	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 15), trackingID, goddd.V400.Number, goddd.SESTO, goddd.Unload)
	chk.Check(err, IsNil)

	c, err = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.SESTO)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.InPort)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.VoyageNumber(""))
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{Type: goddd.Claim, Location: goddd.SESTO})

	// Finally, cargo is claimed in Stockholm. This ends the cargo lifecycle from our perspective.
	err = handlingEventService.RegisterHandlingEvent(toDate(2009, time.March, 16), trackingID, goddd.V400.Number, goddd.SESTO, goddd.Claim)
	chk.Check(err, IsNil)

	c, _ = cargoRepository.Find(trackingID)

	chk.Check(c.Delivery.LastKnownLocation, Equals, goddd.SESTO)
	chk.Check(c.Delivery.TransportStatus, Equals, goddd.Claimed)
	chk.Check(c.Delivery.IsMisdirected, Equals, false)
	chk.Check(c.Delivery.CurrentVoyage, Equals, goddd.VoyageNumber(""))
	chk.Check(c.Delivery.NextExpectedActivity, Equals, goddd.HandlingActivity{})
}

func selectPreferredItinerary(itineraries []goddd.Itinerary) goddd.Itinerary {
	return itineraries[0]
}

func toDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 12, 00, 00, 00, time.UTC)
}

// Stub RoutingService
type stubRoutingService struct{}

func (s *stubRoutingService) FetchRoutesForSpecification(rs goddd.RouteSpecification) []goddd.Itinerary {
	if rs.Origin == goddd.CNHKG {
		return []goddd.Itinerary{
			{Legs: []goddd.Leg{
				goddd.NewLeg("V100", goddd.CNHKG, goddd.USNYC, toDate(2009, time.March, 3), toDate(2009, time.March, 9)),
				goddd.NewLeg("V200", goddd.USNYC, goddd.USCHI, toDate(2009, time.March, 10), toDate(2009, time.March, 14)),
				goddd.NewLeg("V300", goddd.USCHI, goddd.SESTO, toDate(2009, time.March, 7), toDate(2009, time.March, 11)),
			}},
		}
	}

	return []goddd.Itinerary{
		{Legs: []goddd.Leg{
			goddd.NewLeg("V300", goddd.JNTKO, goddd.DEHAM, toDate(2009, time.March, 8), toDate(2009, time.March, 12)),
			goddd.NewLeg("V400", goddd.DEHAM, goddd.SESTO, toDate(2009, time.March, 14), toDate(2009, time.March, 15)),
		}},
	}
}

// Stub HandlingEventHandler
type stubHandlingEventHandler struct {
	InspectionService inspection.Service
}

func (h *stubHandlingEventHandler) CargoWasHandled(event goddd.HandlingEvent) {
	h.InspectionService.InspectCargo(event.TrackingID)
}

// Stub CargoEventHandler
type stubCargoEventHandler struct {
}

func (h *stubCargoEventHandler) CargoWasMisdirected(c *goddd.Cargo) {
}

func (h *stubCargoEventHandler) CargoHasArrived(c *goddd.Cargo) {
}
