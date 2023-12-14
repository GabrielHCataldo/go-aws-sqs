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

const SqsQueueCreateTestName = "SQS_QUEUE_CREATE_TEST_NAME"
const SqsQueueCreateTestUrl = "SQS_QUEUE_CREATE_TEST_URL"

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
}

func initListTestProducer() []testProducer {
	return []testProducer{
		{
			name:     "valid request",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initMessageAttTest()),
				option.NewProducer().SetDelaySeconds(2 * time.Second),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
			async:   false,
		},
		{
			name:     "valid request fifo",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_FIFO"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes("test message string"),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: true,
		},
		{
			name:     "invalid system message attributes",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(messageAttTest{}),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "invalid map message attributes",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initInvalidTestMap()),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "invalid message body",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        "",
			wantErr:  true,
		},
	}
}

func initListTestConsumer[Body, MessageAttributes any]() []testConsumer[Body, MessageAttributes] {
	return []testConsumer[Body, MessageAttributes]{
		{
			name:     "success",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success FIFO",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_FIFO"),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initHandleConsumerWithErr[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "success empty",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_EMPTY"),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumerDefault(),
			async:    false,
			wantErr:  false,
		},
		{
			name:     "failed parse body",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_STRING"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initSimpleHandleConsumerWithErr[Body],
			opts:     initOptionsConsumerDefault(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initSimpleHandleConsumer[Body],
			opts:     initOptionsConsumerDefault(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "failed parse body",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_STRING"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST_EMPTY"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
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
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
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
		QueueName:              os.Getenv("SQS_QUEUE_TEST_NAME"),
		QueueOwnerAWSAccountId: nil,
	}
}

func initGetQueueAttributesInput() GetQueueAttributesInput {
	return GetQueueAttributesInput{
		QueueUrl:       os.Getenv("SQS_QUEUE_TEST"),
		AttributeNames: nil,
	}
}

func initTagsQueue() map[string]string {
	return map[string]string{"tag-test": "test", "go": "test"}
}

func initAttributesQueue(delaySeconds string) map[string]string {
	return map[string]string{"DelaySeconds": delaySeconds}
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

func initMessageString(queueUrl string) {
	if queueUrl != os.Getenv("SQS_QUEUE_TEST_STRING") {
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewProducer().SetMessageAttributes(initMessageAttTest())
	_, err := SendMessage(ctx, queueUrl, "test body string", opt)
	if err != nil {
		logger.Error("error send message:", err)
	}
}

func initMessageStruct(queueUrl string) {
	if queueUrl != os.Getenv("SQS_QUEUE_TEST") && queueUrl != os.Getenv("SQS_QUEUE_TEST_FIFO") {
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewProducer().SetMessageAttributes(initMessageAttTest())
	if queueUrl == os.Getenv("SQS_QUEUE_TEST_FIFO") {
		opt.SetMessageGroupId("group")
		opt.SetMessageDeduplicationId("deduplication")
	}
	_, err := SendMessage(ctx, queueUrl, initTestStruct(), opt)
	if err != nil {
		logger.Error("error send message:", err)
	}
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

func deleteQueueCreateTest() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	_, _ = DeleteQueue(ctx, SqsQueueCreateTestUrl)
}

func getSqsCreateQueueTest() string {
	return "go-aws-sqs-test-create-queue-" + strconv.Itoa(int(time.Now().UnixNano()))
}
