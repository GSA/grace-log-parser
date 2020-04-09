package main

import (
	"context"
	"os"
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

func TestTextBody(t *testing.T) {
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
			body := tc.e.textBody()
			if tc.expected != body {
				t.Errorf("eventHandler() failed. Expected:\n%q\nGot:\n%q\n", tc.expected, body)
			}
		})
	}
}

// nolint: funlen
func TestHtmlBody(t *testing.T) {
	err := os.Setenv("ACCOUNT_ALIAS", "test-account-alias")
	if err != nil {
		t.Fatalf("Unexpected error setting ACCOUNT_ALIAS: %v\n", err)
	}

	tt := map[string]struct {
		e        ConsoleLoginEvent
		expected string
	}{
		"empty event": {
			expected: `<head>
  <style>
    table {border-collapse: collapse;}
    td, th {border: 1px solid Black;}
    th {background: LightGray;}
    tr:nth-child(even) {background: #F3F3F3;}
    tr:nth-child(odd) {background: White;}
    .resource {background-color: RoyalBlue; color: White; font-weight: bold;}
    .blank {background-color: White; border: none;}
    .group {background-color: LightBlue;}
  </style>
</head>
<body>
  <h1> in test-account-alias</h1>
  <table>
    <tr><th colspan="2">Event Details</th></tr>
    <tr><td>EventType</td><td></td></tr>
    <tr><td>EventID</td><td></td></tr>
    <tr><td>EventTime</td><td></td></tr>
    <tr><td>EventName</td><td></td></tr>
    <tr><td>UserAgent</td><td></td></tr>
    <tr><td>AWS Region</td><td></td></tr>
    <tr><td>SourceIPAddress</td><td></td></tr>
  </table>
  &nbsp;
  <table>
    <tr><th colspan="2">UserIdentity</th></tr>
    <tr><td>Type</td><td></td></tr>
    <tr><td>AccountID</td><td></td></tr>
    <tr><td>UserName</td><td></td></tr>
  </table>
</body>
`,
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
			expected: `<head>
  <style>
    table {border-collapse: collapse;}
    td, th {border: 1px solid Black;}
    th {background: LightGray;}
    tr:nth-child(even) {background: #F3F3F3;}
    tr:nth-child(odd) {background: White;}
    .resource {background-color: RoyalBlue; color: White; font-weight: bold;}
    .blank {background-color: White; border: none;}
    .group {background-color: LightBlue;}
  </style>
</head>
<body>
  <h1>test1 in test-account-alias</h1>
  <table>
    <tr><th colspan="2">Event Details</th></tr>
    <tr><td>EventType</td><td>test1</td></tr>
    <tr><td>EventID</td><td>testID</td></tr>
    <tr><td>EventTime</td><td>testTime</td></tr>
    <tr><td>EventName</td><td>testName</td></tr>
    <tr><td>UserAgent</td><td>testAgent</td></tr>
    <tr><td>AWS Region</td><td>testRegion</td></tr>
    <tr><td>SourceIPAddress</td><td>testIP</td></tr>
  </table>
  &nbsp;
  <table>
    <tr><th colspan="2">UserIdentity</th></tr>
    <tr><td>Type</td><td>testType</td></tr>
    <tr><td>AccountID</td><td>testID</td></tr>
    <tr><td>UserName</td><td>testName</td></tr>
  </table>
</body>
`,
		},
	}
	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			body := tc.e.htmlBody()
			if tc.expected != body {
				t.Errorf("eventHandler() failed. Expected:\n%q\nGot:\n%q\n", tc.expected, body)
			}
		})
	}
}
