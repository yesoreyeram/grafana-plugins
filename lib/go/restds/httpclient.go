package restds

import (
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
)

func NewHTTPClient(config *Config) *http.Client {
	hc, _ := httpclient.New()
	return hc
}
