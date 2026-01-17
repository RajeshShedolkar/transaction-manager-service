package cardservice

import "transaction-manager/internal/domain"

// CardEventService handles card-related business events
// and applies ledger + transaction updates.
type CardEventService interface {

	// AUTH_SUCCESS → ledger AUTH → txn AUTHORIZED
	HandleAuthSuccess(event domain.TxEvent) error

	// SETTLEMENT_STARTED → marker → txn PROCESSING
	HandleSettlementStarted(event domain.TxEvent) error

	// DEBIT_CONFIRMED → ledger DEBIT
	HandleDebitConfirmed(event domain.TxEvent) error

	// CREDIT_CONFIRMED → ledger CREDIT → txn COMPLETED
	HandleCreditConfirmed(event domain.TxEvent) error

	// CANCEL_REQUESTED → ledger RELEASE → txn RELEASED
	HandleCancel(event domain.TxEvent) error

	// SETTLEMENT_FAILED → ledger REVERSAL → txn FAILED
	HandleSettlementFailed(event domain.TxEvent) error

	// REFUND_PROCESSED → ledger REFUND → txn REFUNDED
	HandleRefund(event domain.TxEvent) error
}
