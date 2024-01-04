package main

import (
	"github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
	"github.com/GabrielHCataldo/go-aws-sqs-template/sqs/option"
	"github.com/GabrielHCataldo/go-logger/logger"
	"os"
	"os/signal"
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
	simpleReceiveMessage()
	simpleReceiveMessageStruct()
	receiveMessage()
	completeOptions()
	simpleReceiveMessageAsync()
}

func simpleReceiveMessage() {
	sqs.SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST_STRING_URL"), handlerSimple)
}

func simpleReceiveMessageStruct() {
	sqs.SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST_URL"), handler)
}

func receiveMessage() {
	sqs.ReceiveMessage(os.Getenv("SQS_QUEUE_TEST_STRING_URL"), handlerReceiveMessage)
}

func completeOptions() {
	opt := option.NewConsumer().
		// HTTP communication customization options with AWS SQS
		SetHttpClient(option.HttpClient{}).
		// print logs (default: false)
		SetDebugMode(true).
		// If true remove the message from the queue after successfully processed (handler error return is null) (default: false)
		SetDeleteMessageProcessedSuccess(true).
		// Duration time to process the message, timeout applied in the past context. (default: 5 seconds)
		SetConsumerMessageTimeout(5 * time.Second).
		// Delay to run the next search for messages in the queue (default: 0)
		SetDelayQueryLoop(5 * time.Second).
		// The maximum number of messages to return. 1 a 10 (default: 10)
		SetMaxNumberOfMessages(10).
		// The maximum number of messages to return. 1 a 10 (default: 0)
		SetVisibilityTimeout(5 * time.Second).
		// The duration that the received messages are hidden from subsequent
		// retrieve requests after being retrieved by a ReceiveMessage request.
		SetReceiveRequestAttemptId("").
		// The duration for which the call waits for a message to arrive in
		// the queue before returning (default: 0)
		SetWaitTimeSeconds(1 * time.Second)
	sqs.SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST_URL"), handler, opt)
}

func simpleReceiveMessageAsync() {
	sqs.SimpleReceiveMessageAsync(os.Getenv("SQS_QUEUE_TEST_URL"), handler, option.NewConsumer().SetDebugMode(true))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		logger.Info("Stopped application!")
	}
}

func handler(ctx *sqs.SimpleContext[test]) error {
	logger.Debug("ctx simple body struct to process message:", ctx)
	return nil
}

func handlerSimple(ctx *sqs.SimpleContext[string]) error {
	logger.Debug("ctx simple to process message:", ctx)
	return nil
}

func handlerReceiveMessage(ctx *sqs.Context[string, messageAttTest]) error {
	logger.Debug("ctx to process message:", ctx)
	return nil
}
