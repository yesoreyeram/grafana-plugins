package restds

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (ds *pluginHost) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	dsi, err := getInstance(ds.IM, req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	if dsi.Options.HealthCheckURL == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusOk,
			Message: fmt.Sprintf("%s datasource plugin works", dsi.PluginName),
		}, nil
	}
	responseString, meta, err := dsi.RestDS.GetResponse(Query{URL: dsi.Options.HealthCheckURL})
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	if meta.StatusCode != 200 {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Errorf("invalid status code from api. %d", meta.StatusCode).Error(),
		}, nil
	}
	if dsi.Options.CustomHealthCheckValidation != nil {
		return dsi.Options.CustomHealthCheckValidation(responseString, meta), nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("%s datasource plugin works", dsi.PluginName),
	}, nil
}
