package transactions

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KeeganBeuthin/TBV-Go-SDK/pkg/utils"
)

var globalAmount float64

func Execute_credit_leg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Invalid amount value: %s", amount))
	}
	globalAmount = amountFloat

	fmt.Printf("Executing credit leg for amount: %s, account: %s\n", amount, account)

	query := fmt.Sprintf(`PREFIX ex: <http://example.org/>
        SELECT ?balance
        WHERE {
          ex:%s ex:hasBalance ?balance .
        }`, account)

	return utils.StringToPtr(query)
}

func Process_credit_result(resultPtr *byte) *byte {
	result := utils.GoString(resultPtr, -1)
	fmt.Printf("Processing result: %s, amount: %.2f\n", result, globalAmount)

	var data struct {
		Results []struct {
			Balance string `json:"balance"`
		} `json:"results"`
	}

	err := json.Unmarshal([]byte(result), &data)
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Failed to unmarshal result: %v", err))
	}

	if len(data.Results) == 0 {
		return utils.CreateErrorResult("No balance found in result")
	}

	balance, err := strconv.ParseFloat(data.Results[0].Balance, 64)
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Invalid balance value: %s", data.Results[0].Balance))
	}

	newBalance := balance + globalAmount
	responseMessage := fmt.Sprintf("Current balance: %.0f. After credit of %.0f, new balance: %.0f", balance, globalAmount, newBalance)
	fmt.Println(responseMessage)

	return utils.StringToPtr(responseMessage)
}

func Execute_debit_leg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)
	message := fmt.Sprintf("Debiting %s from account %s", amount, account)
	fmt.Printf("Created message: \"%s\"\n", message)
	return utils.StringToPtr(message)
}
