package imps

import "transaction-manager/internal/domain"

type ImpsEventService interface {
	HandleImpsDebitTx(event domain.TxEvent) error
	HandleImpsCreditTx(event domain.TxEvent) error
	// HandleImpsReversalTx(event domain.TxEvent) error
}
