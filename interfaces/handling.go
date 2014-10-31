package interfaces

import (
	"strconv"
	"time"

	"github.com/marcusolsson/goddd/application"
	"github.com/marcusolsson/goddd/domain/cargo"
	"github.com/marcusolsson/goddd/domain/location"
	"github.com/marcusolsson/goddd/domain/voyage"
)

type HandlingEventServiceFacade interface {
	RegisterHandlingEvent(completionTime, trackingID, voyageNumber, unLocode, eventType string) error
}

type handlingEventServiceFacade struct {
	Service application.HandlingEventService
}

func (f *handlingEventServiceFacade) RegisterHandlingEvent(completionTime, trackingID, voyageNumber, unLocode, eventType string) error {
	millis, _ := strconv.ParseInt(completionTime, 10, 64)
	return f.Service.RegisterHandlingEvent(time.Unix(millis/1000, 0), cargo.TrackingID(trackingID), voyage.Number(voyageNumber), location.UNLocode(unLocode), stringToEventType(eventType))
}

func NewHandlingEventServiceFacade(cargoRepository cargo.Repository, locationRepository location.Repository, voyageRepository voyage.Repository, handlingEventRepository cargo.HandlingEventRepository) HandlingEventServiceFacade {
	cargoEventHandler := &cargoEventHandler{}
	cargoInspectionService := application.NewCargoInspectionService(cargoRepository, handlingEventRepository, cargoEventHandler)

	handlingEventFactory := cargo.HandlingEventFactory{
		CargoRepository:    cargoRepository,
		VoyageRepository:   voyageRepository,
		LocationRepository: locationRepository,
	}

	handlingEventHandler := &handlingEventHandler{
		InspectionService: cargoInspectionService,
	}

	handlingEventService := application.NewHandlingEventService(handlingEventRepository, handlingEventFactory, handlingEventHandler)

	return &handlingEventServiceFacade{Service: handlingEventService}
}

type handlingEventHandler struct {
	InspectionService application.CargoInspectionService
}

func (h *handlingEventHandler) CargoWasHandled(event cargo.HandlingEvent) {
	h.InspectionService.InspectCargo(event.TrackingID)
}

type cargoEventHandler struct {
}

func (h *cargoEventHandler) CargoWasMisdirected(c cargo.Cargo) {
}

func (h *cargoEventHandler) CargoHasArrived(c cargo.Cargo) {
}

func stringToEventType(s string) cargo.HandlingEventType {
	types := map[string]cargo.HandlingEventType{
		cargo.Receive.String(): cargo.Receive,
		cargo.Load.String():    cargo.Load,
		cargo.Unload.String():  cargo.Unload,
		cargo.Customs.String(): cargo.Customs,
		cargo.Claim.String():   cargo.Claim,
	}
	return types[s]
}