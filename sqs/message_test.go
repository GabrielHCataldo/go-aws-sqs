package sqs

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDeleteMessage(t *testing.T) {
	initMessageReceiptHandle()
	for _, tt := range initListTestDeleteMessage() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := DeleteMessage(ctx, tt.queueUrl, tt.receiptHandle, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteMessageBatch(t *testing.T) {
	initMessageReceiptHandle()
	for _, tt := range initListTestDeleteMessageBatch() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := DeleteMessageBatch(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMessageBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeMessageVisibility(t *testing.T) {
	initMessageReceiptHandle()
	for _, tt := range initListTestChangeMessageVisibility() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := ChangeMessageVisibility(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeMessageVisibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeMessageVisibilityBatch(t *testing.T) {
	initMessageReceiptHandle()
	for _, tt := range initListTestChangeMessageVisibilityBatch() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := ChangeMessageVisibilityBatch(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeMessageVisibilityBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStartMessageMoveTask(t *testing.T) {
	initMessageString()
	for _, tt := range initListTestStartMessageMoveTask() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			output, err := StartMessageMoveTask(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartMessageMoveTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if output != nil && output.TaskHandle != nil {
				_ = os.Setenv(SqsTaskHandle, *output.TaskHandle)
				cancelMessageMoveTaskTest()
			}
		})
	}
}

func TestCancelMessageMoveTask(t *testing.T) {
	for _, tt := range initListTestCancelMessageMoveTask() {
		t.Run(tt.name, func(t *testing.T) {
			initStartMessageMoveTask()
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := CancelMessageMoveTask(ctx, tt.taskHandle, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelMessageMoveTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListMessageMoveTasks(t *testing.T) {
	initStartMessageMoveTask()
	for _, tt := range initListTestListMessageMoveTask() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := ListMessageMoveTasks(ctx, tt.sourceArn, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListMessageMoveTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
