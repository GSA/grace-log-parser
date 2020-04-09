package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
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

// emailData type for email html template
type emailData struct {
	AccountAlias string
	Event        *ConsoleLoginEvent
}

// eventHandler ... Handles log events by parsing them, filtering and sending
//  emails for select event types
func eventHandler(ctx context.Context, logsEvent events.CloudwatchLogsEvent) error {
	data, err := logsEvent.AWSLogs.Parse()
	if err != nil {
		log.Printf("error parsing log data: %v", err)
		return err
	}

	return handleEvents(data.LogEvents)
}

// handleEvents ... parses log messages out of log events
func handleEvents(logEvents []events.CloudwatchLogsLogEvent) error {
	for _, logEvent := range logEvents {
		var message ConsoleLoginEvent

		err := json.Unmarshal([]byte(logEvent.Message), &message)
		if err != nil {
			log.Printf("error unmarshalling log event message: %v", err)
			return err
		}

		err = handleMessage(&message)
		if err != nil {
			return err
		}
	}

	return nil
}

// handleMessage ... filters log event messages and sends email on matches
func handleMessage(message *ConsoleLoginEvent) error {
	log.Printf("** " + message.EventName + " event **")

	if message.EventName == "ConsoleLogin" {
		return message.sendEmail()
	}

	return nil
}

// textBody makes a text email body from a CloudWatch log event message
func (e *ConsoleLoginEvent) textBody() (b string) {
	b += "EventType: " + e.EventType + "\n"
	b += "EventID: " + e.EventID + "\n"
	b += "EventTime: " + e.EventTime + "\n"
	b += "EventName: " + e.EventName + "\n"
	b += "UserAgent: " + e.UserAgent + "\n"
	b += "AWS Region: " + e.AWSRegion + "\n"
	b += "SourceIPAddress: " + e.SourceIPAddress + "\n"

	b += "\nUserIdentity\n\n"
	b += "Type: " + e.UserIdentity.Type + "\n"
	b += "AccountID: " + e.UserIdentity.AccountID + "\n"
	b += "UserName: " + e.UserIdentity.UserName + "\n"

	return b
}

// htmlBody makes a html email body from a CloudWatch log event message
func (e *ConsoleLoginEvent) htmlBody() (b string) {
	tmpl := template.Must(template.ParseFiles("email.html"))
	data := emailData{
		Event:        e,
		AccountAlias: os.Getenv("ACCOUNT_ALIAS"),
	}
	buf := new(bytes.Buffer)

	err := tmpl.Execute(buf, data)
	if err != nil {
		log.Printf("Error creating html email body: %v\n", err)
		return ""
	}

	return buf.String()
}

// sendEmail sends an email via AWS Simple Email Service (SES)
func (e *ConsoleLoginEvent) sendEmail() error {
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
					Data: aws.String(e.textBody()),
				},
				Html: &ses.Content{
					Data: aws.String(e.htmlBody()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(e.EventType + " " + os.Getenv("ACCOUNT_ALIAS")),
			},
		},
		Source: aws.String(os.Getenv("FROM_EMAIL")),
	}

	resp, err := svc.SendEmail(&input)
	if err != nil {
		log.Printf("error sending email: %v", err)
		return err
	}

	log.Printf("**SES Response:\n%v", resp)

	return nil
}

func main() {
	lambda.Start(eventHandler)
}
