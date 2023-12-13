package sqs

import (
	"errors"
	"github.com/GabrielHCataldo/go-logger/logger"
	"go-aws-sqs/sqs/option"
	"os"
	"os/signal"
	"testing"
	"time"
)

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

type testSimpleConsumer struct {
	name     string
	queueUrl string
	opts     []option.Producer
	wantErr  bool
	async    bool
}

type messageAttTest struct {
	Name      string         `json:"account,omitempty"`
	Text      string         `json:"text,omitempty"`
	Balance   float64        `json:"balance,omitempty"`
	Bool      bool           `json:"bool"`
	Int       int            `json:"int"`
	SubStruct test           `json:"subStruct,omitempty"`
	Map       map[string]any `json:"map,omitempty"`
}

type test struct {
	Name      string         `json:"name,omitempty"`
	BirthDate time.Time      `json:"birthDate,omitempty"`
	Emails    []string       `json:"emails,omitempty"`
	Bank      bank           `json:"bank,omitempty"`
	Map       map[string]any `json:"map,omitempty"`
}

type bank struct {
	Account string  `json:"account,omitempty"`
	Digits  string  `json:"digits,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}

func initListTestProducer() []testProducer {
	return []testProducer{
		{
			name:     "valid request",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initMessageAttTest()),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "valid request FIFO",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_FIFO"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initMessageAttTest()),
				option.NewProducer().SetDebugMode(true),
				option.NewProducer().SetMessageDeduplicationId("test"),
				option.NewProducer().SetMessageGroupId("test"),
			},
			wantErr: false,
		},
		{
			name:     "valid request Async",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        initTestStruct(),
			opts: []option.Producer{
				option.NewProducer().SetMessageAttributes(initMessageAttTest()),
				option.NewProducer().SetDebugMode(true),
			},
			wantErr: false,
			async:   true,
		},
		{
			name:     "invalid queue URL",
			queueUrl: "",
			v:        "test message",
			wantErr:  true,
		},
		{
			name:     "invalid message",
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
			opts:     initOptionsConsumer(),
			wantErr:  false,
		},
		{
			name:     "success error consumer",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initHandleConsumerWithErr[Body, MessageAttributes],
			opts:     initOptionsConsumer(),
			wantErr:  false,
		},
		{
			name:     "success async",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumer(),
			async:    true,
			wantErr:  false,
		},
		{
			name:     "failed",
			queueUrl: "https://google.com/",
			handler:  initHandleConsumer[Body, MessageAttributes],
			opts:     initOptionsConsumer(),
			wantErr:  true,
		},
	}
}

func initTestStruct() test {
	return test{
		Name:      "Test Name",
		BirthDate: time.Now(),
		Emails:    []string{"test@gmail.com", "gabriel@gmail.com", "gabriel.test@gmail.com"},
		Bank: bank{
			Account: "123456",
			Digits:  "2",
			Balance: 200.12,
		},
		Map: initTestMap(),
	}
}

func initMessageAttTest() messageAttTest {
	return messageAttTest{
		Name:      "Name test producer",
		Text:      "Text field",
		Balance:   10.23,
		Bool:      true,
		Int:       3,
		Map:       initTestMap(),
		SubStruct: initTestStruct(),
	}
}

func initTestMap() map[string]any {
	return map[string]any{
		"int":    1,
		"float":  12.2432,
		"bool":   true,
		"slice":  []int{12, 23, 53},
		"string": "text test",
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

func initOptionsConsumer() []option.Consumer {
	return []option.Consumer{
		option.NewConsumer().SetDebugMode(true),
		option.NewConsumer().SetConsumerMessageTimeout(5 * time.Second),
		option.NewConsumer().SetDelayQueryLoop(2 * time.Second),
		option.NewConsumer().SetDeleteMessageProcessedSuccess(true),
		option.NewConsumer().SetMaxNumberOfMessages(10),
		option.NewConsumer().SetReceiveRequestAttemptId(""),
		option.NewConsumer().SetVisibilityTimeout(0 * time.Second),
		option.NewConsumer().SetWaitTimeSeconds(0 * time.Second),
	}
}

func waitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		logger.Info("Stop test!")
	}
}

func failWithoutPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Fail()
	}
}
