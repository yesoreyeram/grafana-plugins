package main

import (
	"context"
	"errors"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

const PluginId = "yesoreyeram-zero-datasource"

type DatasourceInstance struct{}
type PluginHost struct{ IM instancemgmt.InstanceManager }

func (ins *DatasourceInstance) Dispose() {
	backend.Logger.Debug("disposing plugin instance")
}
func (ds *PluginHost) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusError,
		Message: "zero datasource health check not implemented",
	}, nil
}
func (ds *PluginHost) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		response.Responses[q.RefID] = backend.DataResponse{
			Error:  errors.New("zero datasource query not implemented"),
			Status: backend.StatusNotImplemented,
		}
	}
	return response, nil
}
func main() {
	backend.SetupPluginEnvironment(PluginId)
	pluginHost := &PluginHost{
		IM: datasource.NewInstanceManager(func(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
			return &DatasourceInstance{}, nil
		}),
	}
	pluginServer := datasource.ServeOpts{
		QueryDataHandler:   pluginHost,
		CheckHealthHandler: pluginHost,
	}
	if err := datasource.Serve(pluginServer); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
