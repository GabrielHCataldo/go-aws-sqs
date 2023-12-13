package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go-aws-sqs/internal/util"
	"go-aws-sqs/sqs/option"
	"reflect"
	"strconv"
	"time"
)

type MessageReceived[Body, MessageAttributes any] struct {
	// A unique identifier for the message. An Id is considered unique across all
	// Amazon Web Services accounts for an extended period of time.
	Id string
	// An identifier associated with the act of receiving the message. A new receipt
	// handler is returned every time you receive a message. When deleting a message,
	// you provide the last received receipt handler to delete the message.
	ReceiptHandle string
	// The message's contents (not URL-encoded).
	Body                             Body
	ApproximateReceiveCount          *int
	ApproximateFirstReceiveTimestamp *time.Time
	MessageDeduplicationId           *string
	MessageGroupId                   *string
	SenderId                         *string
	SentTimestamp                    *time.Time
	SequenceNumber                   *int
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

type Context[Body, MessageAttributes any] struct {
	context.Context
	QueueUrl string
	Message  MessageReceived[Body, MessageAttributes]
}

type SimpleContext[Body any] Context[Body, map[string]types.MessageAttributeValue]

type HandlerConsumerFunc[Body, MessageAttributes any] func(ctx *Context[Body, MessageAttributes]) error
type HandlerSimpleConsumerFunc[Body any] func(ctx *SimpleContext[Body]) error

var ctxUnitTest context.Context

func ReceiveMessage[Body, MessageAttributes any](
	queueUrl string,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opts ...option.Consumer,
) {
	receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

func ReceiveMessageAsync[Body, MessageAttributes any](
	queueUrl string,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opts ...option.Consumer,
) {
	go receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

func SimpleReceiveMessage[Body any](
	queueUrl string,
	simpleHandle HandlerSimpleConsumerFunc[Body],
	opts ...option.Consumer,
) {
	handler := initHandleConsumerFunc(simpleHandle)
	receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

func SimpleReceiveMessageAsync[Body any](
	queueUrl string,
	simpleHandle HandlerSimpleConsumerFunc[Body],
	opts ...option.Consumer,
) {
	handler := initHandleConsumerFunc(simpleHandle)
	go receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

func receiveMessage[Body, MessageAttributes any](
	queueUrl string,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opt option.Consumer,
) {
	ctx := context.TODO()
	if ctxUnitTest != nil {
		ctx = ctxUnitTest
	}
	ctxClient, cancelCtxClient := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtxClient()
	client, err := getClient(ctxClient, false)
	if err != nil {
		panic(err)
	}
	input := prepareReceiveMessageInput(queueUrl, opt)
	attemptsReceiveMessages := 0
	printLogInitial(opt)
	for {
		if ctx.Err() != nil {
			break
		}
		output, err := client.ReceiveMessage(ctx, &input)
		if err != nil {
			handleError(&attemptsReceiveMessages, err, opt)
			continue
		}
		if len(output.Messages) == 0 {
			loggerInfo(opt.DebugMode, "No msg available to be processed, searching again in", opt.DelayQueryLoop.String())
			time.Sleep(opt.DelayQueryLoop)
			continue
		}
		loggerInfo(opt.DebugMode, "Start process received messages size:", len(output.Messages))
		signal := make(chan struct{}, 1)
		go processMessages[Body, MessageAttributes](queueUrl, *output, handler, opt, &signal)
		select {
		case <-signal:
			time.Sleep(opt.DelayQueryLoop)
		}
	}
}

func prepareReceiveMessageInput(queueUrl string, opt option.Consumer) sqs.ReceiveMessageInput {
	return sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		AttributeNames:      []types.QueueAttributeName{types.QueueAttributeNameAll},
		MaxNumberOfMessages: opt.MaxNumberOfMessages,
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		ReceiveRequestAttemptId: &opt.ReceiveRequestAttemptId,
		VisibilityTimeout:       util.ConvertDurationToInt32(opt.VisibilityTimeout),
		WaitTimeSeconds:         util.ConvertDurationToInt32(opt.WaitTimeSeconds),
	}
}

func printLogInitial(opt option.Consumer) {
	loggerInfo(opt.DebugMode, "Run start find messages with options:", opt)
}

func handleError(attemptsReceiveMessages *int, err error, opt option.Consumer) {
	*attemptsReceiveMessages++
	loggerErr(opt.DebugMode, "Receive message error:", err, " attempt:", attemptsReceiveMessages)
	if *attemptsReceiveMessages >= 3 {
		panic(fmt.Sprintln("Stop consumer: number of failed attempts exceeded 3 err:", err))
	}
	loggerInfo(opt.DebugMode, "Trying again in", opt.DelayQueryLoop.String())
	time.Sleep(opt.DelayQueryLoop)
}

func processMessages[Body, MessageAttributes any](
	queueUrl string,
	output sqs.ReceiveMessageOutput,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opt option.Consumer,
	signal *chan struct{},
) {
	ctx, cancel := context.WithTimeout(context.TODO(), opt.ConsumerMessageTimeout)
	defer cancel()
	var count int
	var mgsS, mgsF []string
	for _, message := range output.Messages {
		nCtx, err := prepareContextConsumer[Body, MessageAttributes](ctx, queueUrl, message)
		if err != nil {
			loggerErr(opt.DebugMode, "error prepare context to consumer:", err)
			return
		}
		err = handler(nCtx)
		appendMessagesByResult(nCtx.Message.Id, err, mgsS, mgsF)
		if opt.DeleteMessageProcessedSuccess {
			go deleteMessage(queueUrl, *message.ReceiptHandle, opt)
		}
		count++
	}
	loggerInfo(opt.DebugMode, "Finish process messages!", "processed:", count, "success:", mgsS, "failed:", mgsF)
	*signal <- struct{}{}
}

func prepareContextConsumer[Body, MessageAttributes any](
	ctx context.Context,
	queueUrl string,
	message types.Message,
) (*Context[Body, MessageAttributes], error) {
	ctxConsumer := &Context[Body, MessageAttributes]{
		Context:  ctx,
		QueueUrl: queueUrl,
	}
	messageReceived := MessageReceived[Body, MessageAttributes]{
		Id:                     *message.MessageId,
		ReceiptHandle:          *message.ReceiptHandle,
		MD5OfBody:              *message.MD5OfBody,
		MD5OfMessageAttributes: message.MD5OfMessageAttributes,
	}
	var body Body
	bodyString := *message.Body
	util.ParseStringToGeneric(bodyString, &body)
	if util.IsZeroReflect(reflect.ValueOf(body)) {
		return nil, ErrParseBody
	}
	messageReceived.Body = body
	if message.MessageAttributes != nil {
		var messagesAttributes MessageAttributes
		reflectMessageAtt := reflect.ValueOf(messagesAttributes)
		if reflectMessageAtt.Kind() == reflect.Map || reflectMessageAtt.Kind() == reflect.Struct {
			if util.IsMapMessageAttributeValues(messagesAttributes) {
				messageReceived.MessageAttributes = any(message.MessageAttributes).(MessageAttributes)
			} else {
				err := convertMessageAttributes[MessageAttributes](message.MessageAttributes, &messagesAttributes)
				if err != nil {
					return nil, err
				} else if !util.IsZeroReflect(reflect.ValueOf(messagesAttributes)) {
					messageReceived.MessageAttributes = messagesAttributes
				}
			}
		}
	}
	fillAttributes[Body, MessageAttributes](message, &messageReceived)
	ctxConsumer.Message = messageReceived
	return ctxConsumer, nil
}

func appendMessagesByResult(messageId string, err error, mgsS, mgsF []string) {
	if err != nil {
		mgsF = append(mgsF, messageId)
	} else {
		mgsS = append(mgsS, messageId)
	}
}

func convertMessageAttributes[T any](messageAttributes map[string]types.MessageAttributeValue, dest *T) error {
	m := map[string]any{}
	for k, v := range messageAttributes {
		var valueProcessed any
		valueString := v.StringValue
		if len(v.BinaryValue) > 0 {
			bs := string(v.BinaryValue)
			valueString = &bs
		}
		if valueString == nil || len(*valueString) == 0 {
			continue
		}
		util.ParseStringToGeneric(*valueString, &valueProcessed)
		if valueProcessed != nil {
			m[k] = valueProcessed
		}
	}
	if len(m) == 0 {
		return nil
	}
	return util.ParseMapToStruct[T](m, dest)
}

func fillAttributes[Body, MessageAttributes any](
	message types.Message,
	messageReceived *MessageReceived[Body, MessageAttributes],
) {
	for k, v := range message.Attributes {
		if len(v) == 0 {
			continue
		}
		switch k {
		case "ApproximateReceiveCount":
			i, err := strconv.Atoi(v)
			if err == nil && i > 0 {
				messageReceived.ApproximateReceiveCount = &i
			}
			break
		case "ApproximateFirstReceiveTimestamp":
			t, err := time.Parse(time.RFC3339, v)
			if err == nil {
				messageReceived.ApproximateFirstReceiveTimestamp = &t
			}
			break
		case "MessageDeduplicationId":
			messageReceived.MessageDeduplicationId = &v
			break
		case "MessageGroupId":
			messageReceived.MessageGroupId = &v
			break
		case "SenderId":
			messageReceived.SenderId = &v
			break
		case "SentTimestamp":
			t, err := time.Parse(time.RFC3339, v)
			if err == nil {
				messageReceived.SentTimestamp = &t
			}
			break
		case "SequenceNumber":
			i, err := strconv.Atoi(v)
			if err == nil && i > 0 {
				messageReceived.SequenceNumber = &i
			}
			break
		}
	}
}

func initHandleConsumerFunc[Body any, MessageAttributes map[string]types.MessageAttributeValue](
	simpleHandle HandlerSimpleConsumerFunc[Body],
) HandlerConsumerFunc[Body, MessageAttributes] {
	return func(ctx *Context[Body, MessageAttributes]) error {
		return simpleHandle(parseContextToSimpleContext[Body, MessageAttributes](ctx))
	}
}

func parseContextToSimpleContext[Body any, MessageAttributes map[string]types.MessageAttributeValue](
	ctx *Context[Body, MessageAttributes],
) *SimpleContext[Body] {
	return &SimpleContext[Body]{
		Context:  ctx.Context,
		QueueUrl: ctx.QueueUrl,
		Message: MessageReceived[Body, map[string]types.MessageAttributeValue]{
			Id:                               ctx.Message.Id,
			ReceiptHandle:                    ctx.Message.ReceiptHandle,
			Body:                             ctx.Message.Body,
			ApproximateReceiveCount:          ctx.Message.ApproximateReceiveCount,
			ApproximateFirstReceiveTimestamp: ctx.Message.ApproximateFirstReceiveTimestamp,
			MessageDeduplicationId:           ctx.Message.MessageDeduplicationId,
			MessageGroupId:                   ctx.Message.MessageGroupId,
			SenderId:                         ctx.Message.SenderId,
			SentTimestamp:                    ctx.Message.SentTimestamp,
			SequenceNumber:                   ctx.Message.SequenceNumber,
			MD5OfBody:                        ctx.Message.MD5OfBody,
			MD5OfMessageAttributes:           ctx.Message.MD5OfMessageAttributes,
			MessageAttributes:                ctx.Message.MessageAttributes,
		},
	}
}

func deleteMessage(queueUrl, receiptHandle string, opt option.Consumer) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := DeleteMessage(ctx, queueUrl, receiptHandle)
	if err != nil {
		loggerErr(opt.DebugMode, "error auto delete message processed success:", err)
	}
}
