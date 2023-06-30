package xmlframer

import (
	"errors"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-plugins/lib/go/jsonframer"
)

type FramerOptions struct {
	FrameName    string
	RootSelector string
	Columns      []jsonframer.ColumnSelector
}

func ToFrame(xmlString string, options FramerOptions) (*data.Frame, error) {
	xml := strings.NewReader(xmlString)
	jsonStr, err := xj.Convert(xml)
	if err != nil {
		return nil, errors.Join(errors.New("error converting xml to grafana data frame"), err)
	}
	framerOptions := jsonframer.FramerOptions{
		FramerType:   jsonframer.FramerTypeGJSON,
		FrameName:    options.FrameName,
		RootSelector: options.RootSelector,
		Columns:      options.Columns,
	}
	return jsonframer.ToFrame(jsonStr.String(), framerOptions)
}
