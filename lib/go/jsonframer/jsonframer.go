package jsonframer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/tidwall/gjson"
	"github.com/xiatechs/jsonata-go"
	"github.com/yesoreyeram/grafana-plugins/lib/go/gframer"
)

type FramerType string

const (
	FramerTypeGJSON   FramerType = "gjson"
	FramerTypeSQLite3 FramerType = "sqlite3"
)

type ColumnsType string

const (
	ColumnsTypeReplace  ColumnsType = "replace"  // Default. Only returns the selected columns
	ColumnsTypeAppend   ColumnsType = "append"   // Append new columns along with the auto-generated ones
	ColumnsTypeOverride ColumnsType = "override" // Override the specific columns with types and alias
)

type FramerOptions struct {
	FramerType   FramerType // `gjson` | `sqlite3`
	SQLite3Query string
	FrameName    string
	RootSelector string
	Columns      []ColumnSelector
	ColumnsType  ColumnsType
}

type ColumnSelector struct {
	Selector   string
	Alias      string
	Type       string
	TimeFormat string
}

func ToFrame(jsonString string, options FramerOptions) (frame *data.Frame, err error) {
	if strings.Trim(jsonString, " ") == "" {
		return frame, errors.New("empty json received")
	}
	if !gjson.Valid(jsonString) {
		return frame, errors.New("invalid json response received")
	}
	outString := jsonString
	switch options.FramerType {
	case "sqlite3":
		outString, err = QueryJSONUsingSQLite3(outString, options.SQLite3Query, options.RootSelector)
		if err != nil {
			return frame, err
		}
		return getFrameFromResponseString(outString, options)
	default:
		outString, err := GetRootData(jsonString, options.RootSelector)
		if err != nil {
			return frame, err
		}
		switch options.ColumnsType {
		case ColumnsTypeAppend:
			defaultFrame, err := getFrameFromResponseString(outString, FramerOptions{
				FrameName:    options.FrameName,
				RootSelector: options.RootSelector,
				Columns:      []ColumnSelector{},
			})
			if err != nil {
				return defaultFrame, err
			}
			newOutString, err := getColumnValuesFromResponseString(outString, options.Columns)
			if err != nil {
				return frame, err
			}
			frameWithNewColumns, err := getFrameFromResponseString(newOutString, options)
			if err != nil {
				return frameWithNewColumns, err
			}
			defaultFrame.Fields = append(defaultFrame.Fields, frameWithNewColumns.Fields...)
			return defaultFrame, nil
		case ColumnsTypeOverride:
			defaultFrame, err := getFrameFromResponseString(outString, FramerOptions{
				FrameName:    options.FrameName,
				RootSelector: options.RootSelector,
				Columns:      []ColumnSelector{},
			})
			if err != nil {
				return defaultFrame, err
			}
			columnsWithOutAlias := []ColumnSelector{}
			for _, c := range options.Columns {
				c.Alias = ""
				columnsWithOutAlias = append(columnsWithOutAlias, c)
			}
			newOutString, err := getColumnValuesFromResponseString(outString, columnsWithOutAlias)
			if err != nil {
				return frame, err
			}
			frameWithNewColumns, err := getFrameFromResponseString(newOutString, options)
			if err != nil {
				return defaultFrame, err
			}
			frameWithNewColumnsMap := map[string]*data.Field{}
			for _, f := range frameWithNewColumns.Fields {
				frameWithNewColumnsMap[f.Name] = f
			}
			for fi, f := range defaultFrame.Fields {
				if nf, ok := frameWithNewColumnsMap[f.Name]; ok {
					defaultFrame.Fields[fi] = nf
				}
			}
			return defaultFrame, err
		}
		outString, err = getColumnValuesFromResponseString(outString, options.Columns)
		if err != nil {
			return frame, err
		}
		return getFrameFromResponseString(outString, options)
	}
}

func GetRootData(jsonString string, rootSelector string) (string, error) {
	if rootSelector != "" {
		r := gjson.Get(string(jsonString), rootSelector)
		if r.Exists() {
			return r.String(), nil
		}
		if e := jsonata.MustCompile(rootSelector); e != nil {
			var data interface{}
			err := json.Unmarshal([]byte(jsonString), &data)
			if err != nil {
				return "", err
			}
			if res, err := e.Eval(data); err == nil {
				if r, err := json.Marshal(res); err == nil {
					return string(r), nil
				}
			}
		}
		return "", errors.New("root object doesn't exist in the response. Root selector:" + rootSelector)

	}
	return jsonString, nil

}

func getColumnValuesFromResponseString(responseString string, columns []ColumnSelector) (string, error) {
	if len(columns) > 0 {
		outString := responseString
		result := gjson.Parse(outString)
		out := []map[string]interface{}{}
		if result.IsArray() {
			result.ForEach(func(key, value gjson.Result) bool {
				oi := map[string]interface{}{}
				for _, col := range columns {
					name := col.Alias
					if name == "" {
						name = col.Selector
					}
					oi[name] = convertFieldValueType(gjson.Get(value.Raw, col.Selector).Value(), col)
				}
				out = append(out, oi)
				return true
			})
		}
		if !result.IsArray() && result.IsObject() {
			oi := map[string]interface{}{}
			for _, col := range columns {
				name := col.Alias
				if name == "" {
					name = col.Selector
				}
				oi[name] = convertFieldValueType(gjson.Get(result.Raw, col.Selector).Value(), col)
			}
			out = append(out, oi)
		}
		a, err := json.Marshal(out)
		if err != nil {
			return "", err
		}
		return string(a), nil
	}
	return responseString, nil
}

func getFrameFromResponseString(responseString string, options FramerOptions) (frame *data.Frame, err error) {
	var out interface{}
	err = json.Unmarshal([]byte(responseString), &out)
	if err != nil {
		return frame, fmt.Errorf("error while un-marshaling response. %s", err.Error())
	}
	columns := []gframer.ColumnSelector{}
	for _, c := range options.Columns {
		columns = append(columns, gframer.ColumnSelector{
			Alias:      c.Alias,
			Selector:   c.Selector,
			Type:       c.Type,
			TimeFormat: c.TimeFormat,
		})
	}
	return gframer.ToDataFrame(out, gframer.FramerOptions{
		FrameName: options.FrameName,
		Columns:   columns,
	})
}

func convertFieldValueType(input interface{}, col ColumnSelector) interface{} {
	return input
}
