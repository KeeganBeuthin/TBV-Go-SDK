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

func ProcessCreditResult(resultPtr *byte) *byte {
	if resultPtr == nil {
		return utils.CreateErrorResult("Error executing RDF query")
	}
	resultStr := utils.GoString(resultPtr, -1)
	fmt.Printf("Processing credit result: %s\n", resultStr)
	var queryResult struct {
		Results []struct {
			Balance json.Number `json:"balance"`
		} `json:"results"`
	}
	err := json.Unmarshal([]byte(resultStr), &queryResult)
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Error parsing JSON: %v", err))
	}
	if len(queryResult.Results) == 0 {
		return utils.CreateErrorResult("Invalid result format")
	}
	balance, err := queryResult.Results[0].Balance.Float64()
	if err != nil {
		return utils.CreateErrorResult(fmt.Sprintf("Invalid balance format: %v", err))
	}
	newBalance := balance + globalAmount
	message := fmt.Sprintf("Current balance: %.2f. After credit of %.2f, new balance: %.2f", balance, globalAmount, newBalance)
	return utils.CreateSuccessResult(message)
}

func ExecuteDebitLeg(amountPtr *byte, amountLen int32, accountPtr *byte, accountLen int32) *byte {
	amount := utils.GoString(amountPtr, amountLen)
	account := utils.GoString(accountPtr, accountLen)
	message := fmt.Sprintf("Debiting %s from account %s", amount, account)
	fmt.Printf("Created message: \"%s\"\n", message)
	return utils.StringToPtr(message)
}
