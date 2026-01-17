package events

import (
	"encoding/json"
	"log"
	"transaction-manager/internal/imps"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/repository"
)

type ImpsEventTransactionHandler struct {
	service imps.ImpsEventService
}

func NewImpsEventTransactionHandler(s imps.ImpsEventService) *ImpsEventTransactionHandler {
	return &ImpsEventTransactionHandler{service: s}
}

func (h *ImpsEventTransactionHandler) HandleImpsEventIdempotent(
	message []byte,
	eventRepo repository.PgxEventRepo,
) {
	var event domain.TxEvent
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

	case domain.EventDebitConfirmed:
		// AUTH → ledger AUTH → txn AUTHORIZED
		h.service.HandleImpsDebitTx(event)

	case domain.EventCreditConfirmed:
		h.service.HandleImpsCreditTx(event)

	default:
		log.Println("Unhandled card event type:", event.EventType)
		return
	}

	if err := eventRepo.MarkProcessed(event.EventID, event.EventType); err != nil {
		log.Println("Failed to mark event processed:", err)
	}
}
