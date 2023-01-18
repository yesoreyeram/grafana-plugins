package macros_test

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/macros"
)

func TestApplyMacros(t *testing.T) {
	from := time.UnixMilli(1500376552001).In(time.UTC) // Tue Jul 18 2017 12:15:52 GMT+0100 (British Summer Time)
	to := time.UnixMilli(1500549352001).In(time.UTC)   // Thu Jul 20 2017 12:15:52 GMT+0100 (British Summer Time)
	tests := []struct {
		name        string
		inputString string
		timeRange   backend.TimeRange
		want        string
		wantErr     bool
	}{
		{inputString: "${__from}", want: "1500376552001"},
		{inputString: "${__from:date}", want: "2017-07-18T11:15:52.001Z"},
		{inputString: "${__from:date:iso}", want: "2017-07-18T11:15:52.001Z"},
		{inputString: "foo ${__from:date:YYYY:MM:DD:hh:mm} bar", want: "foo 2017:07:18:11:15 bar"},
		{inputString: "foo ${__to:date:YYYY-MM-DD:hh,mm} bar", want: "foo 2017-07-20:11,15 bar"},
		{inputString: "from ${__from:date:iso} to ${__to:date:iso}", want: "from 2017-07-18T11:15:52.001Z to 2017-07-20T11:15:52.001Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := macros.ApplyMacros(tt.inputString, backend.TimeRange{From: from, To: to})
			require.Nil(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
