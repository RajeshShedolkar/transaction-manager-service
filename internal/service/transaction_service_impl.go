package service

import (
	//"errors"
	//"time"

	"github.com/google/uuid"

	"transaction-manager/internal/domain"
	"transaction-manager/internal/repository"
	"transaction-manager/internal/state"
)

type TransactionServiceImpl struct {
	txRepo     repository.TransactionRepository
	ledgerRepo repository.LedgerRepository
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	ledgerRepo repository.LedgerRepository,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		txRepo:     txRepo,
		ledgerRepo: ledgerRepo,
	}
}

func (s *TransactionServiceImpl) CreateImmediateTransaction(tx *domain.Transaction) error {

	// 1. Generate Transaction ID
	tx.ID = uuid.New().String()
	tx.Status = domain.StatusInitiated

	// 2. Persist transaction
	err := s.txRepo.Save(tx)
	if err != nil {
		return err
	}

	// 3. Transition to COMPLETED
	err = state.Transition(tx, domain.StatusCompleted)
	if err != nil {
		return err
	}

	// 4. Update status in DB
	err = s.txRepo.UpdateStatus(tx.ID, tx.Status)
	if err != nil {
		return err
	}

	// 5. Append ledger entry
	ledger := &domain.LedgerEntry{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		EntryType:     domain.LedgerDebit,
		Amount:        tx.Amount,
		Source:        "API",
	}

	err = s.ledgerRepo.Append(ledger)
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionServiceImpl) GetTransaction(id string) (*domain.Transaction, []domain.LedgerEntry, error) {

	tx, err := s.txRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	ledger, err := s.ledgerRepo.FindByTransactionID(id)
	if err != nil {
		return nil, nil, err
	}

	return tx, ledger, nil
}

