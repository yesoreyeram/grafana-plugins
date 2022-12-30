package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/yesoreyeram/grafana-plugins/restds"
)

type PetStoreRestDriver struct{}

func (v *PetStoreRestDriver) LoadConfig(settings backend.DataSourceInstanceSettings) (*restds.Config, error) {
	return &restds.Config{
		AuthenticationMethod: restds.AuthTypeNone,
	}, nil
}

func (v *PetStoreRestDriver) LoadSpec() openapi3.Spec {
	specUrl := "https://petstore3.swagger.io/api/v3/openapi.json"
	spec, err := loadSpecFromUrl(specUrl)
	if err != nil {
		return openapi3.Spec{Openapi: "3.0.3"}
	}
	// region Workaround for incorrect server url in spec
	for i, s := range spec.Servers {
		if strings.HasPrefix(s.URL, "/api") {
			spec.Servers[i] = openapi3.Server{URL: "https://petstore3.swagger.io" + s.URL}
		}
	}
	// endregion
	return spec
}

func loadSpecFromUrl(specUrl string) (openapi3.Spec, error) {
	spec := &openapi3.Spec{Openapi: "3.0.3"}
	req, _ := http.NewRequest(http.MethodGet, specUrl, nil)
	hc := http.DefaultClient
	res, err := hc.Do(req)
	if err != nil {
		return *spec, err
	}
	if res != nil {
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return *spec, err
		}
		if res.StatusCode >= http.StatusBadRequest {
			return *spec, errors.New("invalid status code." + res.Status)
		}
		err = json.Unmarshal(bodyBytes, &spec)
		if err != nil {
			return *spec, err
		}
		return *spec, nil
	}
	return *spec, errors.New("invalid/empty status")
}
