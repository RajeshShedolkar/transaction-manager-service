package main

import (
	"fmt"
	"log"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/service"
)

func main() {

	tx := &domain.Transaction{
		PaymentType: "IMMEDIATE",
		PaymentMode: "IMPS",
		Amount:      1000,
		Currency:    "INR",
	}

	fmt.Println("Transaction before:", tx)

	// Here later we will inject real repos
	fmt.Println("Service layer ready")
}
