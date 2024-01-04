package main

import (
	"context"
	"github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
	"github.com/GabrielHCataldo/go-aws-sqs-template/sqs/option"
	"github.com/GabrielHCataldo/go-logger/logger"
	"os"
	"time"
)

type test struct {
	Name      string    `json:"name,omitempty"`
	BirthDate time.Time `json:"birthDate,omitempty"`
	Emails    []string  `json:"emails,omitempty"`
	Bank      bank      `json:"bank,omitempty"`
	Map       map[string]any
}

type bank struct {
	Account string  `json:"account,omitempty"`
	Digits  string  `json:"digits,omitempty"`
	Balance float64 `json:"balance,omitempty"`
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

func main() {
	simple()
	simpleAsync()
	structBody()
	mapBody()
	completeOptions()
}

func simple() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	body := "body test"
	message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_STRING_URL"), body)
	if err != nil {
		logger.Error("error send message:", err)
	} else {
		logger.Info("message sent successfully:", message)
	}
}

func simpleAsync() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	body := "body test"
	sqs.SendMessageAsync(ctx, os.Getenv("SQS_QUEUE_TEST_STRING_URL"), body)
}

func structBody() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	body := initTestStruct()
	message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_URL"), body)
	if err != nil {
		logger.Error("error send message:", err)
	} else {
		logger.Info("message sent successfully:", message)
	}
}

func mapBody() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	body := initTestMap()
	message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_STRING_URL"), body)
	if err != nil {
		logger.Error("error send message:", err)
	} else {
		logger.Info("message sent successfully:", message)
	}
}

func completeOptions() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	body := initTestStruct()
	opt := option.NewProducer().
		// HTTP communication customization options with AWS SQS
		SetHttpClient(option.HttpClient{}).
		// print logs (default: false)
		SetDebugMode(true).
		// delay to delay the availability of message processing (default: 0)
		SetDelaySeconds(5 * time.Second).
		// Message attributes, must be of type Map or Struct, other types are not acceptable.
		SetMessageAttributes(initMessageAttTest()).
		// The message system attribute to send
		SetMessageSystemAttributes(option.MessageSystemAttributes{}).
		// This parameter applies only to FIFO (first-in-first-out) queues. The token used for deduplication of sent messages.
		SetMessageDeduplicationId("").
		// This parameter applies only to FIFO (first-in-first-out) queues. The tag that specifies that a message belongs to a specific message group.
		SetMessageGroupId("")
	message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_URL"), body, opt)
	if err != nil {
		logger.Error("error send message:", err)
	} else {
		logger.Info("message sent successfully:", message)
	}
}

func initTestStruct() test {
	b := bank{
		Account: "123456",
		Digits:  "2",
		Balance: 200.12,
	}
	return test{
		Name:      "Test Name",
		BirthDate: time.Now(),
		Emails:    []string{"test@gmail.com", "gabriel@gmail.com", "gabriel.test@gmail.com"},
		Bank:      b,
		Map:       map[string]any{"int": 1, "bool": true, "float": 1.23, "string": "text test"},
	}
}

func initTestMap() map[string]any {
	return map[string]any{"int": 1, "bool": true, "float": 1.23, "string": "text test"}
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
