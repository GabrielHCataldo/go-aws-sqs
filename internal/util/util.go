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
	"strings"
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

func GetJsonNameByTag(tag string) string {
	result := ""
	splitTag := strings.Split(tag, ",")
	if len(splitTag) > 0 {
		result = splitTag[0]
	}
	if len(result) == 0 || result == "omitempty" {
		return ""
	}
	return result
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
	return v.Kind() == reflect.Invalid || v.IsZero() || IsNilValueReflect(v) ||
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

func ParseStringToGeneric[T any](s string, dest *T) {
	if err := json.Unmarshal([]byte(s), dest); err == nil {
		return
	}
	if i, err := strconv.Atoi(s); err == nil {
		bodyInt, ok := any(i).(T)
		if ok {
			*dest = bodyInt
			return
		}
	}
	if b, err := strconv.ParseBool(s); err == nil {
		bodyBool, ok := any(b).(T)
		if ok {
			*dest = bodyBool
			return
		}
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		bodyFloat, ok := any(f).(T)
		if ok {
			*dest = bodyFloat
			return
		}
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		bodyTime, ok := any(t).(T)
		if ok {
			*dest = bodyTime
			return
		}
	}
	bodyString, ok := any(s).(T)
	if ok {
		*dest = bodyString
	}
}

func ParseMapToStruct[T any](m map[string]any, dest *T) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

func IsMapMessageAttributeValues(a any) bool {
	switch a.(type) {
	case map[string]types.MessageAttributeValue:
		return true
	default:
		return false
	}
}

func IsValidType(a any) bool {
	r := reflect.ValueOf(a)
	if !IsNilValueReflect(r) && r.Kind() != reflect.Invalid {
		return true
	}
	return false
}

func ConvertDurationToInt32(d time.Duration) int32 {
	return int32(d.Seconds())
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
