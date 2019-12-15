package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-lambda-go/lambda"
)

type CloudWatchInput struct {
	PhoneNumber string `json:"phoneNumber"`
	Message string `json:"message"`
}

func getEnv (key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s is not set in environment variables", key)
	}
	return env
}

func Handler (input CloudWatchInput) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Fatalf("Error occurred while initializing Session: %v", err)
	}

	svc := sns.New(sess)
	topicArn := getEnv("TOPIC_ARN")
	sendSMSMessage(svc, &topicArn, &input.PhoneNumber, &input.Message)
}

func main () {
	lambda.Start(Handler)
}

func sendSMSMessage(client *sns.SNS, topicArn, phoneNumber, msg *string) {
	var attrs = map[string]*sns.MessageAttributeValue {
		"endpoint": &sns.MessageAttributeValue{
			DataType: aws.String("String"),
			StringValue: phoneNumber,
		},
	}

	_, err := client.Publish(&sns.PublishInput{
		Message: msg,
		MessageAttributes: attrs,
		TopicArn: topicArn,
	})

	if err != nil {
		log.Printf("Failed to publish for endpoint %s \n%v", *phoneNumber, err)
	}
}