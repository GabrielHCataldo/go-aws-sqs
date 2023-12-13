package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var sqsClient *sqs.Client

func getClient(ctx context.Context, debugMode bool) (*sqs.Client, error) {
	if sqsClient != nil {
		return sqsClient, nil
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	sqsClient = sqs.NewFromConfig(cfg)
	return sqsClient, nil
}
