package anyframer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func SliceToFrame(name string, input []any, columns []Column) (frame *data.Frame, err error) {
	frame = data.NewFrame(name)
	if len(input) < 1 {
		return frame, err
	}
	if len(columns) > 0 {
		for _, column := range columns {
			if column.Selector == "" {
				continue
			}
			fieldName := column.Alias
			if fieldName == "" {
				fieldName = column.Selector
			}
			switch column.Format {
			case ColumnFormatString:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableString, len(input))
				field.Name = fieldName
				for i := 0; i < len(input); i++ {
					currentValue, err := applySelector(input[i], column.Selector)
					if err != nil {
						continue
					}
					switch cvt := currentValue.(type) {
					case string:
						field.Set(i, ToPointer(currentValue.(string)))
					case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
						field.Set(i, ToPointer(fmt.Sprintf("%v", currentValue)))
					case bool:
						field.Set(i, ToPointer(fmt.Sprintf("%v", currentValue.(bool))))
					case time.Time:
						field.Set(i, ToPointer(currentValue.(time.Time).String()))
					default:
						noop(cvt)
						field.Set(i, nil)
					}
				}
				frame.Fields = append(frame.Fields, field)
			case ColumnFormatNumber:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableFloat64, len(input))
				field.Name = fieldName
				for i := 0; i < len(input); i++ {
					currentValue, err := applySelector(input[i], column.Selector)
					if err != nil {
						continue
					}
					switch cvt := currentValue.(type) {
					case string:
						if item, err := strconv.ParseFloat(currentValue.(string), 64); err == nil {
							field.Set(i, ToPointer(item))
						}
					case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
						field.Set(i, ToPointer(getFieldValue(currentValue)))
					default:
						noop(cvt)
						field.Set(i, nil)
					}
				}
				frame.Fields = append(frame.Fields, field)
			case ColumnFormatBoolean:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableBool, len(input))
				field.Name = fieldName
				for i := 0; i < len(input); i++ {
					currentValue, err := applySelector(input[i], column.Selector)
					if err != nil {
						continue
					}
					switch cvt := currentValue.(type) {
					case bool:
						field.Set(i, ToPointer(currentValue))
					case string:
						switch strings.ToLower(strings.TrimSpace(currentValue.(string))) {
						case "true":
							field.Set(i, ToPointer(true))
						case "false":
							field.Set(i, ToPointer(false))
						}

					default:
						noop(cvt)
						field.Set(i, nil)
					}
				}
				frame.Fields = append(frame.Fields, field)
			case ColumnFormatTimeStamp:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableTime, len(input))
				field.Name = fieldName
				for i := 0; i < len(input); i++ {
					currentValue, err := applySelector(input[i], column.Selector)
					if err != nil {
						continue
					}
					switch a := currentValue.(type) {
					case float64:
						if v := fmt.Sprintf("%v", currentValue); v != "" {
							if t, err := time.Parse("2006", v); err == nil {
								field.Set(i, ToPointer(t))
							}
						}
					case string:
						if currentValue.(string) != "" {
							field.Set(i, getTimeFromString(currentValue.(string), column.TimeFormat))
						}
					default:
						noop(a)
						field.Set(i, nil)
					}
				}
				frame.Fields = append(frame.Fields, field)
			case ColumnFormatUnixMsecTimeStamp:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableTime, len(input))
				field.Name = fieldName
				for i := 0; i < len(input); i++ {
					currentValue, err := applySelector(input[i], column.Selector)
					if err != nil {
						continue
					}
					switch cvt := currentValue.(type) {
					case string:
						if item, err := strconv.ParseInt(currentValue.(string), 10, 64); err == nil && currentValue.(string) != "" {
							field.Set(i, ToPointer(time.UnixMilli(item)))
						}
					case float64:
						field.Set(i, ToPointer(time.UnixMilli(int64(currentValue.(float64)))))
					default:
						noop(cvt)
						field.Set(i, nil)
					}
				}
				frame.Fields = append(frame.Fields, field)
			case ColumnFormatUnixSecTimeStamp:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableTime, len(input))
				field.Name = fieldName
				for i := 0; i < len(input); i++ {
					currentValue, err := applySelector(input[i], column.Selector)
					if err != nil {
						continue
					}
					switch cvt := currentValue.(type) {
					case string:
						if item, err := strconv.ParseInt(currentValue.(string), 10, 64); err == nil && currentValue.(string) != "" {
							field.Set(i, ToPointer(time.Unix(item, 0)))
						}
					case float64:
						field.Set(i, ToPointer(time.Unix(int64(currentValue.(float64)), 0)))
					default:
						noop(cvt)
						field.Set(i, nil)
					}
				}
				frame.Fields = append(frame.Fields, field)
			}
		}
		return frame, nil
	}
	for _, item := range input {
		if item != nil {
			switch item.(type) {
			case nil, string, float64, float32, int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint, bool, time.Time, *string, *float64, *float32, *int64, *int32, *int16, *int8, *int, *uint64, *uint32, *uint16, *uint8, *uint, *bool, *time.Time:
				a, _ := getFieldTypeAndValue(item)
				field := data.NewFieldFromFieldType(a, len(input))
				field.Name = name
				for idx, i := range input {
					field.Set(idx, ToPointer(i))
				}
				frame.Fields = append(frame.Fields, field)
			case []any:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableString, len(input))
				field.Name = name
				for idx, i := range input {
					if o, err := json.Marshal(i); err == nil {
						field.Set(idx, ToPointer(string(o)))
					}
				}
				frame.Fields = append(frame.Fields, field)
			default:
				results := map[string]map[int]any{}
				for idx, id := range input {
					if o, ok := id.(map[string]any); ok {
						for k, v := range o {
							if results[k] == nil {
								results[k] = map[int]any{}
							}
							results[k][idx] = getFieldValue(v)
						}
					}
				}
				for _, k := range sortedKeys(results) {
					if results[k] != nil {
						o := []any{}
						for i := 0; i < len(input); i++ {
							o = append(o, results[k][i])
						}
						fieldType := getFieldTypeFromSlice(o)
						if fieldType == data.FieldTypeJSON {
							field := data.NewFieldFromFieldType(data.FieldTypeNullableString, len(input))
							field.Name = k
							for i := 0; i < len(input); i++ {
								if o, err := json.Marshal(o[i]); err == nil {
									field.Set(i, ToPointer(string(o)))
								}
							}
							frame.Fields = append(frame.Fields, field)
						}
						if fieldType != data.FieldTypeJSON {
							if len(columns) < 1 {
								field := data.NewFieldFromFieldType(fieldType, len(input))
								field.Name = k
								for i := 0; i < len(input); i++ {
									field.Set(i, ToPointer(o[i]))
								}
								frame.Fields = append(frame.Fields, field)
							}
						}
					}
				}
			}
			break
		}
	}
	if len(frame.Fields) == 0 {
		field := data.NewFieldFromFieldType(data.FieldTypeNullableString, len(input))
		field.Name = name
		frame.Fields = append(frame.Fields, field)
	}
	return frame, nil
}
