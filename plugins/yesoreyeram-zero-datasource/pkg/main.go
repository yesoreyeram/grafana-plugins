package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
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
		r := &backend.DataResponse{}
		r.Error = errors.New("zero datasource query not implemented")
		response.Responses[q.RefID] = *r
	}
	return response, nil
}
func (ds *PluginHost) GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("oops.. zero datasource resource calls not implemented")); err != nil {
			backend.Logger.Error("error writing resource call response", "path", "/404", "error", err.Error())
		}
	})
	return router
}

func main() {
	backend.SetupPluginEnvironment(PluginId)
	pluginHost := &PluginHost{
		IM: datasource.NewInstanceManager(func(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
			return &DatasourceInstance{}, nil
		}),
	}
	pluginServer := datasource.ServeOpts{
		QueryDataHandler:    pluginHost,
		CheckHealthHandler:  pluginHost,
		CallResourceHandler: httpadapter.New(pluginHost.GetRouter()),
	}
	if err := datasource.Serve(pluginServer); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
