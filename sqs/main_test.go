package sqs

import (
	"context"
	"errors"
	"github.com/GabrielHCataldo/go-logger/logger"
	"go-aws-sqs/sqs/option"
	"os"
	"strconv"
	"testing"
	"time"
)

const SqsQueueTestName = "SQS_QUEUE_TEST_NAME"
const SqsQueueTestUrl = "SQS_QUEUE_TEST_URL"
const SqsQueueTestStringArn = "SQS_QUEUE_TEST_STRING_ARN"
const SqsQueueTestFifoUrl = "SQS_QUEUE_TEST_FIFO_URL"
const SqsQueueTestEmptyUrl = "SQS_QUEUE_TEST_EMPTY_URL"
const SqsQueueTestStringUrl = "SQS_QUEUE_TEST_STRING_URL"
const SqsQueueTestDlqArn = "SQS_QUEUE_TEST_DLQ_ARN"
const SqsQueueCreateTestName = "SQS_QUEUE_CREATE_TEST_NAME"
const SqsQueueCreateTestUrl = "SQS_QUEUE_CREATE_TEST_URL"
const SqsMessageId = "SQS_MESSAGE_ID"
const SqsMessageReceiptHandle = "SQS_MESSAGE_RECEIPT_HANDLE"
const SqsTaskHandle = "SQS_TASK_HANDLE"

type testProducer struct {
	name     string
	queueUrl string
	v        any
	opts     []option.Producer
	wantErr  bool
	async    bool
}

type testConsumer[Body, MessageAttributes any] struct {
	name     string
	queueUrl string
	handler  HandlerConsumerFunc[Body, MessageAttributes]
	opts     []option.Consumer
	wantErr  bool
	async    bool
}

type testSimpleConsumer[Body any] struct {
	name     string
	queueUrl string
	handler  HandlerSimpleConsumerFunc[Body]
	opts     []option.Consumer
	wantErr  bool
	async    bool
}

type testCreateQueue struct {
	name      string
	queueName string
	opts      []option.CreateQueue
	wantErr   bool
}

type testDeleteQueue struct {
	name     string
	queueUrl string
	opts     []option.Default
	wantErr  bool
}

type testTagQueue struct {
	name    string
	input   TagQueueInput
	opts    []option.Default
	wantErr bool
}

type testSetAttributesQueue struct {
	name    string
	input   SetQueueAttributesInput
	opts    []option.Default
	wantErr bool
}

type testUntagQueue struct {
	name    string
	input   UntagQueueInput
	opts    []option.Default
	wantErr bool
}

type testPurgeQueue struct {
	name     string
	queueUrl string
	opts     []option.Default
	wantErr  bool
}

type testGetQueueUrl struct {
	name    string
	input   GetQueueUrlInput
	opts    []option.Default
	wantErr bool
}

type testGetQueueAttributes struct {
	name    string
	input   GetQueueAttributesInput
	opts    []option.Default
	wantErr bool
}

type testListQueueTags struct {
	name     string
	queueUrl string
	opts     []option.Default
	wantErr  bool
}

type testListDeadLetterSourceQueues struct {
	name     string
	queueUrl string
	opts     []option.ListDeadLetterSourceQueues
	wantErr  bool
}

type testListQueue struct {
	name    string
	opts    []option.ListQueues
	wantErr bool
}

type testDeleteMessage struct {
	name          string
	queueUrl      string
	receiptHandle string
	opts          []option.Default
	wantErr       bool
}

type testDeleteMessageBatch struct {
	name    string
	input   DeleteMessageBatchInput
	opts    []option.Default
	wantErr bool
}

type testChangeMessageVisibility struct {
	name    string
	input   ChangeMessageVisibilityInput
	opts    []option.Default
	wantErr bool
}

type testChangeMessageVisibilityBatch struct {
	name    string
	input   ChangeMessageVisibilityBatchInput
	opts    []option.Default
	wantErr bool
}

type testStartMessageMoveTask struct {
	name    string
	input   StartMessageMoveTaskInput
	opts    []option.Default
	wantErr bool
}

type testCancelMessageMoveTask struct {
	name       string
	taskHandle string
	opts       []option.Default
	wantErr    bool
}
type testListMessageMoveTask struct {
	name      string
	sourceArn string
	opts      []option.ListMessageMoveTasks
	wantErr   bool
}

type test struct {
	Name        string    `json:"name,omitempty"`
	BirthDate   time.Time `json:"birthDate,omitempty"`
	Emails      []string  `json:"emails,omitempty"`
	Bank        bank      `json:"bank,omitempty"`
	Map         map[string]any
	EmptyString string `json:"emptyString"`
	PointerBank *bank  `json:"pointerBank,omitempty"`
	Any         any
	HideString  string `json:"-"`
}

type messageAttTest struct {
	Name        string  `json:"account,omitempty"`
	Text        string  `json:"text,omitempty"`
	Balance     float64 `json:"balance,omitempty"`
	Bool        bool    `json:"bool"`
	Int         int     `json:"int"`
	SubStruct   test    `json:"subStruct,omitempty"`
	PointerTest *test   `json:"pointerBank,omitempty"`
	Map         map[string]any
	Any         any
	EmptyString string `json:"emptyString,omitempty"`
	HideString  string `json:"-"`
}

type bank struct {
	Account string  `json:"account,omitempty"`
	Digits  string  `json:"digits,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}

func TestMain(t *testing.M) {
	t.Run()
	deleteQueueCreateTest()
	purgeQueues()
	cancelMessageMoveTaskTest()
}

func initListTestProducer() []testProducer {
	msgAttTest := initMessageAttTest()
	return []testProducer{
		{
			name:     "valid request",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(&msgAttTest),
				option.NewProducer().SetDelaySeconds(2 * time.Second),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
			async:   false,
		},
		{
			name:     "valid request fifo",
			queueUrl: os.Getenv(SqsQueueTestFifoUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetDebugMode(true),
				option.NewProducer().SetOptionHttp(option.Http{}),
				option.NewProducer().SetMessageAttributes(initMessageAttTest()),
				option.NewProducer().SetMessageGroupId("group"),
				option.NewProducer().SetMessageDeduplicationId("deduplication"),
			},
			wantErr: false,
			async:   false,
		},
		{
			name:     "valid request async",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initTestMap()),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
			async:   true,
		},
		{
			name:     "invalid queue URL",
			queueUrl: "https://google.com",
			v:        initTestStruct(),
			wantErr:  true,
		},
		{
			name:     "invalid message attributes",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes("test message string"),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: true,
		},
		{
			name:     "invalid system message attributes",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageSystemAttributes(option.MessageSystemAttributes{
					AWSTraceHeader: "test test",
				}),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: true,
		},
		{
			name:     "empty message attributes",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(struct{}{}),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "invalid map message attributes",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initInvalidTestMap()),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "invalid message body",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			v:        "",
			wantErr:  true,
		},
	}
}

func initListTestConsumer[Body, MessageAttributes any]() []testConsumer[Body, MessageAttributes] {
	return []testConsumer[Body, MessageAttributes]{
		{
			name:     "success",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success fifo",
			queueUrl: os.Getenv(SqsQueueTestFifoUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			handler:  initHandleConsumerWithErr[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "success empty",
			queueUrl: os.Getenv(SqsQueueTestEmptyUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    false,
			wantErr:  false,
		},
		{
			name:     "failed parse body",
			queueUrl: os.Getenv(SqsQueueTestStringUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    false,
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com/",
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerWithErr(),
			wantErr:  true,
		},
	}
}

func initListTestSimpleConsumer[Body any]() []testSimpleConsumer[Body] {
	return []testSimpleConsumer[Body]{
		{
			name:     "success",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			handler:  initSimpleHandleConsumerWithErr[Body],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "failed parse body",
			queueUrl: os.Getenv(SqsQueueTestStringUrl),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			async:    false,
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com/",
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerWithErr(),
			wantErr:  true,
		},
	}
}

func initListTestCreateQueue() []testCreateQueue {
	return []testCreateQueue{
		{
			name:      "success",
			queueName: getSqsCreateQueueTest(),
			opts:      initOptionsCreateQueue(),
			wantErr:   false,
		},
		{
			name:      "failed",
			queueName: "",
			opts:      initOptionsCreateQueue(),
			wantErr:   true,
		},
	}
}

func initListTestDeleteQueue() []testDeleteQueue {
	return []testDeleteQueue{
		{
			name:     "success",
			queueUrl: os.Getenv(SqsQueueCreateTestUrl),
			opts:     initOptionsDefault(),
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com",
			opts:     initOptionsDefault(),
			wantErr:  true,
		},
	}
}

func initListTestTagQueue() []testTagQueue {
	return []testTagQueue{
		{
			name:    "success",
			input:   initTagQueueInput(),
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name:    "failed",
			input:   TagQueueInput{},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestSetAttributesQueue() []testSetAttributesQueue {
	return []testSetAttributesQueue{
		{
			name:    "success",
			input:   initSetQueueAttributesInput(),
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name:    "failed",
			input:   SetQueueAttributesInput{},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestUntagQueue() []testUntagQueue {
	return []testUntagQueue{
		{
			name:    "success",
			input:   initUntagQueueInput(),
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name:    "failed",
			input:   UntagQueueInput{},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestPurgeQueue() []testPurgeQueue {
	return []testPurgeQueue{
		{
			name:     "success",
			queueUrl: os.Getenv(SqsQueueCreateTestUrl),
			opts:     initOptionsDefault(),
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com",
			opts:     initOptionsDefault(),
			wantErr:  true,
		},
	}
}

func initListTestGetQueueUrl() []testGetQueueUrl {
	return []testGetQueueUrl{
		{
			name:    "success",
			input:   initGetQueueUrlInput(),
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name:    "failed",
			input:   GetQueueUrlInput{},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestGetQueueAttributes() []testGetQueueAttributes {
	return []testGetQueueAttributes{
		{
			name:    "success",
			input:   initGetQueueAttributesInput(),
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name:    "failed",
			input:   GetQueueAttributesInput{},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestListQueueTags() []testListQueueTags {
	return []testListQueueTags{
		{
			name:     "success",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			opts:     initOptionsDefault(),
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com",
			opts:     initOptionsDefault(),
			wantErr:  true,
		},
	}
}

func initListTestListDeadLetterSourceQueues() []testListDeadLetterSourceQueues {
	return []testListDeadLetterSourceQueues{
		{
			name:     "success",
			queueUrl: os.Getenv(SqsQueueTestUrl),
			opts:     initOptionsListDeadLetterSourceQueues(),
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com",
			opts:     initOptionsListDeadLetterSourceQueuesWithErr(),
			wantErr:  true,
		},
	}
}

func initListTestListQueue() []testListQueue {
	return []testListQueue{
		{
			name:    "success",
			opts:    initOptionsListQueues(),
			wantErr: false,
		},
		{
			name:    "failed",
			opts:    initOptionsListQueuesWithErr(),
			wantErr: true,
		},
	}
}

func initListTestDeleteMessage() []testDeleteMessage {
	return []testDeleteMessage{
		{
			name:          "success",
			queueUrl:      os.Getenv(SqsQueueTestStringUrl),
			receiptHandle: os.Getenv(SqsMessageReceiptHandle),
			opts:          initOptionsDefault(),
			wantErr:       false,
		},
		{
			name:          "failed",
			queueUrl:      "https://google.com",
			receiptHandle: "",
			opts:          initOptionsDefault(),
			wantErr:       true,
		},
	}
}

func initListTestDeleteMessageBatch() []testDeleteMessageBatch {
	return []testDeleteMessageBatch{
		{
			name: "success",
			input: DeleteMessageBatchInput{
				QueueUrl: os.Getenv(SqsQueueTestStringUrl),
				Entries: []DeleteMessageBatchRequestEntry{
					{
						Id:            os.Getenv(SqsMessageId),
						ReceiptHandle: os.Getenv(SqsMessageReceiptHandle),
					},
				},
			},
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name: "failed",
			input: DeleteMessageBatchInput{
				QueueUrl: "https://google.com/",
				Entries:  []DeleteMessageBatchRequestEntry{},
			},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestChangeMessageVisibility() []testChangeMessageVisibility {
	return []testChangeMessageVisibility{
		{
			name: "success",
			input: ChangeMessageVisibilityInput{
				QueueUrl:          os.Getenv(SqsQueueTestStringUrl),
				ReceiptHandle:     os.Getenv(SqsMessageReceiptHandle),
				VisibilityTimeout: 1,
			},
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name: "failed",
			input: ChangeMessageVisibilityInput{
				QueueUrl: "https://google.com/",
			},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestChangeMessageVisibilityBatch() []testChangeMessageVisibilityBatch {
	return []testChangeMessageVisibilityBatch{
		{
			name: "success",
			input: ChangeMessageVisibilityBatchInput{
				QueueUrl: os.Getenv(SqsQueueTestStringUrl),
				Entries: []ChangeMessageVisibilityBatchRequestEntry{
					{
						Id:                os.Getenv(SqsMessageId),
						ReceiptHandle:     os.Getenv(SqsMessageReceiptHandle),
						VisibilityTimeout: 3,
					},
				},
			},
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name: "failed",
			input: ChangeMessageVisibilityBatchInput{
				QueueUrl: "https://google.com/",
			},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestStartMessageMoveTask() []testStartMessageMoveTask {
	return []testStartMessageMoveTask{
		{
			name:    "success",
			input:   initStartMessageMoveTaskInput(),
			opts:    initOptionsDefault(),
			wantErr: false,
		},
		{
			name: "failed",
			input: StartMessageMoveTaskInput{
				SourceArn: "https://google.com/",
			},
			opts:    initOptionsDefault(),
			wantErr: true,
		},
	}
}

func initListTestCancelMessageMoveTask() []testCancelMessageMoveTask {
	return []testCancelMessageMoveTask{
		{
			name:       "success",
			taskHandle: os.Getenv(SqsTaskHandle),
			opts:       initOptionsDefault(),
			wantErr:    false,
		},
		{
			name:       "failed",
			taskHandle: "",
			opts:       initOptionsDefault(),
			wantErr:    true,
		},
	}
}

func initListTestListMessageMoveTask() []testListMessageMoveTask {
	return []testListMessageMoveTask{
		{
			name:      "success",
			sourceArn: os.Getenv(SqsQueueTestStringArn),
			opts:      initOptionsListMessageMoveTasks(),
			wantErr:   false,
		},
		{
			name:      "failed",
			sourceArn: "",
			opts:      initOptionsListMessageMoveTasksWithErr(),
			wantErr:   true,
		},
	}
}

func initTestStruct() test {
	b := bank{
		Account: "123456",
		Digits:  "2",
		Balance: 200.12,
	}
	return test{
		Name:        "Test Name",
		BirthDate:   time.Now(),
		Emails:      []string{"test@gmail.com", "gabriel@gmail.com", "gabriel.test@gmail.com"},
		Bank:        b,
		Map:         initTestMap(),
		EmptyString: "",
		PointerBank: &b,
		HideString:  "",
	}
}

func initMessageAttTest() messageAttTest {
	t := initTestStruct()
	return messageAttTest{
		Name:        "Name test producer",
		Text:        "Text field",
		Balance:     10.23,
		Bool:        true,
		Int:         3,
		SubStruct:   t,
		PointerTest: &t,
		Map:         initTestMap(),
		HideString:  "hide test",
	}
}

func initTestMap() map[string]any {
	return map[string]any{
		"name":        "Name test producer",
		"text":        "Text field",
		"balance":     10.23,
		"bool":        true,
		"int":         3,
		"emptyString": "",
		"nil":         nil,
	}
}

func initInvalidTestMap() map[string]any {
	var invalid string
	var invalidValue any
	return map[string]any{
		invalid:        "invalid",
		"valueInvalid": invalidValue,
	}
}

func initHandleConsumer[Body, MessageAttributes any](ctx *Context[Body, MessageAttributes]) error {
	logger.Debug("ctx:", ctx)
	return nil
}

func initSimpleHandleConsumer[Body any](ctx *SimpleContext[Body]) error {
	logger.Debug("ctx:", ctx)
	return nil
}

func initHandleConsumerWithErr[Body, MessageAttributes any](ctx *Context[Body, MessageAttributes]) error {
	logger.Debug("ctx:", ctx)
	return initErrorConsumer()
}

func initSimpleHandleConsumerWithErr[Body any](ctx *SimpleContext[Body]) error {
	logger.Debug("ctx:", ctx)
	return initErrorConsumer()
}

func initErrorConsumer() error {
	return errors.New("test error message")
}

func initTagQueueInput() TagQueueInput {
	return TagQueueInput{
		QueueUrl: os.Getenv(SqsQueueCreateTestUrl),
		Tags:     initTagsQueue(),
	}
}

func initSetQueueAttributesInput() SetQueueAttributesInput {
	return SetQueueAttributesInput{
		QueueUrl:   os.Getenv(SqsQueueCreateTestUrl),
		Attributes: initAttributesQueue("900"),
	}
}

func initUntagQueueInput() UntagQueueInput {
	return UntagQueueInput{
		QueueUrl: os.Getenv(SqsQueueCreateTestUrl),
		TagKeys:  []string{"tag-test"},
	}
}

func initGetQueueUrlInput() GetQueueUrlInput {
	return GetQueueUrlInput{
		QueueName:              os.Getenv(SqsQueueTestName),
		QueueOwnerAWSAccountId: nil,
	}
}

func initGetQueueAttributesInput() GetQueueAttributesInput {
	return GetQueueAttributesInput{
		QueueUrl:       os.Getenv(SqsQueueTestUrl),
		AttributeNames: nil,
	}
}

func initTagsQueue() map[string]string {
	return map[string]string{"tag-test": "test", "go": "test"}
}

func initAttributesQueue(delaySeconds string) map[string]string {
	return map[string]string{"DelaySeconds": delaySeconds}
}

func initStartMessageMoveTaskInput() StartMessageMoveTaskInput {
	maxNumberOfMessagesPerSecond := int32(1)
	return StartMessageMoveTaskInput{
		SourceArn:                    os.Getenv(SqsQueueTestDlqArn),
		DestinationArn:               nil,
		MaxNumberOfMessagesPerSecond: &maxNumberOfMessagesPerSecond,
	}
}

func initOptionsConsumerDefault() []option.Consumer {
	return []option.Consumer{
		option.NewConsumer().SetDebugMode(true),
		option.NewConsumer().SetOptionHttp(option.Http{}),
		option.NewConsumer().SetConsumerMessageTimeout(5 * time.Second),
		option.NewConsumer().SetDelayQueryLoop(4 * time.Second),
		option.NewConsumer().SetDeleteMessageProcessedSuccess(true),
		option.NewConsumer().SetMaxNumberOfMessages(1),
		option.NewConsumer().SetReceiveRequestAttemptId(""),
		option.NewConsumer().SetVisibilityTimeout(1 * time.Second),
		option.NewConsumer().SetWaitTimeSeconds(1 * time.Second),
	}
}

func initOptionsConsumerWithErr() []option.Consumer {
	return []option.Consumer{
		option.NewConsumer().SetDebugMode(true),
		option.NewConsumer().SetOptionHttp(option.Http{}),
		option.NewConsumer().SetConsumerMessageTimeout(0),
		option.NewConsumer().SetDelayQueryLoop(0),
		option.NewConsumer().SetDeleteMessageProcessedSuccess(true),
		option.NewConsumer().SetMaxNumberOfMessages(0),
		option.NewConsumer().SetReceiveRequestAttemptId("test"),
		option.NewConsumer().SetVisibilityTimeout(0),
		option.NewConsumer().SetWaitTimeSeconds(0),
	}
}

func initOptionsCreateQueue() []option.CreateQueue {
	return []option.CreateQueue{
		option.NewCreateQueue().SetDebugMode(true),
		option.NewCreateQueue().SetOptionHttp(option.Http{}),
		option.NewCreateQueue().SetAttributes(initAttributesQueue("800")),
		option.NewCreateQueue().SetTags(initTagsQueue()),
	}
}

func initOptionsListQueues() []option.ListQueues {
	return []option.ListQueues{
		option.NewListQueue().SetDebugMode(true),
		option.NewListQueue().SetOptionHttp(option.Http{}),
		option.NewListQueue().SetMaxResults(0),
		option.NewListQueue().SetNextToken(""),
		option.NewListQueue().SetQueueNamePrefix(""),
	}
}

func initOptionsListQueuesWithErr() []option.ListQueues {
	return []option.ListQueues{
		option.NewListQueue().SetDebugMode(true),
		option.NewListQueue().SetOptionHttp(option.Http{}),
		option.NewListQueue().SetMaxResults(10),
		option.NewListQueue().SetNextToken("test"),
		option.NewListQueue().SetQueueNamePrefix("test"),
	}
}

func initOptionsListDeadLetterSourceQueues() []option.ListDeadLetterSourceQueues {
	return []option.ListDeadLetterSourceQueues{
		option.NewListDeadLetterSourceQueues().SetDebugMode(true),
		option.NewListDeadLetterSourceQueues().SetOptionHttp(option.Http{}),
		option.NewListDeadLetterSourceQueues().SetMaxResults(0),
		option.NewListDeadLetterSourceQueues().SetNextToken(""),
	}
}

func initOptionsListDeadLetterSourceQueuesWithErr() []option.ListDeadLetterSourceQueues {
	return []option.ListDeadLetterSourceQueues{
		option.NewListDeadLetterSourceQueues().SetDebugMode(true),
		option.NewListDeadLetterSourceQueues().SetOptionHttp(option.Http{}),
		option.NewListDeadLetterSourceQueues().SetMaxResults(10),
		option.NewListDeadLetterSourceQueues().SetNextToken("token token"),
	}
}

func initOptionsDefault() []option.Default {
	return []option.Default{
		option.NewDefault().SetDebugMode(true),
		option.NewDefault().SetOptionHttp(option.Http{}),
	}
}

func initOptionsListMessageMoveTasks() []option.ListMessageMoveTasks {
	return []option.ListMessageMoveTasks{
		option.NewListMessageMoveTasks().SetDebugMode(true),
		option.NewListMessageMoveTasks().SetOptionHttp(option.Http{}),
		option.NewListMessageMoveTasks().SetMaxResults(10),
	}
}

func initOptionsListMessageMoveTasksWithErr() []option.ListMessageMoveTasks {
	return []option.ListMessageMoveTasks{
		option.NewListMessageMoveTasks().SetDebugMode(true),
		option.NewListMessageMoveTasks().SetOptionHttp(option.Http{}),
		option.NewListMessageMoveTasks().SetMaxResults(0),
	}
}

func initMessageString() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewProducer().SetMessageAttributes(initMessageAttTest())
	output, err := SendMessage(ctx, os.Getenv(SqsQueueTestStringUrl), "test body string", opt)
	if err != nil {
		logger.Error("error send message:", err)
		return
	}
	_ = os.Setenv(SqsMessageId, *output.MessageId)
}

func initMessageStruct(queueUrl string) {
	if queueUrl != os.Getenv(SqsQueueTestUrl) && queueUrl != os.Getenv(SqsQueueTestFifoUrl) {
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewProducer().SetMessageAttributes(initMessageAttTest())
	_, err := SendMessage(ctx, queueUrl, initTestStruct(), opt)
	if err != nil {
		logger.Error("error send message:", err)
	}
}

func initMessageReceiptHandle() {
	initMessageString()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	ctxInterrupt = ctx
	opt := option.NewConsumer().
		SetDebugMode(true).
		SetMaxNumberOfMessages(1).
		SetVisibilityTimeout(0).
		SetDeleteMessageProcessedSuccess(false)
	SimpleReceiveMessage[any](os.Getenv(SqsQueueTestStringUrl), func(ctx *SimpleContext[any]) error {
		_ = os.Setenv(SqsMessageReceiptHandle, ctx.Message.ReceiptHandle)
		cancel()
		return nil
	}, opt)
}

func initQueueCreateTest() {
	deleteQueueCreateTest()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewCreateQueue().SetTags(initTagsQueue())
	name := getSqsCreateQueueTest()
	output, err := CreateQueue(ctx, name, opt)
	if err != nil {
		logger.Error("error create test-queue:", err)
		return
	}
	_ = os.Setenv(SqsQueueCreateTestName, name)
	_ = os.Setenv(SqsQueueCreateTestUrl, *output.QueueUrl)
}

func initStartMessageMoveTask() {
	cancelMessageMoveTaskTest()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	output, err := StartMessageMoveTask(ctx, initStartMessageMoveTaskInput())
	if err != nil {
		logger.Error("error start message move task:", err)
		return
	}
	_ = os.Setenv(SqsTaskHandle, *output.TaskHandle)
}

func deleteQueueCreateTest() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	url := os.Getenv(SqsQueueCreateTestUrl)
	if len(url) == 0 {
		return
	}
	_, _ = DeleteQueue(ctx, os.Getenv(SqsQueueCreateTestUrl))
}

func purgeQueues() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	urlString := os.Getenv(SqsQueueTestStringUrl)
	urlFifo := os.Getenv(SqsQueueTestFifoUrl)
	if len(urlString) != 0 {
		_, _ = PurgeQueue(ctx, os.Getenv(SqsQueueTestStringUrl))
	}
	if len(urlFifo) != 0 {
		_, _ = PurgeQueue(ctx, os.Getenv(SqsQueueTestFifoUrl))
	}
}

func cancelMessageMoveTaskTest() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	taskHandle := os.Getenv(SqsTaskHandle)
	if len(taskHandle) == 0 {
		return
	}
	_, _ = CancelMessageMoveTask(ctx, os.Getenv(SqsTaskHandle))

}

func getSqsCreateQueueTest() string {
	return "go-aws-sqs-test-create-queue-" + strconv.Itoa(int(time.Now().UnixNano()))
}
