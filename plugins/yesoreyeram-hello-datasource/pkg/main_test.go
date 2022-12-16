package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/require"
)

const updateGoldenFile = false

func TestQueryData(t *testing.T) {
	tests := []struct {
		name               string
		queryJson          json.RawMessage
		config             Config
		headersFromGrafana map[string]string
		pluginContext      backend.PluginContext
		want               backend.DataResponse
		wantError          error
	}{
		{
			name:      "empty token should throw error",
			wantError: errors.New("invalid token"),
		},
		{
			name:      "invalid token should throw error",
			config:    Config{APIToken: "welcome"},
			wantError: errors.New("invalid token"),
		},
		{
			name:   "valid token should not throw error",
			config: Config{APIToken: "hello"},
		},
		{
			name:      "valid greeting query",
			queryJson: []byte(`{ "username" : "Sriram"}`),
			config:    Config{APIToken: "hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.name)
			queryJson := tt.queryJson
			if strings.TrimSpace(string(queryJson)) == "" {
				queryJson = []byte(`{}`)
			}
			got := QueryData(context.Background(), backend.DataQuery{JSON: queryJson}, tt.config, tt.headersFromGrafana, tt.pluginContext)
			require.NotNil(t, got)
			if tt.wantError != nil {
				require.NotNil(t, got.Error)
				return
			}
			require.Nil(t, got.Error)
			experimental.CheckGoldenJSONResponse(t, "testdata", t.Name(), &got, updateGoldenFile)
		})
	}
}
