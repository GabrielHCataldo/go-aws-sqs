package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go-aws-sqs/sqs/option"
	"time"
)

type MessageReceived[Body, MessageAttributes any] struct {
	// A unique identifier for the message. A MessageId is considered unique across all
	// Amazon Web Services accounts for an extended period of time.
	MessageId string
	// An identifier associated with the act of receiving the message. A new receipt
	// handle is returned every time you receive a message. When deleting a message,
	// you provide the last received receipt handle to delete the message.
	ReceiptHandle string
	// The message's contents (not URL-encoded).
	Body                             Body
	ApproximateReceiveCount          int
	ApproximateFirstReceiveTimestamp time.Time
	MessageDeduplicationId           *string
	MessageGroupId                   *string
	SenderId                         *string
	SentTimestamp                    time.Time
	SequenceNumber                   int
	// An MD5 digest of the non-URL-encoded message body string.
	MD5OfBody string
	// An MD5 digest of the non-URL-encoded message attribute string. You can use this
	// attribute to verify that Amazon SQS received the message correctly. Amazon SQS
	// URL-decodes the message before creating the MD5 digest. For information about
	// MD5, see RFC1321 (https://www.ietf.org/rfc/rfc1321.txt).
	MD5OfMessageAttributes *string
	// Each message attribute consists of a Name, Type, and Value. For more
	// information, see Amazon SQS message attributes (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-metadata.html#sqs-message-attributes)
	// in the Amazon SQS Developer Guide.
	MessageAttributes MessageAttributes
}

type ContextConsumer[Body, MessageAttributes any] struct {
	context.Context
	QueueUrl string
	Message  MessageReceived[Body, MessageAttributes]
}

type HandleConsumerFunc[Body, MessageAttributes any] func(ctx ContextConsumer[Body, MessageAttributes]) error

func ReceiveMessage[Body, MessageAttributes any](
	queueUrl string,
	handle HandleConsumerFunc[Body, MessageAttributes],
	opts ...*option.OptionsConsumer,
) {
}

func ReceiveMessageAsync[Body, MessageAttributes any](
	queueUrl string,
	handle HandleConsumerFunc[Body, MessageAttributes],
	opts ...*option.OptionsConsumer,
) {
}

func receiveMessage[Body, MessageAttributes any](
	queueUrl string,
	handle HandleConsumerFunc[Body, MessageAttributes],
	opts ...*option.OptionsConsumer,
) error {
	ctxClient, cancelCtxClient := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancelCtxClient()
	client, err := getSqsClient(ctxClient)
	if err != nil {
		return err
	}
	attemptsReceiveMessages := 0
	for {
		input := &sqs.ReceiveMessageInput{
			MessageAttributeNames: []string{
				string(types.QueueAttributeNameAll),
			},
			AttributeNames:      []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:            &queueUrl,
			MaxNumberOfMessages: getMaxNumberOfMessagesByOptsConsumer(opts...),
			VisibilityTimeout:   getVisibilityTimeoutByOptsConsumer(opts...),
		}
		consumerMessageTimeout := getConsumerMessageTimeoutByOptsConsumer(opts...)
		delayRerun := getDelayRerunQueryByOptsConsumer(opts...)
		debugMode := getDebugModeByOptsConsumer(opts...)
		loggerInfo(debugMode, "Run start find receive message MaxNumberOfMessages:", input.MaxNumberOfMessages,
			"VisibilityTimeout:", input.VisibilityTimeout)
		receiveMessageOutput, err := client.ReceiveMessage(context.Background(), input)
		if err != nil {
			attemptsReceiveMessages++
			loggerErr(debugMode, "Receive message error:", err, " attempt:", attemptsReceiveMessages)
			if attemptsReceiveMessages >= 3 {
				loggerInfo(debugMode, "Stop consumer: number of failed attempts exceeded 3")
				return nil
			}
			loggerInfo(debugMode, "Trying again in", delayRerun.String())
			time.Sleep(delayRerun)
			continue
		}
		if len(receiveMessageOutput.Messages) == 0 {
			loggerInfo(debugMode, "No messages available to be processed at this time, searching again in", delayRerun.String())
			time.Sleep(delayRerun)
			continue
		}
		loggerInfo(debugMode, "Start process received messages size:", len(receiveMessageOutput.Messages))
		ctxConsumer, cancelConsumer := context.WithTimeout(context.TODO(), consumerMessageTimeout)
		finish := make(chan struct{}, 1)
		var countMsgProcessed int
		var countMsgSuccess []string
		var countMsgFailed []string
		go func() {
			defer cancelConsumer()
			for _, message := range receiveMessageOutput.Messages {
				// todo: implement me
				//	err = handle()
				//	if err != nil {
				//		countMsgFailed = append(countMsgFailed, *message.MessageId)
				//	} else {
				//		countMsgSuccess = append(countMsgSuccess, *message.MessageId)
				//	}
				countMsgProcessed++
			}
			finish <- struct{}{}
		}()
		select {
		case <-ctxConsumer.Done():
			loggerErr(debugMode, "Timeout context processed:", countMsgProcessed, "success:", countMsgSuccess,
				"failed:", countMsgFailed, "timeout error:", ctxConsumer.Err())
			time.Sleep(delayRerun)
			continue
		case <-finish:
			loggerInfo(debugMode, "Finish process messages!", "processed:", countMsgProcessed,
				"success:", countMsgSuccess, "failed:", countMsgFailed)
			time.Sleep(delayRerun)
			continue
		}
	}
}
