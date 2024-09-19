package transactions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KeeganBeuthin/TBV-Go-SDK/pkg/utils"
)

var globalAmount float64

func Execute_credit_leg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Error: Invalid amount value: %s", amount))
	}
	globalAmount = amountFloat
	fmt.Printf("Executing credit leg for amount: %s, account: %s\n", amount, account)
	query := fmt.Sprintf(`
		PREFIX ex: <http://example.org/>
		SELECT ?balance
		WHERE {
			ex:%s ex:hasBalance ?balance .
		}
	`, account)
	return utils.StringToPtr(query)
}

func Process_credit_result(resultPtr *byte) *byte {
	if resultPtr == nil {
		return utils.CreateErrorResult("Error: Received nil pointer")
	}
	result := utils.GoString(resultPtr, -1)
	fmt.Printf("Processing credit result: %s\n", result)

	// Find the balance value in the JSON string
	balanceStart := strings.Index(result, `"balance":"`) + 11
	balanceEnd := strings.Index(result[balanceStart:], `"`)

	if balanceStart > 10 && balanceEnd > 0 {
		balanceStr := result[balanceStart : balanceStart+balanceEnd]
		balance, err := strconv.ParseFloat(balanceStr, 64)

		if err == nil {
			newBalance := balance + globalAmount
			responseMessage := fmt.Sprintf("Current balance: %.2f. After credit of %.2f, new balance: %.2f", balance, globalAmount, newBalance)
			fmt.Println(responseMessage)
			return utils.StringToPtr(responseMessage)
		} else {
			errorMessage := fmt.Sprintf("Invalid balance value: %s", balanceStr)
			fmt.Println(errorMessage)
			return utils.CreateErrorResult(errorMessage)
		}
	} else {
		errorMessage := "No balance found in result"
		fmt.Println(errorMessage)
		return utils.CreateErrorResult(errorMessage)
	}
}

func Execute_debit_leg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)
	message := fmt.Sprintf("Debiting %s from account %s", amount, account)
	fmt.Printf("Created message: \"%s\"\n", message)
	return utils.StringToPtr(message)
}
