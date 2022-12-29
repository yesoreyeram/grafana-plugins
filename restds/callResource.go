package restds

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (ds *pluginHost) GetRouter(restDriver RestDriver, restDriverOptions RestDriverOptions) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("pong")); err != nil {
			backend.Logger.Error("error writing resource call response", "path", "/ping", "error", err.Error())
		}
	})
	router.HandleFunc("/openapi3", func(w http.ResponseWriter, r *http.Request) {
		spec := restDriver.LoadSpec()
		res, err := spec.MarshalJSON()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if _, err := w.Write(res); err != nil {
			backend.Logger.Error("error writing resource call response", "path", "/openapi3", "error", err.Error())
		}
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
