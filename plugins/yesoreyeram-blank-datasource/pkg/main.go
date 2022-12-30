package main

import (
	"context"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const PluginId = "yesoreyeram-blank-datasource"

type DatasourceInstance struct{}
type PluginHost struct{ IM instancemgmt.InstanceManager }

func (ins *DatasourceInstance) Dispose() {
	backend.Logger.Debug("disposing plugin instance")
}
func (ds *PluginHost) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "blank datasource just works but does nothing",
	}, nil
}
func (ds *PluginHost) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		frame := data.NewFrame(
			q.QueryType, data.NewField("response", nil, []string{"blank response"}),
		).SetMeta(
			&data.FrameMeta{Notices: []data.Notice{{Text: "blank query works. but not fully implemented"}}},
		)
		response.Responses[q.RefID] = backend.DataResponse{
			Frames: data.Frames{frame},
			Status: backend.StatusOK,
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
