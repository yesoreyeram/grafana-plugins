package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/yesoreyeram/grafana-plugins/restds"
)

func main() {
	pluginName := "Pet Store"
	pluginID := "yesoreyeram-petstore-datasource"
	backend.SetupPluginEnvironment(pluginID)
	driver := &PetStoreRestDriver{}
	driverOptions := restds.RestDriverOptions{
		PluginID:       pluginID,
		PluginName:     pluginName,
		HealthCheckURL: "https://petstore.swagger.io/v2/pet/findByStatus?status=available",
	}
	if err := datasource.Serve(restds.NewPlugin(driver, driverOptions)); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
