package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/yesoreyeram/grafana-plugins/restds"
)

type VercelConfig struct {
	BaseURL string `json:"apiUrl,omitempty"`
}

type VercelRestDriver struct{}

func (v *VercelRestDriver) LoadConfig(settings backend.DataSourceInstanceSettings) (*restds.Config, error) {
	config := &VercelConfig{
		BaseURL: "https://api.vercel.com",
	}
	configJson := settings.JSONData
	if configJson == nil {
		configJson = []byte("{}")
	}
	if err := json.Unmarshal(configJson, config); err != nil {
		return nil, fmt.Errorf("error while reading the config. %w", err)
	}
	apiToken, ok := settings.DecryptedSecureJSONData["apiToken"]
	if !ok || apiToken == "" {
		return nil, errors.New("invalid/empty api token")
	}
	return &restds.Config{
		BaseURL:              config.BaseURL,
		AuthenticationMethod: restds.AuthTypeBearerToken,
		BearerToken:          apiToken,
	}, nil
}
