package anyframer_test

import (
	"errors"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	af "github.com/yesoreyeram/grafana-plugins/anyframer"
)

func TestStructToFrame(t *testing.T) {
	tests := []struct {
		name      string
		frameName string
		input     any
		wantErr   error
	}{
		{
			name:    "empty input object should throw error",
			wantErr: errors.New("unable to construct frame"),
		},
		{
			name: "map of mixed objects should parse correctly",
			input: map[string]any{
				"valid string":          "foo",
				"valid string pointer":  af.ToPointer("foo"),
				"valid int64":           int64(123),
				"valid int64 pointer":   af.ToPointer(int64(123)),
				"valid int32":           int32(123),
				"valid int32 pointer":   af.ToPointer(int32(123)),
				"valid int16":           int16(123),
				"valid int16 pointer":   af.ToPointer(int16(123)),
				"valid int8":            int8(123),
				"valid int8 pointer":    af.ToPointer(int8(123)),
				"valid int":             int(123),
				"valid int pointer":     af.ToPointer(int(123)),
				"valid uint64":          uint64(123),
				"valid uint64 pointer":  af.ToPointer(uint64(123)),
				"valid uint32":          uint32(123),
				"valid uint32 pointer":  af.ToPointer(uint32(123)),
				"valid uint16":          uint16(123),
				"valid uint16 pointer":  af.ToPointer(uint16(123)),
				"valid uint8":           uint8(123),
				"valid uint8 pointer":   af.ToPointer(uint8(123)),
				"valid uint":            uint(123),
				"valid uint pointer":    af.ToPointer(uint(123)),
				"valid float64":         float64(123.456),
				"valid float64 pointer": af.ToPointer(float64(123.456)),
				"valid float32":         float32(123.456),
				"valid float32 pointer": af.ToPointer(float32(123.456)),
				"valid bool":            true,
				"valid bool pointer":    af.ToPointer(true),
				"valid time":            time.Unix(1, 0),
				"valid time pointer":    af.ToPointer(time.Unix(1, 0)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.name)
			gotFrame, err := af.StructToFrame(t.Name(), tt.input)
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
