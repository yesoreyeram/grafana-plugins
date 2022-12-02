package anyframer_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/anyframer"
)

func TestGuessType(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		headers http.Header
		input   string
		want    anyframer.InputType
	}{
		{},
		{rawURL: "https://foo"},
		{rawURL: "https://foo.com"},
		{rawURL: "https://foo.com/bar"},
		{rawURL: "https://foo.com/bar.json", want: anyframer.InputTypeJSON},
		{rawURL: "https://foo.com/bar.json?something=nothing", want: anyframer.InputTypeJSON},
		{rawURL: "foo.yaml"},
		{rawURL: "foo.json", want: anyframer.InputTypeJSON},
		{rawURL: "foo/bar.csv", want: anyframer.InputTypeCSV},
		{rawURL: "https://foo.com/bar", headers: http.Header{}},
		{rawURL: "https://foo.com/bar", headers: map[string][]string{"something": {"nothing"}}},
		{rawURL: "https://foo.com/bar", headers: map[string][]string{"Content-Type": {"nothing"}}},
		{rawURL: "https://foo.com/bar", headers: map[string][]string{"Content-Type": {"application/json"}}, want: anyframer.InputTypeJSON},
		{rawURL: "https://foo.com/bar", input: "hello"},
		{rawURL: "https://foo.com/bar", input: " { \"foo\" : 123 } ", want: anyframer.InputTypeJSON},
		{rawURL: "https://foo.com/bar", input: " [1,2,3] ", want: anyframer.InputTypeJSON},
		{rawURL: "https://foo.com/bar", input: "a	b	c\n1	2	3", want: anyframer.InputTypeTSV},
		{rawURL: "https://foo.com/bar", input: "a,b,c\n1,2,3", want: anyframer.InputTypeCSV},
		{rawURL: "https://foo.com/bar", input: "<html ></html>", want: anyframer.InputTypeHTML},
		{rawURL: "https://foo.com/bar", input: "<!DOCTYPE html><html></html>", want: anyframer.InputTypeHTML},
		{rawURL: "https://foo.com/bar", input: "<!doctype html><html></html>", want: anyframer.InputTypeHTML},
		{rawURL: "https://foo.com/bar", input: "<xml></xml>", want: anyframer.InputTypeXML},
		{rawURL: "https://foo.com/bar", input: "<xml ></xml>", want: anyframer.InputTypeXML},
		{rawURL: "https://foo.com/bar", input: "<?xml ></xml>", want: anyframer.InputTypeXML},
		{rawURL: "https://foo.com/bar", input: "<rss></rss>", want: anyframer.InputTypeXML},
	}
	for _, tt := range tests {
		want := tt.want
		if want == "" {
			want = anyframer.InputTypeUnknown
		}
		t.Run(tt.name, func(t *testing.T) {
			framer := anyframer.AnyFramer{RawURL: tt.rawURL, Headers: tt.headers}
			got := framer.GuessType(tt.input)
			require.Equal(t, want, got)
		})
	}
}
