package booking

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/marcusolsson/goddd"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) BookNewCargo(origin goddd.UNLocode, destination goddd.UNLocode, arrivalDeadline time.Time) (id goddd.TrackingID, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "book",
			"origin", origin,
			"destination", destination,
			"arrival_deadline", arrivalDeadline,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.BookNewCargo(origin, destination, arrivalDeadline)
}

func (s *loggingService) LoadCargo(id goddd.TrackingID) (c Cargo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "load",
			"tracking_id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LoadCargo(id)
}

func (s *loggingService) RequestPossibleRoutesForCargo(id goddd.TrackingID) []goddd.Itinerary {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "request_routes",
			"tracking_id", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.RequestPossibleRoutesForCargo(id)
}

func (s *loggingService) AssignCargoToRoute(id goddd.TrackingID, itinerary goddd.Itinerary) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "assign_to_route",
			"tracking_id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AssignCargoToRoute(id, itinerary)
}

func (s *loggingService) ChangeDestination(id goddd.TrackingID, l goddd.UNLocode) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "change_destination",
			"tracking_id", id,
			"destination", l,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.ChangeDestination(id, l)
}

func (s *loggingService) Cargos() []Cargo {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "list_cargos",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Cargos()
}

func (s *loggingService) Locations() []Location {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "list_locations",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Locations()
}
