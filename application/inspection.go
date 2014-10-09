package application

import "github.com/marcusolsson/goddd/domain/cargo"

type CargoInspectionService interface {
	InspectCargo(trackingId cargo.TrackingId)
}

type cargoInspectionService struct {
	cargoRepository         cargo.CargoRepository
	handlingEventRepository cargo.HandlingEventRepository
	eventHandler            EventHandler
}

// TODO: Should be transactional
func (s *cargoInspectionService) InspectCargo(trackingId cargo.TrackingId) {
	c, err := s.cargoRepository.Find(trackingId)

	if err != nil {
		return
	}

	h := s.handlingEventRepository.QueryHandlingHistory(trackingId)

	c.DeriveDeliveryProgress(h)

	if c.Delivery.IsMisdirected {
		s.eventHandler.CargoWasMisdirected(c)
	}

	if c.Delivery.IsUnloadedAtDestination {
		s.eventHandler.CargoHasArrived(c)
	}

	s.cargoRepository.Store(c)
}