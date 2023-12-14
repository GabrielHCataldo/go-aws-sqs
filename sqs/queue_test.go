package sqs

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCreateQueue(t *testing.T) {
	for _, tt := range initListTestCreateQueue() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			output, err := CreateQueue(ctx, tt.queueName, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if output != nil && output.QueueUrl != nil {
				_ = os.Setenv(SqsQueueCreateTestName, tt.queueName)
				_ = os.Setenv(SqsQueueCreateTestUrl, *output.QueueUrl)
				deleteQueueCreateTest()
			}
		})
	}
}

func TestDeleteQueue(t *testing.T) {
	initQueueCreateTest()
	for _, tt := range initListTestDeleteQueue() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := DeleteQueue(ctx, tt.queueUrl, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagQueue(t *testing.T) {
	initQueueCreateTest()
	for _, tt := range initListTestTagQueue() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := TagQueue(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TagQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetAttributesQueue(t *testing.T) {
	initQueueCreateTest()
	for _, tt := range initListTestSetAttributesQueue() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := SetQueueAttributes(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetQueueAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUntagQueue(t *testing.T) {
	initQueueCreateTest()
	for _, tt := range initListTestUntagQueue() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := UntagQueue(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("UntagQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPurgeQueue(t *testing.T) {
	initQueueCreateTest()
	for _, tt := range initListTestPurgeQueue() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := PurgeQueue(ctx, tt.queueUrl, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PurgeQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetQueueUrl(t *testing.T) {
	for _, tt := range initListTestGetQueueUrl() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := GetQueueUrl(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueueUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetQueueAttributes(t *testing.T) {
	for _, tt := range initListTestGetQueueAttributes() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := GetQueueAttributes(ctx, tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueueAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListQueues(t *testing.T) {
	for _, tt := range initListTestListQueue() {
		t.Run(tt.name, func(t *testing.T) {
			d := 5 * time.Second
			if tt.name == "failed" {
				d = 0
			}
			ctx, cancel := context.WithTimeout(context.TODO(), d)
			defer cancel()
			_, err := ListQueues(ctx, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListQueues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListQueueTags(t *testing.T) {
	for _, tt := range initListTestListQueueTags() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := ListQueueTags(ctx, tt.queueUrl, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListQueueTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListDeadLetterSourceQueues(t *testing.T) {
	for _, tt := range initListTestListDeadLetterSourceQueues() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			_, err := ListDeadLetterSourceQueues(ctx, tt.queueUrl, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDeadLetterSourceQueues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
