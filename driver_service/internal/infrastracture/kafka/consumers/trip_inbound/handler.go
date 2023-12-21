package trip_inbound

import (
	"context"
	"driver_service/internal/application/trip"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventHandler struct {
	tripService trip.Service
	logger      *zap.Logger
}

func NewEventHandler(
	tripService trip.Service,
	logger *zap.Logger,
) *EventHandler {
	return &EventHandler{
		tripService: tripService,
		logger:      logger,
	}
}

func (eh *EventHandler) Handle(ctx context.Context, message kafka.Message) error {
	var event Event
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		eh.logger.Error("Error decoding message", zap.Error(err))
		return err
	}

	parsedUUID, err := uuid.Parse(event.Data.TripID)
	if err != nil {
		eh.logger.Error("Error decoding uuid from topic message", zap.Error(err))
		return err
	}

	err = eh.tripService.CreateTrip(ctx, parsedUUID)
	if err != nil {
		eh.logger.Error("Failed to create trip", zap.Error(err))
		return err
	}

	return nil
}