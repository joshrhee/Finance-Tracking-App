package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joshrhee/plaid-go-lambda/CreateLinkToken"
	"github.com/joshrhee/plaid-go-lambda/GetAccessToken"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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
	clientUserId                         = ""

	once     sync.Once
	db       *dynamodb.Client
	dynamoDB *dynamodb.Client

	// secretName = "Plaid-Secret"
	// region = "us-east-1"
	// secretString = ""

	FirstDayOfPreviousMonth = ""
	LastDayOfPreviousMonth  = ""
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
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

var ginLambda *ginadapter.GinLambdaV2

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

	// PLAID_PRODUCTS := [2]string{"auth", "transactions"}
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")
	APP_PORT = os.Getenv("APP_PORT")

	// set defaults
	// if PLAID_PRODUCTS == "" {
	// 	PLAID_PRODUCTS = "transactions"
	// }
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

// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################

// Using secret manager
// func retreiveSecrets() {
// 	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create Secrets Manager client
// 	svc := secretsmanager.NewFromConfig(config)

// 	input := &secretsmanager.GetSecretValueInput{
// 		SecretId:     aws.String(secretName),
// 		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
// 	}

// 	result, err := svc.GetSecretValue(context.TODO(), input)
// 	if err != nil {
// 		// For a list of exceptions thrown, see
// 		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
// 		log.Fatal(err.Error())
// 	}

// 	// Decrypts secret using the associated KMS key.
// 	secretString = *result.SecretString

// 	var secretMap map[string]string
// 	jsonErr := json.Unmarshal([]byte(secretString), &secretMap)
// 	if jsonErr != nil {
// 		fmt.Println("Error json Marshall")
// 	}

// 	PLAID_CLIENT_ID = secretMap["PLAID_CLIENT_ID"]

// 	if (PLAID_ENV == "sandbox") {
// 		PLAID_SECRET = secretMap["PLAID_SECRET_SANDBOX"]
// 	} else if (PLAID_ENV == "development") {
// 		PLAID_SECRET = secretMap["PLAID_SECRET_DEV"]
// 	} else if (PLAID_ENV == "production") {
// 		PLAID_SECRET = secretMap["PLAID_SECRET_PROD"]
// 	}

// 	fmt.Println(`secretMap`, secretMap)
// 	fmt.Println(`secretMap["PLAID_CLIENT_ID"]`, secretMap["PLAID_CLIENT_ID"])
// 	fmt.Println(`PLAID_CLIENT_ID`, PLAID_CLIENT_ID)
// 	fmt.Println(`secretMap["PLAID_SECRET"]`, PLAID_SECRET)

// }

//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################

// Creating Link Token
func createLinkToken(c *gin.Context) {
	CreateLinkToken.CreateLinkToken(c, client, PLAID_COUNTRY_CODES, PLAID_REDIRECT_URI, PLAID_PRODUCTS)
}

// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
func getAccessToken(c *gin.Context) {
	GetAccessToken.GetAccessToken(c, client, &accessToken, &itemID, &transferID, PLAID_PRODUCTS, FirstDayOfPreviousMonth, LastDayOfPreviousMonth)
}

// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################
// ########################################################################################

//	type transaction struct {
//		Date     string   `json:"date"`
//		Amount   float64  `json:"amount"`
//		Category []string `json:"category"`
//		Name     string   `json:"name"`
//	}
func getTransactionDateRange() {
	now := time.Now()

	year, month, _ := now.Date()

	FirstDayOfPreviousMonth = ((time.Date(year, month-1, 1, 0, 0, 0, 0, now.Location())).String())[0:10]
	LastDayOfPreviousMonth = (time.Date(year, month, 0, 0, 0, 0, 0, now.Location())).String()[0:10]
}

//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################
//  ########################################################################################

// func SendMessageToSQS(queueUrl string, messages []string) error {
//     // Load the AWS SDK configuration
//     cfg, err := config.LoadDefaultAWSConfig()
//     if err != nil {
//         return err
//     }

//     // Create a new SQS client
//     svc := sqs.New(cfg)

//     // Loop over the array of messages and send each message to the SQS queue
//     for _, message := range messages {
//         input := &sqs.SendMessageInput{
//             MessageBody: aws.String(message),
//             QueueUrl:    aws.String(queueUrl),
//         }

//         _, err := svc.SendMessage(input)
//         if err != nil {
//             return err
//         }
//     }

//     return nil
// }

// Get DynamoDB
func GetDynamoDB() *dynamodb.Client {
	fmt.Println("GetDynamoDB is called!!!!")
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			panic(err)
		}

		db = dynamodb.NewFromConfig(cfg)
	})

	return db
}

// Write to DynamoDB
func putDynamoDB(clientUserId string) {
	fmt.Println("PutDynamoDB is called!!!!")
	// Define the DynamoDB item to be put
	item := map[string]types.AttributeValue{
		"clientUserID": &types.AttributeValueMemberS{Value: clientUserId},
		"accessToken":  &types.AttributeValueMemberS{Value: accessToken},
	}

	// Create the input object for the PutItem API call
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("PlaidIdTokenTable"),
	}

	// Call the PutItem API
	_, err := db.PutItem(context.Background(), input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("PutItem successful for clientUserID %s\n", clientUserId)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	dynamoDB = GetDynamoDB()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "root endpoint!!!",
		})
	})
	router.POST("/create_link_token", createLinkToken)
	router.POST("/get_access_token", getAccessToken)

	// env := os.Getenv("GIN_MODE")
	// fmt.Println("env: ", env)

	ginLambda = ginadapter.NewV2(router)
	lambda.Start(Handler)

}
