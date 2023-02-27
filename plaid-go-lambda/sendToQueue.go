import (
    "encoding/json"
    "fmt"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Object struct {
    // Define your object structure here
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

func Handler(event events.SQSEvent) error {
    // Create a new AWS SDK configuration
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        return err
    }

    // Create a new SQS client
    svc := sqs.NewFromConfig(cfg)

    // Loop over the messages in the event
    for _, message := range event.Records {
        // Parse the message body as an Object
        var obj Object
        err := json.Unmarshal([]byte(message.Body), &obj)
        if err != nil {
            fmt.Printf("Error parsing message body: %v\n", err)
            continue
        }

        // Process the Object here
        fmt.Printf("Processing Object: %+v\n", obj)
    }

    return nil
}
