package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/swaggest/openapi-go/openapi3"
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

func (v *VercelRestDriver) LoadSpec() openapi3.Spec {
	spec := openapi3.Spec{Openapi: "3.0.3"}
	if spec, err := loadSpecFromUrl("https://openapi.vercel.sh"); err == nil {
		// This part seems to be failing
		// https://github.com/vercel/community/discussions/646#discussioncomment-4553140
		return spec
	}
	schemaTypeString := openapi3.SchemaTypeString
	schemaTypeNumber := openapi3.SchemaTypeNumber
	stringSchema := &openapi3.SchemaOrRef{Schema: &openapi3.Schema{Type: &schemaTypeString}}
	numberSchema := &openapi3.SchemaOrRef{Schema: &openapi3.Schema{Type: &schemaTypeNumber}}
	spec = *spec.WithInfo(openapi3.Info{
		Title:   "Vercel API",
		Version: "0.0.1-dev.1",
		Contact: &openapi3.Contact{
			URL: pointer("https://api.vercel.com"),
		},
	})
	spec = *spec.WithServers(openapi3.Server{
		URL:         "{baseUrl}",
		Description: pointer("All endpoints live under the URL https://api.vercel.com and this is the default URL"),
		Variables: map[string]openapi3.ServerVariable{
			"baseUrl": {
				Default:     "https://api.vercel.com",
				Description: pointer("Full base url of the vercel api. Typically it is https://api.vercel.com"),
			},
		},
	})
	spec = *spec.WithPaths(openapi3.Paths{
		MapOfPathItemValues: map[string]openapi3.PathItem{
			"/v9/projects": {
				Summary:     pointer("Retrieve a list of projects"),
				Description: pointer("Allows to retrieve the list of projects of the authenticated user. The list will be paginated and the provided query parameters allow filtering the returned projects."),
				MapOfOperationValues: map[string]openapi3.Operation{
					"get": {
						Parameters: []openapi3.ParameterOrRef{
							{Parameter: &openapi3.Parameter{Name: "edgeConfigId", In: openapi3.ParameterInQuery, Description: pointer("Filter results by connected Edge Config ID")}},
							{Parameter: &openapi3.Parameter{Name: "edgeConfigTokenId", In: openapi3.ParameterInQuery, Description: pointer("Filter results by connected Edge Config Token ID"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "excludeRepos", In: openapi3.ParameterInQuery, Description: pointer("Filter results by excluding those projects that belong to a repo"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "from", In: openapi3.ParameterInQuery, Description: pointer("Query only projects updated after the given timestamp"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "gitForkProtection", In: openapi3.ParameterInQuery, Description: pointer("Specifies whether PRs from Git forks should require a team member's authorization before it can be deployed"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "limit", In: openapi3.ParameterInQuery, Description: pointer("Limit the number of projects returned"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "repo", In: openapi3.ParameterInQuery, Description: pointer("Filter results by repo. Also used for project count"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "repoId", In: openapi3.ParameterInQuery, Description: pointer("Filter results by Repository ID."), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "repoUrl", In: openapi3.ParameterInQuery, Description: pointer("Filter results by Repository URL."), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "search", In: openapi3.ParameterInQuery, Description: pointer("Search projects by the name field"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "teamId", In: openapi3.ParameterInQuery, Description: pointer("The Team identifier or slug to perform the request on behalf of."), Schema: stringSchema}},
						},
					},
				},
			},
			"/v2/teams": {
				Summary:     pointer("List all teams"),
				Description: pointer("Get a paginated list of all the Teams the authenticated User is a member of."),
				MapOfOperationValues: map[string]openapi3.Operation{
					"get": {
						Parameters: []openapi3.ParameterOrRef{
							{Parameter: &openapi3.Parameter{Name: "limit", In: openapi3.ParameterInQuery, Description: pointer("Maximum number of Teams which may be returned."), Schema: numberSchema}},
							{Parameter: &openapi3.Parameter{Name: "since", In: openapi3.ParameterInQuery, Description: pointer("Timestamp (in milliseconds) to only include Teams created since then."), Schema: numberSchema}},
							{Parameter: &openapi3.Parameter{Name: "until", In: openapi3.ParameterInQuery, Description: pointer("Timestamp (in milliseconds) to only include Teams created until then."), Schema: numberSchema}},
						},
					},
				},
			},
			"/v2/teams/{teamId}/members": {
				Summary:     pointer("List team members"),
				Description: pointer("Get a paginated list of team members for the provided team."),
				Parameters: []openapi3.ParameterOrRef{
					{Parameter: &openapi3.Parameter{Name: "teamId", In: openapi3.ParameterInPath, Description: pointer("Team ID"), Schema: stringSchema, Required: pointer(true)}},
				},
				MapOfOperationValues: map[string]openapi3.Operation{
					"get": {
						Parameters: []openapi3.ParameterOrRef{
							{Parameter: &openapi3.Parameter{Name: "excludeProject", In: openapi3.ParameterInQuery, Description: pointer("Exclude members who belong to the specified project"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "limit", In: openapi3.ParameterInQuery, Description: pointer("Limit how many teams should be returned"), Schema: numberSchema}},
							{Parameter: &openapi3.Parameter{Name: "role", In: openapi3.ParameterInQuery, Description: pointer("Only return members with the specified team role"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "search", In: openapi3.ParameterInQuery, Description: pointer("Search team members by their name, username, and email"), Schema: stringSchema}},
							{Parameter: &openapi3.Parameter{Name: "since", In: openapi3.ParameterInQuery, Description: pointer("Timestamp in milliseconds to only include members added since then"), Schema: numberSchema}},
							{Parameter: &openapi3.Parameter{Name: "until", In: openapi3.ParameterInQuery, Description: pointer("Timestamp in milliseconds to only include members added until then"), Schema: numberSchema}},
						},
					},
				},
			},
		},
	})
	return spec
}

func loadSpecFromUrl(u string) (openapi3.Spec, error) {
	spec := &openapi3.Spec{Openapi: "3.0.3"}
	req, _ := http.NewRequest(http.MethodGet, "https://openapi.vercel.sh", nil)
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

func pointer[T any](input T) *T { return &input }
