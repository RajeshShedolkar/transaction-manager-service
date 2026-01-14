package domain

type EventType string

const (
	EventAuthSuccess       EventType = "AUTH_SUCCESS"
	EventAuthFailed        EventType = "AUTH_FAILED"
	EventSettlementStarted EventType = "SETTLEMENT_STARTED"
	EventDebitConfirmed    EventType = "DEBIT_CONFIRMED"
	EventCreditConfirmed   EventType = "CREDIT_CONFIRMED"
	EventCancelRequested   EventType = "CANCEL_REQUESTED"
	EventSettlementFailed  EventType = "SETTLEMENT_FAILED"
	EventRefundProcessed   EventType = "REFUND_PROCESSED"
)

