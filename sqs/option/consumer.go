package option

import "time"

type OptionsConsumer struct {
	baseOptions
	// default: false
	DeleteMessageProcessedSuccess bool
	// default: 10 seconds
	ConsumerMessageTimeout time.Duration
	// default: 15 seconds
	DelayRerunQuery time.Duration
	// The maximum number of messages to return. Amazon SQS never returns more
	// messages than this value (however, fewer messages might be returned). Valid
	// values: 1 to 10.
	//
	// default: 10
	MaxNumberOfMessages int32
	// This parameter applies only to FIFO (first-in-first-out) queues. The token used
	// for deduplication of ReceiveMessage calls. If a networking issue occurs after a
	// ReceiveMessage action, and instead of a response you receive a generic error, it
	// is possible to retry the same action with an identical ReceiveRequestAttemptId
	// to retrieve the same set of messages, even if their visibility timeout has not
	// yet expired.
	//   - You can use ReceiveRequestAttemptId only for 5 minutes after a
	//   ReceiveMessage action.
	//   - When you set FifoQueue , a caller of the ReceiveMessage action can provide a
	//   ReceiveRequestAttemptId explicitly.
	//   - If a caller of the ReceiveMessage action doesn't provide a
	//   ReceiveRequestAttemptId , Amazon SQS generates a ReceiveRequestAttemptId .
	//   - It is possible to retry the ReceiveMessage action with the same
	//   ReceiveRequestAttemptId if none of the messages have been modified (deleted or
	//   had their visibility changes).
	//   - During a visibility timeout, subsequent calls with the same
	//   ReceiveRequestAttemptId return the same messages and receipt handles. If a
	//   retry occurs within the deduplication interval, it resets the visibility
	//   timeout. For more information, see Visibility Timeout (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html)
	//   in the Amazon SQS Developer Guide. If a caller of the ReceiveMessage action
	//   still processes messages when the visibility timeout expires and messages become
	//   visible, another worker consuming from the same queue can receive the same
	//   messages and therefore process duplicates. Also, if a consumer whose message
	//   processing time is longer than the visibility timeout tries to delete the
	//   processed messages, the action fails with an error. To mitigate this effect,
	//   ensure that your application observes a safe threshold before the visibility
	//   timeout expires and extend the visibility timeout as necessary.
	//   - While messages with a particular MessageGroupId are invisible, no more
	//   messages belonging to the same MessageGroupId are returned until the
	//   visibility timeout expires. You can still receive messages with another
	//   MessageGroupId as long as it is also visible.
	//   - If a caller of ReceiveMessage can't track the ReceiveRequestAttemptId , no
	//   retries work until the original visibility timeout expires. As a result, delays
	//   might occur but the messages in the queue remain in a strict order.
	// The maximum length of ReceiveRequestAttemptId is 128 characters.
	// ReceiveRequestAttemptId can contain alphanumeric characters ( a-z , A-Z , 0-9 )
	// and punctuation ( !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~ ). For best practices of
	// using ReceiveRequestAttemptId , see Using the ReceiveRequestAttemptId Request
	// Parameter (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/using-receiverequestattemptid-request-parameter.html)
	// in the Amazon SQS Developer Guide.
	ReceiveRequestAttemptId *string
	// The duration that the received messages are hidden from subsequent
	// retrieve requests after being retrieved by a ReceiveMessage request.
	//
	// default: 30 seconds
	VisibilityTimeout time.Duration
	// The duration for which the call waits for a message to arrive in
	// the queue before returning. If a message is available, the call returns sooner
	// than WaitTimeSeconds . If no messages are available and the wait time expires,
	// the call returns successfully with an empty list of messages. To avoid HTTP
	// errors, ensure that the HTTP response timeout for ReceiveMessage requests is
	// longer than the WaitTimeSeconds parameter. For example, with the Java SDK, you
	// can set HTTP transport settings using the NettyNioAsyncHttpClient (https://sdk.amazonaws.com/java/api/latest/software/amazon/awssdk/http/nio/netty/NettyNioAsyncHttpClient.html)
	// for asynchronous clients, or the ApacheHttpClient (https://sdk.amazonaws.com/java/api/latest/software/amazon/awssdk/http/apache/ApacheHttpClient.html)
	// for synchronous clients.
	//
	// default: 0
	WaitTime time.Duration
}
