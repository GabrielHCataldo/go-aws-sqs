package client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var sqsClient *sqs.Client

func GetClient(ctx context.Context) *sqs.Client {
	if sqsClient != nil {
		return sqsClient
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	sqsClient = sqs.NewFromConfig(cfg)
	return sqsClient
}
