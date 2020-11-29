package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
	"strings"
	"time"
)

//Refer to https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html
const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "update@prod.bh7cw.me" //"ibh7cw@gmail.com"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

var sess *session.Session
var svc_ses *ses.SES
var svc_db *dynamodb.DynamoDB

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
	log.Printf("Event is NULL: %v", event.Records == nil)

	//log number of records
	log.Printf("Number of Records: %v", len(event.Records))

	//log record message
	for _, m := range event.Records {
		log.Printf("Record Message: %v", m.SNS.Message)

		//send email to the user with sns notification message
		SendSESEmail(m.SNS.Message)
	}

	//log timestamp
	currentTime = time.Now()
	log.Printf("Invocation completed: %v", currentTime.Format("2020-11-26 15:04:05"))

	return nil
}

func initSession() *session.Session {
	log.Println("initialize aws session")
	if sess == nil {
		newSess, err := session.NewSession(&aws.Config{
			Region:aws.String("us-east-1")},
		)
		/*test locally
		newSess, err := session.NewSessionWithOptions(session.Options{
			// Specify profile to load for the session's config
			Profile: "prod",

			// Provide SDK Config options, such as Region.
			Config: aws.Config{
				Region: aws.String("us-east-1"),
			},

			// Force enable Shared Config support
			SharedConfigState: session.SharedConfigEnable,
		})*/

		if err != nil {
			log.Println("can't load the aws session")
			return nil
		} else {
			log.Println("loaded aws session")
			sess = newSess
		}
	}
	return sess
}

func initSESClient() *ses.SES {
	if svc_ses == nil {
		sess = initSession()
		// Create S3 service client
		svc_ses = ses.New(sess)
	}

	return svc_ses
}

func initDBClient() *dynamodb.DynamoDB {
	if svc_db == nil {
		sess = initSession()
		// Create S3 service client
		svc_db = dynamodb.New(sess)
	}

	return svc_db
}

//send email to the user with sns notification message
func SendSESEmail(message string) {
	// Create an SES session.
	svc_ses := initSESClient()

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
	Recipient := strings.TrimLeft(email_context[4], "UserEmail: ")

	Subject := "Notification from bh7cw"
	HtmlBody := "<h1>Notification from bh7cw</h1><p>This email was sent with <a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the <a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
	TextBody := "This email was sent from prod.bh7cw.me with Amazon SES."

	if email_context[0] == "create answer" {
		Subject = fmt.Sprintf("Your question '%v' on bh7cw.me has been answered", strings.TrimLeft(email_context[2], "QuestionText: "))
		HtmlBody = fmt.Sprintf("<h1>Notification from bh7cw.me</h1>"+
			"<p>Hi %v,</p>"+
			"<p>The question %v, %v on bh7cw.me has been answered.</p>"+
			"<p>Answer: %v, %v.</p>"+
			"<p>See more details in %v.</p>"+
			"<p>This email was sent from <a href='%v'>bh7cw.me</a> with <a href='https://aws.amazon.com/ses/'>Amazon SES</a>.</p>",
			strings.TrimLeft(email_context[3], "UserName: "), email_context[1], email_context[2], email_context[5], email_context[6], email_context[7], email_context[7])
		TextBody = fmt.Sprintf("Hi %v,\n"+
			"The question %v, %v on bh7cw.me has been answered.\n"+
			"Answer: %v, %v.\n"+
			"See more details in %v.\n"+
			"This email was sent from bh7cw.me with Amazon SES.\n", strings.TrimLeft(email_context[3], "UserName: "), email_context[1], email_context[2], email_context[5], email_context[6], email_context[7])
	} else if email_context[0] == "update answer" {
		Subject = fmt.Sprintf("The answer '%v' to your question '%v' on bh7cw.me has been updated", strings.TrimLeft(email_context[6], "AnswerText: "), strings.TrimLeft(email_context[2], "QuestionText: "))
		HtmlBody = fmt.Sprintf("<h1>Notification from bh7cw.me</h1>"+
			"<p>Hi %v,</p>"+
			"<p>The answer to your question %v, %v on bh7cw.me has been updated.</p>"+
			"<p>Answer: %v, %v.</p>"+
			"<p>See more details in %v.</p>"+
			"<p>This email was sent from <a href='%v'>bh7cw.me</a> with <a href='https://aws.amazon.com/ses/'>Amazon SES</a>.</p>",
			strings.TrimLeft(email_context[3], "UserName: "), email_context[1], email_context[2], email_context[5], email_context[6], email_context[7], email_context[7])
		TextBody = fmt.Sprintf("Hi %v,\n"+
			"The answer to your question %v, %v on bh7cw.me has been updated.\n"+
			"Answer: %v, %v.\n"+
			"See more details in %v.\n"+
			"This email was sent from bh7cw.me with Amazon SES.\n", strings.TrimLeft(email_context[3], "UserName: "), email_context[1], email_context[2], email_context[5], email_context[6], email_context[7])
	} else if email_context[0] == "delete answer" {
		Subject = fmt.Sprintf("The answer '%v' to your question '%v' on bh7cw.me has been deleted", strings.TrimLeft(email_context[6], "AnswerText: "), strings.TrimLeft(email_context[2], "QuestionText: "))
		HtmlBody = fmt.Sprintf("<h1>Notification from bh7cw.me</h1>"+
			"<p>Hi %v,</p>"+
			"<p>The answer %v, %v to your question %v, %v on bh7cw.me has been deleted.</p>"+
			"<p>This email was sent from bh7cw.me with <a href='https://aws.amazon.com/ses/'>Amazon SES</a>.</p>",
			strings.TrimLeft(email_context[3], "UserName: "), email_context[5], email_context[6], email_context[1], email_context[2])
		TextBody = fmt.Sprintf("Hi %v,\n"+
			"The answer %v, %v to your question %v, %v on bh7cw.me has been deleted.\n"+
			"This email was sent from bh7cw.me with Amazon SES.\n", strings.TrimLeft(email_context[3], "UserName: "), email_context[5], email_context[6], email_context[1], email_context[2])
	} else {
		log.Println("The message is not started as expected")
	}

	//search for email, if already sent, return, otherwise, put in DynamoDB table, and send email
	isExist := searchItemInDynamoDB(TextBody)
	if isExist {
		log.Println("The email has already been sent")
		return
	}

	if err := addItemToDynamoDB(TextBody); err != nil {
		log.Printf("Failed to put email item into DynamoDB table: %v", err)
		return
	}

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
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
	result, err := svc_ses.SendEmail(input)

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

func searchItemInDynamoDB(TextBody string) bool {
	//initialize dynamodb client
	svc_db := initDBClient()

	tableName := "csye6225"

	result, err := svc_db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(TextBody),
			},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if result.Item == nil {
		log.Println("Search email in dynamodb: false")
		return false
	}

	log.Printf("Got item output: %v", result)
	return true
}

//add the email to DynomoDB to avoid sending duplicate emails to users
//refer to https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-create-table-item.html
func addItemToDynamoDB(TextBody string) error {
	//initialize dynamodb client
	svc_db := initDBClient()

	tableName := "csye6225"

	item := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(TextBody),
		},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err := svc_db.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %v\n", err)
		return err
	}

	log.Println("Successfully added email: '" + TextBody + "'")

	return nil
}

func main() {
	lambda.Start(handleRequest)
	//test
	//SendSESEmail("update answer, QuestionID: 1, QuestionText: meaning of cat, UserName: Jane Jenny, UserEmail: jingzhangng20@gmail.com, AnswerID: 1, AnswerText: lovely, Link: http://prod.bh7cw.me:80/v1/question/b1db1852-5c5f-457c-b94d-56b917064eee/answer/931bb982-3573-4187-a8e6-d0870901b880")
	//SendSESEmail("delete answer, QuestionID: 1, QuestionText: meaning of cat, UserName: Jane Jenny, UserEmail: jingzhangng20@gmail.com, AnswerID: 1, AnswerText: lovely")
}
