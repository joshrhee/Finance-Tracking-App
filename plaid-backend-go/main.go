package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v3/plaid"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = ""
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	APP_PORT                             = ""
	client              *plaid.APIClient = nil

	FirstDayOfPreviousMonth = ""
	LastDayOfPreviousMonth = ""
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	// set constants from env
	PLAID_CLIENT_ID = os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET = os.Getenv("PLAID_SECRET")

	if PLAID_CLIENT_ID == "" || PLAID_SECRET == "" {
		log.Fatal("Error: PLAID_SECRET or PLAID_CLIENT_ID is not set. Did you copy .env.example to .env and fill it out?")
	}

	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")
	APP_PORT = os.Getenv("APP_PORT")

	// set defaults
	if PLAID_PRODUCTS == "" {
		PLAID_PRODUCTS = "transactions"
	}
	if PLAID_COUNTRY_CODES == "" {
		PLAID_COUNTRY_CODES = "US"
	}
	if PLAID_ENV == "" {
		PLAID_ENV = "sandbox"
	}
	if APP_PORT == "" {
		APP_PORT = "8000"
	}
	if PLAID_CLIENT_ID == "" {
		log.Fatal("PLAID_CLIENT_ID is not set. Make sure to fill out the .env file")
	}
	if PLAID_SECRET == "" {
		log.Fatal("PLAID_SECRET is not set. Make sure to fill out the .env file")
	}

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)

	getTransactionDateRange()
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}

	router.Use(cors.New(config))

	router.POST("/create_link_token", createLinkToken)
	router.POST("/get_access_token", getAccessToken)
	
	err := router.Run("localhost:8080")
	if err != nil {
		panic("Unable to start the server!")
	}
}

// We store the access_token in memory - in production, store it in a secure
// persistent data store.
var accessToken string
var itemID string

var paymentID string

// The transfer_id is only relevant for the Transfer ACH product.
// We store the transfer_id in memory - in production, store it in a secure
// persistent data store
var transferID string

func createLinkToken(c *gin.Context) {
	linkToken, err := linkTokenCreate(nil)
	if err != nil {
		renderError(c, err)
		return
	}
	fmt.Println("Link token is: ", linkToken)
	c.JSON(200, gin.H{
		"link_token": linkToken,
	})
}

func renderError(c *gin.Context, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.
		c.JSON(http.StatusOK, gin.H{"error": plaidError})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": originalErr.Error()})
}

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func convertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}

	for _, productStr := range productStrs {
		products = append(products, plaid.Products(productStr))
	}

	return products
}

func linkTokenCreate(paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation) (string, error) {
	ctx := context.Background()

	// Institutions from all listed countries will be shown.
	countryCodes := convertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ","))
	redirectURI := PLAID_REDIRECT_URI

	// This should correspond to a unique id for the current user.
	// Typically, this will be a user ID number from your application.
	// Personally identifiable information, such as an email address or phone number, should not be used here.
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		countryCodes,
		user,
	)

	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		// The 'payment_initiation' product has to be the only element in the 'products' list.
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		products := convertProducts(strings.Split(PLAID_PRODUCTS, ","))
		request.SetProducts(products)
	}

	if redirectURI != "" {
		request.SetRedirectUri(redirectURI)
	}

	linkTokenCreateResp, _, err := client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	return linkTokenCreateResp.GetLinkToken(), nil
}

type transaction struct {
	Date string `json:"date"`
	Amount float64 `json:"amount"`
	Category []string `json:"category"`
	Name string `json:"name"`
}

func getTransactionDateRange() {
	now := time.Now()

	year, month, _ := now.Date()

	FirstDayOfPreviousMonth = ((time.Date(year, month - 1, 1, 0, 0, 0, 0, now.Location())).String())[0:10]
	LastDayOfPreviousMonth = (time.Date(year, month, 0, 0, 0, 0, 0, now.Location())).String()[0:10]

	fmt.Println("FirstDayOfPreviousMonth: ", FirstDayOfPreviousMonth)
	fmt.Println("LastDayOfPreviousMonth: ", LastDayOfPreviousMonth)
	
}

func getAccessToken(c *gin.Context) {
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

	publicToken := strings.Join(publicTokenArray,"")

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
		renderError(c, err)
		return
	}

	accessToken = exchangePublicTokenResp.GetAccessToken()
	itemID = exchangePublicTokenResp.GetItemId()
	if itemExists(strings.Split(PLAID_PRODUCTS, ","), "transfer") {
		transferID, err = authorizeAndCreateTransfer(ctx, client, accessToken)
	}

	// fmt.Println("public token: " + publicToken)
	// fmt.Println("access token: " + accessToken)
	// fmt.Println("item ID: " + itemID)

	// c.JSON(http.StatusOK, gin.H{
	// 	"access_token": accessToken,
	// 	"item_id":      itemID,
	// })

	transactionRequest := plaid.NewTransactionsGetRequest(
		accessToken,
		FirstDayOfPreviousMonth,
		LastDayOfPreviousMonth,
	)

	options := plaid.TransactionsGetRequestOptions{
		Count: plaid.PtrInt32(100),
		Offset: plaid.PtrInt32(0),
	}

	transactionRequest.SetOptions(options)

	fmt.Println("TransactionRequest: ", transactionRequest)

	transactionResponse, _, err := client.PlaidApi.TransactionsGet(ctx).TransactionsGetRequest(*transactionRequest).Execute()
	if err != nil {
		fmt.Errorf("Transaction get error: ", err)
	}

	var editedTransactions []transaction
	transactions := transactionResponse.GetTransactions()
	
	for i := 0; i < int(transactionResponse.TotalTransactions); i++ {
		editedTransaction := transaction{
			Date: transactions[i].Date,
			Amount: transactions[i].Amount,
			Category: transactions[i].Category,
			Name: transactions[i].Name,
		}
		fmt.Println("editedTransaction: ", editedTransaction)
		editedTransactions = append(editedTransactions, editedTransaction)
	}


	c.JSON(http.StatusOK, editedTransactions)
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

// Helper function to determine if Transfer is in Plaid product array
func itemExists(array []string, product string) bool {
	for _, item := range array {
		if item == product {
			return true
		}
	}

	return false
}
