// neft_consumer.go
package api

import (
	"transaction-manager/internal/config"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/events"
	"transaction-manager/internal/kfk"
	"transaction-manager/internal/repository"
	"transaction-manager/pkg/logger"

	"go.uber.org/zap"
)

func StartNEFTConsumers(
	handler *TransactionHandler,
	eventRepo *repository.PgxEventRepo,
	eventHandler *events.EventTransactionHandler,
) {

	brokers := config.KAFKA_BROKERS

	// ACCOUNT EVENTS
	kfkAccBalanceBlocked := kfk.NewKafkaReader(brokers, config.KafkaNEFTAccountBalanceBlockedEvt, config.ACC_TM_GROUP)
	kfkAccBalanceDebited := kfk.NewKafkaReader(brokers, config.KafkaNEFTAccountBalanceDebitedEvt, config.ACC_TM_GROUP)
	kfkAccBalanceReleased := kfk.NewKafkaReader(brokers, config.KafkaNEFTAccountBalanceReleasedEvt, config.ACC_TM_GROUP)

	// NEFT NETWORK EVENTS
	kfkNeftSuccess := kfk.NewKafkaReader(brokers, config.KafkaPaymentNEFTDebitSuccessEvt, config.PAYMENT_TM_GROUP)
	kfkNeftFailed := kfk.NewKafkaReader(brokers, config.KafkaPaymentNEFTDebitFailedEvt, config.PAYMENT_TM_GROUP)
	kfkNeftTimeout := kfk.NewKafkaReader(brokers, config.KafkaPaymentNEFTDebitTimeoutEvt, config.PAYMENT_TM_GROUP)

	// -----------------------------
	// ACCOUNT → TM
	// -----------------------------

	go kfk.Consume(kfkAccBalanceBlocked, func(msg []byte) {
		logger.Log.Info("NEFT Account Balance Blocked received", zap.ByteString("msg", msg))

		handler.HandledPayEvent(
			msg,
			domain.StatusBlockRequested,
			domain.SagaBalanceBlocked,
			domain.COMPLETED,
			domain.StatusNetworkRequested,
			domain.SagaNetworkRequested,
			domain.IN_PROGRESS,
			config.KafkaPaymentNEFTDebitCmd,
			*eventRepo)
	})

	// -----------------------------
	// NEFT NETWORK SUCCESS
	// -----------------------------

	go kfk.Consume(kfkNeftSuccess, func(msg []byte) {
		logger.Log.Info("NEFT Network Debit Success received", zap.ByteString("msg", msg))

		handler.HandledPayEvent(
			msg,
			domain.StatusNetworkRequested,
			domain.SagaNetworkRequested,
			domain.COMPLETED,
			domain.StatusNEFTDebited,
			domain.SagaNEFTDebited,
			domain.IN_PROGRESS,
			config.KafkaAccountFinalDebitCmd,
			*eventRepo)
	})

	// -----------------------------
	// ACCOUNT FINAL DEBIT
	// -----------------------------

	go kfk.Consume(kfkAccBalanceDebited, func(msg []byte) {
		logger.Log.Info("NEFT Account Final Debit received", zap.ByteString("msg", msg))

		handler.HandledPayEvent(
			msg,
			domain.StatusFinalDebitFromAcc,
			domain.SagaFinalDebitFromAcc,
			domain.COMPLETED,
			domain.StatusNEFTCompleted,
			"DONE",
			domain.COMPLETED,
			"",
			*eventRepo)
	})

	// -----------------------------
	// NEFT NETWORK FAILURE
	// -----------------------------

	go kfk.Consume(kfkNeftFailed, func(msg []byte) {
		logger.Log.Info("NEFT Network Failed", zap.ByteString("msg", msg))

		handler.HandledPayEvent(
			msg,
			domain.StatusNEFTFailed,
			domain.SagaNetworkRequested,
			domain.FAILED,
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			domain.IN_PROGRESS,
			config.KafkaAccountReleaseHoldCmd,
			*eventRepo)
	})

	go kfk.Consume(kfkNeftTimeout, func(msg []byte) {
		logger.Log.Info("NEFT Network Timeout", zap.ByteString("msg", msg))

		handler.HandledPayEvent(
			msg,
			domain.StatusNetworkTimedOut,
			domain.SagaNetworkRequested,
			domain.FAILED,
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			domain.IN_PROGRESS,
			config.KafkaAccountReleaseHoldCmd,
			*eventRepo)
	})

	// -----------------------------
	// ACCOUNT RELEASE CONFIRMATION
	// -----------------------------

	go kfk.Consume(kfkAccBalanceReleased, func(msg []byte) {
		logger.Log.Info("NEFT Account Balance Released", zap.ByteString("msg", msg))

		handler.HandledPayEvent(
			msg,
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			domain.COMPLETED,
			domain.StatusNEFTFailed,
			"DONE",
			domain.COMPLETED,
			config.KafkaNEFTTransactionFinalEvt,
			*eventRepo)
	})
}
