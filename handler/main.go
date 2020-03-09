package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// eventHandler ... Handles log events by parsing them, filtering and sending
//  emails for select event types
func eventHandler(ctx context.Context, logsEvent events.CloudwatchLogsEvent) {
	eventTypes := []string{
		"ConsoleLogin",
	}
	excludedKeys := []string{
		"accessKeyId",
		"principalId",
	}

	cwData := logsEvent.AWSLogs.Data

	compressedPayload, err := base64.StdEncoding.DecodeString(cwData)
	if err != nil {
		log.Fatalf("error decoding base64 cloudwatch data: %v", err)
		return
	}

	r, err := gzip.NewReader(bytes.NewReader(compressedPayload))
	if err != nil {
		log.Fatalf("error decompressing cloudwatch data: %v", err)
		return
	}

	s, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("error reading decompressed cloudwatch data: %v", err)
		return
	}

	payload := make(map[string]interface{})

	err = json.Unmarshal(s, &payload)
	if err != nil {
		log.Fatalf("error unmarshalling cloudwatch data to map: %v", err)
		return
	}

	log.Printf("***Payload:\n%v", payload)

	logEvents := payload["logEvents"].([]interface{})

	for _, logEvent := range logEvents {
		log.Printf("**Event (%T): %v", logEvent, logEvent)

		if mStr, ok := logEvent.(map[string]interface{})["message"]; ok {
			log.Printf("**message (%T):\n%v", mStr, mStr)

			message := make(map[string]interface{})

			err = json.Unmarshal([]byte(mStr.(string)), &message)
			if err != nil {
				log.Fatalf("error unmarshalling log event message: %v", err)
				return
			}

			if eventType, ok := message["eventName"]; ok {
				log.Printf("eventType (%T):\n%v\n", eventType, eventType)

				if contains(eventTypes, eventType.(string)) {
					log.Println("**matching eventType**")

					body := makeBody(message, excludedKeys)
					sendEmail(body)
				}
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, a := range s {
		if a == str {
			return true
		}
	}

	return false
}

func makeBody(m map[string]interface{}, x []string) (b string) {
	b = "This is the header *Test*\n\n"
	b += "EventType: " + m["eventType"].(string) + "\n"
	b += "EventId: " + m["eventID"].(string) + "\n"
	b += "EventTime: " + m["eventTime"].(string) + "\n"
	b += "EventName: " + m["eventName"].(string) + "\n"
	b += "UserAgent: " + m["userAgent"].(string) + "\n"
	b += "AWS Region: " + m["awsRegion"].(string) + "\n"
	b += "SourceIPAddress: " + m["sourceIPAddress"].(string) + "\n"

	if !contains(x, "userIdentity") {
		if ui, ok := m["userIdentity"].(map[string]interface{}); ok {
			b += "\nUserIdentity \n\n"
			b += parseSesUserIdentity(ui, x)
		}
	}

	return b
}

func parseSesUserIdentity(ui map[string]interface{}, x []string) (s string) {
	for k, i := range ui {
		switch v := i.(type) {
		case string:
			if !contains(x, k) {
				s += k + ": " + v + "\n"
			}
		case map[string]interface{}:
			s += parseSesUserIdentity(v, x)
		default:
			log.Printf("warning: Unknown UserIdentity value type: %T", v)
		}
	}

	return s
}

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
				Data: aws.String(os.Getenv("SUBJECT")),
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
