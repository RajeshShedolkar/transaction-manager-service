package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"transaction-manager/internal/api"
	"transaction-manager/internal/cardservice"
	"transaction-manager/internal/config"
	"transaction-manager/internal/events"
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

	// ---------- Kafka Consumer ----------
	api.StartConsumers(handler, eventRepo, eventHandler)
	api.StartNEFTConsumers(handler, eventRepo, eventHandler)
	logger.Log.Info("Kafka idempotent consumer started")

	// ---------- Start HTTP ----------
	log.Println("Transaction Manager API started on :8080")
	r.Run(":8080")
}
