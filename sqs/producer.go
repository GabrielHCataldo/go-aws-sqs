package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go-aws-sqs/internal/util"
	"go-aws-sqs/sqs/option"
)

// SendMessage sends a message to an SQS queue using the provided context, queue URL, message content, and options.
//
// The function first retrieves an SQS client using the getSqsClient function. It then prepares the message input using
// the prepareMessageInput function. The message input includes the message body, queue URL, delay seconds, message
// attributes, message deduplication ID, message group ID, and message system attributes. The message input is passed to
// the client's SendMessage method, which sends the message to the queue. The function returns the SendMessage output or
// an error if one occurs.
//
// Example usage:
//
//	output, err := SendMessage(ctx, queueUrl, v, opts...)
//
// Parameters:
// - ctx (Context): The context of the request.
// - queueUrl (string): The URL of the SQS queue.
// - v (any): The content of the message.
// - opts (OptionsProducer): Optional options to customize the message. (see OptionsProducer declaration for available options)
//
// Returns:
// - *sqs.SendMessageOutput: The result of the SendMessage operation.
// - error: An error if one occurs during the SendMessage operation.
func SendMessage(ctx context.Context, queueUrl string, v any, opts ...*option.OptionsProducer) (
	*sqs.SendMessageOutput, error) {
	debugMode := getDebugModeByOptsProducer(opts...)
	loggerInfo(debugMode, "getting client sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	loggerInfo(debugMode, "preparing message input..")
	input, err := prepareMessageInput(queueUrl, v, opts...)
	if err != nil {
		loggerErr(debugMode, "error prepare message input:", err)
		return nil, err
	}
	loggerInfo(debugMode, "sending message..")
	output, err := client.SendMessage(ctx, input, optionsHttp(getOptionsHttpByOptsProducer(opts...)))
	if err != nil {
		loggerErr(debugMode, "error send message:", err)
	} else {
		loggerInfo(debugMode, "message sent successfully:", output)
	}
	return output, err
}

// SendMessageAsync sends a message asynchronously to an SQS queue using the provided context, queue URL,
// message content, and options.
//
// The function spawns a goroutine and calls the sendMessageAsync function internally to send the message.
// It first checks if the debug mode is enabled by checking the options. It then calls the SendMessage function,
// which sends the message to the queue. If an error occurs and debug mode is enabled, it logs the error.
// If debug mode is enabled, it logs the successful message delivery.
//
// Parameters:
// - ctx (Context): The context of the request.
// - queueUrl (string): The URL of the SQS queue.
// - v (any): The content of the message.
// - opts (OptionsProducer): Optional options to customize the message. (see OptionsProducer declaration for available options)
//
// Example usage:
// SendMessageAsync(ctx, queueUrl, v, opts...)
func SendMessageAsync(ctx context.Context, queueUrl string, v any, opts ...*option.OptionsProducer) {
	go sendMessageAsync(ctx, queueUrl, v, opts...)
}

func sendMessageAsync(ctx context.Context, queueUrl string, v any, opts ...*option.OptionsProducer) {
	_, _ = SendMessage(ctx, queueUrl, v, opts...)
}

func prepareMessageInput(queueUrl string, v any, opts ...*option.OptionsProducer) (*sqs.SendMessageInput, error) {
	body, err := util.ConvertToString(v)
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		return nil, ErrMessageEmpty
	}
	messageAttByOpt, err := getMessageAttByOptsProducer(opts...)
	if err != nil {
		return nil, err
	}
	messageSystemAttByOpt, err := getMessageSystemAttByOptsProducer(opts...)
	if err != nil {
		return nil, err
	}
	delaySeconds := getDelaySecondsByOptsProducer(opts...)
	messageDeduplicationId := getMessageDeduplicationIdByOptsProducer(opts...)
	messageGroupIdByOpts := getMessageGroupIdByOptsProducer(opts...)
	return &sqs.SendMessageInput{
		MessageBody:             aws.String(body),
		QueueUrl:                &queueUrl,
		DelaySeconds:            delaySeconds,
		MessageAttributes:       messageAttByOpt,
		MessageDeduplicationId:  messageDeduplicationId,
		MessageGroupId:          messageGroupIdByOpts,
		MessageSystemAttributes: messageSystemAttByOpt,
	}, nil
}
