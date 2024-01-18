package sqs

import (
	"context"
	"errors"
	"github.com/GabrielHCataldo/go-aws-sqs-template/sqs/option"
	"github.com/GabrielHCataldo/go-logger/logger"
	"os"
	"strconv"
	"testing"
	"time"
)

const sqsQueueTestName = "SQS_QUEUE_TEST_NAME"
const sqsQueueTestUrl = "SQS_QUEUE_TEST_URL"
const sqsQueueTestStringArn = "SQS_QUEUE_TEST_STRING_ARN"
const sqsQueueTestFifoUrl = "SQS_QUEUE_TEST_FIFO_URL"
const sqsQueueTestEmptyUrl = "SQS_QUEUE_TEST_EMPTY_URL"
const sqsQueueTestStringUrl = "SQS_QUEUE_TEST_STRING_URL"
const sqsQueueTestDlqArn = "SQS_QUEUE_TEST_DLQ_ARN"
const sqsQueueCreateTestName = "SQS_QUEUE_CREATE_TEST_NAME"
const sqsQueueCreateTestUrl = "SQS_QUEUE_CREATE_TEST_URL"
const sqsMessageId = "SQS_MESSAGE_ID"
const sqsMessageReceiptHandle = "SQS_MESSAGE_RECEIPT_HANDLE"
const sqsTaskHandle = "SQS_TASK_HANDLE"

type testProducer struct {
	name     string
	queueUrl string
	v        any
	opts     []*option.Producer
	wantErr  bool
	async    bool
}

type testConsumer[Body, MessageAttributes any] struct {
	name     string
	queueUrl string
	handler  HandlerConsumerFunc[Body, MessageAttributes]
	opts     []*option.Consumer
	wantErr  bool
	async    bool
}

type testSimpleConsumer[Body any] struct {
	name     string
	queueUrl string
	handler  HandlerSimpleConsumerFunc[Body]
	opts     []*option.Consumer
	wantErr  bool
	async    bool
}

type testCreateQueue struct {
	name      string
	queueName string
	opts      []*option.CreateQueue
	wantErr   bool
}

type testDeleteQueue struct {
	name     string
	queueUrl string
	opts     []*option.Default
	wantErr  bool
}

type testTagQueue struct {
	name    string
	input   TagQueueInput
	opts    []*option.Default
	wantErr bool
}

type testSetAttributesQueue struct {
	name    string
	input   SetQueueAttributesInput
	opts    []*option.Default
	wantErr bool
}

type testUntagQueue struct {
	name    string
	input   UntagQueueInput
	opts    []*option.Default
	wantErr bool
}

type testPurgeQueue struct {
	name     string
	queueUrl string
	opts     []*option.Default
	wantErr  bool
}

type testGetQueueUrl struct {
	name    string
	input   GetQueueUrlInput
	opts    []*option.Default
	wantErr bool
}

type testGetQueueAttributes struct {
	name    string
	input   GetQueueAttributesInput
	opts    []*option.Default
	wantErr bool
}

type testListQueueTags struct {
	name     string
	queueUrl string
	opts     []*option.Default
	wantErr  bool
}

type testListDeadLetterSourceQueues struct {
	name     string
	queueUrl string
	opts     []*option.ListDeadLetterSourceQueues
	wantErr  bool
}

type testListQueue struct {
	name    string
	opts    []*option.ListQueues
	wantErr bool
}

type testDeleteMessage struct {
	name          string
	queueUrl      string
	receiptHandle string
	opts          []*option.Default
	wantErr       bool
}

type testDeleteMessageBatch struct {
	name    string
	input   DeleteMessageBatchInput
	opts    []*option.Default
	wantErr bool
}

type testChangeMessageVisibility struct {
	name    string
	input   ChangeMessageVisibilityInput
	opts    []*option.Default
	wantErr bool
}

type testChangeMessageVisibilityBatch struct {
	name    string
	input   ChangeMessageVisibilityBatchInput
	opts    []*option.Default
	wantErr bool
}

type testStartMessageMoveTask struct {
	name    string
	input   StartMessageMoveTaskInput
	opts    []*option.Default
	wantErr bool
}

type testCancelMessageMoveTask struct {
	name       string
	taskHandle string
	opts       []*option.Default
	wantErr    bool
}
type testListMessageMoveTask struct {
	name      string
	sourceArn string
	opts      []*option.ListMessageMoveTasks
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
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
				option.NewProducer().SetMessageAttributes(&msgAttTest),
				option.NewProducer().SetDelaySeconds(2 * time.Second),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
			async:   false,
		},
		{
			name:     "valid request fifo",
			queueUrl: os.Getenv(sqsQueueTestFifoUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
				option.NewProducer().SetDebugMode(true),
				option.NewProducer().SetHttpClient(option.HttpClient{}),
				option.NewProducer().SetMessageAttributes(initMessageAttTest()),
				option.NewProducer().SetMessageGroupId("group"),
				option.NewProducer().SetMessageDeduplicationId("deduplication"),
			},
			wantErr: false,
			async:   false,
		},
		{
			name:     "valid request async",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
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
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
				option.NewProducer().SetMessageAttributes("test message string"),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: true,
		},
		{
			name:     "invalid system message attributes",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
				option.NewProducer().SetMessageSystemAttributes(option.MessageSystemAttributes{
					AWSTraceHeader: "test test",
				}),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: true,
		},
		{
			name:     "empty message attributes",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
				option.NewProducer().SetMessageAttributes(struct{}{}),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "invalid map message attributes",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        initTestStruct(),
			opts: []*option.Producer{
				option.NewProducer().SetMessageAttributes(initInvalidTestMap()),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "invalid message body",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			v:        "",
			wantErr:  true,
		},
	}
}

func initListTestConsumer[Body, MessageAttributes any]() []testConsumer[Body, MessageAttributes] {
	return []testConsumer[Body, MessageAttributes]{
		{
			name:     "success",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success fifo",
			queueUrl: os.Getenv(sqsQueueTestFifoUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initHandleConsumerWithErr[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "success empty",
			queueUrl: os.Getenv(sqsQueueTestEmptyUrl),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    false,
			wantErr:  false,
		},
		{
			name:     "failed parse body",
			queueUrl: os.Getenv(sqsQueueTestStringUrl),
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
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initSimpleHandleConsumerWithErr[Body],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "failed parse body",
			queueUrl: os.Getenv(sqsQueueTestStringUrl),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			async:    false,
			wantErr:  false,
		},
		{
			name:     "failed timeout",
			queueUrl: os.Getenv(sqsQueueTestUrl),
			handler:  initSimpleHandleConsumerTimeout[Body],
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
			queueUrl: os.Getenv(sqsQueueCreateTestUrl),
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
			queueUrl: os.Getenv(sqsQueueCreateTestUrl),
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
			queueUrl: os.Getenv(sqsQueueTestUrl),
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
			queueUrl: os.Getenv(sqsQueueTestUrl),
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
			queueUrl:      os.Getenv(sqsQueueTestStringUrl),
			receiptHandle: os.Getenv(sqsMessageReceiptHandle),
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
				QueueUrl: os.Getenv(sqsQueueTestStringUrl),
				Entries: []DeleteMessageBatchRequestEntry{
					{
						Id:            os.Getenv(sqsMessageId),
						ReceiptHandle: os.Getenv(sqsMessageReceiptHandle),
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
				QueueUrl:          os.Getenv(sqsQueueTestStringUrl),
				ReceiptHandle:     os.Getenv(sqsMessageReceiptHandle),
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
				QueueUrl: os.Getenv(sqsQueueTestStringUrl),
				Entries: []ChangeMessageVisibilityBatchRequestEntry{
					{
						Id:                os.Getenv(sqsMessageId),
						ReceiptHandle:     os.Getenv(sqsMessageReceiptHandle),
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
			taskHandle: os.Getenv(sqsTaskHandle),
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
			sourceArn: os.Getenv(sqsQueueTestStringArn),
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

func initSimpleHandleConsumerTimeout[Body any](ctx *SimpleContext[Body]) error {
	logger.Debug("ctx:", ctx)
	time.Sleep(30 * time.Second)
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
		QueueUrl: os.Getenv(sqsQueueCreateTestUrl),
		Tags:     initTagsQueue(),
	}
}

func initSetQueueAttributesInput() SetQueueAttributesInput {
	return SetQueueAttributesInput{
		QueueUrl:   os.Getenv(sqsQueueCreateTestUrl),
		Attributes: initAttributesQueue("900"),
	}
}

func initUntagQueueInput() UntagQueueInput {
	return UntagQueueInput{
		QueueUrl: os.Getenv(sqsQueueCreateTestUrl),
		TagKeys:  []string{"tag-test"},
	}
}

func initGetQueueUrlInput() GetQueueUrlInput {
	return GetQueueUrlInput{
		QueueName:              os.Getenv(sqsQueueTestName),
		QueueOwnerAWSAccountId: nil,
	}
}

func initGetQueueAttributesInput() GetQueueAttributesInput {
	return GetQueueAttributesInput{
		QueueUrl:       os.Getenv(sqsQueueTestUrl),
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
		SourceArn:                    os.Getenv(sqsQueueTestDlqArn),
		DestinationArn:               nil,
		MaxNumberOfMessagesPerSecond: &maxNumberOfMessagesPerSecond,
	}
}

func initOptionsConsumerDefault() []*option.Consumer {
	return []*option.Consumer{
		nil,
		option.NewConsumer().SetDebugMode(true),
		option.NewConsumer().SetHttpClient(option.HttpClient{}),
		option.NewConsumer().SetConsumerMessageTimeout(5 * time.Second),
		option.NewConsumer().SetDelayQueryLoop(4 * time.Second),
		option.NewConsumer().SetDeleteMessageProcessedSuccess(true),
		option.NewConsumer().SetMaxNumberOfMessages(1),
		option.NewConsumer().SetReceiveRequestAttemptId(""),
		option.NewConsumer().SetVisibilityTimeout(1 * time.Second),
		option.NewConsumer().SetWaitTimeSeconds(1 * time.Second),
	}
}

func initOptionsConsumerWithErr() []*option.Consumer {
	return []*option.Consumer{
		nil,
		option.NewConsumer().SetDebugMode(true),
		option.NewConsumer().SetHttpClient(option.HttpClient{}),
		option.NewConsumer().SetConsumerMessageTimeout(0),
		option.NewConsumer().SetDelayQueryLoop(0),
		option.NewConsumer().SetDeleteMessageProcessedSuccess(true),
		option.NewConsumer().SetMaxNumberOfMessages(0),
		option.NewConsumer().SetReceiveRequestAttemptId("test"),
		option.NewConsumer().SetVisibilityTimeout(0),
		option.NewConsumer().SetWaitTimeSeconds(0),
	}
}

func initOptionsCreateQueue() []*option.CreateQueue {
	return []*option.CreateQueue{
		nil,
		option.NewCreateQueue().SetDebugMode(true),
		option.NewCreateQueue().SetHttpClient(option.HttpClient{}),
		option.NewCreateQueue().SetAttributes(initAttributesQueue("800")),
		option.NewCreateQueue().SetTags(initTagsQueue()),
	}
}

func initOptionsListQueues() []*option.ListQueues {
	return []*option.ListQueues{
		nil,
		option.NewListQueue().SetDebugMode(true),
		option.NewListQueue().SetHttpClient(option.HttpClient{}),
		option.NewListQueue().SetMaxResults(0),
		option.NewListQueue().SetNextToken(""),
		option.NewListQueue().SetQueueNamePrefix(""),
	}
}

func initOptionsListQueuesWithErr() []*option.ListQueues {
	return []*option.ListQueues{
		nil,
		option.NewListQueue().SetDebugMode(true),
		option.NewListQueue().SetHttpClient(option.HttpClient{}),
		option.NewListQueue().SetMaxResults(10),
		option.NewListQueue().SetNextToken("test"),
		option.NewListQueue().SetQueueNamePrefix("test"),
	}
}

func initOptionsListDeadLetterSourceQueues() []*option.ListDeadLetterSourceQueues {
	return []*option.ListDeadLetterSourceQueues{
		nil,
		option.NewListDeadLetterSourceQueues().SetDebugMode(true),
		option.NewListDeadLetterSourceQueues().SetHttpClient(option.HttpClient{}),
		option.NewListDeadLetterSourceQueues().SetMaxResults(0),
		option.NewListDeadLetterSourceQueues().SetNextToken(""),
	}
}

func initOptionsListDeadLetterSourceQueuesWithErr() []*option.ListDeadLetterSourceQueues {
	return []*option.ListDeadLetterSourceQueues{
		nil,
		option.NewListDeadLetterSourceQueues().SetDebugMode(true),
		option.NewListDeadLetterSourceQueues().SetHttpClient(option.HttpClient{}),
		option.NewListDeadLetterSourceQueues().SetMaxResults(10),
		option.NewListDeadLetterSourceQueues().SetNextToken("token token"),
	}
}

func initOptionsDefault() []*option.Default {
	return []*option.Default{
		nil,
		option.NewDefault().SetDebugMode(true),
		option.NewDefault().SetHttpClient(option.HttpClient{}),
	}
}

func initOptionsListMessageMoveTasks() []*option.ListMessageMoveTasks {
	return []*option.ListMessageMoveTasks{
		nil,
		option.NewListMessageMoveTasks().SetDebugMode(true),
		option.NewListMessageMoveTasks().SetHttpClient(option.HttpClient{}),
		option.NewListMessageMoveTasks().SetMaxResults(10),
	}
}

func initOptionsListMessageMoveTasksWithErr() []*option.ListMessageMoveTasks {
	return []*option.ListMessageMoveTasks{
		nil,
		option.NewListMessageMoveTasks().SetDebugMode(true),
		option.NewListMessageMoveTasks().SetHttpClient(option.HttpClient{}),
		option.NewListMessageMoveTasks().SetMaxResults(0),
	}
}

func initMessageString() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewProducer().SetMessageAttributes(initMessageAttTest())
	output, err := SendMessage(ctx, os.Getenv(sqsQueueTestStringUrl), "test body string", opt)
	if err != nil {
		logger.Error("error send message:", err)
		return
	}
	_ = os.Setenv(sqsMessageId, *output.MessageId)
}

func initMessageStruct(queueUrl string) {
	if queueUrl != os.Getenv(sqsQueueTestUrl) && queueUrl != os.Getenv(sqsQueueTestFifoUrl) {
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
	SimpleReceiveMessage[any](os.Getenv(sqsQueueTestStringUrl), func(ctx *SimpleContext[any]) error {
		_ = os.Setenv(sqsMessageReceiptHandle, ctx.Message.ReceiptHandle)
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
	_ = os.Setenv(sqsQueueCreateTestName, name)
	_ = os.Setenv(sqsQueueCreateTestUrl, *output.QueueUrl)
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
	_ = os.Setenv(sqsTaskHandle, *output.TaskHandle)
}

func deleteQueueCreateTest() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	url := os.Getenv(sqsQueueCreateTestUrl)
	if len(url) == 0 {
		return
	}
	_, _ = DeleteQueue(ctx, os.Getenv(sqsQueueCreateTestUrl))
}

func purgeQueues() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	urlString := os.Getenv(sqsQueueTestStringUrl)
	urlFifo := os.Getenv(sqsQueueTestFifoUrl)
	if len(urlString) != 0 {
		_, _ = PurgeQueue(ctx, os.Getenv(sqsQueueTestStringUrl))
	}
	if len(urlFifo) != 0 {
		_, _ = PurgeQueue(ctx, os.Getenv(sqsQueueTestFifoUrl))
	}
}

func cancelMessageMoveTaskTest() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	taskHandle := os.Getenv(sqsTaskHandle)
	if len(taskHandle) == 0 {
		return
	}
	_, _ = CancelMessageMoveTask(ctx, os.Getenv(sqsTaskHandle))

}

func getSqsCreateQueueTest() string {
	return "go-aws-sqs-template-test-create-queue-" + strconv.Itoa(int(time.Now().UnixNano()))
}
