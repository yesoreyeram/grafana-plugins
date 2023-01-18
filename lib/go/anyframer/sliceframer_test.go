package anyframer_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	af "github.com/yesoreyeram/grafana-plugins/anyframer"
)

func TestSliceToFrame(t *testing.T) {
	tests := []struct {
		name      string
		frameName string
		input     []any
		columns   []af.Column
		wantErr   error
	}{
		{
			name:  "array of strings",
			input: []any{"hello", "world"},
		},
		{
			name:  "array of numbers",
			input: []any{float64(123.45), float64(6)},
		},
		{
			name:  "array of boolean",
			input: []any{false, true},
		},
		{
			name:  "array of array",
			input: []any{[]any{"hello", "world"}, []any{float64(123.45), float64(6)}},
		},
		{
			name: "array of objects",
			input: []any{
				map[string]any{"username": "foo", "age": 123, "hobbies": []string{"swimming"}, "salaried": true, "country": "uk"},
				map[string]any{"username": "bar", "age": 456.789, "hobbies": []string{"painting"}, "salaried": false, "city": "chennai"},
			},
		},
		{
			name: "array of objects with columns",
			input: []any{
				map[string]any{"username": "foo", "age": 123, "hobbies": []string{"swimming"}, "salaried": true, "country": "uk"},
				map[string]any{"username": "bar", "age": 456.789, "hobbies": []string{"painting"}, "salaried": false, "city": "chennai"},
			},
			columns: []af.Column{
				{Selector: "username", Format: af.ColumnFormatString, Alias: "User Name"},
				{Selector: "age", Format: af.ColumnFormatNumber},
				{Selector: "salaried", Format: af.ColumnFormatBoolean},
			},
		},
		{
			name: "array of objects with columns and nested selector",
			input: []any{
				map[string]any{"username": "foo", "age": 123, "hobbies": []string{"swimming"}, "salaried": true, "country": "uk", "address": map[string]any{"country": "uk", "postcode": 123}},
				map[string]any{"username": "bar", "age": 456.789, "hobbies": []string{"painting"}, "salaried": false, "city": "chennai", "address": map[string]any{"country": "india", "postcode": 567.890}},
			},
			columns: []af.Column{
				{Selector: "username", Format: af.ColumnFormatString, Alias: "User Name"},
				{Selector: "age", Format: af.ColumnFormatNumber},
				{Selector: "salaried", Format: af.ColumnFormatBoolean},
				{Selector: "address.postcode", Format: af.ColumnFormatNumber, Alias: "Post Code"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.name)
			gotFrame, err := af.SliceToFrame(t.Name(), tt.input, tt.columns)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, gotFrame)
			experimental.CheckGoldenJSONFrame(t, "testdata/golden", t.Name(), gotFrame, updateGoldenFile)
		})
	}
}
