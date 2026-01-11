package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"transaction-manager/internal/api"
	"transaction-manager/internal/config"
	"transaction-manager/internal/repository"
	"transaction-manager/internal/service"
)

func main() {

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}

	db, err := config.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	txRepo := repository.NewPgxTransactionRepo(db)
	ledgerRepo := repository.NewPgxLedgerRepo(db)

	txService := service.NewTransactionService(txRepo, ledgerRepo)
	handler := api.NewTransactionHandler(txService)

	r := gin.Default()
	r.POST("/api/v1/transactions", handler.CreateTransaction)

	log.Println("Transaction Manager API started on :8080")
	r.Run(":8080")
}
