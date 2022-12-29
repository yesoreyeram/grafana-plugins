package restds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/yesoreyeram/grafana-plugins/restds"
)

func TestSpecFromJson(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
		want       openapi3.Spec
		wantErr    error
		test       func(t *testing.T, got openapi3.Spec)
	}{
		{
			jsonString: `{
				"openapi"	: 	"3.0.3",
				"info"		:	{
					"title" 	: "Sample API",
					"version" 	: "0.0.1"
				},
				"servers"	: 	[
				],
				"paths": {}
			}`,
			test: func(t *testing.T, got openapi3.Spec) {
				require.Equal(t, "0.0.1", got.Info.Version)
				require.Equal(t, "Sample API", got.Info.Title)
				require.Equal(t, "3.0.3", got.Openapi)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := restds.SpecFromJson(tt.jsonString)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, got)
			if tt.test != nil {
				tt.test(t, got)
			}
		})
	}
}
