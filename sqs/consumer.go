package sqs

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-aws-sqs-template/internal/client"
	"github.com/GabrielHCataldo/go-aws-sqs-template/internal/util"
	"github.com/GabrielHCataldo/go-aws-sqs-template/sqs/option"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"reflect"
	"time"
)

// MessageReceived represents a received message from Amazon Simple Queue Service (SQS).
// It contains information about the message, such as its ID, receipt handle, body, attributes,
type MessageReceived[Body, MessageAttributes any] struct {
	// A unique identifier for the message. An Id is considered unique across all
	// Amazon Web Services accounts for an extended period of time.
	Id string
	// An identifier associated with the act of receiving the message. A new receipt
	// handler is returned every time you receive a message. When deleting a message,
	// you provide the last received receipt handler to delete the message.
	ReceiptHandle string
	// Message body converted
	Body Body
	// Supported attributes:
	//   - ApproximateReceiveCount
	//   - ApproximateFirstReceiveTimestamp
	//   - MessageDeduplicationId
	//   - MessageGroupId
	//   - SenderId
	//   - SentTimestamp
	//   - SequenceNumber
	// ApproximateFirstReceiveTimestamp and SentTimestamp are each returned as an
	// integer representing the epoch time (http://en.wikipedia.org/wiki/Unix_time) in
	// milliseconds.
	Attributes Attributes
	// An MD5 digest of the non-URL-encoded message body string.
	MD5OfBody string
	// An MD5 digest of the non-URL-encoded message attribute string. You can use this
	// attribute to verify that Amazon SQS received the message correctly. Amazon SQS
	// URL-decodes the message before creating the MD5 digest. For information about
	// MD5, see RFC1321 (https://www.ietf.org/rfc/rfc1321.txt).
	MD5OfMessageAttributes *string
	// Converted message attributes
	MessageAttributes MessageAttributes
}

// Attributes represents the various attributes associated with a received message from Amazon Simple Queue Service (SQS).
// It includes information such as the approximate receipt count, the approximate timestamp of the first receipt,
// message deduplication ID (FIFO), message group ID (FIFO), sender ID, sent timestamp, and sequence number.
type Attributes struct {
	// Approximate receipt count
	ApproximateReceiveCount int
	// Approximate timestamp of first receipt
	ApproximateFirstReceiveTimestamp time.Time
	// Message deduplication ID (FIFO)
	MessageDeduplicationId string
	// Message group ID (FIFO)
	MessageGroupId string
	// sender ID
	SenderId string
	// Timestamp sent
	SentTimestamp time.Time
	// Sequence Number
	SequenceNumber int
}

// Context represents the execution context of a consumer handler in the package.
// It contains the necessary information for handling the message in the queue.
type Context[Body, MessageAttributes any] struct {
	// context used to process the message, may have a timeout if you pass the option.Consumer.ConsumerMessageTimeout
	// as a parameter
	context.Context `json:"-"`
	// message queue url
	QueueUrl string
	// converted message to process
	Message MessageReceived[Body, MessageAttributes]
}

// SimpleContext represents the execution context of a consumer handler in the package.
// It contains the necessary information for handling the message in the queue.
type SimpleContext[Body any] Context[Body, map[string]types.MessageAttributeValue]

// HandlerConsumerFunc is a function that consumes a message and returns an error if a failure occurs while processing the message.
type HandlerConsumerFunc[Body, MessageAttributes any] func(ctx *Context[Body, MessageAttributes]) error

// HandlerSimpleConsumerFunc is a function that consumes a message and returns an error if a failure occurs while processing the message.
type HandlerSimpleConsumerFunc[Body any] func(ctx *SimpleContext[Body]) error

type channelMessageProcessed struct {
	Err    error
	Signal *chan struct{}
}

var ctxInterrupt context.Context

// ReceiveMessage Works as a repeating job, when triggered, it will fetch messages from the indicated queue
// in the queueUrl parameter, if it does not find any messages it will reprocess the function looking for new messages again,
// otherwise, it will process these messages converting them into Context with the types of Body and MessageAttributes
// passed explicitly in the function, converted to the Context.Message.Body and Context.Message.MessageAttributes field,
// if you don't use MessageAttributes, you can use the SimpleReceiveMessage function. After conversion,
// we call the handler parameter so that this function processes message by message, in this handler we expect a return
// error, if no error occurred when processing the message, we have the auto-delete option
// (option.Consumer.DeleteMessageProcessedSuccess) if you have it enabled, we will remove the message from the file for you.
// finally, when processing all messages, we return to the initial flow looking for new messages.
//
// # Parameters
//
// - queueUrl: url of the queue where you want to fetch messages
// - handler: function to process the received message
// - opts: list of option.Consumer to customize the job
//
// # Panic
//
// If 3 errors occur when obtaining the message from the queue, it will trigger a panic informing the error returned
// from AWS SQS.
func ReceiveMessage[Body, MessageAttributes any](
	queueUrl string,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opts ...option.Consumer,
) {
	receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

// ReceiveMessageAsync Works like a repeating job, when triggered, it will fetch messages from the indicated queue
// in the queueUrl parameter, if it does not find any messages it will reprocess the function looking for new messages again,
// otherwise, it will process these messages converting them into Context with the types of Body and MessageAttributes
// passed explicitly in the function, converted to the Context.Message.Body and Context.Message.MessageAttributes field,
// if you don't use MessageAttributes, you can use the SimpleReceiveMessageAsync function. After conversion,
// we call the handler parameter so that this function processes message by message, in this handler we expect a return
// error, if no error occurred when processing the message, we have the auto-delete option
// (option.Consumer.DeleteMessageProcessedSuccess) where if enabled, we remove the message from the queue for you.
// finally, when processing all messages, we return to the initial flow looking for new messages.
//
// Unlike ReceiveMessage, this function will be processed asynchronously.
//
// # Parameters
//
// - queueUrl: url of the queue where you want to fetch messages
// - handler: function to process the received message
// - opts: list of options.Consumer to customize the job
//
// # Panic
//
// If 3 errors occur when obtaining the message from the queue, it will trigger a panic informing the error returned
// from AWS SQS.
func ReceiveMessageAsync[Body, MessageAttributes any](
	queueUrl string,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opts ...option.Consumer,
) {
	go receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

// SimpleReceiveMessage Works as a repeating job, when triggered, it will fetch messages from the indicated queue
// in the queueUrl parameter, if it does not find any messages it will reprocess the function looking for new messages
// again, otherwise, it will process these messages converting them to Context with the Body type passed explicitly
// in the function, converted to the Context.Message.Body field. After the conversion, we call the handler parameter
// so that this function processes message by message, in this handler we expect an error return, if none has occurred
// error processing the message, we have the auto-delete option (option.Consumer.DeleteMessageProcessedSuccess) in which case
// is enabled, we remove the message from the queue for you. Finally, after processing all messages, we return to the flow
// initial searching for new messages.
//
// # Parameters
//
// - queueUrl: url of the queue where you want to fetch messages
// - handler: function to process the received message
// - opts: list of options.Consumer to customize the job
//
// # Panic
//
// If 3 errors occur when obtaining the message from the queue, it will trigger a panic informing the error returned
// from AWS SQS.
func SimpleReceiveMessage[Body any](
	queueUrl string,
	simpleHandle HandlerSimpleConsumerFunc[Body],
	opts ...option.Consumer,
) {
	handler := initHandleConsumerFunc(simpleHandle)
	receiveMessage(queueUrl, handler, option.GetConsumerByParams(opts))
}

// SimpleReceiveMessageAsync Works like a repeating job, when triggered, it will fetch messages from the indicated queue
// in the queueUrl parameter, if it does not find any messages it will reprocess the function looking for new messages again,
// otherwise, it will process these messages converting them to Context with the Body type passed explicitly
// in the function, converted to the Context.Message.Body field. After the conversion, we call the handler parameter so
// that this function processes message by message, in this handler we expect an error return, if none has occurred
// error processing the message, we have the auto-delete option (option.Consumer.DeleteMessageProcessedSuccess) in which case
// is enabled, we remove the message from the queue for you. Finally, after processing all messages, we return to the flow
// initial searching for new messages.
//
// Unlike SimpleReceiveMessage, this function will be processed asynchronously.
//
// # Parameters
//
// - queueUrl: url of the queue where you want to fetch messages
// - handler: function to process the received message
// - opts: list of options.Consumer to customize the job
//
// # Panic
//
// If 3 errors occur when obtaining the message from the queue, it will trigger a panic informing the error returned
// from AWS SQS.
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
	if ctxInterrupt != nil {
		ctx = ctxInterrupt
	}
	ctxClient, cancelCtxClient := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtxClient()
	sqsClient := client.GetClient(ctxClient)
	input := prepareReceiveMessageInput(queueUrl, opt)
	attemptsReceiveMessages := 0
	printLogInitial(opt)
	for {
		if ctx.Err() != nil {
			break
		}
		output, err := sqsClient.ReceiveMessage(ctx, &input)
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
		processMessages[Body, MessageAttributes](queueUrl, output, handler, opt)
		time.Sleep(opt.DelayQueryLoop)
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
		ReceiveRequestAttemptId: opt.ReceiveRequestAttemptId,
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
	output *sqs.ReceiveMessageOutput,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opt option.Consumer,
) {
	var count int
	var mgsS, mgsF []string
	for _, message := range output.Messages {
		mgsS, mgsF = processMessage(queueUrl, handler, message, opt)
		count++
	}
	loggerInfo(opt.DebugMode, "Finish process messages!", "processed:", count, "success:", mgsS, "failed:", mgsF)
}

func processMessage[Body, MessageAttributes any](
	queueUrl string,
	handler HandlerConsumerFunc[Body, MessageAttributes],
	message types.Message,
	opt option.Consumer,
) (mgsS, mgsF []string) {
	ctx, cancel := context.WithTimeout(context.TODO(), opt.ConsumerMessageTimeout)
	defer cancel()
	ctxConsumer, err := prepareContextConsumer[Body, MessageAttributes](ctx, queueUrl, message)
	if err != nil {
		loggerErr(opt.DebugMode, "error prepare context to consumer:", err)
		return
	}
	signal := make(chan struct{}, 1)
	channel := channelMessageProcessed{
		Signal: &signal,
	}
	go processHandler(ctxConsumer, handler, opt, &channel)
	select {
	case <-ctx.Done():
		appendMessagesByResult(ctxConsumer.Message.Id, ctx.Err(), &mgsS, &mgsF)
		break
	case <-*channel.Signal:
		appendMessagesByResult(*message.MessageId, channel.Err, &mgsS, &mgsF)
		break
	}
	return mgsS, mgsF
}

func processHandler[Body, MessageAttributes any](
	ctx *Context[Body, MessageAttributes],
	handler HandlerConsumerFunc[Body, MessageAttributes],
	opt option.Consumer,
	channel *channelMessageProcessed,
) {
	err := handler(ctx)
	if ctx.Err() != nil {
		return
	}
	if err == nil && opt.DeleteMessageProcessedSuccess {
		go deleteMessage(ctx.QueueUrl, ctx.Message.ReceiptHandle)
	}
	channel.Err = err
	*channel.Signal <- struct{}{}
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
				convertMessageAttributes[MessageAttributes](message.MessageAttributes, &messagesAttributes)
				messageReceived.MessageAttributes = messagesAttributes
			}
		}
	}
	fillAttributes[Body, MessageAttributes](message, &messageReceived)
	ctxConsumer.Message = messageReceived
	return ctxConsumer, nil
}

func appendMessagesByResult(messageId string, err error, mgsS, mgsF *[]string) {
	if err != nil {
		*mgsF = append(*mgsF, messageId)
	} else {
		*mgsS = append(*mgsS, messageId)
	}
}

func convertMessageAttributes[T any](messageAttributes map[string]types.MessageAttributeValue, dest *T) {
	m := map[string]any{}
	for k, v := range messageAttributes {
		var valueProcessed any
		valueString := v.StringValue
		util.ParseStringToGeneric(*valueString, &valueProcessed)
		if valueProcessed != nil {
			m[k] = valueProcessed
		}
	}
	_ = util.ParseMapToStruct[T](m, dest)
}

func fillAttributes[Body, MessageAttributes any](
	message types.Message,
	messageReceived *MessageReceived[Body, MessageAttributes],
) {
	m := map[string]any{}
	for k, v := range message.Attributes {
		var valueProcessed any
		util.ParseStringToGeneric(v, &valueProcessed)
		if valueProcessed != nil {
			m[k] = valueProcessed
		}
	}
	_ = util.ParseMapToStruct[Attributes](m, &messageReceived.Attributes)
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
			Id:                     ctx.Message.Id,
			ReceiptHandle:          ctx.Message.ReceiptHandle,
			Body:                   ctx.Message.Body,
			Attributes:             ctx.Message.Attributes,
			MD5OfBody:              ctx.Message.MD5OfBody,
			MD5OfMessageAttributes: ctx.Message.MD5OfMessageAttributes,
			MessageAttributes:      ctx.Message.MessageAttributes,
		},
	}
}

func deleteMessage(queueUrl, receiptHandle string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = DeleteMessage(ctx, queueUrl, receiptHandle)
}
