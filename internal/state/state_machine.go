package state

import "transaction-manager/internal/domain"

var allowedTransitions = map[domain.TransactionStatus][]domain.TransactionStatus{
	domain.StatusInitiated: {
		domain.StatusCompleted,
		domain.StatusPending,
		domain.StatusFailed,
		domain.StatusAuthorized,
	},
	domain.StatusPending: {
		domain.StatusCompleted,
		domain.StatusFailed,
		domain.StatusTimedOut,
	},
	domain.StatusAuthorized: {
		domain.StatusCompleted,
		domain.StatusReleased,
	},
}

func CanTransition(from, to domain.TransactionStatus) bool {
	nextStates, ok := allowedTransitions[from]
	if !ok {
		return false
	}

	for _, state := range nextStates {
		if state == to {
			return true
		}
	}
	return false
}
