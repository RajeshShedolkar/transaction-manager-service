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

func StartConsumers(
	handler *TransactionHandler,
	eventRepo *repository.PgxEventRepo,
	eventHandler *events.EventTransactionHandler,
) {
	brokers := config.KAFKA_BROKERS
	cardAuthReader := kfk.NewKafkaReader(brokers, config.KAFKA_CARD_EVENT_TOPIC, "tm-card-group")
	KfkAccountBalanceBlockedReader := kfk.NewKafkaReader(brokers, config.KafkaAccountBalanceBlockedEvt, config.ACC_TM_GROUP)
	kfkNetworkRequestedReader := kfk.NewKafkaReader(brokers, config.KafkaPaymentIMPSDebitSuccessEvt, config.PAYMENT_TM_GROUP)
	KafkaAccountBalanceDebitedEvtReader := kfk.NewKafkaReader(brokers, config.KafkaAccountBalanceDebitedEvt, config.ACC_TM_GROUP)

	//  payment.events.debit-failed / timeout
	kfkNetworkRequestedFailReader := kfk.NewKafkaReader(brokers, config.KafkaPaymentIMPSDebitFailedEvt, config.PAYMENT_TM_GROUP)
	kfkNetworkRequestedTimedOutReader := kfk.NewKafkaReader(brokers, config.KafkaPaymentIMPSDebitTimeoutEvt, config.PAYMENT_TM_GROUP)

	// payment relesed
	kfkAccBalancedeRealsedReader := kfk.NewKafkaReader(brokers, config.KafkaAccountBalanceReleasedEvt, config.PAYMENT_TM_GROUP)

	go kfk.Consume(cardAuthReader, func(msg []byte) {

		eventHandler.HandleCardEventIdempotent(msg, *eventRepo)
	})

	go kfk.Consume(KfkAccountBalanceBlockedReader, func(msg []byte) {
		logger.Log.Info("Account Balance Blocked event received in TM %s", zap.ByteString("message", msg))
		// handler.HandleAccBalBlocked(msg, *eventRepo)
		handler.HandledPayEvent(
			msg,
			domain.StatusBlockRequested,
			domain.SagaBalanceBlocked,
			domain.COMPLETED,
			domain.StatusNetworkRequested,
			domain.SagaNetworkRequested,
			domain.IN_PROGRESS,
			config.KafkaPaymentIMPSDebitCmd,
			*eventRepo)
	})
	// CONSUME FROM PAYMENT NETWORK
	// ---- SUCCESS CASE-----

	go kfk.Consume(kfkNetworkRequestedReader, func(msg []byte) { // "payment.events.debit-success"
		logger.Log.Info("Payment Network confirms Success event and ready for debit received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusNetworkRequested,
			domain.SagaNetworkRequested,
			domain.COMPLETED,
			domain.StatusIMPSDebited,
			domain.SagaIMPSDebited,
			domain.IN_PROGRESS,
			config.KafkaPaymentIMPSDebitCmd,
			*eventRepo)
	})

	go kfk.Consume(kfkNetworkRequestedReader, func(msg []byte) {
		logger.Log.Info("Payment Network Debit Success event received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusIMPSDebited,
			domain.SagaIMPSDebited,
			domain.COMPLETED,
			domain.StatusFinalDebitFromAcc,
			domain.SagaFinalDebitFromAcc,
			domain.IN_PROGRESS,
			config.KafkaAccountFinalDebitCmd,
			*eventRepo)
	})

	go kfk.Consume(KafkaAccountBalanceDebitedEvtReader, func(msg []byte) {
		logger.Log.Info("Account Balance Debited event received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusCompleted,
			domain.SagaFinalDebitFromAcc,
			domain.COMPLETED,
			domain.StatusCompleted,
			"DONE",
			domain.COMPLETED,
			"",
			*eventRepo)
	})

	// PAYMENT NETWORK FAIL CASE

	go kfk.Consume(kfkNetworkRequestedFailReader, func(msg []byte) {
		logger.Log.Info("Payment Network is emiting debit-failed event and ready for debit received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusNetworkDebitFailed,
			domain.SagaNetworkRequested,
			"",
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			"",
			config.KafkaAccountReleaseHoldCmd,
			*eventRepo)
	})

	go kfk.Consume(kfkNetworkRequestedTimedOutReader, func(msg []byte) {
		logger.Log.Info("Payment Network is emiting debit-failed and ready for debit received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusNetworkTimedOut,
			domain.SagaNetworkRequested,
			"",
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			"",
			config.KafkaAccountReleaseHoldCmd,
			*eventRepo)
	})

	// ACCOUNT SERVICE CONFIRM ABOUT RELEASED FUND
	go kfk.Consume(kfkAccBalancedeRealsedReader, func(msg []byte) {
		logger.Log.Info("ACCOUNT SERVICE CONFIRM ABOUT RELEASED FUND event and ready to make entry in Ledge by TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			"",
			domain.StatusFailed,
			"",
			"",
			"Notify to downstream consumer / account service",
			*eventRepo)
	})

	go kfk.Consume(kfkAccBalancedeRealsedReader, func(msg []byte) {
		logger.Log.Info("ACCOUNT SERVICE CONFIRM ABOUT RELEASED FUND event and ready to make entry in Ledge by TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusReleasedRequested,
			domain.SagaRelease,
			"",
			domain.StatusFailed,
			"",
			"",
			"Notify to downstream consumner/ acc service",
			*eventRepo)
	})

}
