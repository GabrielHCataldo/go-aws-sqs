package sqs

import (
	"context"
	"testing"
	"time"
)

func TestReceiveMessage(t *testing.T) {
	for _, tt := range initListTestConsumer[test, messageAttTest]() {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantErr {
					t.Errorf("ReceiveMessage() error = %v, wantErr %v", r, tt.wantErr)
				}
			}()
			initMessageString()
			initMessageStruct(tt.queueUrl)
			d := 5 * time.Second
			if tt.name == "failed" {
				d = 20 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.TODO(), d)
			defer cancel()
			ctxInterrupt = ctx
			if tt.async {
				ReceiveMessageAsync(tt.queueUrl, tt.handler, tt.opts...)
				select {
				case <-ctx.Done():
				}
			} else {
				ReceiveMessage(tt.queueUrl, tt.handler, tt.opts...)
			}
		})
	}
}

func TestSimpleReceiveMessage(t *testing.T) {
	for _, tt := range initListTestSimpleConsumer[test]() {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantErr {
					t.Errorf("SimpleReceiveMessage() error = %v, wantErr %v", r, tt.wantErr)
				}
			}()
			initMessageString()
			initMessageStruct(tt.queueUrl)
			d := 5 * time.Second
			if tt.name == "failed" {
				d = 20 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.TODO(), d)
			defer cancel()
			ctxInterrupt = ctx
			if tt.async {
				SimpleReceiveMessageAsync(tt.queueUrl, tt.handler, tt.opts...)
				select {
				case <-ctx.Done():
				}
			} else {
				SimpleReceiveMessage(tt.queueUrl, tt.handler, tt.opts...)
			}
		})
	}
}
