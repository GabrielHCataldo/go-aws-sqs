package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go-aws-sqs/sqs/option"
)

type TagQueueInput struct {
	// The URL of the queue.
	//
	// This member is required.
	QueueUrl string
	// The list of tags to be added to the specified queue.
	//
	// This member is required.
	Tags map[string]string
}

type SetQueueAttributesInput struct {
	// A map of attributes to set. The following lists the names, descriptions, and
	// values of the special request parameters that the SetQueueAttributes action
	// uses:
	//   - DelaySeconds – The length of time, in seconds, for which the delivery of all
	//   messages in the queue is delayed. Valid values: An integer from 0 to 900 (15
	//   minutes). Default: 0.
	//   - MaximumMessageSize – The limit of how many bytes a message can contain
	//   before Amazon SQS rejects it. Valid values: An integer from 1,024 bytes (1 KiB)
	//   up to 262,144 bytes (256 KiB). Default: 262,144 (256 KiB).
	//   - MessageRetentionPeriod – The length of time, in seconds, for which Amazon
	//   SQS retains a message. Valid values: An integer representing seconds, from 60 (1
	//   minute) to 1,209,600 (14 days). Default: 345,600 (4 days). When you change a
	//   queue's attributes, the change can take up to 60 seconds for most of the
	//   attributes to propagate throughout the Amazon SQS system. Changes made to the
	//   MessageRetentionPeriod attribute can take up to 15 minutes and will impact
	//   existing messages in the queue potentially causing them to be expired and
	//   deleted if the MessageRetentionPeriod is reduced below the age of existing
	//   messages.
	//   - Policy – The queue's policy. A valid Amazon Web Services policy. For more
	//   information about policy structure, see Overview of Amazon Web Services IAM
	//   Policies (https://docs.aws.amazon.com/IAM/latest/UserGuide/PoliciesOverview.html)
	//   in the Identity and Access Management User Guide.
	//   - ReceiveMessageWaitTimeSeconds – The length of time, in seconds, for which a
	//   ReceiveMessage action waits for a message to arrive. Valid values: An integer
	//   from 0 to 20 (seconds). Default: 0.
	//   - VisibilityTimeout – The visibility timeout for the queue, in seconds. Valid
	//   values: An integer from 0 to 43,200 (12 hours). Default: 30. For more
	//   information about the visibility timeout, see Visibility Timeout (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html)
	//   in the Amazon SQS Developer Guide.
	// The following attributes apply only to dead-letter queues: (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
	//   - RedrivePolicy – The string that includes the parameters for the dead-letter
	//   queue functionality of the source queue as a JSON object. The parameters are as
	//   follows:
	//   - deadLetterTargetArn – The Amazon Resource Name (ARN) of the dead-letter
	//   queue to which Amazon SQS moves messages after the value of maxReceiveCount is
	//   exceeded.
	//   - maxReceiveCount – The number of times a message is delivered to the source
	//   queue before being moved to the dead-letter queue. Default: 10. When the
	//   ReceiveCount for a message exceeds the maxReceiveCount for a queue, Amazon SQS
	//   moves the message to the dead-letter-queue.
	//   - RedriveAllowPolicy – The string that includes the parameters for the
	//   permissions for the dead-letter queue redrive permission and which source queues
	//   can specify dead-letter queues as a JSON object. The parameters are as follows:
	//   - redrivePermission – The permission type that defines which source queues can
	//   specify the current queue as the dead-letter queue. Valid values are:
	//   - allowAll – (Default) Any source queues in this Amazon Web Services account
	//   in the same Region can specify this queue as the dead-letter queue.
	//   - denyAll – No source queues can specify this queue as the dead-letter queue.
	//   - byQueue – Only queues specified by the sourceQueueArns parameter can specify
	//   this queue as the dead-letter queue.
	//   - sourceQueueArns – The Amazon Resource Names (ARN)s of the source queues that
	//   can specify this queue as the dead-letter queue and redrive messages. You can
	//   specify this parameter only when the redrivePermission parameter is set to
	//   byQueue . You can specify up to 10 source queue ARNs. To allow more than 10
	//   source queues to specify dead-letter queues, set the redrivePermission
	//   parameter to allowAll .
	// The dead-letter queue of a FIFO queue must also be a FIFO queue. Similarly, the
	// dead-letter queue of a standard queue must also be a standard queue. The
	// following attributes apply only to server-side-encryption (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-server-side-encryption.html)
	// :
	//   - KmsMasterKeyId – The ID of an Amazon Web Services managed customer master
	//   key (CMK) for Amazon SQS or a custom CMK. For more information, see Key Terms (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-server-side-encryption.html#sqs-sse-key-terms)
	//   . While the alias of the AWS-managed CMK for Amazon SQS is always
	//   alias/aws/sqs , the alias of a custom CMK can, for example, be alias/MyAlias
	//   . For more examples, see KeyId (https://docs.aws.amazon.com/kms/latest/APIReference/API_DescribeKey.html#API_DescribeKey_RequestParameters)
	//   in the Key Management Service API Reference.
	//   - KmsDataKeyReusePeriodSeconds – The length of time, in seconds, for which
	//   Amazon SQS can reuse a data key (https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#data-keys)
	//   to encrypt or decrypt messages before calling KMS again. An integer representing
	//   seconds, between 60 seconds (1 minute) and 86,400 seconds (24 hours). Default:
	//   300 (5 minutes). A shorter time period provides better security but results in
	//   more calls to KMS which might incur charges after Free Tier. For more
	//   information, see How Does the Data Key Reuse Period Work? (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-server-side-encryption.html#sqs-how-does-the-data-key-reuse-period-work)
	//   .
	//   - SqsManagedSseEnabled – Enables server-side queue encryption using SQS owned
	//   encryption keys. Only one server-side encryption option is supported per queue
	//   (for example, SSE-KMS (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-sse-existing-queue.html)
	//   or SSE-SQS (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-sqs-sse-queue.html)
	//   ).
	// The following attribute applies only to FIFO (first-in-first-out) queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues.html)
	// :
	//   - ContentBasedDeduplication – Enables content-based deduplication. For more
	//   information, see Exactly-once processing (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-exactly-once-processing.html)
	//   in the Amazon SQS Developer Guide. Note the following:
	//   - Every message must have a unique MessageDeduplicationId .
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
	// The following attributes apply only to high throughput for FIFO queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/high-throughput-fifo.html)
	// :
	//   - DeduplicationScope – Specifies whether message deduplication occurs at the
	//   message group or queue level. Valid values are messageGroup and queue .
	//   - FifoThroughputLimit – Specifies whether the FIFO queue throughput quota
	//   applies to the entire queue or per message group. Valid values are perQueue
	//   and perMessageGroupId . The perMessageGroupId value is allowed only when the
	//   value for DeduplicationScope is messageGroup .
	// To enable high throughput for FIFO queues, do the following:
	//   - Set DeduplicationScope to messageGroup .
	//   - Set FifoThroughputLimit to perMessageGroupId .
	// If you set these attributes to anything other than the values shown for
	// enabling high throughput, normal throughput is in effect and deduplication
	// occurs as specified. For information on throughput quotas, see Quotas related
	// to messages (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/quotas-messages.html)
	// in the Amazon SQS Developer Guide.
	//
	// This member is required.
	Attributes map[string]string
	// The URL of the Amazon SQS queue whose attributes are set. Queue URLs and names
	// are case-sensitive.
	//
	// This member is required.
	QueueUrl string
}

type UntagQueueInput struct {
	// The URL of the queue.
	//
	// This member is required.
	QueueUrl *string
	// The list of tags to be removed from the specified queue.
	//
	// This member is required.
	TagKeys []string
}

type GetQueueUrlInput struct {
	// The name of the queue whose URL must be fetched. Maximum 80 characters. Valid
	// values: alphanumeric characters, hyphens ( - ), and underscores ( _ ). Queue
	// URLs and names are case-sensitive.
	//
	// This member is required.
	QueueName string
	// The Amazon Web Services account ID of the account that created the queue.
	QueueOwnerAWSAccountId string
}

type GetQueueAttributesInput struct {
	// The URL of the Amazon SQS queue whose attribute information is retrieved. Queue
	// URLs and names are case-sensitive.
	//
	// This member is required.
	QueueUrl string
	// A list of attributes for which to retrieve information. The AttributeNames
	// parameter is optional, but if you don't specify values for this parameter, the
	// request returns empty results. In the future, new attributes might be added. If
	// you write code that calls this action, we recommend that you structure your code
	// so that it can handle new attributes gracefully. The following attributes are
	// supported: The ApproximateNumberOfMessagesDelayed ,
	// ApproximateNumberOfMessagesNotVisible , and ApproximateNumberOfMessages metrics
	// may not achieve consistency until at least 1 minute after the producers stop
	// sending messages. This period is required for the queue metadata to reach
	// eventual consistency.
	//   - All – Returns all values.
	//   - ApproximateNumberOfMessages – Returns the approximate number of messages
	//   available for retrieval from the queue.
	//   - ApproximateNumberOfMessagesDelayed – Returns the approximate number of
	//   messages in the queue that are delayed and not available for reading
	//   immediately. This can happen when the queue is configured as a delay queue or
	//   when a message has been sent with a delay parameter.
	//   - ApproximateNumberOfMessagesNotVisible – Returns the approximate number of
	//   messages that are in flight. Messages are considered to be in flight if they
	//   have been sent to a client but have not yet been deleted or have not yet reached
	//   the end of their visibility window.
	//   - CreatedTimestamp – Returns the time when the queue was created in seconds (
	//   epoch time (http://en.wikipedia.org/wiki/Unix_time) ).
	//   - DelaySeconds – Returns the default delay on the queue in seconds.
	//   - LastModifiedTimestamp – Returns the time when the queue was last changed in
	//   seconds ( epoch time (http://en.wikipedia.org/wiki/Unix_time) ).
	//   - MaximumMessageSize – Returns the limit of how many bytes a message can
	//   contain before Amazon SQS rejects it.
	//   - MessageRetentionPeriod – Returns the length of time, in seconds, for which
	//   Amazon SQS retains a message. When you change a queue's attributes, the change
	//   can take up to 60 seconds for most of the attributes to propagate throughout the
	//   Amazon SQS system. Changes made to the MessageRetentionPeriod attribute can
	//   take up to 15 minutes and will impact existing messages in the queue potentially
	//   causing them to be expired and deleted if the MessageRetentionPeriod is
	//   reduced below the age of existing messages.
	//   - Policy – Returns the policy of the queue.
	//   - QueueArn – Returns the Amazon resource name (ARN) of the queue.
	//   - ReceiveMessageWaitTimeSeconds – Returns the length of time, in seconds, for
	//   which the ReceiveMessage action waits for a message to arrive.
	//   - VisibilityTimeout – Returns the visibility timeout for the queue. For more
	//   information about the visibility timeout, see Visibility Timeout (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html)
	//   in the Amazon SQS Developer Guide.
	// The following attributes apply only to dead-letter queues: (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
	//   - RedrivePolicy – The string that includes the parameters for the dead-letter
	//   queue functionality of the source queue as a JSON object. The parameters are as
	//   follows:
	//   - deadLetterTargetArn – The Amazon Resource Name (ARN) of the dead-letter
	//   queue to which Amazon SQS moves messages after the value of maxReceiveCount is
	//   exceeded.
	//   - maxReceiveCount – The number of times a message is delivered to the source
	//   queue before being moved to the dead-letter queue. Default: 10. When the
	//   ReceiveCount for a message exceeds the maxReceiveCount for a queue, Amazon SQS
	//   moves the message to the dead-letter-queue.
	//   - RedriveAllowPolicy – The string that includes the parameters for the
	//   permissions for the dead-letter queue redrive permission and which source queues
	//   can specify dead-letter queues as a JSON object. The parameters are as follows:
	//   - redrivePermission – The permission type that defines which source queues can
	//   specify the current queue as the dead-letter queue. Valid values are:
	//   - allowAll – (Default) Any source queues in this Amazon Web Services account
	//   in the same Region can specify this queue as the dead-letter queue.
	//   - denyAll – No source queues can specify this queue as the dead-letter queue.
	//   - byQueue – Only queues specified by the sourceQueueArns parameter can specify
	//   this queue as the dead-letter queue.
	//   - sourceQueueArns – The Amazon Resource Names (ARN)s of the source queues that
	//   can specify this queue as the dead-letter queue and redrive messages. You can
	//   specify this parameter only when the redrivePermission parameter is set to
	//   byQueue . You can specify up to 10 source queue ARNs. To allow more than 10
	//   source queues to specify dead-letter queues, set the redrivePermission
	//   parameter to allowAll .
	// The dead-letter queue of a FIFO queue must also be a FIFO queue. Similarly, the
	// dead-letter queue of a standard queue must also be a standard queue. The
	// following attributes apply only to server-side-encryption (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-server-side-encryption.html)
	// :
	//   - KmsMasterKeyId – Returns the ID of an Amazon Web Services managed customer
	//   master key (CMK) for Amazon SQS or a custom CMK. For more information, see
	//   Key Terms (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-server-side-encryption.html#sqs-sse-key-terms)
	//   .
	//   - KmsDataKeyReusePeriodSeconds – Returns the length of time, in seconds, for
	//   which Amazon SQS can reuse a data key to encrypt or decrypt messages before
	//   calling KMS again. For more information, see How Does the Data Key Reuse
	//   Period Work? (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-server-side-encryption.html#sqs-how-does-the-data-key-reuse-period-work)
	//   .
	//   - SqsManagedSseEnabled – Returns information about whether the queue is using
	//   SSE-SQS encryption using SQS owned encryption keys. Only one server-side
	//   encryption option is supported per queue (for example, SSE-KMS (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-sse-existing-queue.html)
	//   or SSE-SQS (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-sqs-sse-queue.html)
	//   ).
	// The following attributes apply only to FIFO (first-in-first-out) queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues.html)
	// :
	//   - FifoQueue – Returns information about whether the queue is FIFO. For more
	//   information, see FIFO queue logic (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-understanding-logic.html)
	//   in the Amazon SQS Developer Guide. To determine whether a queue is FIFO (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues.html)
	//   , you can check whether QueueName ends with the .fifo suffix.
	//   - ContentBasedDeduplication – Returns whether content-based deduplication is
	//   enabled for the queue. For more information, see Exactly-once processing (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-exactly-once-processing.html)
	//   in the Amazon SQS Developer Guide.
	// The following attributes apply only to high throughput for FIFO queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/high-throughput-fifo.html)
	// :
	//   - DeduplicationScope – Specifies whether message deduplication occurs at the
	//   message group or queue level. Valid values are messageGroup and queue .
	//   - FifoThroughputLimit – Specifies whether the FIFO queue throughput quota
	//   applies to the entire queue or per message group. Valid values are perQueue
	//   and perMessageGroupId . The perMessageGroupId value is allowed only when the
	//   value for DeduplicationScope is messageGroup .
	// To enable high throughput for FIFO queues, do the following:
	//   - Set DeduplicationScope to messageGroup .
	//   - Set FifoThroughputLimit to perMessageGroupId .
	// If you set these attributes to anything other than the values shown for
	// enabling high throughput, normal throughput is in effect and deduplication
	// occurs as specified. For information on throughput quotas, see Quotas related
	// to messages (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/quotas-messages.html)
	// in the Amazon SQS Developer Guide.
	AttributeNames []types.QueueAttributeName
}

// Creates a new standard or FIFO queue. You can pass one or more attributes in
// the request. Keep the following in mind:
//   - If you don't specify the FifoQueue attribute, Amazon SQS creates a standard
//     queue. You can't change the queue type after you create it and you can't convert
//     an existing standard queue into a FIFO queue. You must either create a new FIFO
//     queue for your application or delete your existing standard queue and recreate
//     it as a FIFO queue. For more information, see Moving From a Standard Queue to
//     a FIFO Queue (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues.html#FIFO-queues-moving)
//     in the Amazon SQS Developer Guide.
//   - If you don't provide a value for an attribute, the queue is created with
//     the default value for the attribute.
//   - If you delete a queue, you must wait at least 60 seconds before creating a
//     queue with the same name.
//
// To successfully create a new queue, you must provide a queue name that adheres
// to the limits related to queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/limits-queues.html)
// and is unique within the scope of your queues. After you create a queue, you
// must wait at least one second after the queue is created to be able to use the
// queue. To get the queue URL, use the GetQueueUrl action. GetQueueUrl requires
// only the QueueName parameter. be aware of existing queue names:
//   - If you provide the name of an existing queue along with the exact names and
//     values of all the queue's attributes, CreateQueue returns the queue URL for
//     the existing queue.
//   - If the queue name, attribute names, or attribute values don't match an
//     existing queue, CreateQueue returns an error.
//
// Cross-account permissions don't apply to this action. For more information, see
// Grant cross-account permissions to a role and a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon SQS Developer Guide.
func CreateQueue(ctx context.Context, queueName string, opts ...*option.OptionsCreateQueue) (*sqs.CreateQueueOutput,
	error) {
	debugMode := getDebugModeByOptsCreateQueue(opts...)
	loggerInfo(debugMode, "creating queue sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  &queueName,
		Attributes: getAttributesByOptsCreateQueue(opts...),
		Tags:       getTagsByOptsCreateQueue(opts...),
	}, optionsHttp(getOptionsHttpByOptsCreateQueue(opts...)))
	if err != nil {
		loggerErr(debugMode, "error create queue sqs:", err)
	} else {
		loggerInfo(debugMode, "queue created successfully:", output)
	}
	return output, err
}

// Add cost allocation tags to the specified Amazon SQS queue. For an overview,
// see Tagging Your Amazon SQS Queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-queue-tags.html)
// in the Amazon SQS Developer Guide. When you use queue tags, keep the following
// guidelines in mind:
//   - Adding more than 50 tags to a queue isn't recommended.
//   - Tags don't have any semantic meaning. Amazon SQS interprets tags as
//     character strings.
//   - Tags are case-sensitive.
//   - A new tag with a key identical to that of an existing tag overwrites the
//     existing tag.
//
// For a full list of tag restrictions, see Quotas related to queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-limits.html#limits-queues)
// in the Amazon SQS Developer Guide. Cross-account permissions don't apply to this
// action. For more information, see Grant cross-account permissions to a role and
// a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon SQS Developer Guide.
func TagQueue(ctx context.Context, input TagQueueInput, opts ...*option.OptionsDefault) (*sqs.TagQueueOutput, error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "tag queue sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.TagQueue(ctx, &sqs.TagQueueInput{
		QueueUrl: &input.QueueUrl,
		Tags:     input.Tags,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error tag queue sqs:", err)
	} else {
		loggerInfo(debugMode, "queue deleted successfully:", output)
	}
	return output, err
}

// Sets the value of one or more queue attributes. When you change a queue's
// attributes, the change can take up to 60 seconds for most of the attributes to
// propagate throughout the Amazon SQS system. Changes made to the
// MessageRetentionPeriod attribute can take up to 15 minutes and will impact
// existing messages in the queue potentially causing them to be expired and
// deleted if the MessageRetentionPeriod is reduced below the age of existing
// messages.
//   - In the future, new attributes might be added. If you write code that calls
//     this action, we recommend that you structure your code so that it can handle new
//     attributes gracefully.
//   - Cross-account permissions don't apply to this action. For more information,
//     see Grant cross-account permissions to a role and a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
//     in the Amazon SQS Developer Guide.
//   - To remove the ability to change queue permissions, you must deny permission
//     to the AddPermission , RemovePermission , and SetQueueAttributes actions in
//     your IAM policy.
func SetAttributeQueue(ctx context.Context, input SetQueueAttributesInput, opts ...*option.OptionsDefault) (
	*sqs.SetQueueAttributesOutput, error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "set attributes queue sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl:   &input.QueueUrl,
		Attributes: input.Attributes,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error set attributes queue sqs:", err)
	} else {
		loggerInfo(debugMode, "queue updated successfully:", output)
	}
	return output, err
}

// Remove cost allocation tags from the specified Amazon SQS queue. For an
// overview, see Tagging Your Amazon SQS Queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-queue-tags.html)
// in the Amazon SQS Developer Guide. Cross-account permissions don't apply to this
// action. For more information, see Grant cross-account permissions to a role and
// a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon SQS Developer Guide.
func UntagQueue(ctx context.Context, input UntagQueueInput, opts ...*option.OptionsDefault) (*sqs.UntagQueueOutput,
	error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "untag queue sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.UntagQueue(ctx, &sqs.UntagQueueInput{
		QueueUrl: input.QueueUrl,
		TagKeys:  input.TagKeys,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error untag queue sqs:", err)
	} else {
		loggerInfo(debugMode, "queue untagged successfully:", output)
	}
	return output, err
}

// Deletes available messages in a queue (including in-flight messages) specified
// by the QueueURL parameter. When you use the PurgeQueue action, you can't
// retrieve any messages deleted from a queue. The message deletion process takes
// up to 60 seconds. We recommend waiting for 60 seconds regardless of your queue's
// size. Messages sent to the queue before you call PurgeQueue might be received
// but are deleted within the next minute. Messages sent to the queue after you
// call PurgeQueue might be deleted while the queue is being purged.
func PurgeQueue(ctx context.Context, queueUrl string, opts ...*option.OptionsDefault) (*sqs.PurgeQueueOutput, error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "purge queue sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
		QueueUrl: &queueUrl,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error purge queue sqs:", err)
	} else {
		loggerInfo(debugMode, "queue purged successfully:", output)
	}
	return output, err
}

// Deletes the queue specified by the QueueUrl , regardless of the queue's
// contents. Be careful with the DeleteQueue action: When you delete a queue, any
// messages in the queue are no longer available. When you delete a queue, the
// deletion process takes up to 60 seconds. Requests you send involving that queue
// during the 60 seconds might succeed. For example, a SendMessage request might
// succeed, but after 60 seconds the queue and the message you sent no longer
// exist. When you delete a queue, you must wait at least 60 seconds before
// creating a queue with the same name. Cross-account permissions don't apply to
// this action. For more information, see Grant cross-account permissions to a
// role and a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon SQS Developer Guide. The delete operation uses the HTTP GET verb.
func DeleteQueue(ctx context.Context, queueUrl string, opts ...*option.OptionsDefault) (*sqs.DeleteQueueOutput, error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "deleting queue sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.DeleteQueue(ctx, &sqs.DeleteQueueInput{
		QueueUrl: &queueUrl,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error delete queue sqs:", err)
	} else {
		loggerInfo(debugMode, "queue deleted successfully:", output)
	}
	return output, err
}

// GetQueueUrl Returns the URL of an existing Amazon SQS queue. To access a queue that belongs
// to another AWS account, use the QueueOwnerAWSAccountId parameter to specify the
// account ID of the queue's owner. The queue's owner must grant you permission to
// access the queue. For more information about shared queue access, see
// AddPermission or see Allow Developers to Write Messages to a Shared Queue (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-writing-an-sqs-policy.html#write-messages-to-shared-queue)
// in the Amazon SQS Developer Guide.
func GetQueueUrl(ctx context.Context, input GetQueueUrlInput, opts ...*option.OptionsDefault) (*sqs.GetQueueUrlOutput,
	error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "get queue url sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName:              &input.QueueName,
		QueueOwnerAWSAccountId: &input.QueueOwnerAWSAccountId,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error get queue url sqs:", err)
	} else {
		loggerInfo(debugMode, "get queue successfully:", output)
	}
	return output, err
}

// GetQueueAttributes Gets attributes for the specified queue.
// To determine whether a queue is FIFO (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues.html)
// you can check whether QueueName ends with the .fifo suffix.
func GetQueueAttributes(ctx context.Context, input GetQueueAttributesInput, opts ...*option.OptionsDefault) (
	*sqs.GetQueueAttributesOutput, error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "get queue attributes sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       &input.QueueUrl,
		AttributeNames: input.AttributeNames,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error queue attributes sqs:", err)
	} else {
		loggerInfo(debugMode, "get queue attributes successfully:", output)
	}
	return output, err
}

// ListQueues Returns a list of your queues in the current region. The response includes a
// maximum of 1,000 results. If you specify a value for the optional
// QueueNamePrefix parameter, only queues with a name that begins with the
// specified value are returned. The listQueues methods supports pagination. Set
// parameter MaxResults in the request to specify the maximum number of results to
// be returned in the response. If you do not set MaxResults , the response
// includes a maximum of 1,000 results. If you set MaxResults and there are
// additional results to display, the response includes a value for NextToken . Use
// NextToken as a parameter in your next request to listQueues to receive the next
// page of results. Cross-account permissions don't apply to this action. For more
// information, see Grant cross-account permissions to a role and a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon SQS Developer Guide.
func ListQueues(ctx context.Context, opts ...*option.OptionsListQueues) (*sqs.ListQueuesOutput, error) {
	debugMode := getDebugModeByOptsListQueues(opts...)
	loggerInfo(debugMode, "list queue tags sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.ListQueues(ctx, &sqs.ListQueuesInput{
		MaxResults:      getMaxResultsByOptsListQueues(opts...),
		NextToken:       getNextTokenByOptsListQueues(opts...),
		QueueNamePrefix: getQueueNamePrefixByOptsListQueues(opts...),
	}, optionsHttp(getOptionsHttpByOptsListQueues(opts...)))
	if err != nil {
		loggerErr(debugMode, "error list queue tags sqs:", err)
	} else {
		loggerInfo(debugMode, "get queue tags successfully:", output)
	}
	return output, err
}

// ListQueueTags List all cost allocation tags added to the specified Amazon SQS queue. For an
// overview, see Tagging Your Amazon SQS Queues (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-queue-tags.html)
// in the Amazon SQS Developer Guide. Cross-account permissions don't apply to this
// action. For more information, see Grant cross-account permissions to a role and
// a username (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon SQS Developer Guide.
func ListQueueTags(ctx context.Context, queueUrl string, opts ...*option.OptionsDefault) (
	*sqs.ListQueueTagsOutput, error) {
	debugMode := getDebugModeByOptsDefault(opts...)
	loggerInfo(debugMode, "list queue tags sqs..")
	client, err := getSqsClient(ctx)
	if err != nil {
		loggerErr(debugMode, "error get client sqs:", err)
		return nil, err
	}
	output, err := client.ListQueueTags(ctx, &sqs.ListQueueTagsInput{
		QueueUrl: &queueUrl,
	}, optionsHttp(getOptionsHttpByOptsDefault(opts...)))
	if err != nil {
		loggerErr(debugMode, "error list queue tags sqs:", err)
	} else {
		loggerInfo(debugMode, "get queue tags successfully:", output)
	}
	return output, err
}
