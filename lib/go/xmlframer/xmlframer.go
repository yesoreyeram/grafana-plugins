package xmlframer

import (
	"errors"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-plugins/lib/go/jsonframer"
)

type XMLFramerOptions struct {
	FrameName    string
	RootSelector string
	Columns      []jsonframer.ColumnSelector
}

func XmlStringToFrame(xmlString string, options XMLFramerOptions) (*data.Frame, error) {
	xml := strings.NewReader(xmlString)
	jsonStr, err := xj.Convert(xml)
	if err != nil {
		return nil, errors.Join(errors.New("error converting xml to grafana data frame"), err)
	}
	framerOptions := jsonframer.JSONFramerOptions{
		FramerType:   jsonframer.FramerTypeGJSON,
		FrameName:    options.FrameName,
		RootSelector: options.RootSelector,
		Columns:      options.Columns,
	}
	return jsonframer.JsonStringToFrame(jsonStr.String(), framerOptions)
}
