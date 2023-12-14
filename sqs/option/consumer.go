package option

import "time"

type Consumer struct {
	Default
	// default: false
	DeleteMessageProcessedSuccess bool
	// default: 5 seconds
	ConsumerMessageTimeout time.Duration
	// default: 5 seconds
	DelayQueryLoop time.Duration
	// The maximum number of messages to return. Amazon SQS never returns more
	// messages than this value (however, fewer messages might be returned). Valid
	// values: 1 to 10.
	//
	// default: 10
	MaxNumberOfMessages int32
	// The duration that the received messages are hidden from subsequent
	// retrieve requests after being retrieved by a ReceiveMessage request.
	//
	// default: 0 seconds
	VisibilityTimeout time.Duration
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
	// default: 0 seconds
	WaitTimeSeconds time.Duration
}

func NewConsumer() Consumer {
	return Consumer{}
}

func (o Consumer) SetDeleteMessageProcessedSuccess(b bool) Consumer {
	o.DeleteMessageProcessedSuccess = b
	return o
}

func (o Consumer) SetConsumerMessageTimeout(d time.Duration) Consumer {
	o.ConsumerMessageTimeout = d
	return o
}

func (o Consumer) SetMaxNumberOfMessages(i int32) Consumer {
	o.MaxNumberOfMessages = i
	return o
}

func (o Consumer) SetDelayQueryLoop(d time.Duration) Consumer {
	o.DelayQueryLoop = d
	return o
}

func (o Consumer) SetReceiveRequestAttemptId(s string) Consumer {
	o.ReceiveRequestAttemptId = &s
	return o
}

func (o Consumer) SetVisibilityTimeout(d time.Duration) Consumer {
	o.VisibilityTimeout = d
	return o
}

func (o Consumer) SetWaitTimeSeconds(d time.Duration) Consumer {
	o.WaitTimeSeconds = d
	return o
}

func (o Consumer) SetDebugMode(b bool) Consumer {
	o.DebugMode = b
	return o
}

func (o Consumer) SetOptionHttp(opt Http) Consumer {
	o.OptionHttp = &opt
	return o
}

func GetConsumerByParams(opts []Consumer) Consumer {
	var result Consumer
	for _, opt := range opts {
		fillDefaultFields(opt.Default, &result.Default)
		if opt.DeleteMessageProcessedSuccess {
			result.DeleteMessageProcessedSuccess = opt.DeleteMessageProcessedSuccess
		}
		if opt.ConsumerMessageTimeout.Seconds() > 0 {
			result.ConsumerMessageTimeout = opt.ConsumerMessageTimeout
		}
		if opt.DelayQueryLoop.Seconds() > 0 {
			result.DelayQueryLoop = opt.DelayQueryLoop
		}
		if opt.MaxNumberOfMessages > 0 {
			result.MaxNumberOfMessages = opt.MaxNumberOfMessages
		}
		if opt.VisibilityTimeout.Seconds() > 0 {
			result.VisibilityTimeout = opt.VisibilityTimeout
		}
		if opt.ReceiveRequestAttemptId != nil && len(*opt.ReceiveRequestAttemptId) != 0 {
			result.ReceiveRequestAttemptId = opt.ReceiveRequestAttemptId
		}
		if opt.WaitTimeSeconds.Seconds() > 0 {
			result.WaitTimeSeconds = opt.WaitTimeSeconds
		}
	}
	if result.MaxNumberOfMessages <= 0 {
		result.MaxNumberOfMessages = 10
	}
	if result.ConsumerMessageTimeout.Seconds() == 0 {
		result.ConsumerMessageTimeout = 5 * time.Second
	}
	if result.DelayQueryLoop.Seconds() == 0 {
		result.DelayQueryLoop = 5 * time.Second
	}
	return result
}
