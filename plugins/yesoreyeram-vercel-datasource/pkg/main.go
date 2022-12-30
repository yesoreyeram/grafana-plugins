package main

import (
	"fmt"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/yesoreyeram/grafana-plugins/anyframer"
	"github.com/yesoreyeram/grafana-plugins/restds"
)

func main() {
	pluginName := "Vercel"
	pluginID := "yesoreyeram-vercel-datasource"
	backend.SetupPluginEnvironment(pluginID)
	driver := &VercelRestDriver{}
	driverOptions := restds.RestDriverOptions{
		PluginID:                    pluginID,
		PluginName:                  pluginName,
		HealthCheckURL:              "https://api.vercel.com/v9/projects",
		CustomHealthCheckValidation: CustomHealthCheck,
	}
	if err := datasource.Serve(restds.NewPlugin(driver, driverOptions)); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}

func CustomHealthCheck(responseString string, meta restds.ResponseMeta) *backend.CheckHealthResult {
	if responseString == "" {
		return &backend.CheckHealthResult{
			Message: "invalid response received",
			Status:  backend.HealthStatusError,
		}
	}
	framer := &anyframer.AnyFramer{
		InputType:    anyframer.InputTypeJSON,
		RootSelector: "projects",
	}
	frame, err := framer.ToFrame(responseString)
	if err != nil {
		return &backend.CheckHealthResult{
			Message: err.Error(),
			Status:  backend.HealthStatusError,
		}
	}
	return &backend.CheckHealthResult{
		Message: fmt.Sprintf("Connected to Vercel. %d projects found.", frame.Rows()),
		Status:  backend.HealthStatusOk,
	}
}
