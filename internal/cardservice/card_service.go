package cardservice

import "transaction-manager/internal/domain"

// CardEventService handles card-related business events
// and applies ledger + transaction updates.
type CardEventService interface {

	// AUTH_SUCCESS → ledger AUTH → txn AUTHORIZED
	HandleAuthSuccess(event domain.CardEvent) error

	// SETTLEMENT_STARTED → marker → txn PROCESSING
	HandleSettlementStarted(event domain.CardEvent) error

	// DEBIT_CONFIRMED → ledger DEBIT
	HandleDebitConfirmed(event domain.CardEvent) error

	// CREDIT_CONFIRMED → ledger CREDIT → txn COMPLETED
	HandleCreditConfirmed(event domain.CardEvent) error

	// CANCEL_REQUESTED → ledger RELEASE → txn RELEASED
	HandleCancel(event domain.CardEvent) error

	// SETTLEMENT_FAILED → ledger REVERSAL → txn FAILED
	HandleSettlementFailed(event domain.CardEvent) error

	// REFUND_PROCESSED → ledger REFUND → txn REFUNDED
	HandleRefund(event domain.CardEvent) error
}
