package macros_test

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/lib/go/macros"
)

func TestApplyMacros(t *testing.T) {
	// https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#__from-and-__to
	from := time.UnixMilli(1594671549254).In(time.UTC)
	to := time.UnixMilli(1500549352001).In(time.UTC)
	tests := []struct {
		name        string
		inputString string
		timeRange   backend.TimeRange
		pluginCtx   backend.PluginContext
		want        string
		wantErr     bool
	}{
		{inputString: "${__from}", want: "1594671549254"},
		{inputString: "${__from:date}", want: "2020-07-13T20:19:09.254Z"},
		{inputString: "${__from:date:seconds}", want: "1594671549"},
		{inputString: "${__from:date:iso}", want: "2020-07-13T20:19:09.254Z"},
		{inputString: "foo ${__from:date:YYYY:MM:DD:hh:mm} bar", want: "foo 2020:07:13:08:19 bar"},
		{inputString: "foo ${__from:date:YYYY:MM:DD:HH:mm} bar", want: "foo 2020:07:13:20:19 bar"},
		{inputString: "foo ${__to:date:YYYY-MM-DD:hh,mm} bar", want: "foo 2017-07-20:11,15 bar"},
		{inputString: "from ${__from:date:iso} to ${__to:date:iso}", want: "from 2020-07-13T20:19:09.254Z to 2017-07-20T11:15:52.001Z"},

		{inputString: "${__timeFrom}", want: "1594671549254"},
		{inputString: "${__timeFrom:date} ${__timeFrom:date}", want: "2020-07-13T20:19:09.254Z 2020-07-13T20:19:09.254Z"},
		{inputString: "from ${__timeFrom:date:iso} to ${__timeTo:date:iso}", want: "from 2020-07-13T20:19:09.254Z to 2017-07-20T11:15:52.001Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := macros.ApplyMacros(
				tt.inputString,
				backend.DataQuery{TimeRange: backend.TimeRange{From: from, To: to}},
				tt.pluginCtx,
			)
			require.Nil(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
