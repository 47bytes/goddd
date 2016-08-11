package handling

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

func (s *instrumentingService) RegisterHandlingEvent(completionTime time.Time, trackingID goddd.TrackingID, voyage goddd.VoyageNumber,
	loc goddd.UNLocode, eventType goddd.HandlingEventType) error {

	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "register_incident"}
		s.requestCount.With(methodField).Add(1)
		s.requestLatency.With(methodField).Observe(time.Since(begin))
	}(time.Now())

	return s.Service.RegisterHandlingEvent(completionTime, trackingID, voyage, loc, eventType)
}
