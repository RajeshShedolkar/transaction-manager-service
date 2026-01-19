// NEFT – ACCOUNT EVENTS
package config

const (
	KafkaNEFTAccountBalanceBlockedEvt  = "neft.account.events.balance-blocked"
	KafkaNEFTAccountBalanceDebitedEvt  = "neft.account.events.balance-debited"
	KafkaNEFTAccountBalanceReleasedEvt = "neft.account.events.balance-released"
)

// NEFT – NETWORK COMMAND
const (
	KafkaPaymentNEFTDebitCmd = "neft.payment.commands.debit"
)

// NEFT – NETWORK EVENTS
const (
	KafkaPaymentNEFTDebitSuccessEvt = "neft.payment.events.debit-success"
	KafkaPaymentNEFTDebitFailedEvt  = "neft.payment.events.debit-failed"
	KafkaPaymentNEFTDebitTimeoutEvt = "neft.payment.events.debit-timeout"
)

// NEFT – FINAL EVENT
const (
	KafkaNEFTTransactionFinalEvt = "neft.events.transaction-final"
)

// ======================================================
// NEFT – ALL TOPICS
// ======================================================

var AllNEFTTopics = []string{
	// Account events
	KafkaNEFTAccountBalanceBlockedEvt,
	KafkaNEFTAccountBalanceDebitedEvt,
	KafkaNEFTAccountBalanceReleasedEvt,

	// Network command
	KafkaPaymentNEFTDebitCmd,

	// Network events
	KafkaPaymentNEFTDebitSuccessEvt,
	KafkaPaymentNEFTDebitFailedEvt,
	KafkaPaymentNEFTDebitTimeoutEvt,

	// Final event
	KafkaNEFTTransactionFinalEvt,
}
