package anyframer

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func StructToFrame(name string, input any) (frame *data.Frame, err error) {
	frame = data.NewFrame(name)
	if in, ok := input.(map[string]any); ok {
		fields := map[string]*data.Field{}
		for key, value := range in {
			newKey := fmt.Sprintf("%v", key)
			switch x := value.(type) {
			case nil, string, float64, float32, int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint, bool, time.Time:
				a, b := getFieldTypeAndValue(value)
				field := data.NewFieldFromFieldType(a, 1)
				field.Name = newKey
				field.Set(0, ToPointer(b))
				fields[newKey] = field
			case *string, *float64, *float32, *int64, *int32, *int16, *int8, *int, *uint64, *uint32, *uint16, *uint8, *uint, *bool, *time.Time:
				a, b := getFieldTypeAndValue(value)
				field := data.NewFieldFromFieldType(a, 1)
				field.Name = newKey
				field.Set(0, ToPointer(b))
				fields[newKey] = field
			default:
				fieldType, b := getFieldTypeAndValue(value)
				if fieldType == data.FieldTypeJSON {
					fieldType = data.FieldTypeNullableString
				}
				field := data.NewFieldFromFieldType(fieldType, 1)
				field.Name = newKey
				if o, err := json.Marshal(b); err == nil {
					field.Set(0, ToPointer(string(o)))
					fields[newKey] = field
				}
				noop(x)
			}
		}
		for _, key := range sortedKeys(in) {
			if f, ok := fields[key]; ok && f != nil {
				frame.Fields = append(frame.Fields, f)
			}
		}
		return frame, err
	}
	err = errors.New("unable to construct frame")
	return frame, err
}
