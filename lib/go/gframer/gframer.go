package gframer

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type ColumnSelector struct {
	Selector   string
	Alias      string
	Type       string
	TimeFormat string
}

type FramerOptions struct {
	FrameName           string
	ExecutedQueryString string
	Columns             []ColumnSelector
}

func noOperation(x interface{}) {}

func ToDataFrame(input interface{}, options FramerOptions) (frame *data.Frame, err error) {
	switch x := input.(type) {
	case nil, string, float64, float32, int64, int32, int16, int, bool:
		return structToFrame(options.FrameName, map[string]interface{}{options.FrameName: input}, options.ExecutedQueryString)
	case []interface{}:
		return sliceToFrame(options.FrameName, input.([]interface{}), options)
	default:
		noOperation(x)
		return structToFrame(options.FrameName, input, options.ExecutedQueryString)
	}
}

func structToFrame(name string, input interface{}, executedQueryString string) (frame *data.Frame, err error) {
	frame = data.NewFrame(name)
	if executedQueryString != "" {
		frame.Meta = &data.FrameMeta{
			ExecutedQueryString: executedQueryString,
		}
	}
	if in, ok := input.(map[string]interface{}); ok {
		fields := map[string]*data.Field{}
		for key, value := range in {
			switch x := value.(type) {
			case nil, string, float64, float32, int64, int32, int16, int, bool:
				noOperation(x)
				a, b := getFieldTypeAndValue(value)
				field := data.NewFieldFromFieldType(a, 1)
				field.Name = key
				field.Set(0, ToPointer(b))
				fields[key] = field
			default:
				fieldType, b := getFieldTypeAndValue(value)
				if fieldType == data.FieldTypeJSON {
					fieldType = data.FieldTypeNullableString
				}
				field := data.NewFieldFromFieldType(fieldType, 1)
				field.Name = key
				if o, err := json.Marshal(b); err == nil {
					field.Set(0, ToPointer(string(o)))
					fields[key] = field
				}
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

func sliceToFrame(name string, input []interface{}, options FramerOptions) (frame *data.Frame, err error) {
	frame = data.NewFrame(name)
	if options.ExecutedQueryString != "" {
		frame.Meta = &data.FrameMeta{
			ExecutedQueryString: options.ExecutedQueryString,
		}
	}
	if len(input) < 1 {
		return frame, err
	}
	for _, item := range input {
		if item != nil {
			switch item.(type) {
			case string, float64, float32, int64, int32, int16, int, bool:
				a, _ := getFieldTypeAndValue(item)
				field := data.NewFieldFromFieldType(a, len(input))
				field.Name = name
				for idx, i := range input {
					field.Set(idx, ToPointer(i))
				}
				frame.Fields = append(frame.Fields, field)
			case []interface{}:
				field := data.NewFieldFromFieldType(data.FieldTypeNullableString, len(input))
				field.Name = name
				for idx, i := range input {
					if o, err := json.Marshal(i); err == nil {
						field.Set(idx, ToPointer(string(o)))
					}
				}
				frame.Fields = append(frame.Fields, field)
			default:
				results := map[string]map[int]interface{}{}
				for idx, id := range input {
					if o, ok := id.(map[string]interface{}); ok {
						for k, v := range o {
							if results[k] == nil {
								results[k] = map[int]interface{}{}
							}
							results[k][idx] = v
						}
					}
				}
				for _, k := range sortedKeys(results) {
					if results[k] != nil {
						o := []interface{}{}
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
							if len(options.Columns) > 0 {
								for _, c := range options.Columns {
									if c.Alias == k || (c.Alias == "" && c.Selector == k) {
										switch c.Type {
										case "string":
											field := anyToNullableString(input, k, nil, o)
											frame.Fields = append(frame.Fields, field)
										case "boolean":
											field := anyToNullableBool(input, k, nil, o)
											frame.Fields = append(frame.Fields, field)
										case "number":
											field := anyToNullableNumber(input, k, nil, o)
											frame.Fields = append(frame.Fields, field)
										case "timestamp":
											field := anyToNullableTimestamp(input, k, nil, o, c.TimeFormat)
											frame.Fields = append(frame.Fields, field)
										case "timestamp_epoch":
											field := anyToNullableTimestampEpoch(input, k, nil, o)
											frame.Fields = append(frame.Fields, field)
										case "timestamp_epoch_s":
											field := anyToNullableTimestampEpochSecond(input, k, nil, o)
											frame.Fields = append(frame.Fields, field)
										default:
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
							if len(options.Columns) < 1 {
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

func getFieldTypeAndValue(value interface{}) (t data.FieldType, out interface{}) {
	switch x := value.(type) {
	case nil:
		return data.FieldTypeNullableString, value
	case string:
		return data.FieldTypeNullableString, value
	case float64:
		return data.FieldTypeNullableFloat64, value
	case float32:
		return data.FieldTypeNullableFloat64, float64(value.(float32))
	case int64:
		return data.FieldTypeNullableFloat64, float64(value.(int64))
	case int32:
		return data.FieldTypeNullableFloat64, float64(value.(int32))
	case int16:
		return data.FieldTypeNullableFloat64, float64(value.(int16))
	case int:
		return data.FieldTypeNullableFloat64, float64(value.(int))
	case bool:
		return data.FieldTypeNullableBool, value
	case interface{}:
		return data.FieldTypeJSON, value
	default:
		noOperation(x)
		return data.FieldTypeNullableString, value
	}
}

func getFieldTypeFromSlice(value []interface{}) (t data.FieldType) {
	for _, item := range value {
		if item != nil {
			a, _ := getFieldTypeAndValue(item)
			return a
		}
	}
	return data.FieldTypeNullableString
}

func sortedKeys(in interface{}) []string {
	if input, ok := in.(map[string]interface{}); ok {
		keys := make([]string, len(input))
		var idx int
		for key := range input {
			keys[idx] = key
			idx++
		}
		sort.Strings(keys)
		return keys
	}
	if input, ok := in.(map[string]map[int]interface{}); ok {
		keys := make([]string, len(input))
		var idx int
		for key := range input {
			keys[idx] = key
			idx++
		}
		sort.Strings(keys)
		return keys
	}
	return []string{}
}

func ToPointer(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case int8:
		return &v
	case *int8:
		return value
	case int16:
		return &v
	case *int16:
		return value
	case int32:
		return &v
	case *int32:
		return value
	case int64:
		return &v
	case *int64:
		return value
	case uint8:
		return &v
	case *uint8:
		return value
	case uint16:
		return &v
	case *uint16:
		return value
	case uint32:
		return &v
	case *uint32:
		return value
	case uint64:
		return &v
	case *uint64:
		return value
	case float32:
		return &v
	case *float32:
		return value
	case float64:
		return &v
	case *float64:
		return value
	case string:
		return &v
	case *string:
		return value
	case bool:
		return &v
	case *bool:
		return value
	case time.Time:
		return &v
	case *time.Time:
		return value
	default:
		return nil
	}
}
