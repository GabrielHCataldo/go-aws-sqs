package sqs

import (
	"context"
	"go-aws-sqs/sqs/option"
	"testing"
)

func TestListQueues(t *testing.T) {
	// Use the test cases you think are necessary
	testCases := []struct {
		name    string
		ctx     context.Context
		opts    []option.ListQueues
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.TODO(),
			wantErr: false,
			opts: []option.ListQueues{
				option.NewListQueue().SetDebugMode(true),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ListQueues(tt.ctx, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListQueues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
