package restds

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// region DataSource Instance
type datasourceInstance struct {
	PluginName string
	PluginID   string
	Config     Config
}

func getInstance(im instancemgmt.InstanceManager, ctx backend.PluginContext) (*datasourceInstance, error) {
	instance, err := im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return instance.(*datasourceInstance), nil
}

func (ins *datasourceInstance) Dispose() {
	backend.Logger.Debug("disposing plugin instance")
}

//endregion

// region Plugin Host
type pluginHost struct{ IM instancemgmt.InstanceManager }

func (ds *pluginHost) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	dsi, err := getInstance(ds.IM, req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("%s datasource plugin works", dsi.PluginName),
	}, nil
}
func (ds *pluginHost) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
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
func (ds *pluginHost) GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("pong")); err != nil {
			backend.Logger.Error("error writing resource call response", "path", "/ping", "error", err.Error())
		}
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backend.Logger.Debug("resource call received", "url", r.URL.String())
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("oops.. resource not found")); err != nil {
			backend.Logger.Error("error writing resource call response", "path", "/404", "error", err.Error())
		}
	})
	return router
}

//endregion

// region Rest Driver
type RestDriverOptions struct {
	PluginName string
	PluginID   string
}
type RestDriver interface {
	LoadConfig(settings backend.DataSourceInstanceSettings) (*Config, error)
}

//endregion

func NewPlugin(restDriver RestDriver, restDriverOptions RestDriverOptions) datasource.ServeOpts {
	pluginHost := &pluginHost{
		IM: datasource.NewInstanceManager(func(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
			config, err := restDriver.LoadConfig(settings)
			if err != nil {
				return nil, fmt.Errorf("error loading config. %w", err)
			}
			return &datasourceInstance{
				PluginID:   restDriverOptions.PluginID,
				PluginName: restDriverOptions.PluginName,
				Config:     *config,
			}, nil
		}),
	}
	return datasource.ServeOpts{
		QueryDataHandler:    pluginHost,
		CheckHealthHandler:  pluginHost,
		CallResourceHandler: httpadapter.New(pluginHost.GetRouter()),
	}
}
