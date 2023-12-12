package util

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"reflect"
	"strconv"
	"time"
)

func ConvertToString(a any) (string, error) {
	if a == nil {
		return "", nil
	}
	v := reflect.ValueOf(a)
	if v.Type().Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.CanInterface() || IsZeroReflect(v) {
		return "", nil
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		b, err := json.Marshal(v.Interface())
		return string(b), err
	default:
		return convertToStringByType(a)
	}
}

func ConvertToBytes(a any) ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	t := reflect.TypeOf(a)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Interface:
		return json.Marshal(a)
	default:
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err := enc.Encode(a)
		return buffer.Bytes(), err
	}
}

func GetDataType(a any) (*string, error) {
	if a == nil {
		return nil, nil
	}
	t := reflect.TypeOf(a)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Interface, reflect.String, reflect.Bool:
		return aws.String("String"), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return aws.String("Number"), nil
	default:
		return nil, nil
	}
}

func IsNilValueReflect(v reflect.Value) bool {
	if (v.Kind() == reflect.Interface ||
		v.Kind() == reflect.Pointer ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Func ||
		v.Kind() == reflect.UnsafePointer ||
		v.Kind() == reflect.Map) && v.IsNil() {
		return true
	}
	return false
}

func IsZeroReflect(v reflect.Value) bool {
	return v.IsZero() || IsNilValueReflect(v) ||
		(v.Kind() == reflect.Map && len(v.MapKeys()) == 0) ||
		(v.Kind() == reflect.Struct && v.NumField() == 0) ||
		(v.Kind() == reflect.Slice && v.Len() == 0) ||
		(v.Kind() == reflect.Array && v.Len() == 0)
}

func IsNonZeroMessageAttValue(v *types.MessageAttributeValue) bool {
	return v != nil && v.DataType != nil &&
		(v.StringValue != nil || len(v.StringListValues) != 0) ||
		(v.BinaryValue != nil || len(v.BinaryListValues) != 0)
}

func convertToStringByType(a any) (string, error) {
	switch t := a.(type) {
	case int:
		return strconv.Itoa(t), nil
	case int8:
		return strconv.Itoa(int(t)), nil
	case uint8:
		return strconv.Itoa(int(t)), nil
	case int16:
		return strconv.Itoa(int(t)), nil
	case uint16:
		return strconv.Itoa(int(t)), nil
	case int32:
		return strconv.Itoa(int(t)), nil
	case uint32:
		return strconv.Itoa(int(t)), nil
	case int64:
		return strconv.Itoa(int(t)), nil
	case uint64:
		return strconv.Itoa(int(t)), nil
	case bool:
		return strconv.FormatBool(t), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case time.Time:
		b, err := t.MarshalJSON()
		return string(b), err
	case string:
		return t, nil
	}
	return fmt.Sprintf(`"%s"`, a), nil
}
