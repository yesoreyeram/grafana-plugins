package restds

import (
	"context"
	"errors"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

type pluginHost struct{ IM instancemgmt.InstanceManager }
type datasourceInstance struct {
	PluginName string
	PluginID   string
	Options    RestDriverOptions
	RestDS     RestDS
}

func getInstance(ctx context.Context, pluginCtx backend.PluginContext, im instancemgmt.InstanceManager) (*datasourceInstance, error) {
	instance, err := im.Get(ctx, pluginCtx)
	if err != nil {
		return nil, err
	}
	return instance.(*datasourceInstance), nil
}

func (ins *datasourceInstance) Dispose() {
	backend.Logger.Debug("disposing plugin instance")
}

func NewPlugin(restDriver RestDriver, restDriverOptions RestDriverOptions) datasource.ServeOpts {
	pluginHost := &pluginHost{
		IM: datasource.NewInstanceManager(func(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
			config, err := restDriver.LoadConfig(settings)
			if err != nil {
				return nil, fmt.Errorf("error loading config. %w", err)
			}
			if config == nil {
				return nil, fmt.Errorf("error loading config. %w", errors.New("invalid/empty config"))
			}
			restDs := &RestDS{Config: *config, HTTPClient: NewHTTPClient(config)}
			return &datasourceInstance{
				PluginID:   restDriverOptions.PluginID,
				PluginName: restDriverOptions.PluginName,
				Options:    restDriverOptions,
				RestDS:     *restDs,
			}, nil
		}),
	}
	return datasource.ServeOpts{
		QueryDataHandler:    pluginHost,
		CheckHealthHandler:  pluginHost,
		CallResourceHandler: httpadapter.New(pluginHost.GetRouter(restDriver, restDriverOptions)),
	}
}
