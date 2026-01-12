package events

import "transaction-manager/internal/service"

func HandleNEFTEvent(txService service.TransactionService, txID string, status string) {
	txService.HandleNEFTSettlement(txID, status)
}
