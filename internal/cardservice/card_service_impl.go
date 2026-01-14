package cardservice

import (
	"transaction-manager/internal/domain"
	"transaction-manager/internal/repository"

	"github.com/google/uuid"
)

type EventTransactionServiceImpl struct {
	txRepo     repository.TransactionRepository
	ledgerRepo repository.LedgerRepository
	eventRepo  repository.EventRepository
}

func NewEventTransactionService(
	txRepo repository.TransactionRepository,
	ledgerRepo repository.LedgerRepository,
	eventRepo repository.EventRepository,
) *EventTransactionServiceImpl {
	return &EventTransactionServiceImpl{
		txRepo:     txRepo,
		ledgerRepo: ledgerRepo,
		eventRepo:  eventRepo,
	}
}

func (s *EventTransactionServiceImpl) HandleCardAuth(event domain.CardEvent) error {

	tx := &domain.Transaction{
		ID:          event.NetworkTxn,
		PaymentType: "CARD",
		PaymentMode: "CARD",
		Status:      domain.StatusAuthorized,
		Amount:      event.Amount,
		Currency:    "INR",
	}

	if err := s.txRepo.Save(tx); err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		EntryType:     domain.LedgerAuth,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	return s.ledgerRepo.Append(ledger)
}
