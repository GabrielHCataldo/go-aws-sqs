package option

import (
	"github.com/GabrielHCataldo/go-aws-sqs-template/internal/util"
	"reflect"
	"time"
)

type Producer struct {
	Default
	// The length of time, in seconds, for which to delay a specific message. Maximum: 15 minutes.
	// Messages with a positive DelaySeconds value become available for processing after the delay period is finished.
	// If you don't specify a value, the default value for the queue applies. When you set
	// FifoQueue , you can't set DelaySeconds per message. You can set this parameter
	// only on a queue level.
	DelaySeconds time.Duration `json:"delaySeconds,omitempty"`
	// Message attributes, must be of type Map or Struct, other types are not acceptable.
	MessageAttributes any `json:"messageAttributes,omitempty"`
	// The message system attribute to send.
	//   - Currently, the only supported message system attribute is AWSTraceHeader .
	//   Its type must be String and its value must be a correctly formatted X-Ray
	//   trace header string.
	//   - The size of a message system attribute doesn't count towards the total size
	//   of a message.
	MessageSystemAttributes *MessageSystemAttributes `json:"messageSystemAttributes,omitempty"`
	// This parameter applies only to FIFO (first-in-first-out) queues. The token used
	// for deduplication of sent messages. If a message with a particular
	// MessageDeduplicationId is sent successfully, any messages sent with the same
	// MessageDeduplicationId are accepted successfully but aren't delivered during the
	// 5-minute deduplication interval. For more information, see Exactly-once
	// processing (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-exactly-once-processing.html)
	// in the Amazon SQS Developer Guide.
	//   - Every message must have a unique MessageDeduplicationId ,
	//   - You may provide a MessageDeduplicationId explicitly.
	//   - If you aren't able to provide a MessageDeduplicationId and you enable
	//   ContentBasedDeduplication for your queue, Amazon SQS uses an SHA-256 hash to
	//   generate the MessageDeduplicationId using the body of the message (but not the
	//   attributes of the message).
	//   - If you don't provide a MessageDeduplicationId and the queue doesn't have
	//   ContentBasedDeduplication set, the action fails with an error.
	//   - If the queue has ContentBasedDeduplication set, your MessageDeduplicationId
	//   overrides the generated one.
	//   - When ContentBasedDeduplication is in effect, messages with identical content
	//   sent within the deduplication interval are treated as duplicates and only one
	//   copy of the message is delivered.
	//   - If you send one message with ContentBasedDeduplication enabled and then
	//   another message with a MessageDeduplicationId that is the same as the one
	//   generated for the first MessageDeduplicationId , the two messages are treated
	//   as duplicates and only one copy of the message is delivered.
	// The MessageDeduplicationId is available to the consumer of the message (this
	// can be useful for troubleshooting delivery issues). If a message is sent
	// successfully but the acknowledgement is lost and the message is resent with the
	// same MessageDeduplicationId after the deduplication interval, Amazon SQS can't
	// detect duplicate messages. Amazon SQS continues to keep track of the message
	// deduplication ID even after the message is received and deleted. The maximum
	// length of MessageDeduplicationId is 128 characters. MessageDeduplicationId can
	// contain alphanumeric characters ( a-z , A-Z , 0-9 ) and punctuation (
	// !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~ ). For best practices of using
	// MessageDeduplicationId , see Using the MessageDeduplicationId Property (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/using-messagededuplicationid-property.html)
	// in the Amazon SQS Developer Guide.
	MessageDeduplicationId *string `json:"messageDeduplicationId,omitempty"`
	// This parameter applies only to FIFO (first-in-first-out) queues. The tag that
	// specifies that a message belongs to a specific message group. Messages that
	// belong to the same message group are processed in a FIFO manner (however,
	// messages in different message groups might be processed out of order). To
	// interleave multiple ordered streams within a single queue, use MessageGroupId
	// values (for example, session data for multiple users). In this scenario,
	// multiple consumers can process the queue, but the session data of each user is
	// processed in a FIFO fashion.
	//   - You must associate a non-empty MessageGroupId with a message. If you don't
	//   provide a MessageGroupId , the action fails.
	//   - ReceiveMessage might return messages with multiple MessageGroupId values.
	//   For each MessageGroupId , the messages are sorted by time sent. The caller
	//   can't specify a MessageGroupId .
	// The length of MessageGroupId is 128 characters. Valid values: alphanumeric
	// characters and punctuation (!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~) . For best
	// practices of using MessageGroupId , see Using the MessageGroupId Property (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/using-messagegroupid-property.html)
	// in the Amazon SQS Developer Guide. MessageGroupId is required for FIFO queues.
	// You can't use it for Standard queues.
	MessageGroupId *string `json:"messageGroupId,omitempty"`
}

type MessageSystemAttributes struct {
	// its value must be a correctly formatted X-Ray
	// trace header string.
	//
	// this member is required
	AWSTraceHeader string `json:"AWSTraceHeader,omitempty"`
}

func NewProducer() Producer {
	return Producer{}
}

func (o Producer) SetDelaySeconds(t time.Duration) Producer {
	o.DelaySeconds = t
	return o
}

func (o Producer) SetMessageAttributes(m any) Producer {
	o.MessageAttributes = m
	return o
}

func (o Producer) SetMessageSystemAttributes(m MessageSystemAttributes) Producer {
	o.MessageSystemAttributes = &m
	return o
}

func (o Producer) SetMessageDeduplicationId(s string) Producer {
	o.MessageDeduplicationId = &s
	return o
}

func (o Producer) SetMessageGroupId(s string) Producer {
	o.MessageGroupId = &s
	return o
}

func (o Producer) SetDebugMode(b bool) Producer {
	o.DebugMode = b
	return o
}

func (o Producer) SetHttpClient(opt HttpClient) Producer {
	o.HttpClient = &opt
	return o
}

func GetProducerByParams(opts []Producer) Producer {
	var result Producer
	for _, opt := range opts {
		fillDefaultFields(opt.Default, &result.Default)
		if opt.DelaySeconds > 0 {
			result.DelaySeconds = opt.DelaySeconds
		}
		if util.IsValidType(opt.MessageAttributes) {
			result.MessageAttributes = opt.MessageAttributes
		}
		if util.IsValidType(opt.MessageSystemAttributes) &&
			!util.IsZeroReflect(reflect.ValueOf(opt.MessageSystemAttributes)) &&
			len(opt.MessageSystemAttributes.AWSTraceHeader) != 0 {
			result.MessageSystemAttributes = opt.MessageSystemAttributes
		}
		if opt.MessageDeduplicationId != nil && len(*opt.MessageDeduplicationId) != 0 {
			result.MessageDeduplicationId = opt.MessageDeduplicationId
		}
		if opt.MessageGroupId != nil && len(*opt.MessageGroupId) != 0 {
			result.MessageGroupId = opt.MessageGroupId
		}
	}
	return result
}
