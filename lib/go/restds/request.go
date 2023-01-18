package restds

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

func GetRequest(config Config, query Query, headersFromGrafana map[string]string) (*http.Request, error) {
	acceptHeader := "application/json"
	contentTypeHeader := "application/json"
	var body io.Reader
	if strings.EqualFold(http.MethodPost, string(query.Method)) {
		switch query.BodyType {
		case BodyTypeRaw:
			body = strings.NewReader(query.Body)
			if strings.TrimSpace(query.BodyContentType) != "" {
				contentTypeHeader = query.BodyContentType
			}
		case BodyTypeFormData:
			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			defer writer.Close()
			for _, f := range query.BodyForm {
				_ = writer.WriteField(f.Key, f.Value)
			}
			if writer != nil {
				contentTypeHeader = writer.FormDataContentType()
			}
			body = payload
		case BodyTypeFormReloaded:
			form := url.Values{}
			for _, f := range query.BodyForm {
				if strings.TrimSpace(f.Key) != "" {
					form.Set(f.Key, f.Value)
				}
			}
			body = strings.NewReader(form.Encode())
			contentTypeHeader = "application/x-www-form-urlencoded"
		case BodyTypeGraphQL:
			jsonData := map[string]string{"query": query.BodyGraphQLQuery}
			jsonValue, _ := json.Marshal(jsonData)
			body = strings.NewReader(string(jsonValue))
		default:
			body = strings.NewReader(query.Body)
		}
	}
	req, err := http.NewRequest(strings.ToUpper(string(query.Method)), normalizeURL(query.URL), body)
	if err != nil {
		return req, err
	}
	switch query.QueryType {
	case QueryTypeJSON:
		req.Header.Add(HeaderKeyAccept, acceptHeader)
	default:
		req.Header.Add(HeaderKeyAccept, acceptHeader)
	}
	req.Header.Add(HeaderKeyContentType, contentTypeHeader)
	for k, v := range config.Headers {
		if k != "" {
			req.Header.Add(k, v)
			if strings.EqualFold(k, HeaderKeyAccept) || strings.EqualFold(k, HeaderKeyContentType) {
				req.Header.Set(k, v)
			}
		}
	}
	for _, header := range query.Headers {
		if header.Key != "" {
			req.Header.Add(header.Key, header.Value)
			if strings.EqualFold(header.Key, HeaderKeyAccept) || strings.EqualFold(header.Key, HeaderKeyContentType) {
				req.Header.Set(header.Key, header.Value)
			}
		}
	}
	if config.AuthenticationMethod == AuthTypeBasic {
		req.Header.Set(HeaderKeyAuthorization, fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(config.BasicAuthUser+":"+config.BasicAuthPassword))))
	}
	if config.AuthenticationMethod == AuthTypeBearerToken {
		req.Header.Set(HeaderKeyAuthorization, fmt.Sprintf("Bearer %s", config.BearerToken))
	}
	if config.AuthenticationMethod == AuthTypeApiKey && config.ApiKeyType != ApiKeyTypeQuery {
		req.Header.Set(config.ApiKeyKey, config.ApiKeyValue)
	}
	if config.AuthenticationMethod == AuthTypeForwardOauth {
		req.Header.Set(HeaderKeyAuthorization, headersFromGrafana[HeaderKeyAuthorization])
		if headersFromGrafana[headerKeyIdToken] != "" {
			req.Header.Set(headerKeyIdToken, headersFromGrafana[headerKeyIdToken])
		}
	}
	q := req.URL.Query()
	for k, v := range config.QueryParams {
		if strings.TrimSpace(k) != "" {
			q.Add(k, v)
		}
	}
	if config.AuthenticationMethod == AuthTypeApiKey && config.ApiKeyType == ApiKeyTypeQuery {
		if config.ApiKeyKey != "" {
			q.Set(config.ApiKeyKey, config.ApiKeyValue)
		}
	}
	req.URL.RawQuery = q.Encode()
	return req, err
}

func normalizeURL(u string) string {
	urlArray := strings.Split(u, "/")
	if strings.HasPrefix(u, "https://github.com") && len(urlArray) > 5 && urlArray[5] == "blob" && urlArray[4] != "blob" && urlArray[3] != "blob" {
		u = strings.Replace(u, "https://github.com", "https://raw.githubusercontent.com", 1)
		u = strings.Replace(u, "/blob/", "/", 1)
	}
	return u
}
