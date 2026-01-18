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
	sagaRepo   repository.SagaRepository
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	ledgerRepo repository.LedgerRepository,
	sagaRepo repository.SagaRepository,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		txRepo:     txRepo,
		ledgerRepo: ledgerRepo,
		sagaRepo:   sagaRepo,
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
	// err = state.Transition(tx, domain.StatusCompleted)
	// if err != nil {
	// 	logger.Log.Error("STATE_TRANSITION_FAILED", zap.Error(err))
	// 	return err
	// }

	// 4. Update status in DB
	// err = s.txRepo.UpdateStatus(tx.ID, tx.Status)
	// if err != nil {
	// 	logger.Log.Error("UPDATING_TRANSACTION_STATUS_FAILED", zap.Error(err))
	// 	return err
	// }

	// 5. Append ledger entry
	ledger := &domain.LedgerEntry{
		ID:            uuid.New().String(),
		AccountRefId:  tx.SourceRefId,
		TransactionID: tx.ID,
		DcFlag:        tx.DcFlag,
		EntryType:     domain.LedgerInit,
		Amount:        tx.Amount,
		Source:        "API",
	}

	err = s.ledgerRepo.Append(ledger)
	if err != nil {
		logger.Log.Error("APPENDING_LEDGER_ENTRY_FAILED", zap.Error(err))
		return err
	}
	logger.Log.Info("IMMEDIATE_TRANSACTION_INIT_SUCCESSFULLY", zap.String("transaction_id", tx.ID))
	// 6. Append saga step
	s.RecordSagaStep(tx, string(domain.SagaInit), domain.SagaStatusStarted)

	return nil
}

func (s *TransactionServiceImpl) GetTransaction(id string) (*domain.Transaction, []domain.LedgerEntry, error) {

	tx, err := s.txRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	// ledger, err := s.ledgerRepo.FindByTransactionID(id)
	// if err != nil {
	// 	return nil, nil, err
	// }

	return tx, nil, nil
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
		AccountRefId:  tx.SourceRefId,
		TransactionID: tx.ID,
		DcFlag:        tx.DcFlag,
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
		AccountRefId:  tx.SourceRefId,
		DcFlag:        tx.DcFlag,
		EntryType:     ledgerType,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	return s.ledgerRepo.Append(ledger)
}

func (s *TransactionServiceImpl) RecordSagaStep(tx *domain.Transaction, step, status string) {
	err := s.sagaRepo.AddStep(&domain.SagaStep{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		StepName:      step,
		Status:        status,
		TxState:       string(tx.Status),
	})
	if err != nil {
		logger.Log.Error("RECORDING_SAGA_STEP_FAILED", zap.Error(err))
		return
	}
	logger.Log.Info("SAGA_STEP_RECORDED", zap.String("transaction_id", tx.ID), zap.String("step", step), zap.String("status", status))
}

func (s *TransactionServiceImpl) UpdateSagaStatus(txID, status string) {
	_ = s.sagaRepo.UpdateSagaStatus(txID, status)
}

func (s *TransactionServiceImpl) UpdateTransactionWithSaga(txID *domain.Transaction, status domain.TransactionStatus, sagaCurrState string) error {
	logger.Log.Info("UPDATING_TRANSACTION_WITH_SAGA", zap.String("transaction_id", txID.ID), zap.String("new_status", string(status)), zap.String("saga_state", sagaCurrState))
	err := s.txRepo.UpdateStatusWithSaga(txID.ID, status, sagaCurrState)
	if err != nil {
		logger.Log.Error("UPDATING_TRANSACTION_WITH_SAGA_FAILED", zap.Error(err))
		return err
	}
	l_err := s.ledgerRepo.Append(&domain.LedgerEntry{
		ID:            uuid.New().String(),
		TransactionID: txID.ID,
		AccountRefId:  txID.SourceRefId,
		DcFlag:        txID.DcFlag,
		EntryType:     domain.LedgerBlock,
		Amount:        txID.Amount,
		Source:        "Event",
	})
	if l_err != nil {
		logger.Log.Error("APPENDING_LEDGER_ENTRY_FAILED", zap.Error(l_err))
		return l_err
	}
	return nil
}
