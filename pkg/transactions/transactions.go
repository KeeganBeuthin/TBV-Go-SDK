//go:build js && wasm
// +build js,wasm

package transactions

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KeeganBeuthin/TBV-Go-SDK/pkg/utils"
)

var globalAmount float64

func ExecuteCreditLeg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Invalid amount value \"%s\"", amount))
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

func ProcessCreditResult(result *string) *string {
	if result == nil {
		fmt.Println("ProcessCreditResult: Received nil pointer")
		return nil
	}
	fmt.Printf("ProcessCreditResult called with result: %s\n", *result)
	var data struct {
		Results []struct {
			Balance json.Number `json:"balance"`
		} `json:"results"`
	}
	err := json.Unmarshal([]byte(*result), &data)
	if err != nil {
		fmt.Printf("Error unmarshaling result: %v\n", err)
		return nil
	}
	if len(data.Results) == 0 {
		return nil
	}
	balanceStr := data.Results[0].Balance.String()
	if balanceStr == "" || balanceStr == "null" {
		return nil
	}
	balance, err := data.Results[0].Balance.Float64()
	if err != nil {
		return nil
	}
	message := fmt.Sprintf("Current balance: %.2f", balance)
	return &message
}

func ExecuteDebitLeg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)
	message := fmt.Sprintf("Debiting %s from account %s", amount, account)
	fmt.Printf("Created message: \"%s\"\n", message)
	return utils.StringToPtr(message)
}
