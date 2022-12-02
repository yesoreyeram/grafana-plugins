package anyframer

import (
	"fmt"
	"strings"
	"time"
)

var possibleDateFormats = []string{"2006-01-02", "2006-1-2", "2006/01/02", "2006/1/2", "01/02/2006", "1/2/2006"}
var possibleDateTimeSeparators = []string{"T", " "}
var possibleTimeFormats = []string{"", "15:04", "15:04:05.999999", "15:04:05.999999Z", "15:04:05.999999 -07:00", "15:04:05 MST"}

func getTimeFromString(input string, timeFormat string) *time.Time {
	if timeFormat == "auto" {
		timeFormat = ""
	}
	var possibleLayouts = []string{time.RFC3339, "2006", timeFormat}
	for _, d := range possibleDateFormats {
		for _, t := range possibleTimeFormats {
			for _, s := range possibleDateTimeSeparators {
				l := strings.TrimSpace(fmt.Sprintf("%s%s%s", d, s, t))
				if l != "" && !strings.HasSuffix(l, " T") {
					possibleLayouts = append(possibleLayouts, l)
				}
			}
		}
	}
	possibleLayouts = append(possibleLayouts, "2006-01", "2006/01", "01-2006", "01/2006")
	for _, layout := range possibleLayouts {
		if t, err := time.Parse(layout, input); err == nil && layout != "" {
			return &t
		}
	}
	return nil
}
