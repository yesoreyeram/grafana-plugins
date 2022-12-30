package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// region Constants
const PluginId = "yesoreyeram-hello-datasource"
const APIToken = "hello"

// endregion

// region Models - Config
type NameTransformMode string

const (
	NameTransformModeNone      NameTransformMode = "none"
	NameTransformModeUpperCase NameTransformMode = "upper_case"
	NameTransformModeLower     NameTransformMode = "lower_case"
)

type Config struct {
	NameTransformMode NameTransformMode `json:"nameTransformMode,omitempty"`
	APIToken          string            `json:"-"`
}

func (config *Config) Validate() error {
	if config.APIToken != APIToken {
		return errors.New("invalid token")
	}
	return nil
}

func LoadConfig(s backend.DataSourceInstanceSettings) (*Config, error) {
	config := &Config{}
	configJson := s.JSONData
	if configJson == nil {
		configJson = []byte("{}")
	}
	if err := json.Unmarshal(configJson, config); err != nil {
		return nil, fmt.Errorf("error while reading the config. %w", err)
	}
	if apiToken, ok := s.DecryptedSecureJSONData["apiToken"]; ok {
		config.APIToken = apiToken
	}
	config = applyDefaultsToConfig(config)
	return config, nil
}

func applyDefaultsToConfig(config *Config) *Config {
	if config.NameTransformMode == "" {
		config.NameTransformMode = NameTransformModeNone
	}
	return config
}

// endregion

// region Models - Query
type QueryType string

const (
	QueryTypeGreet QueryType = "greet"
)

type Query struct {
	RefID     string    `json:"refId"`
	QueryType QueryType `json:"queryType"`
	Greeting  string    `json:"greeting"`
	UserName  string    `json:"username"`
}

func LoadQuery(backendQuery backend.DataQuery, pluginContext backend.PluginContext) (*Query, error) {
	query := &Query{}
	queryJson := backendQuery.JSON
	if queryJson == nil {
		queryJson = []byte("{}")
	}
	if err := json.Unmarshal(queryJson, query); err != nil {
		return nil, fmt.Errorf("error while reading the query. %w", err)
	}
	if query.RefID == "" {
		query.RefID = backendQuery.RefID
	}
	query = applyDefaultsToQuery(query)
	query = applyMacros(query, backendQuery.TimeRange, pluginContext)
	return query, nil
}

func applyDefaultsToQuery(query *Query) *Query {
	if query.QueryType == "" {
		query.QueryType = QueryTypeGreet
	}
	if query.QueryType == QueryTypeGreet {
		if query.UserName == "" {
			query.UserName = "Grafana User"
		}
		if query.Greeting == "" {
			query.Greeting = "Hello"
		}
	}
	return query
}

func applyMacros(query *Query, timeRange backend.TimeRange, pluginContext backend.PluginContext) *Query {
	return query
}

// endregion

// region Datasource Instance
type DatasourceInstance struct {
	Config Config
}

func (is *DatasourceInstance) Dispose() {}

func getInstance(im instancemgmt.InstanceManager, ctx backend.PluginContext) (*DatasourceInstance, error) {
	instance, err := im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return instance.(*DatasourceInstance), nil
}

func GetPluginInstance(backendSettings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	config, err := LoadConfig(backendSettings)
	if err != nil {
		return nil, fmt.Errorf("error getting plugin instance. %w", err)
	}
	return &DatasourceInstance{
		Config: *config,
	}, nil
}

// endregion

// region Plugin Host
type PluginHost struct {
	IM instancemgmt.InstanceManager
}

// endregion

// region Plugin Host - Check Health
func (ds *PluginHost) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	dsi, err := getInstance(ds.IM, req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	if err = dsi.Config.Validate(); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}

// endregion

// region Plugin Host - Query Data
func (ds *PluginHost) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	instance, err := getInstance(ds.IM, req.PluginContext)
	if err != nil {
		backend.Logger.Error("error getting datasource instance from plugin context")
		return nil, fmt.Errorf("error getting datasource instance. %w", err)
	}
	return QueryDataResponse(ctx, instance, req)
}

func QueryDataResponse(ctx context.Context, instance *DatasourceInstance, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	if instance == nil {
		return response, errors.New("invalid instance received")
	}
	for _, q := range req.Queries {
		response.Responses[q.RefID] = QueryData(ctx, q, instance.Config, req.Headers, req.PluginContext)
	}
	return response, nil
}

func QueryData(ctx context.Context, backendQuery backend.DataQuery, config Config, headersFromGrafana map[string]string, pluginContext backend.PluginContext) backend.DataResponse {
	response := &backend.DataResponse{}
	query, err := LoadQuery(backendQuery, pluginContext)
	if err != nil {
		response.Error = err
		return *response
	}
	if err = config.Validate(); err != nil {
		response.Error = err
		return *response
	}
	switch query.QueryType {
	case QueryTypeGreet:
		greeting := fmt.Sprintf("%s %s!", query.Greeting, query.UserName)
		switch config.NameTransformMode {
		case NameTransformModeLower:
			greeting = strings.ToLower(greeting)
		case NameTransformModeUpperCase:
			greeting = strings.ToUpper(greeting)
		}
		greeting = strings.TrimSpace(greeting)
		frame := data.NewFrame("Hello", data.NewField("greeting", nil, []string{greeting}))
		response.Frames = append(response.Frames, frame)
	default:
		response.Error = errors.New("unknown query type")
		return *response
	}
	return *response
}

// endregion

// region Plugin Host - Resource Calls
func (ds *PluginHost) GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("pong")); err != nil {
			backend.Logger.Error("error writing resource call response", "path", "/ping", "error", err.Error())
		}
	})
	router.HandleFunc("/greeting-list", func(w http.ResponseWriter, r *http.Request) {
		greetingList := []string{"Hello", "Welcome"}
		writeResponse(greetingList, nil, w)
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

func writeResponse(resp interface{}, err error, rw http.ResponseWriter) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := rw.Write(b); err != nil {
		backend.Logger.Error("error writing resource call response", "message", err.Error())
	}
}

// endregion

func main() {
	backend.SetupPluginEnvironment(PluginId)
	pluginHost := &PluginHost{IM: datasource.NewInstanceManager(GetPluginInstance)}
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
