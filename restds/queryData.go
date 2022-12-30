package restds

import (
	"context"
	"errors"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/yesoreyeram/grafana-plugins/anyframer"
)

func (ds *pluginHost) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	instance, err := getInstance(ds.IM, req.PluginContext)
	if err != nil {
		backend.Logger.Error("error getting datasource instance from plugin context")
		return nil, fmt.Errorf("error getting datasource instance. %w", err)
	}
	return QueryDataResponse(ctx, instance, req)
}

func QueryDataResponse(ctx context.Context, instance *datasourceInstance, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	if instance == nil {
		return response, errors.New("invalid instance received")
	}
	for _, q := range req.Queries {
		response.Responses[q.RefID] = QueryData(ctx, q, *instance, req.Headers, req.PluginContext)
	}
	return response, nil
}

func QueryData(ctx context.Context, backendQuery backend.DataQuery, instance datasourceInstance, headersFromGrafana map[string]string, pluginContext backend.PluginContext) backend.DataResponse {
	response := &backend.DataResponse{}
	query, err := LoadQuery(backendQuery, pluginContext)
	if err != nil {
		response.Error = err
		return *response
	}
	switch query.QueryType {
	default:
		body, meta, err := instance.RestDS.GetResponse(*query)
		if err != nil {
			response.Error = err
			return *response
		}
		framer := anyframer.AnyFramer{
			InputType:    anyframer.InputTypeJSON,
			RootSelector: query.RootSelector,
			Headers:      meta.Headers,
		}
		f, err := framer.ToFrame(body)
		if err != nil {
			response.Error = err
			return *response
		}
		response.Frames = append(response.Frames, f)
	}
	return *response
}
