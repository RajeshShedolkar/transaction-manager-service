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

	processed, err := eventRepo.IsProcessed(event.EventID)
	if err != nil {
		log.Println("Idempotency check failed:", err)
		return
	}
	if processed {
		log.Println("Duplicate event ignored:", event.EventID)
		return
	}

	log.Println("Processing card event:", event.EventID, event.EventType)
	// et = AUTH_SUCCESS/AUTH_FAILED/SETTLEMENT_STARTED/DEBIT_CONFIRMED/CREDIT_CONFIRMED/
	// et = CANCEL_REQUESTED/SETTLEMENT_FAILED/REFUND_PROCESSED
	et := domain.EventType(event.EventType)
	switch et {

	case domain.EventAuthSuccess:
		// AUTH → ledger AUTH → txn AUTHORIZED
		h.service.HandleAuthSuccess(event)

	case domain.EventSettlementStarted:
		// SETTLEMENT marker → txn PROCESSING
		h.service.HandleSettlementStarted(event)

	case domain.EventDebitConfirmed:
		// DEBIT ledger entry
		h.service.HandleDebitConfirmed(event)

	case domain.EventCreditConfirmed:
		// CREDIT ledger entry → txn COMPLETED
		h.service.HandleCreditConfirmed(event)

	case domain.EventCancelRequested:
		// RELEASE ledger entry → txn RELEASED
		h.service.HandleCancel(event)

	case domain.EventSettlementFailed:
		// REVERSAL ledger entry → txn FAILED
		h.service.HandleSettlementFailed(event)

	case domain.EventRefundProcessed:
		// REFUND ledger entry → txn REFUNDED
		h.service.HandleRefund(event)

	default:
		log.Println("Unhandled card event type:", event.EventType)
		return
	}

	if err := eventRepo.MarkProcessed(event.EventID, event.EventType); err != nil {
		log.Println("Failed to mark event processed:", err)
	}
}
