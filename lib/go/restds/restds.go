package restds

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	HeaderKeyAccept        = "Accept"
	HeaderKeyContentType   = "Content-Type"
	HeaderKeyAuthorization = "Authorization"
	headerKeyIdToken       = "X-ID-Token"
)

type RestDS struct {
	Config     Config
	HTTPClient *http.Client
}

type ResponseMeta struct {
	RawURL     string
	Status     string
	StatusCode int
	Headers    http.Header
}

func (restds *RestDS) GetResponse(query Query) (responseBody string, meta ResponseMeta, err error) {
	req, err := GetRequest(restds.Config, query, map[string]string{})
	if err != nil {
		return responseBody, meta, err
	}
	res, err := restds.HTTPClient.Do(req)
	if err != nil {
		return responseBody, meta, err
	}
	if res != nil {
		defer res.Body.Close()
		meta.RawURL = req.URL.String()
		meta.Headers = res.Header
		meta.Status = res.Status
		meta.StatusCode = res.StatusCode
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", meta, fmt.Errorf("error reading the url response. %w", err)
		}
		if res.StatusCode >= http.StatusBadRequest {
			return "", meta, fmt.Errorf("invalid response received. status code: HTTP %d %s", res.StatusCode, http.StatusText(res.StatusCode))
		}
		return string(bodyBytes), meta, nil
	}
	return responseBody, meta, errors.New("unexpected error while getting data")
}
