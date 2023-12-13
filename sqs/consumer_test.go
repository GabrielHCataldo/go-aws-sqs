package sqs

import (
	"context"
	"testing"
	"time"
)

func TestReceiveMessage(t *testing.T) {
	for _, tt := range initListTestConsumer[test, messageAttTest]() {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				defer failWithoutPanic(t)
			}
			ctxProducer, cancelProducer := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancelProducer()
			_, _ = SendMessage(ctxProducer, tt.queueUrl, initTestStruct())
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			ctxUnitTest = ctx
			ReceiveMessage(tt.queueUrl, initHandleConsumer[test, messageAttTest], tt.opts...)
		})
	}
}

func TestSimpleReceiveMessage(t *testing.T) {
	//SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST"), initSimpleHandleConsumer, option.NewConsumer().SetDebugMode(true))
}
