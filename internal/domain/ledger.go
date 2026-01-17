package domain

type LedgerEntryType string

const (
	// System / lifecycle
	LedgerInit LedgerEntryType = "INIT"

	// Authorization / hold
	LedgerAuth  LedgerEntryType = "AUTH"  // for card (HOLD)
	LedgerBlock LedgerEntryType = "BLOCK" // balance blocked (HOLD)

	// Processing (optional / intermediate)
	LedgerProcess LedgerEntryType = "PROCESS"

	// Final money movement
	LedgerDebit  LedgerEntryType = "DEBIT"  // customer debited
	LedgerCredit LedgerEntryType = "CREDIT" // used only for internal settlement, NOT beneficiary

	// Settlement / clearing
	LedgerSettlement LedgerEntryType = "SETTLEMENT"

	// Compensation flows
	LedgerRelease  LedgerEntryType = "RELEASE"  // unblock hold
	LedgerReversal LedgerEntryType = "REVERSAL" // debit reversed
	LedgerRefund   LedgerEntryType = "REFUND"   // user-visible refund
)

type LedgerEntry struct {
	ID            string
	TransactionID string
	AccountRefId  string
	DcFlag        string          // D or C
	EntryType     LedgerEntryType // AUTH, SETTLEMENT, RELEASE, REVERSAL, DEBIT, CREDIT, REFUND
	Amount        int64
	Source        string // API, EVENT, SYSTEM
	Msg           string
}
