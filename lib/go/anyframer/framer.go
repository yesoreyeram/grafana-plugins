package anyframer

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	xj "github.com/basgys/goxml2json"
)

func (options *AnyFramer) toJSONObject(input string) (any, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("invalid/empty input")
	}
	switch options.InputType {
	case InputTypeJSON:
		return toObjectFromJSONString(input, *options)
	case InputTypeTSV:
		options.CSVOptions.Delimiter = "\t"
		return toObjectFromCSVString(input, *options)
	case InputTypeCSV:
		return toObjectFromCSVString(input, *options)
	case InputTypeXML:
		return toObjectFromXMLString(input, *options)
	case InputTypeHTML:
		return toObjectFromHTMLString(input, *options)
	case InputTypeUnknown:
		var o any
		if err := json.Unmarshal([]byte(input), &o); err != nil {
			return input, nil
		}
		return o, nil
	default:
		return nil, errors.New("error converting input to JSON")
	}
}

func toObjectFromJSONString(jsonString string, options AnyFramer) (jsonObject any, err error) {
	err = json.Unmarshal([]byte(jsonString), &jsonObject)
	if err != nil {
		return nil, err
	}
	return jsonObject, nil
}

func toObjectFromCSVString(csvString string, options AnyFramer) (jsonObject any, err error) {
	out := []any{}
	delimiter := options.CSVOptions.Delimiter
	if delimiter == "" {
		delimiter = ","
	}
	if len(options.CSVOptions.Headers) > 0 {
		headers := []string{}
		for _, h := range options.CSVOptions.Headers {
			h = strings.TrimSpace(h)
			if strings.HasPrefix(h, `"`) && strings.HasSuffix(h, `"`) {
				headers = append(headers, h)
				continue
			}
			headers = append(headers, `"`+h+`"`)
		}
		csvString = strings.Join(headers, delimiter) + "\n" + csvString
	}
	r := csv.NewReader(strings.NewReader(csvString))
	r.LazyQuotes = true
	if options.CSVOptions.Comment != "" {
		r.Comment = rune(options.CSVOptions.Comment[0])
	}
	if delimiter != "" {
		r.Comma = rune(delimiter[0])
	}
	if options.CSVOptions.RelaxColumnCount {
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
		if !options.CSVOptions.SkipLinesWithError {
			return nil, fmt.Errorf("error reading csv response. %w, %v", err, record)
		}
	}
	if len(parsedCSV) == 0 {
		return nil, errors.New("invalid/empty csv")
	}
	header := []string{}
	records := [][]string{}
	if !options.CSVOptions.NoHeaders {
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
	if options.CSVOptions.NoHeaders {
		records = parsedCSV
		if len(records) > 0 {
			for i := 0; i < len(records[0]); i++ {
				header = append(header, fmt.Sprintf("%d", i+1))
			}
		}
	}
	for _, row := range records {
		item := map[string]any{}
		for hID, h := range header {
			if hID < len(row) {
				item[h] = row[hID]
			}
		}
		out = append(out, item)
	}
	return out, nil
}

func toObjectFromHTMLString(htmlString string, options AnyFramer) (jsonObject any, err error) {
	return toObjectFromXMLString(htmlString, options)
}

func toObjectFromXMLString(xmlString string, options AnyFramer) (jsonObject any, err error) {
	o, err := xmlContentToJSONString(xmlString, options)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(o), &jsonObject)
	if err != nil {
		return nil, err
	}
	return jsonObject, nil
}

func xmlContentToJSONString(xmlString string, options AnyFramer) (jsonString string, err error) {
	xml := strings.NewReader(xmlString)
	json, err := xj.Convert(xml)
	if err != nil {
		return "", err
	}
	o := json.String()
	if strings.HasSuffix(o, "\n") {
		return strings.TrimSuffix(o, "\n"), nil
	}
	return json.String(), nil
}
