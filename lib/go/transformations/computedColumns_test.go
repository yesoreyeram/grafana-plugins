package transformations_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/lib/go/transformations"
)

func TestGetFrameWithComputedColumns(t *testing.T) {
	tests := []struct {
		name    string
		input   *data.Frame
		columns []transformations.ComputedColumn
		output  *data.Frame
	}{
		{
			input:   data.NewFrame("sample", data.NewField("id", nil, []*float64{toFP(1), toFP(2)})),
			columns: []transformations.ComputedColumn{{Selector: "id + id", Text: "twice"}},
			output:  data.NewFrame("sample", data.NewField("id", nil, []*float64{toFP(1), toFP(2)}), data.NewField("twice", nil, []*float64{toFP(2), toFP(4)})),
		},
		{
			input:   data.NewFrame("sample", data.NewField("Cylinders", nil, []*float64{toFP(1), toFP(2)}), data.NewField("Horsepower", nil, []*float64{toFP(2.3), toFP(4.5)})),
			columns: []transformations.ComputedColumn{{Selector: "[Cylinders] + horsepower", Text: "power"}},
			output:  data.NewFrame("sample", data.NewField("Cylinders", nil, []*float64{toFP(1), toFP(2)}), data.NewField("Horsepower", nil, []*float64{toFP(2.3), toFP(4.5)}), data.NewField("power", nil, []*float64{toFP(3.3), toFP(6.5)})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformations.GetFrameWithComputedColumns(tt.input, tt.columns)
			require.Nil(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.output, got)
		})
	}
}

func toFP(v float64) *float64 {
	return &v
}
