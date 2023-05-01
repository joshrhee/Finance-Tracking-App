package GetTransactions

import (
	"context"
	"fmt"

	"github.com/plaid/plaid-go/v3/plaid"
)

type Transaction struct {
	Date     string   `json:"date"`
	Amount   float64  `json:"amount"`
	Category []string `json:"category"`
	Name     string   `json:"name"`
}

func GetTransactions(accessToken *string, client *plaid.APIClient, FirstDayOfPreviousMonth string, LastDayOfPreviousMonth string, ctx context.Context) []Transaction {
	// Getting Transaction
	transactionRequest := plaid.NewTransactionsGetRequest(
		*accessToken,
		FirstDayOfPreviousMonth,
		LastDayOfPreviousMonth,
	)

	//fmt.Println("accessToken: ", *accessToken)
	//fmt.Println("FirstDayOfPreviousMonth: ", FirstDayOfPreviousMonth)
	//fmt.Println("LastDayOfPreviousMonth: ", LastDayOfPreviousMonth)

	//fmt.Println("TransactionRequest: ", *transactionRequest)

	options := plaid.TransactionsGetRequestOptions{
		Count:  plaid.PtrInt32(100),
		Offset: plaid.PtrInt32(0),
	}

	transactionRequest.SetOptions(options)

	//fmt.Println("After SetOptions, TransactionRequest: ", *transactionRequest)

	transactionResponse, _, err := client.PlaidApi.TransactionsGet(ctx).TransactionsGetRequest(*transactionRequest).Execute()
	if err != nil {
		fmt.Errorf("Transaction get error: ", err)
	}

	//fmt.Println("transactionResponse: ", transactionResponse)

	var editedTransactions []Transaction
	transactions := transactionResponse.GetTransactions()

	//fmt.Println("transactions: ", transactions)
	//fmt.Println("transactionResponse.TotalTransactions: ", transactionResponse.TotalTransactions)

	for i := 0; i < int(transactionResponse.TotalTransactions); i++ {
		editedTransaction := Transaction{
			Date:     transactions[i].Date,
			Amount:   transactions[i].Amount,
			Category: transactions[i].Category,
			Name:     transactions[i].Name,
		}
		editedTransactions = append(editedTransactions, editedTransaction)
	}

	fmt.Println("editedTransactions: ", editedTransactions)
	return editedTransactions
}