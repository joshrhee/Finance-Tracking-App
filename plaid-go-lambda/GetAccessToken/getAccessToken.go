package GetAccessToken

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"

	"github.com/plaid/plaid-go/v3/plaid"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	PLAID_PRODUCTS                           = ""
	client                  *plaid.APIClient = nil
	accessToken             *string          = nil
	itemID                  *string          = nil
	transferID              *string          = nil
	FirstDayOfPreviousMonth                  = ""
	LastDayOfPreviousMonth                   = ""
)

type transaction struct {
	Date     string   `json:"date"`
	Amount   float64  `json:"amount"`
	Category []string `json:"category"`
	Name     string   `json:"name"`
}

func RenderError(c *gin.Context, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.
		c.JSON(http.StatusOK, gin.H{"error": plaidError})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": originalErr.Error()})
}

// Getting Access token
func GetAccessToken(c *gin.Context, plaidClient *plaid.APIClient, plaidAccessToken *string, plaidItemID *string, plaidTransferID *string, plaidProducts string, firstDayOfPreviousMonth string, lastDayOfPreviousMonth string) {

	client = plaidClient
	accessToken = plaidAccessToken
	itemID = plaidItemID
	transferID = plaidTransferID

	PLAID_PRODUCTS = plaidProducts
	FirstDayOfPreviousMonth = firstDayOfPreviousMonth
	LastDayOfPreviousMonth = lastDayOfPreviousMonth

	encodedRequestBody, _ := ioutil.ReadAll(c.Request.Body)
	stringPublicTokenObject := string(encodedRequestBody)
	// fmt.Println("stringPublicTokenObject: ", stringPublicTokenObject)

	splitedString := strings.Split(stringPublicTokenObject, "")

	publicTokenArray := []string{}
	isEqualFound := false
	for i := 0; i < len(splitedString); i++ {
		if splitedString[i] == "=" {
			isEqualFound = true
			continue
		}

		if !isEqualFound {
			continue
		}

		publicTokenArray = append(publicTokenArray, splitedString[i])
	}

	publicToken := strings.Join(publicTokenArray, "")

	// publicToken := c.Query("public_token")
	// uid := c.PostForm("uid")
	fmt.Println("publicToken: ", publicToken)
	// fmt.Println("uid: ", uid)
	if publicToken == "" {
		fmt.Println("public token is not exist!!")
	}
	ctx := context.Background()

	// exchange the public_token for an access_token
	exchangePublicTokenResp, _, err := client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	if err != nil {
		RenderError(c, err)
		return
	}

	*accessToken = exchangePublicTokenResp.GetAccessToken()

	// Add DynamoDB for {clientUserID: accessToken}
	//putDynamoDB(clientUserId)

	*itemID = exchangePublicTokenResp.GetItemId()
	if itemExists(strings.Split(PLAID_PRODUCTS, ","), "transfer") {
		*transferID, err = authorizeAndCreateTransfer(ctx, client, *accessToken)
	}

	// Getting Transaction
	transactionRequest := plaid.NewTransactionsGetRequest(
		*accessToken,
		FirstDayOfPreviousMonth,
		LastDayOfPreviousMonth,
	)

	fmt.Println("accessToken: ", *accessToken)
	fmt.Println("FirstDayOfPreviousMonth: ", FirstDayOfPreviousMonth)
	fmt.Println("LastDayOfPreviousMonth: ", LastDayOfPreviousMonth)

	fmt.Println("TransactionRequest: ", *transactionRequest)

	options := plaid.TransactionsGetRequestOptions{
		Count:  plaid.PtrInt32(100),
		Offset: plaid.PtrInt32(0),
	}

	transactionRequest.SetOptions(options)

	fmt.Println("After SetOptions, TransactionRequest: ", *transactionRequest)

	transactionResponse, _, err := client.PlaidApi.TransactionsGet(ctx).TransactionsGetRequest(*transactionRequest).Execute()
	if err != nil {
		fmt.Errorf("Transaction get error: ", err)
	}

	fmt.Println("transactionResponse: ", transactionResponse)

	var editedTransactions []transaction
	transactions := transactionResponse.GetTransactions()

	fmt.Println("transactions: ", transactions)
	fmt.Println("transactionResponse.TotalTransactions: ", transactionResponse.TotalTransactions)

	for i := 0; i < int(transactionResponse.TotalTransactions); i++ {
		editedTransaction := transaction{
			Date:     transactions[i].Date,
			Amount:   transactions[i].Amount,
			Category: transactions[i].Category,
			Name:     transactions[i].Name,
		}
		editedTransactions = append(editedTransactions, editedTransaction)
	}

	fmt.Println("editedTransactions: ", editedTransactions)
	c.JSON(http.StatusOK, editedTransactions)
}

// Helper function to determine if Transfer is in Plaid product array
func itemExists(array []string, product string) bool {
	for _, item := range array {
		if item == product {
			return true
		}
	}

	return false
}

// This is a helper function to authorize and create a Transfer after successful
// exchange of a public_token for an access_token. The transfer_id is then used
// to obtain the data about that particular Transfer.
func authorizeAndCreateTransfer(ctx context.Context, client *plaid.APIClient, accessToken string) (string, error) {
	// We call /accounts/get to obtain first account_id - in production,
	// account_id's should be persisted in a data store and retrieved
	// from there.
	accountsGetResp, _, _ := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()

	accountID := accountsGetResp.GetAccounts()[0].AccountId

	transferAuthorizationCreateUser := plaid.NewTransferUserInRequest("FirstName LastName")
	transferAuthorizationCreateRequest := plaid.NewTransferAuthorizationCreateRequest(
		accessToken,
		accountID,
		"credit",
		"ach",
		"1.34",
		"ppd",
		*transferAuthorizationCreateUser,
	)
	transferAuthorizationCreateResp, _, err := client.PlaidApi.TransferAuthorizationCreate(ctx).TransferAuthorizationCreateRequest(*transferAuthorizationCreateRequest).Execute()
	if err != nil {
		return "", err
	}
	authorizationID := transferAuthorizationCreateResp.GetAuthorization().Id

	transferCreateRequest := plaid.NewTransferCreateRequest(
		accessToken,
		accountID,
		authorizationID,
		"credit",
		"ach",
		"1.34",
		"Payment",
		"ppd",
		*transferAuthorizationCreateUser,
	)
	transferCreateResp, _, err := client.PlaidApi.TransferCreate(ctx).TransferCreateRequest(*transferCreateRequest).Execute()
	if err != nil {
		return "", err
	}

	return transferCreateResp.GetTransfer().Id, nil
}
