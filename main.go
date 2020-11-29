package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"
	"os"
	"strings"
	"time"
)

//Refer to https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html
const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "update@prod.bh7cw.me"//"ibh7cw@gmail.com"

	// The character encoding for the email.
	CharSet = "UTF-8"
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

		//send email to the user with sns notification message
		SendSESEmail(m.SNS.Message)
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

//send email to the user with sns notification message
func SendSESEmail(message string) {
	sess, err := session.NewSession(&aws.Config{
		Region:aws.String("us-east-1")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	//message content:
	//create/update(8): create/update/delete answer, QuestionID: %v, QuestionText: %v, UserName: %v %v, UserEmail: %v, AnswerID: %v, AnswerText: %v, link: %v
	//delete(7): QuestionID: %v, QuestionText: %v, UserName: %v %v, UserEmail: %v, AnswerID: %v, AnswerText: %v
	//create/update/delete answer: 0
	//QuestionID: 1
	//QuestionText: 2
	//UserName: 3
	//UserEmail: 4
	//AnswerID: 5
	//AnswerText: 6
	//Link: 7
	email_context := strings.Split(message, ",")

	if len(email_context) != 7 && len(email_context) != 8 {
		log.Println("Message is not as expected")
		return
	}

	//prepare information before assemble the email
	Recipient := email_context[4]

	Subject := "Notification from bh7cw"
	HtmlBody :=  "<h1>Notification from bh7cw</h1><p>This email was sent with <a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the <a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
	TextBody := "This email was sent from prod.bh7cw.me with Amazon SES."

	if email_context[0] == "create answer" {
		Subject = fmt.Sprintf("Your question %v on bh7cw.me has been answered", email_context[2])
		HtmlBody = fmt.Sprintf("<h1>Notification from bh7cw.me</h1>" +
			"<p>Hi %v,</p>" +
			"<p>The question QuestionID: %v, QuestionText: %v on bh7cw.me has been answered.</p>" +
			"<p>Answer: AnswerID: %v, AnswerText: %v.</p>" +
			"<p>See more details in %v.</p>" +
			"<p>This email was sent from <a href='%v'>bh7cw.me</a> with <a href='https://aws.amazon.com/ses/'>Amazon SES</a>.</p>",
			email_context[3], email_context[1], email_context[2], email_context[5], email_context[6], email_context[7], email_context[7])
		TextBody = fmt.Sprintf("Hi %v,\n" +
			"The question QuestionID: %v, QuestionText: %v on bh7cw.me has been answered.\n" +
			"Answer: AnswerID: %v, AnswerText: %v.\n" +
			"See more details in %v.\n" +
			"This email was sent from bh7cw.me with Amazon SES.\n", email_context[3], email_context[1], email_context[2], email_context[5], email_context[6], email_context[7])
	} else if email_context[0] == "update answer" {
		Subject = fmt.Sprintf("The answer %v to your question %v on bh7cw.me has been updated", email_context[6], email_context[2])
		HtmlBody = fmt.Sprintf("<h1>Notification from bh7cw.me</h1>" +
			"<p>Hi %v,</p>" +
			"<p>The answer to your question QuestionID: %v, QuestionText: %v on bh7cw.me has been updated.</p>" +
			"<p>Answer: AnswerID: %v, AnswerText: %v.</p>" +
			"<p>See more details in %v.</p>" +
			"<p>This email was sent from <a href='%v'>bh7cw.me</a> with <a href='https://aws.amazon.com/ses/'>Amazon SES</a>.</p>",
			email_context[3], email_context[1], email_context[2], email_context[5], email_context[6], email_context[7], email_context[7])
		TextBody = fmt.Sprintf("Hi %v,\n" +
			"The answer to your question QuestionID: %v, QuestionText: %v on bh7cw.me has been updated.\n" +
			"Answer: AnswerID: %v, AnswerText: %v.\n" +
			"See more details in %v.\n" +
			"This email was sent from bh7cw.me with Amazon SES.\n", email_context[3], email_context[1], email_context[2], email_context[5], email_context[6], email_context[7])
	} else if email_context[0] == "delete answer" {
		Subject = fmt.Sprintf("The answer %v to your question %v on bh7cw.me has been deleted", email_context[6], email_context[2])
		HtmlBody = fmt.Sprintf("<h1>Notification from bh7cw.me</h1>" +
			"<p>Hi %v,</p>" +
			"<p>The answer AnswerID: %v, AnswerText: %v to your question QuestionID: %v, QuestionText: %v on bh7cw.me has been deleted.</p>" +
			"<p>This email was sent from bh7cw.me with <a href='https://aws.amazon.com/ses/'>Amazon SES</a>.</p>",
			email_context[3], email_context[5], email_context[6], email_context[1], email_context[2])
		TextBody = fmt.Sprintf("Hi %v,\n" +
			"The answer AnswerID: %v, AnswerText: %v to your question QuestionID: %v, QuestionText: %v on bh7cw.me has been deleted.\n" +
			"This email was sent from bh7cw.me with Amazon SES.\n", email_context[3], email_context[5], email_context[6], email_context[1], email_context[2])
	} else {
		log.Println("The message is not started as expected")
	}

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{
			},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}

		return
	}

	log.Println("Email Sent to address: " + Recipient)
	log.Println(result)
}

func main() {
	lambda.Start(handleRequest)
}
