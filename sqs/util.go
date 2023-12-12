package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go-aws-sqs/internal/util"
	"go-aws-sqs/sqs/option"
	"reflect"
)

func getDebugModeByOptsProducer(opts ...*option.OptionsProducer) bool {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.DebugMode {
			return true
		}
	}
	return false
}

func getOptionsHttpByOptsProducer(opts ...*option.OptionsProducer) *option.OptionsHttp {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.OptionsHttp != nil {
			return opt.OptionsHttp
		}
	}
	return nil
}

func getDelaySecondsByOptsProducer(opts ...*option.OptionsProducer) int32 {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.DelaySeconds > 0 {
			return int32(opt.DelaySeconds)
		}
	}
	return 0
}

func getMessageAttByOptsProducer(opts ...*option.OptionsProducer) (map[string]types.MessageAttributeValue, error) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		v := reflect.ValueOf(opt.MessageAttributes)
		if util.IsZeroReflect(v) {
			continue
		}
		return convertToMessageAttValue(v)
	}
	return nil, nil
}

func getMessageDeduplicationIdByOptsProducer(opts ...*option.OptionsProducer) *string {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if len(opt.MessageDeduplicationId) != 0 {
			return &opt.MessageDeduplicationId
		}
	}
	return nil
}

func getMessageGroupIdByOptsProducer(opts ...*option.OptionsProducer) *string {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if len(opt.MessageGroupId) != 0 {
			return &opt.MessageGroupId
		}
	}
	return nil
}

func getMessageSystemAttByOptsProducer(opts ...*option.OptionsProducer) (map[string]types.MessageSystemAttributeValue, error) {
	result := map[string]types.MessageSystemAttributeValue{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		v := reflect.ValueOf(opt.MessageSystemAttributes)
		if util.IsZeroReflect(v) {
			continue
		}
		resultMessageAttValue, err := convertToMessageAttValue(v)
		if err != nil {
			return nil, err
		}
		for k, mv := range resultMessageAttValue {
			result[k] = types.MessageSystemAttributeValue{
				DataType:    mv.DataType,
				StringValue: mv.StringValue,
			}
		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func convertToMessageAttValue(v reflect.Value) (map[string]types.MessageAttributeValue, error) {
	result := map[string]types.MessageAttributeValue{}
	for _, key := range v.MapKeys() {
		mKey := key.Convert(v.Type().Key())
		mValue := v.MapIndex(mKey)
		if !mKey.CanInterface() {
			continue
		}
		mKeyString, err := util.ConvertToString(mKey.Interface())
		if err != nil {
			return nil, err
		} else if len(mKeyString) != 0 && util.IsNilValueReflect(mValue) {
			continue
		}
		mValueProcessed, err := processMessageAttValue(mValue)
		if err != nil {
			return nil, err
		}
		if len(mKeyString) != 0 && util.IsNonZeroMessageAttValue(mValueProcessed) {
			result[mKeyString] = *mValueProcessed
		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func processMessageAttValue(v reflect.Value) (*types.MessageAttributeValue, error) {
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	result := types.MessageAttributeValue{}
	dataType, err := util.GetDataType(v.Interface())
	if err != nil {
		return nil, err
	} else if dataType == nil {
		return nil, nil
	}
	result.DataType = dataType
	switch *dataType {
	case "String", "Number":
		strValue, err := util.ConvertToString(v.Interface())
		if err != nil {
			return nil, err
		} else if len(strValue) == 0 {
			return nil, nil
		}
		result.StringValue = &strValue
		break
	default:
		byteValue, err := util.ConvertToBytes(v.Interface())
		if err != nil {
			return nil, err
		} else if len(byteValue) == 0 {
			return nil, nil
		}
		result.BinaryValue = byteValue
		break
	}
	return &result, nil
}

func getDebugModeByOptsCreateQueue(opts ...*option.OptionsCreateQueue) bool {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.DebugMode {
			return true
		}
	}
	return false
}

func getOptionsHttpByOptsCreateQueue(opts ...*option.OptionsCreateQueue) *option.OptionsHttp {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.OptionsHttp != nil {
			return opt.OptionsHttp
		}
	}
	return nil
}

func getAttributesByOptsCreateQueue(opts ...*option.OptionsCreateQueue) map[string]string {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.Attributes != nil {
			return opt.Attributes
		}
	}
	return nil
}

func getTagsByOptsCreateQueue(opts ...*option.OptionsCreateQueue) map[string]string {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.Tags != nil {
			return opt.Tags
		}
	}
	return nil
}

func getDebugModeByOptsDefault(opts ...*option.OptionsDefault) bool {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.DebugMode {
			return true
		}
	}
	return false
}

func getOptionsHttpByOptsDefault(opts ...*option.OptionsDefault) *option.OptionsHttp {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.OptionsHttp != nil {
			return opt.OptionsHttp
		}
	}
	return nil
}

func getDebugModeByOptsListQueues(opts ...*option.OptionsListQueues) bool {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.DebugMode {
			return true
		}
	}
	return false
}

func getOptionsHttpByOptsListQueues(opts ...*option.OptionsListQueues) *option.OptionsHttp {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.OptionsHttp != nil {
			return opt.OptionsHttp
		}
	}
	return nil
}

func getMaxResultsByOptsListQueues(opts ...*option.OptionsListQueues) *int32 {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.MaxResults != nil {
			return opt.MaxResults
		}
	}
	return nil
}

func getNextTokenByOptsListQueues(opts ...*option.OptionsListQueues) *string {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.NextToken != nil {
			return opt.NextToken
		}
	}
	return nil
}

func getQueueNamePrefixByOptsListQueues(opts ...*option.OptionsListQueues) *string {
	for _, opt := range opts {
		if opt == nil {
			continue
		} else if opt.QueueNamePrefix != nil {
			return opt.QueueNamePrefix
		}
	}
	return nil
}
