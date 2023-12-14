package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go-aws-sqs/internal/client"
	"go-aws-sqs/sqs/option"
)

type DeleteMessageBatchInput struct {
	// The URL of the Amazon SQS queue from which messages are deleted. Queue URLs and
	// names are case-sensitive.
	//
	// This member is required.
	QueueUrl string
	// Lists the receipt handles for the messages to be deleted.
	//
	// This member is required.
	Entries []DeleteMessageBatchRequestEntry
}

type DeleteMessageBatchRequestEntry struct {
	// The identifier for this particular receipt handler. This is used to communicate
	// the result. The Id s of a batch request need to be unique within a request. This
	// identifier can have up to 80 characters. The following characters are accepted:
	// alphanumeric characters, hyphens(-), and underscores (_).
	//
	// This member is required.
	Id string
	// A receipt handler.
	//
	// This member is required.
	ReceiptHandle string
}

type ChangeMessageVisibilityInput struct {
	// The URL of the Amazon SQS queue whose message's visibility is changed. Queue
	// URLs and names are case-sensitive.
	//
	// This member is required.
	QueueUrl string
	// The receipt handler associated with the message, whose visibility timeout is
	// changed. This parameter is returned by the ReceiveMessage action.
	//
	// This member is required.
	ReceiptHandle *string
	// The new value for the message's visibility timeout (in seconds). Values range: 0
	// to 43200 . Maximum: 12 hours.
	//
	// This member is required.
	VisibilityTimeout int32
}

type ChangeMessageVisibilityBatchInput struct {
	// The URL of the Amazon SQS queue whose messages' visibility is changed. Queue
	// URLs and names are case-sensitive.
	//
	// This member is required.
	QueueUrl string
	// Lists the receipt handles of the messages for which the visibility timeout must
	// be changed.
	//
	// This member is required.
	Entries []ChangeMessageVisibilityBatchRequestEntry
}

type ChangeMessageVisibilityBatchRequestEntry struct {
	// An identifier for this particular receipt handler used to communicate the
	// result. The Id s of a batch request need to be unique within a request. This
	// identifier can have up to 80 characters. The following characters are accepted:
	// alphanumeric characters, hyphens(-), and underscores (_).
	//
	// This member is required.
	Id string
	// A receipt handler.
	//
	// This member is required.
	ReceiptHandle string
	// The new value (in seconds) for the message's visibility timeout.
	VisibilityTimeout int32
}

type StartMessageMoveTaskInput struct {
	// The ARN of the queue that contains the messages to be moved to another queue.
	// Currently, only ARNs of dead-letter queues (DLQs) whose sources are other Amazon
	// SQS queues are accepted. DLQs whose sources are non-SQS queues, such as Lambda
	// or Amazon SNS topics, are not currently supported.
	//
	// This member is required.
	SourceArn string
	// The ARN of the queue that receives the moved messages. You can use this field
	// to specify the destination queue where you would like to redrive messages. If
	// this field is left blank, the messages will be redriven back to their respective
	// original source queues.
	DestinationArn *string
	// The number of messages to be moved per second (the message movement rate). You
	// can use this field to define a fixed message movement rate. The maximum value
	// for messages per second is 500. If this field is left blank, the system will
	// optimize the rate based on the queue message backlog size, which may vary
	// throughout the duration of the message movement task.
	MaxNumberOfMessagesPerSecond *int32
}

// DeleteMessage Deletes the specified message from the specified queue. To select the message
// to delete, use the ReceiptHandle of the message (not the MessageId which you
// receive when you send the message). Amazon SQS can delete a message from a queue
// even if a visibility timeout setting causes the message to be locked by another
// consumer. Amazon SQS automatically deletes messages left in a queue longer than
// the retention period configured for the queue. The ReceiptHandle is associated
// with a specific instance of receiving a message. If you receive a message more
// than once, the ReceiptHandle is different each time you receive a message. When
// you use the DeleteMessage action, you must provide the most recently received
// ReceiptHandle for the message (otherwise, the request succeeds, but the message
// will not be deleted). For standard queues, it is possible to receive a message
// even after you delete it. This might happen on rare occasions if one of the
// servers which stores a copy of the message is unavailable when you send the
// request to delete the message. The copy remains on the server and might be
// returned to you during a subsequent receive request. You should ensure that your
// application is idempotent, so that receiving a message more than once does not
// cause issues.
func DeleteMessage(ctx context.Context, queueUrl, receiptHandle string, opts ...option.Default) (
	*sqs.DeleteMessageOutput, error) {
	opt := option.GetDefaultByParams(opts)
	loggerInfo(opt.DebugMode, "deleting message..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: &receiptHandle,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error delete message:", err)
	} else {
		loggerInfo(opt.DebugMode, "message deleted successfully:", output)
	}
	return output, err
}

// DeleteMessageBatch Deletes up to ten messages from the specified queue. This is a batch version of
// DeleteMessage . The result of the action on each message is reported
// individually in the response. Because the batch request can result in a
// combination of successful and unsuccessful actions, you should check for batch
// errors even when the call returns an HTTP status code of 200 .
func DeleteMessageBatch(ctx context.Context, input DeleteMessageBatchInput, opts ...option.Default) (
	*sqs.DeleteMessageBatchOutput, error) {
	opt := option.GetDefaultByParams(opts)
	loggerInfo(opt.DebugMode, "deleting messages batch..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		Entries:  prepareEntriesDeleteMessageBatch(input.Entries),
		QueueUrl: &input.QueueUrl,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error delete messages batch:", err)
	} else {
		loggerInfo(opt.DebugMode, "messages deleted successfully:", output)
	}
	return output, err
}

// ChangeMessageVisibility Changes the visibility timeout of a specified message in a queue to a new
// value. The default visibility timeout for a message is 30 seconds. The minimum
// is 0 seconds. The maximum is 12 hours. For more information, see Visibility
// Timeout (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html)
// in the Amazon SQS Developer Guide. For example, if the default timeout for a
// queue is 60 seconds, 15 seconds have elapsed since you received the message, and
// you send a ChangeMessageVisibility call with VisibilityTimeout set to 10
// seconds, the 10 seconds begin to count from the time that you make the
// ChangeMessageVisibility call. Thus, any attempt to change the visibility timeout
// or to delete that message 10 seconds after you initially change the visibility
// timeout (a total of 25 seconds) might result in an error. An Amazon SQS message
// has three basic states:
//   - Sent to a queue by a producer.
//   - Received from the queue by a consumer.
//   - Deleted from the queue.
//
// A message is considered to be stored after it is sent to a queue by a producer,
// but not yet received from the queue by a consumer (that is, between states 1 and
// 2). There is no limit to the number of stored messages. A message is considered
// to be in flight after it is received from a queue by a consumer, but not yet
// deleted from the queue (that is, between states 2 and 3). There is a limit to
// the number of in flight messages. Limits that apply to in flight messages are
// unrelated to the unlimited number of stored messages. For most standard queues
// (depending on queue traffic and message backlog), there can be a maximum of
// approximately 120,000 in flight messages (received from a queue by a consumer,
// but not yet deleted from the queue). If you reach this limit, Amazon SQS returns
// the OverLimit error message. To avoid reaching the limit, you should delete
// messages from the queue after they're processed. You can also increase the
// number of queues you use to process your messages. To request a limit increase,
// file a support request (https://console.aws.amazon.com/support/home#/case/create?issueType=service-limit-increase&limitType=service-code-sqs)
// . For FIFO queues, there can be a maximum of 20,000 in flight messages (received
// from a queue by a consumer, but not yet deleted from the queue). If you reach
// this limit, Amazon SQS returns no error messages. If you attempt to set the
// VisibilityTimeout to a value greater than the maximum time left, Amazon SQS
// returns an error. Amazon SQS doesn't automatically recalculate and increase the
// timeout to the maximum remaining time. Unlike with a queue, when you change the
// visibility timeout for a specific message the timeout value is applied
// immediately but isn't saved in memory for that message. If you don't delete a
// message after it is received, the visibility timeout for the message reverts to
// the original timeout value (not to the value you set using the
// ChangeMessageVisibility action) the next time the message is received.
func ChangeMessageVisibility(ctx context.Context, input ChangeMessageVisibilityInput, opts ...option.Default) (
	*sqs.ChangeMessageVisibilityOutput, error) {
	opt := option.GetDefaultByParams(opts)
	loggerInfo(opt.DebugMode, "changing messages visibility..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          input.ReceiptHandle,
		ReceiptHandle:     input.ReceiptHandle,
		VisibilityTimeout: input.VisibilityTimeout,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error charge message visibility:", err)
	} else {
		loggerInfo(opt.DebugMode, "message charge visibility successfully:", output)
	}
	return output, err
}

// ChangeMessageVisibilityBatch Changes the visibility timeout of multiple messages. This is a batch version of
// ChangeMessageVisibility . The result of the action on each message is reported
// individually in the response. You can send up to 10 ChangeMessageVisibility
// requests with each ChangeMessageVisibilityBatch action. Because the batch
// request can result in a combination of successful and unsuccessful actions, you
// should check for batch errors even when the call returns an HTTP status code of
// 200 .
func ChangeMessageVisibilityBatch(ctx context.Context, input ChangeMessageVisibilityBatchInput, opts ...option.Default) (
	*sqs.ChangeMessageVisibilityBatchOutput, error) {
	opt := option.GetDefaultByParams(opts)
	loggerInfo(opt.DebugMode, "changing message visibility batch..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.ChangeMessageVisibilityBatch(ctx, &sqs.ChangeMessageVisibilityBatchInput{
		Entries:  prepareEntriesChangeMessageVisibilityBatch(input.Entries),
		QueueUrl: &input.QueueUrl,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error charge message visibility batch:", err)
	} else {
		loggerInfo(opt.DebugMode, "messages charge visibility successfully:", output)
	}
	return output, err
}

// StartMessageMoveTask Starts an asynchronous task to move messages from a specified source queue to a
// specified destination queue.
//   - This action is currently limited to supporting message redrive from queues
//     that are configured as dead-letter queues (DLQs) (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
//     of other Amazon SQS queues only. Non-SQS queue sources of dead-letter queues,
//     such as Lambda or Amazon SNS topics, are currently not supported.
//   - In dead-letter queues redrive context, the StartMessageMoveTask the source
//     queue is the DLQ, while the destination queue can be the original source queue
//     (from which the messages were driven to the dead-letter-queue), or a custom
//     destination queue.
//   - Currently, only standard queues support redrive. FIFO queues don't support
//     redrive.
//   - Only one active message movement task is supported per queue at any given
//     time.
func StartMessageMoveTask(ctx context.Context, input StartMessageMoveTaskInput, opts ...option.Default) (
	*sqs.StartMessageMoveTaskOutput, error) {
	opt := option.GetDefaultByParams(opts)
	loggerInfo(opt.DebugMode, "starting message move task..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.StartMessageMoveTask(ctx, &sqs.StartMessageMoveTaskInput{
		SourceArn:                    &input.SourceArn,
		DestinationArn:               input.DestinationArn,
		MaxNumberOfMessagesPerSecond: input.MaxNumberOfMessagesPerSecond,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error start message move task:", err)
	} else {
		loggerInfo(opt.DebugMode, "start message move task successfully:", output)
	}
	return output, err
}

// CancelMessageMoveTask Cancels a specified message movement task. A message movement can only be
// cancelled when the current status is RUNNING. Cancelling a message movement task
// does not revert the messages that have already been moved. It can only stop the
// messages that have not been moved yet.
//   - This action is currently limited to supporting message redrive from
//     dead-letter queues (DLQs) (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
//     only. In this context, the source queue is the dead-letter queue (DLQ), while
//     the destination queue can be the original source queue (from which the messages
//     were driven to the dead-letter-queue), or a custom destination queue.
//   - Currently, only standard queues are supported.
//   - Only one active message movement task is supported per queue at any given
//     time.
func CancelMessageMoveTask(ctx context.Context, taskHandle string, opts ...option.Default) (
	*sqs.CancelMessageMoveTaskOutput, error) {
	opt := option.GetDefaultByParams(opts)
	loggerInfo(opt.DebugMode, "canceling message move task..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.CancelMessageMoveTask(ctx, &sqs.CancelMessageMoveTaskInput{
		TaskHandle: &taskHandle,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error cancel message move task:", err)
	} else {
		loggerInfo(opt.DebugMode, "cancel message move task successfully:", output)
	}
	return output, err
}

// ListMessageMoveTasks Gets the most recent message movement tasks (up to 10) under a specific source
// queue.
//   - This action is currently limited to supporting message redrive from
//     dead-letter queues (DLQs) (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
//     only. In this context, the source queue is the dead-letter queue (DLQ), while
//     the destination queue can be the original source queue (from which the messages
//     were driven to the dead-letter-queue), or a custom destination queue.
//   - Currently, only standard queues are supported.
//   - Only one active message movement task is supported per queue at any given
//     time.
func ListMessageMoveTasks(ctx context.Context, sourceArn string, opts ...option.ListMessageMoveTasks) (
	*sqs.ListMessageMoveTasksOutput, error) {
	opt := option.GetListMessageMoveTaskByParams(opts)
	loggerInfo(opt.DebugMode, "listing message move tasks..")
	sqsClient := client.GetClient(ctx)
	output, err := sqsClient.ListMessageMoveTasks(ctx, &sqs.ListMessageMoveTasksInput{
		SourceArn:  &sourceArn,
		MaxResults: &opt.MaxResults,
	}, option.FuncByOptionHttp(opt.OptionHttp))
	if err != nil {
		loggerErr(opt.DebugMode, "error list message move tasks:", err)
	} else {
		loggerInfo(opt.DebugMode, "list message move tasks successfully:", output)
	}
	return output, err
}

func prepareEntriesDeleteMessageBatch(entries []DeleteMessageBatchRequestEntry) []types.DeleteMessageBatchRequestEntry {
	var result []types.DeleteMessageBatchRequestEntry
	for _, v := range entries {
		if len(v.Id) == 0 && len(v.ReceiptHandle) == 0 {
			continue
		}
		rEntry := types.DeleteMessageBatchRequestEntry{
			Id:            nil,
			ReceiptHandle: nil,
		}
		if len(v.Id) > 0 {
			rEntry.Id = &v.Id
		}
		if len(v.ReceiptHandle) > 0 {
			rEntry.ReceiptHandle = &v.ReceiptHandle
		}
		result = append(result, rEntry)
	}
	return result
}

func prepareEntriesChangeMessageVisibilityBatch(
	entries []ChangeMessageVisibilityBatchRequestEntry,
) []types.ChangeMessageVisibilityBatchRequestEntry {
	var result []types.ChangeMessageVisibilityBatchRequestEntry
	for _, v := range entries {
		if len(v.Id) == 0 && len(v.ReceiptHandle) == 0 {
			continue
		}
		rEntry := types.ChangeMessageVisibilityBatchRequestEntry{
			Id:            nil,
			ReceiptHandle: nil,
		}
		if len(v.Id) > 0 {
			rEntry.Id = &v.Id
		}
		if len(v.ReceiptHandle) > 0 {
			rEntry.ReceiptHandle = &v.ReceiptHandle
		}
		rEntry.VisibilityTimeout = v.VisibilityTimeout
		result = append(result, rEntry)
	}
	return result
}
