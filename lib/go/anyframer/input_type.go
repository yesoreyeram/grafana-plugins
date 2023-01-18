package anyframer

import (
	"net/http"
	"net/url"
	"strings"
)

// InputType ...
type InputType string

const (
	// InputTypeUnknown ...
	InputTypeUnknown InputType = "unknown"
	// InputTypeJSON ...
	InputTypeJSON InputType = "json"
	// InputTypeCSV ...
	InputTypeCSV InputType = "csv"
	// InputTypeTSV ...
	InputTypeTSV InputType = "tsv"
	// InputTypeHTML ...
	InputTypeHTML InputType = "html"
	// InputTypeXML ...
	InputTypeXML InputType = "xml"
)

// GuessType guesses the framer type from input string and other framer parameters such as RawURL, headers
func (options *AnyFramer) GuessType(input string) InputType {
	o := options.InputType
	if o == "" || o == InputTypeUnknown {
		o = guessInputTypeFromURL(options.RawURL)
		if o == InputTypeUnknown && len(options.Headers) > 0 {
			o = guessInputTypeFromResponseHeaders(options.Headers)
		}
		if o == InputTypeUnknown && input != "" {
			o = guessInputTypeFromInput(input)
		}
	}
	return o
}

func guessInputTypeFromURL(rawURL string) InputType {
	rawURL = strings.ToLower(rawURL)
	if u, err := url.Parse(rawURL); err == nil {
		rawURL = strings.ToLower(u.Path)
	}
	return guessInputTypeFromFileName(rawURL)
}

func guessInputTypeFromFileName(fileName string) InputType {
	fileName = strings.ToLower(fileName)
	if strings.HasSuffix(fileName, ".json") {
		return InputTypeJSON
	}
	if strings.HasSuffix(fileName, ".csv") {
		return InputTypeCSV
	}
	if strings.HasSuffix(fileName, ".tsv") {
		return InputTypeTSV
	}
	if strings.HasSuffix(fileName, ".xml") {
		return InputTypeXML
	}
	if strings.HasSuffix(fileName, ".html") {
		return InputTypeHTML
	}
	return InputTypeUnknown
}

func guessInputTypeFromResponseHeaders(headers http.Header) InputType {
	contentType := strings.ToLower(headers.Get("Content-Type"))
	if strings.HasPrefix(contentType, "application/json") {
		return InputTypeJSON
	}
	if strings.HasPrefix(contentType, "text/csv") {
		return InputTypeCSV
	}
	if strings.HasPrefix(contentType, "text/tab-separated-values") {
		return InputTypeTSV
	}
	if strings.HasPrefix(contentType, "application/xml") || strings.HasPrefix(contentType, "text/xml") {
		return InputTypeXML
	}
	return InputTypeUnknown
}

func guessInputTypeFromInput(input string) InputType {
	input = strings.TrimSpace(strings.ToLower(input))
	if strings.HasPrefix(input, "{") && strings.HasSuffix(input, "}") {
		return InputTypeJSON
	}
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		return InputTypeJSON
	}
	if (strings.HasPrefix(input, "<!doctype html") || strings.HasPrefix(input, "<html")) && strings.HasSuffix(input, ">") {
		return InputTypeHTML
	}
	if strings.HasPrefix(input, "<") && strings.HasSuffix(input, ">") {
		return InputTypeXML
	}
	csvArray := strings.Split(input, "\n")
	if len(csvArray) > 1 && strings.Contains(csvArray[0], "\t") {
		return InputTypeTSV
	}
	if len(csvArray) > 1 && strings.Contains(csvArray[0], ",") {
		return InputTypeCSV
	}
	return InputTypeUnknown
}
