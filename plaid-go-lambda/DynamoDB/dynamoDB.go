package DynamoDB

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"sync"
)

var (
	once sync.Once
	db   *dynamodb.Client
)

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
func PutDynamoDB(clientUserId string, accessToken string) {
	fmt.Println("PutDynamoDB is called!!!!")

	fmt.Println("clientUserId: ", clientUserId)
	fmt.Println("accessToken: ", accessToken)

	// Define the DynamoDB item to be put
	item := map[string]types.AttributeValue{
		"clientUserID": &types.AttributeValueMemberS{Value: clientUserId},
		"accessToken":  &types.AttributeValueMemberS{Value: accessToken},
	}

	fmt.Println("DynamoDB, item: ", item)

	// Create the input object for the PutItem API call
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("PlaidIdTokenTable"),
	}

	fmt.Println("DynamoDB, input: ", input)

	// Call the PutItem API
	_, err := db.PutItem(context.Background(), input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("PutItem successful for clientUserID %s\n", clientUserId)
}
