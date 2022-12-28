package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/yesoreyeram/grafana-plugins/restds"
)

func main() {
	pluginName := "Vercel"
	pluginID := "yesoreyeram-vercel-datasource"
	backend.SetupPluginEnvironment(pluginID)
	pluginServer := restds.NewPlugin(&VercelRestDriver{}, restds.RestDriverOptions{
		PluginID:   pluginID,
		PluginName: pluginName,
	})
	if err := datasource.Serve(pluginServer); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
