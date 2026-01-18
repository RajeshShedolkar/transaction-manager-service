package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"transaction-manager/internal/api"
	"transaction-manager/internal/cardservice"
	"transaction-manager/internal/config"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/events"
	"transaction-manager/internal/kfk"
	"transaction-manager/internal/repository"
	"transaction-manager/internal/service"
	"transaction-manager/pkg/logger"
)

func main() {

	// ---------- Logger ----------
	logger.Init(config.ENV)
	defer logger.Log.Sync()
	logger.Log.Info("service started")

	// ---------- Env ----------
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	// ---------- DB ----------
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}

	db, err := config.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ---------- Repositories ----------
	txRepo := repository.NewPgxTransactionRepo(db)
	ledgerRepo := repository.NewPgxLedgerRepo(db)
	eventRepo := repository.NewPgxEventRepo(db)
	sagaRepo := repository.NewPgxSagaRepo(db)

	// ---------- Services ----------
	txService := service.NewTransactionService(txRepo, ledgerRepo, sagaRepo)
	eventTxService := cardservice.NewEventTransactionService(txRepo, ledgerRepo, eventRepo)
	accountEventPublisher := service.NewKafkaAccountEventPublisher()
	// ---------- API ----------
	handler := api.NewTransactionHandler(txService, accountEventPublisher)
	eventHandler := events.NewEventTransactionHandler(eventTxService)

	r := gin.Default()
	r.POST("/api/v1/transactions", handler.CreateTransaction)
	r.GET("/api/v1/transactions/:id", handler.GetTransaction)

	//kafka producer
	events.InitKafkaTopics()

	// ---------- Kafka Consumer ----------
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
		handler.HandleAccBalBlocked(msg, *eventRepo)
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
			domain.SagaFinalDebited,
			domain.StatusCompleted,
			"DONE",
			"",
			*eventRepo)
	})
	logger.Log.Info("Kafka idempotent consumer started")

	// ---------- Start HTTP ----------
	log.Println("Transaction Manager API started on :8080")
	r.Run(":8080")
}
