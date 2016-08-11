package booking

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/marcusolsson/goddd"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.TimeHistogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) BookNewCargo(origin, destination goddd.UNLocode, arrivalDeadline time.Time) (goddd.TrackingID, error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "book"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.BookNewCargo(origin, destination, arrivalDeadline)
}

func (s *instrumentingService) LoadCargo(id goddd.TrackingID) (c Cargo, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "load"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.LoadCargo(id)
}

func (s *instrumentingService) RequestPossibleRoutesForCargo(id goddd.TrackingID) []goddd.Itinerary {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "request_routes"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.RequestPossibleRoutesForCargo(id)
}

func (s *instrumentingService) AssignCargoToRoute(id goddd.TrackingID, itinerary goddd.Itinerary) (err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "assign_to_route"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.AssignCargoToRoute(id, itinerary)
}

func (s *instrumentingService) ChangeDestination(id goddd.TrackingID, l goddd.UNLocode) (err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "change_destination"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.ChangeDestination(id, l)
}

func (s *instrumentingService) Cargos() []Cargo {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "list_cargos"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.Cargos()
}

func (s *instrumentingService) Locations() []Location {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "list_locations"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.Locations()
}
