package csvframer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-plugins/lib/go/gframer"
)

type FramerOptions struct {
	FrameName          string
	Columns            []gframer.ColumnSelector
	Delimiter          string
	SkipLinesWithError bool
	Comment            string
	RelaxColumnCount   bool
	NoHeaders          bool
}

func ToFrame(csvString string, options FramerOptions) (frame *data.Frame, err error) {
	if strings.TrimSpace(csvString) == "" {
		return frame, errors.New("empty/invalid csv")
	}
	r := csv.NewReader(strings.NewReader(csvString))
	r.LazyQuotes = true
	if options.Comment != "" {
		r.Comment = rune(options.Comment[0])
	}
	if options.Delimiter != "" {
		r.Comma = rune(options.Delimiter[0])
	}
	if options.RelaxColumnCount {
		r.FieldsPerRecord = -1
	}
	parsedCSV := [][]string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err == nil {
			parsedCSV = append(parsedCSV, record)
			continue
		}
		if !options.SkipLinesWithError {
			return frame, fmt.Errorf("error reading csv response. %w, %v", err, record)
		}
	}
	out := []interface{}{}
	header := []string{}
	records := [][]string{}
	if !options.NoHeaders {
		header = parsedCSV[0]
		for idx, hItem := range header {
			for _, col := range options.Columns {
				if col.Selector == hItem && col.Alias != "" {
					header[idx] = col.Alias
				}
			}
		}
		records = parsedCSV[1:]
	}
	if options.NoHeaders {
		records = parsedCSV
		if len(records) > 0 {
			for i := 0; i < len(records[0]); i++ {
				header = append(header, fmt.Sprintf("%d", i+1))
			}
		}
	}
	for _, row := range records {
		item := map[string]interface{}{}
		for colId, col := range header {
			if colId < len(row) {
				item[col] = row[colId]
			}
		}
		out = append(out, item)
	}
	framerOptions := gframer.FramerOptions{
		FrameName: options.FrameName,
		Columns:   options.Columns,
	}
	return gframer.ToDataFrame(out, framerOptions)
}
