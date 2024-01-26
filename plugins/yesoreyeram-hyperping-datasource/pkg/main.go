package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-plugins/lib/go/jsonframer"
)

const pluginId = "yesoreyeram-hyperping-datasource"

func main() {
	backend.SetupPluginEnvironment(pluginId)
	if err := datasource.Manage(pluginId, New, datasource.ManageOpts{}); err != nil {
		backend.Logger.Error("error starting plugin", "pluginId", pluginId, "error", err.Error())
		return
	}
}

//#region Datasource

type DataSource struct {
	HyperPingClient HyperPingClient
	Settings        *Settings
	backend.CallResourceHandler
}

func New(ctx context.Context, dis backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var err error
	settings, err := LoadSettings(dis)
	if err != nil {
		return nil, err
	}
	hyperpingClient, err := NewClient(settings.API_TOKEN)
	if err != nil {
		return nil, err
	}
	ds := &DataSource{
		CallResourceHandler: newResourceHandler(settings, hyperpingClient),
		HyperPingClient:     hyperpingClient,
		Settings:            settings,
	}
	return ds, nil
}

//#endregion

//#region Models

type Settings struct {
	API_TOKEN string `json:"-"`
}

func (s *Settings) ApplyDefaults() error {
	var err error
	return err
}

func (s *Settings) Validate() error {
	var err error
	if s.API_TOKEN == "" {
		return errors.New("invalid/empty api token. get one from https://app.hyperping.io/project/developers")
	}
	return err
}

func LoadSettings(dis backend.DataSourceInstanceSettings) (s *Settings, err error) {
	jsonData := dis.JSONData
	if jsonData == nil {
		jsonData = []byte(`{}`)
	}
	err = json.Unmarshal(jsonData, &s)
	if err != nil {
		return s, err
	}
	s.API_TOKEN = dis.DecryptedSecureJSONData["api_token"]
	if err = s.ApplyDefaults(); err != nil {
		return s, err
	}
	return s, err
}

type QueryType string

const (
	QueryTypeMonitors           QueryType = "monitors"
	QueryTypeMaintenanceWindows QueryType = "maintenance-windows"
)

type Query struct {
	RefID     string            `json:"-"`
	TimeRange backend.TimeRange `json:"-"`
	QueryType QueryType         `json:"queryType,omitempty"`
}

func (q *Query) ApplyDefaults() error {
	var err error
	if q.QueryType == "" {
		q.QueryType = QueryTypeMonitors
	}
	return err
}

func (q *Query) Validate() error {
	var err error
	return err
}

func LoadQuery(input backend.DataQuery) (Query, error) {
	q := Query{}
	jsonData := input.JSON
	if jsonData == nil {
		jsonData = []byte(`{}`)
	}
	err := json.Unmarshal(jsonData, &q)
	if err != nil {
		return q, err
	}
	q.RefID = input.RefID
	q.TimeRange = input.TimeRange
	return q, err
}

//#endregion

//#region Health Check

func (ds *DataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	err := ds.Settings.Validate()
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	res, err := ds.HyperPingClient.GetMonitors(Query{RefID: "health-check"})
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("hyperping datasource working. %d monitors found", res.Rows()),
	}, nil
}

//#endregion

//#region Query

func (ds *DataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		query, err := LoadQuery(q)
		if err != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusInternal, err.Error())
			continue
		}
		switch query.QueryType {
		case QueryTypeMaintenanceWindows:
			frame, err := ds.HyperPingClient.GetMaintenanceWindows(query)
			if err != nil {
				response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusInternal, err.Error())
				continue
			}
			for _, field := range frame.Fields {
				if field.Config == nil {
					field.Config = &data.FieldConfig{}
				}
				if field.Name == "uuid" {
					field.Config.Links = []data.DataLink{
						{
							Title:       "Open in hyperping",
							URL:         "https://app.hyperping.io/maintenance/${__data.fields.uuid}",
							TargetBlank: true,
						},
					}
				}
			}
			response.Responses[query.RefID] = backend.DataResponse{Frames: data.Frames{frame}}
		case QueryTypeMonitors:
			frame, err := ds.HyperPingClient.GetMonitors(query)
			if err != nil {
				response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusInternal, err.Error())
				continue
			}
			for _, field := range frame.Fields {
				if field.Config == nil {
					field.Config = &data.FieldConfig{}
				}
				if field.Name == "uuid" {
					field.Config.Links = []data.DataLink{
						{
							Title:       "Open in hyperping",
							URL:         "https://app.hyperping.io/report/${__data.fields.uuid}",
							TargetBlank: true,
						},
					}
				}
			}
			response.Responses[query.RefID] = backend.DataResponse{Frames: data.Frames{frame}}
		default:
			response.Responses[query.RefID] = backend.ErrDataResponse(backend.StatusInternal, "invalid query type")
		}
	}
	return response, nil
}

//#endregion

//#region Resource Calls

func newResourceHandler(settings *Settings, client HyperPingClient) backend.CallResourceHandler {
	router := mux.NewRouter()
	router.HandleFunc("/ping", WithDatasource(client, GetPingHandler)).Methods("GET")
	return httpadapter.New(router)
}

func WithDatasource(client HyperPingClient, getHandler func(client HyperPingClient) http.HandlerFunc) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		h := getHandler(client)
		h.ServeHTTP(rw, r)
	}
}

func GetPingHandler(client HyperPingClient) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		writeResponse("pong", nil, rw)
	}
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
	_, _ = rw.Write(b)
}

//#endregion

//#region HyperPing Client

type HyperPingClient struct {
	HTTP_CLIENT *http.Client
}

func NewClient(apiToken string) (HyperPingClient, error) {
	hc, err := httpclient.New(httpclient.Options{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", apiToken),
		},
	})
	client := HyperPingClient{HTTP_CLIENT: hc}
	return client, err
}

func (hpc *HyperPingClient) GetMonitors(query Query) (*data.Frame, error) {
	res, err := hpc.get("https://api.hyperping.io/v1/monitors")
	if err != nil {
		return nil, err
	}
	return jsonframer.ToFrame(res, jsonframer.FramerOptions{FramerType: jsonframer.FramerTypeGJSON})
}

func (hpc *HyperPingClient) GetMaintenanceWindows(query Query) (*data.Frame, error) {
	res, err := hpc.get("https://api.hyperping.io/v1/maintenance-windows")
	if err != nil {
		return nil, err
	}
	return jsonframer.ToFrame(res, jsonframer.FramerOptions{
		FramerType:   jsonframer.FramerTypeGJSON,
		FrameName:    query.RefID,
		RootSelector: "maintenanceWindows",
		Columns: []jsonframer.ColumnSelector{
			{Selector: "name", Type: "string"},
			{Selector: "uuid", Type: "string"},
			{Selector: "start_date", Type: "timestamp"},
			{Selector: "end_date", Type: "timestamp"},
		}})
}

func (hpc *HyperPingClient) get(u string) (string, error) {
	if u == "" {
		u = "https://api.hyperping.io/v1/monitors"
	}
	res, err := hpc.HTTP_CLIENT.Get(u)
	if err != nil {
		return "", err
	}
	if res != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code received from hyperping. status code: %s", res.Status)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil

}

//#endregion
