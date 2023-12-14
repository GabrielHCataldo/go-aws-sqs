package sqs

import (
	"context"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	for _, tt := range initListTestProducer() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			var err error
			if tt.async {
				SendMessageAsync(ctx, tt.queueUrl, tt.v, tt.opts...)
				select {
				case <-ctx.Done():
				}
			} else {
				_, err = SendMessage(ctx, tt.queueUrl, tt.v, tt.opts...)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
