package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"time"
)

func handleRequest(ctx context.Context, event events.SNSEvent) error {
	log.Print("Event started: ")

	// event
	records, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("EVENT: %s", records)

	//log timestamp
	currentTime := time.Now()
	log.Printf("Invocation started: %v", currentTime.Format("2020-11-26 15:04:05"))

	//log sns event is null
	log.Printf("Event is NULL: %v", event.Records==nil)

	//log number of records
	log.Printf("Number of Records: %v", len(event.Records))

	//log record message
	for i, m := range event.Records {
		log.Printf("Record Message No.%v: %v", i, m.SNS.Message)
	}

	//log timestamp
	currentTime = time.Now()
	log.Printf("Invocation completed: %v", currentTime.Format("2020-11-26 15:04:05"))

	// environment variables
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
		log.Println(element)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
