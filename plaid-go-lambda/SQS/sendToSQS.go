package SQS

//func SendMessageToSQS(queueUrl string, messages []string) error {
//	// Load the AWS SDK configuration
//	cfg, err := config.LoadDefaultConfig()
//	if err != nil {
//		return err
//	}
//
//	// Create a new SQS client
//	svc := sqs.New(cfg)
//
//	// Loop over the array of messages and send each message to the SQS queue
//	for _, message := range messages {
//		input := &sqs.SendMessageInput{
//			MessageBody: aws.String(message),
//			QueueUrl:    aws.String(queueUrl),
//		}
//
//		_, err := svc.SendMessage(input)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
