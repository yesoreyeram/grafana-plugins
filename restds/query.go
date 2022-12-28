package restds

import (
	"net/http"
)

type Query struct {
	RefID            string         `json:"refId"`
	QueryType        QueryType      `json:"type"`
	Data             string         `json:"data,omitempty"`
	URL              string         `json:"url,omitempty"`
	Method           QueryURLMethod `json:"method,omitempty"`
	Headers          []KV           `json:"headers,omitempty"`
	BodyType         BodyType       `json:"bodyType,omitempty"`
	Body             string         `json:"body,omitempty"`
	BodyContentType  string         `json:"bodyContentType,omitempty"`
	BodyForm         []KV           `json:"bodyForm,omitempty"`
	BodyGraphQLQuery string         `json:"bodyGraphQLQuery,omitempty"`
}

type QueryType string

const (
	QueryTypeAuto QueryType = "auto"
	QueryTypeJSON QueryType = "json"
	QueryTypeCSV  QueryType = "csv"
	QueryTypeTSV  QueryType = "tsv"
	QueryTypeXML  QueryType = "xml"
	QueryTypeHTML QueryType = "html"
)

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
