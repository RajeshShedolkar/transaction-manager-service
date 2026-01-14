package utils

import (
	"strings"
)

func GetPaymentTypeViaAPI(dcFlag string) string {
	if strings.ToUpper(dcFlag) == "D" {
		return "DEBIT"
	}
	return "CREDIT"
}
