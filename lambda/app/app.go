package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	logs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type logConfig struct {
	LastEventID string `json:"last_event_id"`
}

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		return fmt.Errorf("failed to connect to AWS: %v", err)
	}
	c, err := getConfig(sess, os.Getenv("SECRET_ID"))
	if err != nil {
		return err
	}

	return nil
}

func getEventsSinceID(cfg client.ConfigProvider, eventID string) error {
	svc := logs.New(cfg)
	svc.FilterLogEventsPages(&logs.FilterLogEventsInput{})
}

func getConfig(cfg client.ConfigProvider, secretID string) (*logConfig, error) {
	svc := secretsmanager.New(cfg)
	out, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret value for: %s -> %v", secretID, err)
	}

	var c logConfig
	err = json.Unmarshal([]byte(aws.StringValue(out.SecretString)), &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal log config: %v", err)
	}

	return &c, nil
}
