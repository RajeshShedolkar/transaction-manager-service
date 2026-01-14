package events

import (
	"encoding/json"
	"log"

	"transaction-manager/internal/cardservice"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/repository"
)

type EventTransactionHandler struct {
	service cardservice.CardEventService
}

func NewEventTransactionHandler(s cardservice.CardEventService) *EventTransactionHandler {
	return &EventTransactionHandler{service: s}
}

func (h *EventTransactionHandler) HandleCardEventIdempotent(
	message []byte,
	eventRepo repository.PgxEventRepo,
) {

	var event domain.CardEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}
	processed, _ := eventRepo.IsProcessed(event.EventID)
	if processed {
		log.Println("Duplicate event ignored:", event.EventID)
		return
	}

	log.Println("Processing event:", event.EventID, event.EventType)

	switch event.EventType {
	case "CARD_AUTH":
		h.service.HandleCardAuth(event)
		// case "CARD_SETTLEMENT":
		// 	handler.HandleCardSettlement(event.NetworkTxn)
		// case "CARD_AUTH_RELEASE":
		// 	handler.HandleCardRelease(event.NetworkTxn)
	}

	eventRepo.MarkProcessed(event.EventID, event.EventType)
}
