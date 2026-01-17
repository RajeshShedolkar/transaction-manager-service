package imps

import (
	"errors"
	//"transaction-manager/internal/api"

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

// ---------- DEBIT CONFIRMED ----------
func (s *EventTransactionServiceImpl) HandleDebitConfirmed(event domain.TxEvent) error {

	if event.TxId == "" {
		// TODO: Handle DLQ publish
		//api.PublishToDlqAccountService(event, config.KAFKA_ACCOUNT_TOPIC, "missing TxId in event", context.Background())
		return errors.New("missing TxId in event")
	}

	tx, err := s.txRepo.FindByID(event.TxId)
	if err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: tx.ID,
		AccountRefId:  tx.SourceRefId,
		DcFlag:        "D",
		EntryType:     domain.LedgerDebit,
		Amount:        tx.Amount,
		Source:        "EVENT",
		Msg:           "Debit confirmed from source account",
	}

	return s.ledgerRepo.Append(ledger)
}
