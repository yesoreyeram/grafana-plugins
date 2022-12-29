package restds

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/swaggest/openapi-go/openapi3"
)

type RestDriverOptions struct {
	PluginName                  string
	PluginID                    string
	HealthCheckURL              string
	CustomHealthCheckValidation func(responseString string, meta ResponseMeta) *backend.CheckHealthResult
}

type RestDriver interface {
	LoadConfig(settings backend.DataSourceInstanceSettings) (*Config, error)
	LoadSpec() openapi3.Spec
}
