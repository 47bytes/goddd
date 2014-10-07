package handling

import (
	"testing"
	"time"

	"github.com/marcusolsson/goddd/domain/cargo"
	"github.com/marcusolsson/goddd/domain/location"
	. "gopkg.in/check.v1"
)

// Hook gocheck up to the "go test" runner
func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestRegisterHandlingEvent(c *C) {

	var (
		cargoRepository         = cargo.NewCargoRepository()
		handlingEventRepository = &cargo.HandlingEventRepositoryInMem{}
		handlingEventFactory    = cargo.HandlingEventFactory{
			CargoRepository: cargoRepository,
		}
	)

	var service HandlingEventService = &handlingEventService{
		handlingEventRepository: handlingEventRepository,
		handlingEventFactory:    handlingEventFactory,
	}

	var (
		completionTime = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
		trackingId     = cargo.TrackingId("ABC123")
		voyageNumber   = "CM001"
		unLocode       = location.Stockholm.UNLocode
		eventType      = cargo.Load
	)

	service.RegisterHandlingEvent(completionTime, trackingId, voyageNumber, unLocode, eventType)
}