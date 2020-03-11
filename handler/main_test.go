package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestEventHandler(t *testing.T) {
	tt := map[string]struct {
		ctx         context.Context
		logsEvent   events.CloudwatchLogsEvent
		expectedErr string
	}{
		"empty event": {
			logsEvent:   events.CloudwatchLogsEvent{},
			expectedErr: "EOF",
		},
		"no events": {
			logsEvent: events.CloudwatchLogsEvent{
				AWSLogs: events.CloudwatchLogsRawData{
					Data: "H4sIAI0AaF4AA6tWUMrJT3ctS80rKVayUoiOVajlAgCBPI/bFAAAAA==",
				},
			},
			expectedErr: "",
		},
		"no matching event": {
			logsEvent: events.CloudwatchLogsEvent{
				AWSLogs: events.CloudwatchLogsRawData{
					Data: "H4sIADIGaF4AA6tWUMrJT3ctS80rKVayUoiuVspNLS5OTE8FcpSqFWKUUkFSfom5qTFAkRglv3zfxJLkjBilWqXaWIVaLgDh8yUHQAAAAA==",
				},
			},
			expectedErr: "",
		},
		"bad message": {
			logsEvent: events.CloudwatchLogsEvent{
				AWSLogs: events.CloudwatchLogsRawData{
					Data: "H4sIAOTqaF4AA6tWUMrJT3ctS80rKVayUoiuVspNLS5OTE8FcpSqlWpjFWq5AICclQgkAAAA",
				},
			},
			expectedErr: "unexpected end of JSON input",
		},
		"matching event": {
			logsEvent: events.CloudwatchLogsEvent{
				AWSLogs: events.CloudwatchLogsRawData{
					Data: "H4sIAGUGaF4AA6tWUMrJT3ctS80rKVayUoiuVspNLS5OTE8FcpSqFWKUUkFSfom5qTFAkRgl5/y84vycVJ/89My8GKVapdpYhVouAJbAiGdFAAAA",
				},
			},
			expectedErr: "MissingRegion: could not find region configuration",
		},
	}
	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			err := eventHandler(tc.ctx, tc.logsEvent)
			if tc.expectedErr == "" && err != nil {
				t.Errorf("eventHandler() failed. Got error: %v\n", err)
			} else if err != nil && tc.expectedErr != err.Error() {
				t.Errorf("eventHandler() failed. Expected error: %s. Got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMakeBody(t *testing.T) {
	tt := map[string]struct {
		e        ConsoleLoginEvent
		expected string
	}{
		"empty event": {
			expected: "EventType: \nEventID: \nEventTime: \nEventName: \nUserAgent: \n" +
				"AWS Region: \nSourceIPAddress: \n\nUserIdentity\n\nType: \nAccountID: \n" +
				"UserName: \n",
		},
		"event": {
			e: ConsoleLoginEvent{
				EventType:       "test1",
				EventID:         "testID",
				EventTime:       "testTime",
				EventName:       "testName",
				UserAgent:       "testAgent",
				AWSRegion:       "testRegion",
				SourceIPAddress: "testIP",

				UserIdentity: IAMUserIdentity{
					Type:      "testType",
					AccountID: "testID",
					UserName:  "testName",
				},
			},
			expected: `EventType: test1
EventID: testID
EventTime: testTime
EventName: testName
UserAgent: testAgent
AWS Region: testRegion
SourceIPAddress: testIP

UserIdentity

Type: testType
AccountID: testID
UserName: testName
`,
		},
	}
	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			body := makeBody(&tc.e)
			if tc.expected != body {
				t.Errorf("eventHandler() failed. Expected:\n%q\nGot:\n%q\n", tc.expected, body)
			}
		})
	}
}
