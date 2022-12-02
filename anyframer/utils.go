package anyframer

import (
	"sort"
	"time"
)

func sortedKeys(in any) []string {
	if input, ok := in.(map[string]any); ok {
		keys := make([]string, len(input))
		var idx int
		for key := range input {
			keys[idx] = key
			idx++
		}
		sort.Strings(keys)
		return keys
	}
	if input, ok := in.(map[string]map[int]any); ok {
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

func ToPointer(value any) any {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case int:
		return &v
	case *int:
		return value
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
	case uint:
		return &v
	case *uint:
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
		noop(value)
		return nil
	}
}

func noop(x any) {}
