package sqs

import (
	"context"
	"go-aws-sqs/sqs/option"
	"os"
	"testing"
	"time"
)

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

type mapTest struct {
	Int    int     `json:"int,omitempty"`
	Float  float64 `json:"float,omitempty"`
	Bool   bool    `json:"bool,omitempty"`
	Slice  []int   `json:"slice,omitempty"`
	String string  `json:"string,omitempty"`
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name     string
		queueUrl string
		v        any
		opts     []*option.OptionsProducer
		wantErr  bool
		async    bool
	}{
		{
			name:     "valid request",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        getTestStruct(),
			opts: []*option.OptionsProducer{
				option.Producer().SetMessageAttributes(option.MessageAttribute{
					"int":        1,
					"float":      12.2432,
					"string":     "text test",
					"bool":       true,
					"struct":     getTestStruct(),
					"anotherMap": getTestMap(),
				}),
				option.Producer().SetDebugMode(true),
			},
			wantErr: false,
		},
		{
			name:     "valid request FIFO",
			queueUrl: os.Getenv("SQS_QUEUE_TEST_FIFO"),
			v:        getTestStruct(),
			opts: []*option.OptionsProducer{
				option.Producer().SetMessageAttributes(option.MessageAttribute{
					"int":        1,
					"float":      12.2432,
					"string":     "text test",
					"bool":       true,
					"struct":     getTestStruct(),
					"anotherMap": getTestMap(),
				}),
				option.Producer().SetDebugMode(true),
				option.Producer().SetMessageDeduplicationId("test"),
				option.Producer().SetMessageGroupId("test"),
			},
			wantErr: false,
		},
		{
			name:     "valid request Async",
			queueUrl: os.Getenv("SQS_QUEUE_TEST"),
			v:        getTestStruct(),
			opts: []*option.OptionsProducer{
				option.Producer().SetMessageAttributes(option.MessageAttribute{
					"int":        1,
					"float":      12.2432,
					"string":     "text test",
					"bool":       true,
					"struct":     getTestStruct(),
					"anotherMap": getTestMap(),
				}),
				option.Producer().SetDebugMode(true),
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			var err error
			if tt.async {
				SendMessageAsync(ctx, tt.queueUrl, tt.v, tt.opts...)
			} else {
				_, err = SendMessage(ctx, tt.queueUrl, tt.v, tt.opts...)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestStruct() test {
	return test{
		Name:      "Test Name",
		BirthDate: time.Now(),
		Emails:    []string{"test@gmail.com", "gabriel@gmail.com", "gabriel.test@gmail.com"},
		Bank: bank{
			Account: "123456",
			Digits:  "2",
			Balance: 200.12,
		},
		Map: getTestMap(),
	}
}

func getTestMap() map[string]any {
	return map[string]any{
		"int":    1,
		"float":  12.2432,
		"bool":   true,
		"slice":  []int{12, 23, 53},
		"string": "text test",
	}
}
