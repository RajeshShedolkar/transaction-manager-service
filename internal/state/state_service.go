package state

import (
	"errors"
	"transaction-manager/internal/domain"
)

func Transition(tx *domain.Transaction, newState domain.TransactionStatus) error {
	if !CanTransition(tx.Status, newState) {
		return errors.New("invalid state transition")
	}

	tx.Status = newState
	return nil
}
