package restds

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type QueryType string

const (
	QueryTypeOpenAPI QueryType = "openApi3"
	QueryTypeAuto    QueryType = "auto"
	QueryTypeJSON    QueryType = "json"
	QueryTypeCSV     QueryType = "csv"
	QueryTypeTSV     QueryType = "tsv"
	QueryTypeXML     QueryType = "xml"
	QueryTypeHTML    QueryType = "html"
)

type Query struct {
	RefID            string         `json:"refId"`
	QueryType        QueryType      `json:"type"`
	URL              string         `json:"url,omitempty"`
	Method           QueryURLMethod `json:"method,omitempty"`
	Headers          []KV           `json:"headers,omitempty"`
	BodyType         BodyType       `json:"bodyType,omitempty"`
	Body             string         `json:"body,omitempty"`
	BodyContentType  string         `json:"bodyContentType,omitempty"`
	BodyForm         []KV           `json:"bodyForm,omitempty"`
	BodyGraphQLQuery string         `json:"bodyGraphQLQuery,omitempty"`
	RootSelector     string         `json:"rootSelector,omitempty"`
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
	return query, nil
}

type QueryURLMethod string

const (
	QueryURLMethodGet  QueryURLMethod = http.MethodGet
	QueryURLMethodPost QueryURLMethod = http.MethodPost
)

type BodyType string

const (
	BodyTypeNone         BodyType = "none"
	BodyTypeRaw          BodyType = "raw"
	BodyTypeFormData     BodyType = "form-data"
	BodyTypeFormReloaded BodyType = "x-www-form-urlencoded"
	BodyTypeGraphQL      BodyType = "graphql"
)

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
