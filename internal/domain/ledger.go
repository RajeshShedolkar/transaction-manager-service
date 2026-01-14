package domain

type LedgerEntryType string

const (
	LedgerDebit      LedgerEntryType = "DEBIT"
	LedgerCredit     LedgerEntryType = "CREDIT"
	LedgerAuth       LedgerEntryType = "AUTH"
	LedgerSettlement LedgerEntryType = "SETTLEMENT"
	LedgerRelease    LedgerEntryType = "RELEASE"
	LedgerReversal   LedgerEntryType = "REVERSAL"
)

type LedgerEntry struct {
	ID            string
	TransactionID string
	EntryType     LedgerEntryType
	Amount        int64
	Source        string // API / EVENT
}
