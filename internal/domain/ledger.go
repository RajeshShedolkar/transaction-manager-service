package domain

type LedgerEntryType string

const (
	LedgerInit       LedgerEntryType = "INIT"
	LedgerDebit      LedgerEntryType = "DEBIT"
	LedgerCredit     LedgerEntryType = "CREDIT"
	LedgerAuth       LedgerEntryType = "AUTH"
	LedgerProcess    LedgerEntryType = "PROCESS"
	LedgerSettlement LedgerEntryType = "SETTLEMENT"
	LedgerRelease    LedgerEntryType = "RELEASE"
	LedgerReversal   LedgerEntryType = "REVERSAL"
	LedgerRefund     LedgerEntryType = "REFUND"
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
