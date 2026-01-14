package cardservice

import (
	"errors"
	"fmt"
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

func (s *EventTransactionServiceImpl) HandleAuthSuccess(event domain.CardEvent) error {

	tx := &domain.Transaction{
		ID:               uuid.New().String(),
		UserRefId:        event.UserRefId,
		SourceRefId:      event.SourceRefId, //user accout id
		DestinationRefId: event.DestinationRefId,
		DcFlag:           event.DcFlag,
		PaymentType:      "DEEFERRED",
		PaymentMode:      "CARD",
		Amount:           event.Amount,
		Currency:         "INR",
		Status:           domain.StatusAuthorized,
		NetworkTxnId:     event.NetworkTxnId,
		GatewayTxnId:     event.GatewayTxnId,
	}

	if err := s.txRepo.Save(tx); err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		AccountRefId:  tx.SourceRefId,
		DcFlag:        tx.DcFlag,
		EntryType:     domain.LedgerAuth,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	return s.ledgerRepo.Append(ledger)
}

// ---------- SETTLEMENT STARTED ----------
func (s *EventTransactionServiceImpl) HandleSettlementStarted(event domain.CardEvent) error {

	tx, err := s.txRepo.FindByNetworkTxnID(event.NetworkTxnId)
	if err != nil {
		return err
	}
	fmt.Println("Found transaction for settlement:", tx.SourceRefId)
	ledger := &domain.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: tx.ID,
		AccountRefId:  tx.SourceRefId,
		DcFlag:        "D",
		EntryType:     domain.LedgerSettlement,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	if err := s.ledgerRepo.Append(ledger); err != nil {
		return err
	}

	tx.Status = domain.StatusProcessing
	return s.txRepo.UpdateStatus(tx.ID, tx.Status)
}

// ---------- DEBIT CONFIRMED ----------
func (s *EventTransactionServiceImpl) HandleDebitConfirmed(event domain.CardEvent) error {

	tx, err := s.txRepo.FindByNetworkTxnID(event.NetworkTxnId)
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
	}

	return s.ledgerRepo.Append(ledger)
}

// ---------- CREDIT CONFIRMED ----------
func (s *EventTransactionServiceImpl) HandleCreditConfirmed(event domain.CardEvent) error {

	tx, err := s.txRepo.FindByNetworkTxnID(event.NetworkTxnId)
	if err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: tx.ID,
		AccountRefId:  tx.DestinationRefId,
		DcFlag:        "C",
		EntryType:     domain.LedgerCredit,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	if err := s.ledgerRepo.Append(ledger); err != nil {
		return err
	}

	tx.Status = domain.StatusCompleted
	return s.txRepo.UpdateStatus(tx.ID, tx.Status)
}

// ---------- CANCEL REQUESTED (RELEASE) ----------
func (s *EventTransactionServiceImpl) HandleCancel(event domain.CardEvent) error {

	tx, err := s.txRepo.FindByNetworkTxnID(event.NetworkTxnId)
	if err != nil {
		return err
	}

	if tx.Status != domain.StatusAuthorized {
		return errors.New("cancel not allowed in current state")
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: tx.ID,
		AccountRefId:  tx.SourceRefId,
		EntryType:     domain.LedgerRelease,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	if err := s.ledgerRepo.Append(ledger); err != nil {
		return err
	}

	tx.Status = domain.StatusReleased
	return s.txRepo.UpdateStatus(tx.ID, tx.Status)
}

// ---------- SETTLEMENT FAILED (REVERSAL) ----------
func (s *EventTransactionServiceImpl) HandleSettlementFailed(event domain.CardEvent) error {

	tx, err := s.txRepo.FindByNetworkTxnID(event.NetworkTxnId)
	if err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: tx.ID,
		AccountRefId:  tx.SourceRefId,
		EntryType:     domain.LedgerReversal,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	if err := s.ledgerRepo.Append(ledger); err != nil {
		return err
	}

	tx.Status = domain.StatusFailed
	return s.txRepo.UpdateStatus(tx.ID, tx.Status)
}

// ---------- REFUND PROCESSED ----------
func (s *EventTransactionServiceImpl) HandleRefund(event domain.CardEvent) error {

	tx, err := s.txRepo.FindByNetworkTxnID(event.NetworkTxnId)
	if err != nil {
		return err
	}

	ledger := &domain.LedgerEntry{
		ID:            uuid.NewString(),
		TransactionID: tx.ID,
		AccountRefId:  tx.SourceRefId,
		DcFlag:        "C",
		EntryType:     domain.LedgerRefund,
		Amount:        tx.Amount,
		Source:        "EVENT",
	}

	if err := s.ledgerRepo.Append(ledger); err != nil {
		return err
	}

	tx.Status = domain.StatusRefunded
	return s.txRepo.UpdateStatus(tx.ID, tx.Status)
}
