package CreateLinkToken

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/v3/plaid"
)

var (
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	PLAID_PRODUCTS                       = ""
	client              *plaid.APIClient = nil
	clientUserId                         = ""
)

// Creating Link Token
func CreateLinkToken(c *gin.Context, plaidClient *plaid.APIClient, plaidCountryCode string, plaidPredictUri string, plaidProducts string, mainClientUserId string) {
	// retreiveSecrets

	fmt.Println("CreateLinkToken is started!!!!")

	client = plaidClient
	clientUserId = mainClientUserId

	PLAID_COUNTRY_CODES = plaidCountryCode
	PLAID_REDIRECT_URI = plaidPredictUri
	PLAID_PRODUCTS = plaidProducts

	linkToken, err := linkTokenCreate(nil)
	if err != nil {
		renderError(c, err)
		return
	}
	fmt.Println("Success!!!! Link token: ", linkToken)

	c.JSON(200, gin.H{
		"link_token": linkToken,
	})
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
		ClientUserId: clientUserId,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		countryCodes,
		user,
	)

	fmt.Println("Request!!!!!!!!!: ", request)

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

	fmt.Println("linkTokenCreateResp!!!!!!!: ", linkTokenCreateResp)

	return linkTokenCreateResp.GetLinkToken(), nil
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
