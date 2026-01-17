package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"transaction-manager/internal/config"
// 	"transaction-manager/internal/domain"
// 	"transaction-manager/internal/repository"
// 	"transaction-manager/internal/service"
// 	"transaction-manager/pkg/logger"
// )

// func main() {
// 	if err := logger.Init("dev"); err != nil {
// 		panic("failed to initialize logger: " + err.Error())
// 	}
// 	defer logger.Log.Sync()

// 	dsn := os.Getenv("DB_DSN")
// 	if dsn == "" {
// 		log.Fatal("DB_DSN not set")
// 	}

// 	// Connect DB
// 	db, err := config.NewDB(dsn)
// 	if err != nil {
// 		log.Fatal("DB connection failed:", err)
// 	}
// 	defer db.Close()

// 	// Create repositories
// 	txRepo := repository.NewPgxTransactionRepo(db)
// 	ledgerRepo := repository.NewPgxLedgerRepo(db)

// 	// Create service
// 	txService := service.NewTransactionService(txRepo, ledgerRepo)

// 	// Create transaction object
// 	tx := &domain.Transaction{
// 		PaymentType: "IMMEDIATE",
// 		PaymentMode: "IMPS",
// 		Amount:      1500,
// 		Currency:    "INR",
// 	}

// 	fmt.Println("Before service call:", tx)

// 	// Call service
// 	err = txService.CreateImmediateTransaction(tx)
// 	if err != nil {
// 		log.Fatal("Transaction failed:", err)
// 	}

// 	fmt.Println("After service call:", tx)
// 	fmt.Println("Transaction flow completed successfully")
// }
