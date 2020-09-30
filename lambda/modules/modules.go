package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var modulesMap map[string]Module

// Register stores the provided type in the map of registered types
func Register(module Module) {
	disabled := os.Getenv("DISABLED_MODULES")
	modules := strings.Split(disabled, ",")
	for _, m := range modules {
		if strings.EqualFold(module.Name(), m) {
			fmt.Printf("disabled %s module", module.Name())
			return
		}
	}
	if modulesMap == nil {
		modulesMap = make(map[string]Module)
	}
	modulesMap[module.Name()] = module
	fmt.Printf("registered %s module", module.Name())
}

// Module is the interface type for all Modules
type Module interface {
	// Name returns the name of the Module
	Name() string
	// Filter takes the start time (last known event time) and the end time (when the lambda started executing)
	// then returns a *cloudwatchlogs.FilterLogEventsInput which will be used to query for events
	Filter(startTime time.Time, endTime time.Time) []*cloudwatchlogs.FilterLogEventsInput
	// Begin is called before events are processed
	Begin() error
	// Process is called passing an event, true should be
	// returned if this event is acceptable for processing
	Process(*Event) error
	// End is called after all events have been processed
	End() error
}

// NotApplicableErr should be returned from Process() if the event is not applicable
type NotApplicableErr struct {
}

func (n NotApplicableErr) Error() string {
	return "not applicable"
}

// IsNotApplicable returns true if the error is an IsNotApplicable error
func IsNotApplicable(err error) bool {
	_, ok := err.(NotApplicableErr)
	return ok
}

// Event contains all event information
type Event struct {
	Version            string                   `json:"eventVersion"`
	Time               string                   `json:"eventTime"`
	Name               string                   `json:"eventName"`
	Source             string                   `json:"eventSource"`
	Region             string                   `json:"awsRegion"`
	IPAddr             string                   `json:"sourceIPAddress"`
	UserAgent          string                   `json:"userAgent"`
	UserIdentity       *UserIdentity            `json:"userIdentity"`
	RequestParameters  json.RawMessage          `json:"requestParameters"`
	ResponseElements   map[string]interface{}   `json:"responseElements"`
	RequestID          string                   `json:"requestID"`
	ID                 string                   `json:"eventID"`
	ReadOnly           bool                     `json:"readOnly"`
	Resources          []map[string]interface{} `json:"resources"`
	Type               string                   `json:"eventType"`
	RecipientAccountID string                   `json:"recipientAccountId"`
	SharedEventID      string                   `json:"sharedEventID"`
	APIVersion         string                   `json:"apiVersion"`
}

// UserIdentity contains the contents of the UserIdentity property of an event
type UserIdentity struct {
	Type           string          `json:"type"`
	PrincipalID    string          `json:"principalId"`
	Arn            string          `json:"arn"`
	AccountID      string          `json:"accountId"`
	AccessKeyID    string          `json:"accessKeyId"`
	UserName       string          `json:"userName"`
	SessionContext *SessionContext `json:"sessionContext"`
}

// SessionContext contains the contents of the SessionContext property of an event
type SessionContext struct {
	SessionIssuer       *SessionIssuer     `json:"sessionIssuer"`
	WebIDFederationData json.RawMessage    `json:"webIdFederationData"`
	Attributes          *SessionAttributes `json:"attributes"`
}

// SessionIssuer contains the contents of the SessionIssuer property of a session context
type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountID   string `json:"accountId"`
	UserName    string `json:"userName"`
}

// SessionAttributes contains the contents of the SessionAttributes property of a session context
type SessionAttributes struct {
	MFAAuthenticated string `json:"mfaAuthenticated"`
	CreationDate     string `json:"creationDate"`
}
