package anyframer

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func getFieldTypeAndValue(value any) (t data.FieldType, out any) {
	switch x := value.(type) {
	case nil:
		return data.FieldTypeNullableString, value
	case string:
		return data.FieldTypeNullableString, value
	case *string:
		return data.FieldTypeNullableString, value
	case float64:
		return data.FieldTypeNullableFloat64, value
	case *float64:
		return data.FieldTypeNullableFloat64, value
	case float32:
		return data.FieldTypeNullableFloat64, float64(value.(float32))
	case *float32:
		a := value.(*float32)
		return data.FieldTypeNullableFloat64, float64(*a)
	case int64:
		return data.FieldTypeNullableFloat64, float64(value.(int64))
	case *int64:
		a := value.(*int64)
		return data.FieldTypeNullableFloat64, float64(*a)
	case int32:
		return data.FieldTypeNullableFloat64, float64(value.(int32))
	case *int32:
		a := value.(*int32)
		return data.FieldTypeNullableFloat64, float64(*a)
	case int16:
		return data.FieldTypeNullableFloat64, float64(value.(int16))
	case *int16:
		a := value.(*int16)
		return data.FieldTypeNullableFloat64, float64(*a)
	case int8:
		return data.FieldTypeNullableFloat64, float64(value.(int8))
	case *int8:
		a := value.(*int8)
		return data.FieldTypeNullableFloat64, float64(*a)
	case int:
		return data.FieldTypeNullableFloat64, float64(value.(int))
	case *int:
		a := value.(*int)
		return data.FieldTypeNullableFloat64, float64(*a)
	case uint64:
		return data.FieldTypeNullableFloat64, float64(value.(uint64))
	case *uint64:
		a := value.(*uint64)
		return data.FieldTypeNullableFloat64, float64(*a)
	case uint32:
		return data.FieldTypeNullableFloat64, float64(value.(uint32))
	case *uint32:
		a := value.(*uint32)
		return data.FieldTypeNullableFloat64, float64(*a)
	case uint16:
		return data.FieldTypeNullableFloat64, float64(value.(uint16))
	case *uint16:
		a := value.(*uint16)
		return data.FieldTypeNullableFloat64, float64(*a)
	case uint8:
		return data.FieldTypeNullableFloat64, float64(value.(uint8))
	case *uint8:
		a := value.(*uint8)
		return data.FieldTypeNullableFloat64, float64(*a)
	case uint:
		return data.FieldTypeNullableFloat64, float64(value.(uint))
	case *uint:
		a := value.(*uint)
		return data.FieldTypeNullableFloat64, float64(*a)
	case bool:
		return data.FieldTypeNullableBool, value
	case *bool:
		return data.FieldTypeNullableBool, value
	case time.Time:
		return data.FieldTypeNullableTime, value
	case *time.Time:
		return data.FieldTypeNullableTime, value
	case any:
		return data.FieldTypeJSON, value
	default:
		noop(x)
		return data.FieldTypeNullableString, value
	}
}

func getFieldTypeFromSlice(value []any) (t data.FieldType) {
	for _, item := range value {
		if item != nil {
			a, _ := getFieldTypeAndValue(item)
			return a
		}
	}
	return data.FieldTypeNullableString
}

func getFieldValue(item any) (v any) {
	_, out := getFieldTypeAndValue(item)
	return out
}
