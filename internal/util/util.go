package util

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ConvertToString(a any) string {
	if a == nil {
		return ""
	}
	v := reflect.ValueOf(a)
	if v.Type().Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.CanInterface() || IsZeroReflect(v) {
		return ""
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		b, _ := json.Marshal(v.Interface())
		return string(b)
	default:
		return convertToStringByType(a)
	}
}

func GetDataType(a any) string {
	t := reflect.TypeOf(a)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Interface, reflect.String, reflect.Bool:
		return "String"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "Number"
	default:
		return ""
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
	return v.Kind() == reflect.Invalid || v.IsZero() || IsNilValueReflect(v) || !v.CanInterface() ||
		(v.Kind() == reflect.Map && len(v.MapKeys()) == 0) ||
		(v.Kind() == reflect.Struct && v.NumField() == 0) ||
		(v.Kind() == reflect.Slice && v.Len() == 0) ||
		(v.Kind() == reflect.Array && v.Len() == 0)
}

func IsNonZeroMessageAttValue(v *types.MessageAttributeValue) bool {
	return v != nil && v.DataType != nil &&
		(v.StringValue != nil && len(*v.StringValue) != 0)
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

func convertToStringByType(a any) string {
	switch t := a.(type) {
	case int:
		return strconv.Itoa(t)
	case int8:
		return strconv.Itoa(int(t))
	case uint8:
		return strconv.Itoa(int(t))
	case int16:
		return strconv.Itoa(int(t))
	case uint16:
		return strconv.Itoa(int(t))
	case int32:
		return strconv.Itoa(int(t))
	case uint32:
		return strconv.Itoa(int(t))
	case int64:
		return strconv.Itoa(int(t))
	case uint64:
		return strconv.Itoa(int(t))
	case bool:
		return strconv.FormatBool(t)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case time.Time:
		b, _ := t.MarshalJSON()
		return string(b)
	case string:
		return t
	}
	return fmt.Sprintf(`"%s"`, a)
}
