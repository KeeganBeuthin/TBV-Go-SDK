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
	if resultPtr == nil {
		fmt.Println("Error: Received nil pointer in Process_credit_result")
		return utils.CreateErrorResult("Error: Received nil pointer")
	}

	result := utils.GoString(resultPtr, -1)
	fmt.Printf("Raw result string: %s\n", result)
	fmt.Printf("Processing result, global amount: %.2f\n", globalAmount)

	var data struct {
		Results []struct {
			Balance string `json:"balance"`
		} `json:"results"`
	}

	err := json.Unmarshal([]byte(result), &data)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal result: %v", err)
		fmt.Println(errMsg)
		return utils.CreateErrorResult(errMsg)
	}

	fmt.Printf("Unmarshaled data: %+v\n", data)

	if len(data.Results) == 0 {
		errMsg := "No balance found in result"
		fmt.Println(errMsg)
		return utils.CreateErrorResult(errMsg)
	}

	balance, err := strconv.ParseFloat(data.Results[0].Balance, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid balance value: %s", data.Results[0].Balance)
		fmt.Println(errMsg)
		return utils.CreateErrorResult(errMsg)
	}

	fmt.Printf("Parsed balance: %.2f\n", balance)

	newBalance := balance + globalAmount
	responseMessage := fmt.Sprintf("Current balance: %.2f. After credit of %.2f, new balance: %.2f", balance, globalAmount, newBalance)
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
