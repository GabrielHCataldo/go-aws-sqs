package sqs

import (
	"context"
	"errors"
	"github.com/GabrielHCataldo/go-aws-sqs/internal/client"
	"github.com/GabrielHCataldo/go-aws-sqs/internal/util"
	"github.com/GabrielHCataldo/go-aws-sqs/sqs/option"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"reflect"
)

// SendMessage sends a message to an SQS queue using the provided context, queue URL, message content, and options.
//
// The function first retrieves an SQS client using the getSqsClient function. It then prepares the message input using
// the prepareMessageInput function. The message input includes the message body, queue URL, delay seconds, message
// attributes, message deduplication ID, message group ID, and message system attributes. The message input is passed to
// the client's SendMessage method, which sends the message to the queue. The function returns the SendMessage output or
// an error if one occurs.
//
// Example usage:
//
//	output, err := SendMessage(ctx, queueUrl, v, opts...)
//
// # Parameters:
//
// - ctx (Context): The context of the request.
// - queueUrl (string): The URL of the SQS queue.
// - body (any): The content of the message.
// - opts (option.Producer): Optional options to customize the message. (see OptionsProducer declaration for available options)
//
// # Returns:
//
// - *sqs.SendMessageOutput: The result of the SendMessage operation.
// - error: An error if one occurs during the SendMessage operation.
func SendMessage(ctx context.Context, queueUrl string, body any, opts ...option.Producer) (*sqs.SendMessageOutput, error) {
	opt := option.GetProducerByParams(opts)
	loggerInfo(opt.DebugMode, "getting client sqs..")
	sqsClient := client.GetClient(ctx)
	loggerInfo(opt.DebugMode, "preparing message input..")
	input, err := prepareMessageInput(queueUrl, body, opt)
	if err != nil {
		loggerErr(opt.DebugMode, "error prepare message input:", err)
		return nil, err
	}
	loggerInfo(opt.DebugMode, "sending message..")
	output, err := sqsClient.SendMessage(ctx, input, option.FuncByHttpClient(opt.HttpClient))
	if err != nil {
		loggerErr(opt.DebugMode, "error send message:", err)
	} else {
		loggerInfo(opt.DebugMode, "message sent successfully:", output)
	}
	return output, err
}

// SendMessageAsync sends a message asynchronously to an SQS queue using the provided context, queue URL,
// message content, and options.
//
// The function spawns a goroutine and calls the sendMessageAsync function internally to send the message.
// It first checks if the debug mode is enabled by checking the options. It then calls the SendMessage function,
// which sends the message to the queue. If an error occurs and debug mode is enabled, it logs the error.
// If debug mode is enabled, it logs the successful message delivery.
//
// # Parameters:
//
// - ctx (Context): The context of the request.
// - queueUrl (string): The URL of the SQS queue.
// - body (any): The content of the message.
// - opts (OptionsProducer): Optional options to customize the message. (see OptionsProducer declaration for available options)
//
// # Example usage:
//
// SendMessageAsync(ctx, queueUrl, v, opts...)
func SendMessageAsync(ctx context.Context, queueUrl string, body any, opts ...option.Producer) {
	go sendMessageAsync(ctx, queueUrl, body, opts...)
}

func sendMessageAsync(ctx context.Context, queueUrl string, body any, opts ...option.Producer) {
	_, _ = SendMessage(ctx, queueUrl, body, opts...)
}

func prepareMessageInput(queueUrl string, v any, opt option.Producer) (*sqs.SendMessageInput, error) {
	body := util.ConvertToString(v)
	if len(body) == 0 {
		return nil, ErrMessageBodyEmpty
	}
	messageAttByOpt, err := getMessageAttValueByOpt(opt)
	if err != nil {
		return nil, err
	}
	return &sqs.SendMessageInput{
		MessageBody:             aws.String(body),
		QueueUrl:                &queueUrl,
		DelaySeconds:            util.ConvertDurationToInt32(opt.DelaySeconds),
		MessageAttributes:       messageAttByOpt,
		MessageDeduplicationId:  opt.MessageDeduplicationId,
		MessageGroupId:          opt.MessageGroupId,
		MessageSystemAttributes: getMessageSystemAttValueByOpt(opt),
	}, nil
}

func getMessageAttValueByOpt(opt option.Producer) (map[string]types.MessageAttributeValue, error) {
	if opt.MessageAttributes == nil {
		return nil, nil
	}
	v := reflect.ValueOf(opt.MessageAttributes)
	t := reflect.TypeOf(opt.MessageAttributes)
	if t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface {
		v = v.Elem()
		t = t.Elem()
	}
	if util.IsZeroReflect(v) {
		return nil, nil
	}
	return convertToMessageAttValue(t, v)
}

func getMessageSystemAttValueByOpt(opt option.Producer) map[string]types.MessageSystemAttributeValue {
	v := reflect.ValueOf(opt.MessageSystemAttributes)
	if util.IsZeroReflect(v) {
		return nil
	}
	return map[string]types.MessageSystemAttributeValue{
		"AWSTraceHeader": {
			DataType:    aws.String("String"),
			StringValue: aws.String(opt.MessageSystemAttributes.AWSTraceHeader),
		},
	}
}

func convertToMessageAttValue(t reflect.Type, v reflect.Value) (map[string]types.MessageAttributeValue, error) {
	switch v.Kind() {
	case reflect.Map:
		return convertMapToMessageAttValue(v)
	case reflect.Struct:
		return convertStructToMessageAttValue(t, v)
	default:
		return nil, errors.New("messageAttributes or messageSystemAttributes must be map or structure")
	}
}

func convertMapToMessageAttValue(v reflect.Value) (map[string]types.MessageAttributeValue, error) {
	result := map[string]types.MessageAttributeValue{}
	for _, key := range v.MapKeys() {
		mKey := key.Convert(v.Type().Key())
		mValue := v.MapIndex(mKey)
		if util.IsZeroReflect(mKey) {
			continue
		}
		mKeyString := util.ConvertToString(mKey.Interface())
		if util.IsZeroReflect(mValue) {
			continue
		}
		mValueProcessed := convertReflectToMessageAttributeValue(mValue)
		if len(mKeyString) != 0 && util.IsNonZeroMessageAttValue(mValueProcessed) {
			result[mKeyString] = *mValueProcessed
		}
	}
	return result, nil
}

func convertStructToMessageAttValue(t reflect.Type, v reflect.Value) (map[string]types.MessageAttributeValue, error) {
	result := map[string]types.MessageAttributeValue{}
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldStruct := t.Field(i)
		fieldName := util.GetJsonNameByTag(fieldStruct.Tag.Get("json"))
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = fieldStruct.Name
		}
		if !fieldValue.CanInterface() || util.IsZeroReflect(fieldValue) {
			continue
		}
		valueConverted := convertReflectToMessageAttributeValue(fieldValue)
		if len(fieldName) != 0 && util.IsNonZeroMessageAttValue(valueConverted) {
			result[fieldName] = *valueConverted
		}
	}
	return result, nil
}

func convertReflectToMessageAttributeValue(v reflect.Value) *types.MessageAttributeValue {
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	result := types.MessageAttributeValue{}
	dataType := util.GetDataType(v.Interface())
	result.DataType = &dataType
	if dataType == "String" || dataType == "Number" {
		strValue := util.ConvertToString(v.Interface())
		if len(strValue) == 0 {
			return nil
		}
		result.StringValue = &strValue
		return &result
	}
	return nil
}
