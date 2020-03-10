package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// ConsoleLoginEvent ... type for parsing event message json string
type ConsoleLoginEvent struct {
	EventVersion    string          `json:"eventVersion"`
	UserIdentity    IAMUserIdentity `json:"userIdentity"`
	EventTime       string          `json:"eventTime"`
	EventSource     string          `json:"eventSource"`
	EventName       string          `json:"eventName"`
	AWSRegion       string          `json:"awsRegion"`
	SourceIPAddress string          `json:"sourceIPAddress"`
	UserAgent       string          `json:"userAgent"`
	// "requestParameters": null,
	// "responseElements": {
	//    "ConsoleLogin": "Success"
	// },
	// "additionalEventData": {
	EventID            string `json:"eventID"`
	EventType          string `json:"eventType"`
	RecipientAccountID string `json:"recipientAccountId"`
}

// IAMUserIdentity ... type for userIdentity in ConsoleLoginEvent
type IAMUserIdentity struct {
	Type string `json:"type"`
	// *ignore* principalId string
	// *ignore* arn string
	AccountID string `json:"accountId"`
	UserName  string `json:"userName"`
}

// eventHandler ... Handles log events by parsing them, filtering and sending
//  emails for select event types
func eventHandler(ctx context.Context, logsEvent events.CloudwatchLogsEvent) {
	data, err := logsEvent.AWSLogs.Parse()
	if err != nil {
		log.Fatalf("error parsing log data: %v", err)
		return
	}

	handleEvents(data.LogEvents)
}

// handleEvents ... parses log messages out of log events
func handleEvents(logEvents []events.CloudwatchLogsLogEvent) {
	for _, logEvent := range logEvents {
		var message ConsoleLoginEvent

		err := json.Unmarshal([]byte(logEvent.Message), &message)
		if err != nil {
			log.Fatalf("error unmarshalling log event message: %v", err)
			return
		}

		handleMessage(&message)
	}
}

// handleMessage ... filters log event messages and sends email on matches
func handleMessage(message *ConsoleLoginEvent) {
	if message.EventName == "ConsoleLogin" {
		body := makeBody(message)
		sendEmail(body)
	}
}

// makeBody ... makes a text email body from a CloudWatch log event message
func makeBody(e *ConsoleLoginEvent) (b string) {
	b += "EventType: " + e.EventType + "\n"
	b += "EventID: " + e.EventID + "\n"
	b += "EventTime: " + e.EventTime + "\n"
	b += "EventName: " + e.EventName + "\n"
	b += "UserAgent: " + e.UserAgent + "\n"
	b += "AWS Region: " + e.AWSRegion + "\n"
	b += "SourceIPAddress: " + e.SourceIPAddress + "\n"

	b += "\nUserIdentity \n\n"
	b += "Type: " + e.UserIdentity.Type + "\n"
	b += "AccountID:" + e.UserIdentity.AccountID + "\n"
	b += "UserName:" + e.UserIdentity.UserName + "\n"

	return b
}

// sendEmail ... Sends an email via AWS Simple Email Service (SES)
func sendEmail(b string) {
	sess := session.Must(session.NewSession())
	svc := ses.New(sess)
	input := ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(strings.Split(os.Getenv("TO_EMAIL"), ",")),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(b),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Testing Go Port refactor 3"),
			},
		},
		Source: aws.String(os.Getenv("FROM_EMAIL")),
	}

	resp, err := svc.SendEmail(&input)
	if err != nil {
		log.Fatalf("error sending email: %v", err)
	}

	log.Printf("**SES Response:\n%v", resp)
}

func main() {
	lambda.Start(eventHandler)
}
