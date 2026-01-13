package service

import (
	//"errors"
	//"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"transaction-manager/internal/domain"
	"transaction-manager/internal/repository"
	"transaction-manager/internal/state"
	"transaction-manager/pkg/logger"
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
	logger.Log.Info("GENERATING_TRANSACTION_ID")
	tx.ID = uuid.New().String()
	tx.Status = domain.StatusInitiated

	// 2. Persist transaction
	logger.Log.Info("PERSISTING_TRANSACTION")
	err := s.txRepo.Save(tx)
	if err != nil {
		logger.Log.Error("PERSISTING_TRANSACTION_FAILED", zap.Error(err))
		return err
	}

	// 3. Transition to COMPLETED
	err = state.Transition(tx, domain.StatusCompleted)
	if err != nil {
		logger.Log.Error("STATE_TRANSITION_FAILED", zap.Error(err))
		return err
	}

	// 4. Update status in DB
	err = s.txRepo.UpdateStatus(tx.ID, tx.Status)
	if err != nil {
		logger.Log.Error("UPDATING_TRANSACTION_STATUS_FAILED", zap.Error(err))
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
		logger.Log.Error("APPENDING_LEDGER_ENTRY_FAILED", zap.Error(err))
		return err
	}
	logger.Log.Info("IMMEDIATE_TRANSACTION_CREATED_SUCCESSFULLY", zap.String("transaction_id", tx.ID))
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

func (s *TransactionServiceImpl) CreateNEFTTransaction(tx *domain.Transaction) error {

	tx.ID = uuid.New().String()
	tx.Status = domain.StatusInitiated

	// Save initial transaction
	err := s.txRepo.Save(tx)
	if err != nil {
		return err
	}

	// Move to PENDING
	err = state.Transition(tx, domain.StatusPending)
	if err != nil {
		return err
	}

	err = s.txRepo.UpdateStatus(tx.ID, tx.Status)
	if err != nil {
		return err
	}

	// Ledger entry for debit
	ledger := &domain.LedgerEntry{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		EntryType:     domain.LedgerDebit,
		Amount:        tx.Amount,
		Source:        "API",
	}

	return s.ledgerRepo.Append(ledger)
}

func (s *TransactionServiceImpl) HandleNEFTSettlement(txID string, status string) error {

	tx, _, err := s.GetTransaction(txID)
	if err != nil {
		return err
	}

	var newState domain.TransactionStatus
	var ledgerType domain.LedgerEntryType

	if status == "COMPLETED" {
		newState = domain.StatusCompleted
		ledgerType = domain.LedgerSettlement
	} else {
		newState = domain.StatusFailed
		ledgerType = domain.LedgerReversal
	}

	err = state.Transition(tx, newState)
	if err != nil {
		return err
	}

	err = s.txRepo.UpdateStatus(tx.ID, tx.Status)
	if err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		EntryType:     ledgerType,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	return s.ledgerRepo.Append(ledger)
}
