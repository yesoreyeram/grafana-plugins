package gframer_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/lib/go/gframer"
)

func TestToDataFrame(t *testing.T) {
	updateGoldenText := false
	t.Run("nil", func(t *testing.T) {
		var input interface{}
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: ""}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("string", func(t *testing.T) {
		input := `foo`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "foo"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
		require.Nil(t, err)
	})
	t.Run("integer", func(t *testing.T) {
		input := 21
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "21"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("float", func(t *testing.T) {
		input := 21.43
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "21.43"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("bool", func(t *testing.T) {
		input := true
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "true"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("object", func(t *testing.T) {
		input := `{ "name":"foo", "age": 12, "hobbies":["cricket","music"], "isPrimeUser": true, "fullname": { "first": "foo", "last":"bar" } }`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("string-array", func(t *testing.T) {
		input := `["foo","bar"]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("number-array", func(t *testing.T) {
		input := `[12,14.56,0,30]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("array-inside-array", func(t *testing.T) {
		input := `[["one","two"],["three"]]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("all-null-array", func(t *testing.T) {
		input := `[null,null]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
}

func TestToDataFrameSlices(t *testing.T) {
	updateGoldenText := false
	files, err := os.ReadDir("./testdata/slices")
	if err != nil {
		require.Nil(t, err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			t.Run(f.Name(), func(t *testing.T) {
				fileContent, err := os.ReadFile("./testdata/slices/" + f.Name())
				require.Nil(t, err)
				options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: ""}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := gframer.ToDataFrame(out, options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				goldenFileName := strings.TrimSuffix(f.Name(), ".json")
				experimental.CheckGoldenJSONFrame(t, "testdata/slices", goldenFileName, gotFrame, updateGoldenText)
			})
		}
	}
}
