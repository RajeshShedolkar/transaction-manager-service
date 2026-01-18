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
			domain.StatusNetworkRequested,
			domain.SagaNetworkRequested,
			config.KafkaPaymentIMPSDebitCmd,
			*eventRepo)
	})

	go kfk.Consume(kfkNetworkRequestedReader, func(msg []byte) {
		logger.Log.Info("Payment Network confirms Success event and ready for debit received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusNetworkRequested,
			domain.SagaNetworkRequested,
			domain.StatusIMPSDebited,
			domain.SagaIMPSDebited,
			config.KafkaPaymentIMPSDebitCmd,
			*eventRepo)
	})

	go kfk.Consume(kfkNetworkRequestedReader, func(msg []byte) {
		logger.Log.Info("Payment Network Debit Success event received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusIMPSDebited,
			domain.SagaIMPSDebited,
			domain.StatusFinalDebitFromAcc,
			domain.SagaFinalDebitFromAcc,
			config.KafkaAccountFinalDebitCmd,
			*eventRepo)
	})

	go kfk.Consume(KafkaAccountBalanceDebitedEvtReader, func(msg []byte) {
		logger.Log.Info("Account Balance Debited event received in TM %s", zap.ByteString("message", msg))
		handler.HandledPayEvent(
			msg,
			domain.StatusCompleted,
			domain.SagaFinalDebitFromAcc,
			domain.StatusCompleted,
			"DONE",
			"",
			*eventRepo)
	})

}
