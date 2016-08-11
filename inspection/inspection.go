// Package inspection provides means to inspect cargos.
package inspection

import "github.com/marcusolsson/goddd"

// EventHandler provides means of subscribing to inspection events.
type EventHandler interface {
	CargoWasMisdirected(*goddd.Cargo)
	CargoHasArrived(*goddd.Cargo)
}

// Service provides cargo inspection operations.
type Service interface {
	// InspectCargo inspects cargo and send relevant notifications to
	// interested parties, for example if a cargo has been misdirected, or
	// unloaded at the final destination.
	InspectCargo(trackingID goddd.TrackingID)
}

type service struct {
	cargoRepository         goddd.CargoRepository
	handlingEventRepository goddd.HandlingEventRepository
	cargoEventHandler       EventHandler
}

// TODO: Should be transactional
func (s *service) InspectCargo(trackingID goddd.TrackingID) {
	c, err := s.cargoRepository.Find(trackingID)
	if err != nil {
		return
	}

	h := s.handlingEventRepository.QueryHandlingHistory(trackingID)

	c.DeriveDeliveryProgress(h)

	if c.Delivery.IsMisdirected {
		s.cargoEventHandler.CargoWasMisdirected(c)
	}

	if c.Delivery.IsUnloadedAtDestination {
		s.cargoEventHandler.CargoHasArrived(c)
	}

	s.cargoRepository.Store(c)
}

// NewService creates a inspection service with necessary dependencies.
func NewService(cargoRepository goddd.CargoRepository, handlingEventRepository goddd.HandlingEventRepository, eventHandler EventHandler) Service {
	return &service{cargoRepository, handlingEventRepository, eventHandler}
}
