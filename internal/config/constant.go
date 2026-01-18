package config

var KAFKA_BROKERS = []string{"localhost:9092"}
var ENV = "prod"
var ACC_TM_GROUP = "account-svc"
var PAYMENT_TM_GROUP = "payment-network"

// TOPIC to publish account service commands
var KAFKA_ACCOUNT_TOPIC = "account.commands"
var KAFKA_CARD_EVENT_TOPIC = "card-auth-events"
var KAFKA_DLQ_TOPIC = "tm.events.dlq"

// ======================================================
// IMPS – ACCOUNT SERVICE (COMMANDS)
// TM → Account Service
// ======================================================

const (
	KafkaAccountBalanceBlockCmd = "account.commands.balance-block"
	KafkaAccountFinalDebitCmd   = "account.commands.final-debit"
	KafkaAccountReleaseHoldCmd  = "account.commands.release-hold"
)

var AccountCommandTopics = []string{
	KafkaAccountBalanceBlockCmd,
	KafkaAccountFinalDebitCmd,
	KafkaAccountReleaseHoldCmd,
}

// ======================================================
// IMPS – ACCOUNT SERVICE (EVENTS)
// Account Service → TM
// ======================================================

const (
	KafkaAccountBalanceBlockedEvt  = "account.events.balance-blocked"
	KafkaAccountBalanceDebitedEvt  = "account.events.balance-debited"
	KafkaAccountBalanceReleasedEvt = "account.events.balance-released"
)

var AccountEventTopics = []string{
	KafkaAccountBalanceBlockedEvt,
	KafkaAccountBalanceDebitedEvt,
	KafkaAccountBalanceReleasedEvt,
}

// ======================================================
// IMPS – PAYMENT NETWORK (COMMANDS)
// TM → Internal Payment Network / IMPS Adapter
// ======================================================

const (
	KafkaPaymentNetwork = "payment.commands.debit"
	KafkaPaymentIMPSDebitCmd = "payment.commands.debit"
)

var PaymentIMPSCommandTopics = []string{
	KafkaPaymentIMPSDebitCmd,
}

// ======================================================
// IMPS – PAYMENT NETWORK (EVENTS)
// IMPS Adapter / NPCI → TM
// ======================================================

const (
	KafkaPaymentIMPSDebitSuccessEvt = "payment.events.debit-success"
	KafkaPaymentIMPSDebitFailedEvt  = "payment.events.debit-failed"
	KafkaPaymentIMPSDebitTimeoutEvt = "payment.events.debit-timeout"
)

var PaymentIMPSEventTopics = []string{
	KafkaPaymentIMPSDebitSuccessEvt,
	KafkaPaymentIMPSDebitFailedEvt,
	KafkaPaymentIMPSDebitTimeoutEvt,
}

// ======================================================
// IMPS – FINAL TRANSACTION EVENT
// TM → Downstream consumers (Ledger, Notification, etc.)
// ======================================================

const (
	KafkaIMPSTransactionFinalEvt = "imps.events.transaction-final"
)

var FinalTransactionTopics = []string{
	KafkaIMPSTransactionFinalEvt,
}

// ======================================================
// IMPS – ALL TOPICS (Useful for bootstrap / health / ACLs)
// ======================================================

var AllIMPSTopics = []string{
	// Account commands
	KafkaAccountBalanceBlockCmd,
	KafkaAccountFinalDebitCmd,
	KafkaAccountReleaseHoldCmd,

	// Account events
	KafkaAccountBalanceBlockedEvt,
	KafkaAccountBalanceDebitedEvt,
	KafkaAccountBalanceReleasedEvt,

	// Payment network
	KafkaPaymentIMPSDebitCmd,
	KafkaPaymentIMPSDebitSuccessEvt,
	KafkaPaymentIMPSDebitFailedEvt,
	KafkaPaymentIMPSDebitTimeoutEvt,

	// Final event
	KafkaIMPSTransactionFinalEvt,
}
