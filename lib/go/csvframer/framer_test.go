package csvframer_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/lib/go/csvframer"
	"github.com/yesoreyeram/grafana-plugins/lib/go/gframer"
)

func TestCsvStringToFrame(t *testing.T) {
	tests := []struct {
		name      string
		csvString string
		options   csvframer.CSVFramerOptions
		wantError error
	}{
		{
			name:      "empty csv should return error",
			wantError: errors.New("empty/invalid csv"),
		},
		{
			name:      "valid csv should not return error",
			csvString: strings.Join([]string{`a,b,c`, `1,2,3`, `11,12,13`, `21,22,23`}, "\n"),
		},
		{
			name:      "valid csv without headers should not return error",
			csvString: strings.Join([]string{`1,2,3`, `11,12,13`, `21,22,23`}, "\n"),
			options:   csvframer.CSVFramerOptions{NoHeaders: true},
		},
		{
			name:      "framer options should be respected",
			csvString: strings.Join([]string{`a	b	c`, `1	2	3`, `11	12	13`, `21	22	23`}, "\n"),
			options: csvframer.CSVFramerOptions{FrameName: "foo", Delimiter: "\t", RelaxColumnCount: true, Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "timestamp_epoch"},
			}},
		},
		{
			name:      "relax column count",
			csvString: strings.Join([]string{`a	b	c`, `1	2	3`, `11	12`, `21	22	23`}, "\n"),
			options: csvframer.CSVFramerOptions{FrameName: "foo", Delimiter: "\t", SkipLinesWithError: true, Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "timestamp_epoch"},
			}},
		},
		{
			name:      "Skip empty lines",
			csvString: strings.Join([]string{`a	b	c`, `1	2	3`, ``, `21	22	23`}, "\n"),
			options: csvframer.CSVFramerOptions{FrameName: "foo", Delimiter: "\t", Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "timestamp_epoch_s"},
			}},
		},
		{
			name:      "relax column count",
			csvString: strings.Join([]string{`a;b;c`, `1;2;3`, `11;13`, `21;22;23`}, "\n"),
			options: csvframer.CSVFramerOptions{FrameName: "foo", Delimiter: ";", RelaxColumnCount: true, Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "string"},
			}},
		},
		{
			name:      "comment",
			csvString: strings.Join([]string{`# foo`, `a,b,c`, `#01,02,03`, `1,2,3`, `11,12,13`, `21,22,23`, `#`}, "\n"),
			options:   csvframer.CSVFramerOptions{Comment: "#"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrame, err := csvframer.CsvStringToFrame(tt.csvString, tt.options)
			if tt.wantError != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantError, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, gotFrame)
			if tt.wantError == nil {
				goldenFileName := strings.Replace(t.Name(), "TestCsvStringToFrame/", "", 1)
				experimental.CheckGoldenJSONFrame(t, "testdata", goldenFileName, gotFrame, false)
			}
		})
	}
}
