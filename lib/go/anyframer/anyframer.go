package anyframer

import (
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// AnyFramer defines the framer options
type AnyFramer struct {
	Name         string      `json:"name,omitempty"`
	InputType    InputType   `json:"inputType,omitempty"`
	RawURL       string      `json:"rawUrl,omitempty"`
	Headers      http.Header `json:"headers,omitempty"`
	RootSelector string      `json:"rootSelector,omitempty"`
	Columns      []Column    `json:"columns,omitempty"`
	CSVOptions   CSVOptions  `json:"csvOptions,omitempty"`
}

// ToFrame converts the given input string or input interface to data frame
func (framerOptions *AnyFramer) ToFrame(input any) (*data.Frame, error) {
	if framerOptions.Name == "" {
		framerOptions.Name = "response"
	}
	switch t := input.(type) {
	case string:
		if framerOptions.InputType == "" || framerOptions.InputType == InputTypeUnknown {
			framerOptions.InputType = framerOptions.GuessType(input.(string))
		}
		jsonObject, err := framerOptions.toJSONObject(input.(string))
		if err != nil {
			return nil, err
		}
		return toFrameFromInterface(jsonObject, *framerOptions)
	default:
		noop(t)
		return toFrameFromInterface(input, *framerOptions)
	}
}

func toFrameFromInterface(input any, options AnyFramer) (*data.Frame, error) {
	input, err := applySelector(input, options.RootSelector)
	if err != nil {
		return nil, err
	}
	switch x := input.(type) {
	case nil, string, float64, float32, int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint, bool, time.Time:
		return StructToFrame(options.Name, map[string]any{options.Name: input})
	case []any:
		return SliceToFrame(options.Name, input.([]any), options.Columns)
	default:
		noop(x)
		return StructToFrame(options.Name, input)
	}
}

// CSVOptions ...
type CSVOptions struct {
	Delimiter          string   `json:"delimiter,omitempty"`
	Comment            string   `json:"comment,omitempty"`
	RelaxColumnCount   bool     `json:"relaxColumnCount,omitempty"`
	SkipLinesWithError bool     `json:"skipLinesWithError,omitempty"`
	NoHeaders          bool     `json:"noHeaders,omitempty"`
	Headers            []string `json:"headers,omitempty"`
}
