package jsonframer_test

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/lib/go/jsonframer"
)

func TestJsonStringToFrame(t *testing.T) {
	updateTestData := false
	tests := []struct {
		name           string
		responseString string
		refId          string
		rootSelector   string
		columns        []jsonframer.ColumnSelector
		columnsType    jsonframer.ColumnsType
		wantFrame      *data.Frame
		wantErr        error
	}{
		{
			name:           "empty string should throw error",
			responseString: "",
			wantErr:        errors.New("empty json received"),
		},
		{
			name:           "invalid json should throw error",
			responseString: "{",
			wantErr:        errors.New("invalid json response received"),
		},
		{
			name:           "valid json object should not throw error",
			responseString: "{}",
		},
		{
			name:           "valid json array should not throw error",
			responseString: "[]",
		},
		{
			name:           "valid string array should not throw error",
			responseString: `["foo", "bar"]`,
		},
		{
			name:           "valid numeric array should not throw error",
			responseString: `[123, 123.45]`,
		},
		{
			name:           "valid json object with data should not throw error",
			responseString: `{ "username": "foo", "age": 1, "height" : 123.45,  "isPremium": true, "hobbies": ["reading","swimming"] }`,
		},
		{
			name:           "valid json array with data should not throw error",
			responseString: `[{ "username": "foo", "age": 1, "height" : 123.45,  "isPremium": true, "hobbies": ["reading","swimming"] }]`,
		},
		{
			name: "valid json array with multiple rows should not throw error",
			responseString: `[
				{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
				{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["reading","swimming"], "occupation": "student" }
			]`,
		},
		{
			name: "without root data and valid json array with multiple rows should not throw error",
			responseString: `{
				"meta" : {},
				"data" : [
					{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
					{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["reading","swimming"], "occupation": "student" }
				]
			}`,
		},
		{
			name: "with root data and valid json array with multiple rows should not throw error",
			responseString: `{
				"meta" : {},
				"data" : [
					{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
					{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["reading","swimming"], "occupation": "student" }
				]
			}`,
			rootSelector: "data",
		},
		{
			name: "with root data and selectors should produce valid frame",
			responseString: `{
				"meta" : {},
				"data" : [
					{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
					{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["reading","swimming"], "occupation": "student" }
				]
			}`,
			rootSelector: "data",
			columns: []jsonframer.ColumnSelector{
				{Selector: "username", Alias: "user-name"},
				{Selector: "occupation"},
			},
		},
		{
			name: "with root data and selectors should produce valid frame for non array object",
			responseString: `{
				"meta" : {},
				"data" : { "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["reading","swimming"], "occupation": "student" }
			}`,
			rootSelector: "data",
			columns: []jsonframer.ColumnSelector{
				{Selector: "username", Alias: "user-name"},
				{Selector: "occupation"},
			},
		},
		{
			name: "column values",
			responseString: `[
				{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
				{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["reading","swimming"], "occupation": "student" }
			]`,
			rootSelector: "",
			columns: []jsonframer.ColumnSelector{
				{Selector: "age"},
				{Selector: "occupation"},
			},
		},
		{
			name: "column values with append",
			responseString: `[
				{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
				{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["writing","swimming"], "occupation": "student" }
			]`,
			refId:        "M",
			rootSelector: "",
			columns: []jsonframer.ColumnSelector{
				{Selector: "age"},
				{Selector: "hobbies.0", Alias: "Primary Hobby"},
				{Selector: "occupation"},
			},
			columnsType: jsonframer.ColumnsTypeAppend,
		},
		{
			name: "column values with overwrite",
			responseString: `[
				{ "username": "foo", "age": 1, "height" : 123,  "isPremium": true, "hobbies": ["reading","swimming"] },
				{ "username": "bar", "age": 2, "height" : 123.45,  "isPremium": false, "hobbies": ["writing","swimming"], "occupation": "student" }
			]`,
			rootSelector: "",
			columns: []jsonframer.ColumnSelector{
				{Selector: "age", Type: "string", Alias: "Age"},
			},
			columnsType: jsonframer.ColumnsTypeOverride,
		},
		{
			name: "string",
			responseString: `{
				"sss": [
					{ "foo" : "1.2", "bar1": 4, "baz" : true },
					{ "foo" : "3", "bar1": 5.6, "baz" : false }
				]
			}`,
			rootSelector: "sss",
			columns: []jsonframer.ColumnSelector{
				{Selector: "foo", Type: "string"},
				{Selector: "bar1", Alias: "bar", Type: "string"},
				{Selector: "baz", Type: "string"},
			},
		},
		{
			name: "number",
			responseString: `{
				"sss": [
					{ "foo" : "1.2", "bar1": 4, "baz" : true },
					{ "foo" : "3", "bar1": 5.6, "baz" : false }
				]
			}`,
			rootSelector: "sss",
			columns: []jsonframer.ColumnSelector{
				{Selector: "foo", Type: "number"},
				{Selector: "bar1", Alias: "bar", Type: "number"},
				{Selector: "baz", Type: "number"},
			},
		},
		{
			name: "timestamp",
			responseString: `[
				{ "foo" : "2011-01-01T00:00:00.000Z", "bar1": 1325376000000, "baz" : true },
				{ "foo" : "2012-01-01T00:00:00.000Z", "bar1": 1356998400000, "baz" : false }
			]`,
			rootSelector: "",
			columns: []jsonframer.ColumnSelector{
				{Selector: "foo", Type: "timestamp"},
				{Selector: "bar1", Alias: "bar", Type: "timestamp"},
				{Selector: "baz", Type: "timestamp"},
			},
		},
		{
			name: "timestamp_epoch",
			responseString: `{
				"sss": [
					{ "foo" : "1262304000000", "bar1": 1325376000000, "baz" : true },
					{ "foo" : "1293840000000", "bar1": 1356998400000, "baz" : false }
				]
			}`,
			rootSelector: "sss",
			columns: []jsonframer.ColumnSelector{
				{Selector: "foo", Type: "timestamp_epoch"},
				{Selector: "bar1", Alias: "bar", Type: "timestamp_epoch"},
				{Selector: "baz", Type: "timestamp_epoch"},
			},
		},
		{
			name: "timestamp_epoch_s",
			responseString: `[
				{ "foo" : "1262304000", "bar1": 1325376000, "baz" : true },
				{ "foo" : "1293840000", "bar1": 1356998400, "baz" : false }
			]`,
			rootSelector: "",
			columns: []jsonframer.ColumnSelector{
				{Selector: "foo", Type: "timestamp_epoch_s"},
				{Selector: "bar1", Alias: "bar", Type: "timestamp_epoch_s"},
				{Selector: "baz", Type: "timestamp_epoch_s"},
			},
		},
		{
			name: "string with jsonata",
			responseString: `{
				"sss": [
					{ "foo" : "1.2", "bar1": 4, "baz" : true },
					{ "foo" : "3", "bar1": 5.6, "baz" : false }
				]
			}`,
			rootSelector: "sss.foo",
		},
		{
			name: "string with jsonata numbers",
			responseString: `{
				"sss": [
					{ "foo" : "1.2", "bar1": 4, "baz" : true },
					{ "foo" : "3", "bar1": 5.6, "baz" : false }
				]
			}`,
			rootSelector: "sss.bar1",
		},
		{
			name: "string with jsonata summarize",
			responseString: `{
				"sss": [
					{ "foo" : "1.2", "bar1": 4, "baz" : true },
					{ "foo" : "3", "bar1": 5.6, "baz" : false }
				]
			}`,
			rootSelector: "$sum(sss.bar1)",
		},
		{
			name: "eval function",
			responseString: `{
				"inputs" : [
					{
						"a" : 1,
						"b" : "{\"c\":11}"
					},
					{
						"a" : 2,
						"b": "{\"c\":22}"
					}
				]
			}`,
			rootSelector: `$map(inputs,function($v){{
				"a": $v.a,
				"b": $v.b,
				"c": $eval("", $v.b).c
			  }})`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrame, err := jsonframer.ToFrame(tt.responseString, jsonframer.FramerOptions{
				FrameName:    tt.refId,
				RootSelector: tt.rootSelector,
				Columns:      tt.columns,
				ColumnsType:  tt.columnsType,
			})
			if tt.wantErr != nil {
				require.NotNil(t, err)
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, gotFrame)
			goldenFileName := strings.Replace(t.Name(), "TestJsonStringToFrame/", "", 1)
			experimental.CheckGoldenJSONFrame(t, "testdata", goldenFileName, gotFrame, updateTestData)
		})
	}
}

func TestAzureFrame(t *testing.T) {
	fileContent, err := os.ReadFile("./testdata/azure/cost-management-daily.json")
	require.Nil(t, err)
	options := jsonframer.FramerOptions{
		RootSelector: "properties.rows",
		Columns: []jsonframer.ColumnSelector{
			{Selector: "0", Type: "number"},
			{Selector: "1", Type: "number"},
			{Selector: "2", Type: "timestamp", TimeFormat: "20060102"},
			{Selector: "3"},
		},
	}
	var out interface{}
	err = json.Unmarshal(fileContent, &out)
	require.Nil(t, err)
	gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
	require.Nil(t, err)
	require.NotNil(t, gotFrame)
	experimental.CheckGoldenJSONFrame(t, "testdata/azure", "cost-management-daily", gotFrame, false)
}
