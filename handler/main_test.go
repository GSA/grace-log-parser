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
		/* Need to be able to mock SES service
		"matching event": {
			logsEvent: events.CloudwatchLogsEvent{
				AWSLogs: events.CloudwatchLogsRawData{
					Data: "H4sIAGUGaF4AA6tWUMrJT3ctS80rKVayUoiuVspNLS5OTE8FcpSqFWKUUkFSfom5qTFAkRgl5/y84vycVJ/89My8GKVapdpYhVouAJbAiGdFAAAA",
				},
			},
			expectedErr: "",
		},
		*/
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
