package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-lambda-go/lambda"
)

type CloudWatchInput struct {
	PhoneNumber string `json:"phoneNumber"`
	Message string `json:"message"`
}

func Handler (input CloudWatchInput) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Fatalf("Error occurred while initializing Session: %v", err)
	}

	svc := sns.New(sess)
	sendSMSMessage(svc, &input.PhoneNumber, &input.Message)
}

func main () {
	lambda.Start(Handler)
}

func sendSMSMessage(client *sns.SNS, phoneNumber, msg *string) {
	log.Printf("Publishing an SNS message for phone number: %s", *phoneNumber)
	_, err := client.Publish(&sns.PublishInput{
		Message: msg,
		PhoneNumber: phoneNumber,
	})

	if err != nil {
		log.Printf("Failed to publish to SNS for phone number %s \n%v", *phoneNumber, err)
	}
}