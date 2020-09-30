package email

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/GSA/grace-log-parser/handler/modules"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ses"
)

var typesMap map[string]Type

// Register stores the provided type in the map of registered types
func Register(typ Type) {
	disabled := os.Getenv("DISABLED_MODULES")
	modules := strings.Split(disabled, ",")
	for _, t := range modules {
		if strings.EqualFold(typ.Name(), t) {
			fmt.Printf("disabled %s sub-module", typ.Name())
			return
		}
	}
	if typesMap == nil {
		typesMap = make(map[string]Type)
	}
	typesMap[typ.Name()] = typ
	fmt.Printf("registered %s sub-module", typ.Name())
}

// Type interface provides functionality for matching events
type Type interface {
	Name() string
	Begin() error
	Filter(startTime time.Time, endTime time.Time) []*cloudwatchlogs.FilterLogEventsInput
	Process(*modules.Event) (string, error)
	End() error
}

// Email provides the capability for sending emails for different events
type Email struct{}

// New returns a newly instantiated Email object
func New() *Email {
	return &Email{}
}

// Begin calls Type.Begin() for all know Type objects passed as typ
func (e *Email) Begin() error {
	for name, typ := range typesMap {
		if err := typ.Begin(); err != nil {
			return fmt.Errorf("failed Begin() for: %s -> %v", name, err)
		}
	}
	return nil
}

// End calls Type.End() for all known Type objects passed as typ
func (e *Email) End() error {
	for name, typ := range typesMap {
		if err := typ.End(); err != nil {
			return fmt.Errorf("failed End() for: %s -> %v", name, err)
		}
	}
	return nil
}

func (e *Email) Filter(startTime time.Time, endTime time.Time) []*cloudwatchlogs.FilterLogEventsInput {
	if startTime
	for name, typ := range typesMap {
		if err := typ.End(); err != nil {
			return fmt.Errorf("failed End() for: %s -> %v", name, err)
		}
	}
	return nil
}

// Process returns true if the event is an acceptable event for emailing
func (e *Email) Process(evt *modules.Event) error {
	var (
		detail string
		err    error
	)
	for name, typ := range typesMap {
		if detail, err = typ.Process(evt); err != nil {
			if modules.IsNotApplicable(err) {
				continue // skip this Type
			}
			return fmt.Errorf("failed Process() for: %s -> %v", name, err)
		}
		err = e.send(evt, detail)
		if err != nil {
			return err
		}
	}
	return nil
}

// body makes an html email body from a CloudWatch log event message
func (e *Email) body(evt *modules.Event, detail string) (string, error) {
	tmpl := template.Must(template.ParseFiles("email.html"))
	buf := &bytes.Buffer{}

	err := tmpl.Execute(buf, struct {
		Alias  string
		Event  *modules.Event
		Detail string
	}{
		Alias:  os.Getenv("ACCOUNT_ALIAS"),
		Event:  evt,
		Detail: detail,
	})
	if err != nil {
		return "", fmt.Errorf("error creating email body: %v", err)
	}

	return buf.String(), nil
}

// send sends an email via AWS Simple Email Service (SES)
func (e *Email) send(evt *modules.Event, detail string) error {
	sess := session.Must(session.NewSession())
	svc := ses.New(sess)

	body, err := e.body(evt, detail)
	if err != nil {
		return err
	}

	input := ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(strings.Split(os.Getenv("TO_EMAIL"), ",")),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(""),
				},
				Html: &ses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(evt.Type + " " + os.Getenv("ACCOUNT_ALIAS")),
			},
		},
		Source: aws.String(os.Getenv("FROM_EMAIL")),
	}

	resp, err := svc.SendEmail(&input)
	if err != nil {
		return fmt.Errorf("failed to send email for type: %v", err)
	}

	log.Printf("dispatched email: %v", resp)
	return nil
}
