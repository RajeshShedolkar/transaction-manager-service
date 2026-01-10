package main

import (
	"log"
	// "net/http"
	"os"
	"transaction-manager/internal/config"
	"transaction-manager/internal/repository"
)

func main() {
	log.Println("Starting Transaction Manager Service...")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}

	db, err := config.NewDB(dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	txRepo := repository.NewPgxTransactionRepo(db)
	ledgerRepo := repository.NewPgxLedgerRepo(db)

	log.Println("Transaction Manager started with pgx DB")

	_ = txRepo
	_ = ledgerRepo

	select {}
	
}
