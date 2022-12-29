package restds

import (
	"encoding/json"

	"github.com/swaggest/openapi-go/openapi3"
)

func SpecFromJson(jsonString string) (openapi3.Spec, error) {
	spec := openapi3.Spec{}
	err := json.Unmarshal([]byte(jsonString), &spec)
	return spec, err
}
