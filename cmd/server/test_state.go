package main

import (
	"fmt"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/state"
)

func main() {
	tx := domain.Transaction{
		ID:     "TXN1",
		Status: domain.StatusInitiated,
	}

	err := state.Transition(&tx, domain.StatusCompleted)
	fmt.Println("Status:", tx.Status, "Error:", err)

	err = state.Transition(&tx, domain.StatusPending)
	fmt.Println("Status:", tx.Status, "Error:", err)
}
