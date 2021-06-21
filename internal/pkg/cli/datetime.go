package cli

import (
	"fmt"
	"math"
	"time"
)

type prettyduration struct {
	secondsthreshold float64
	handler          func(time.Duration) string
}

var prettydurations = []prettyduration{
	{
		secondsthreshold: 0,
		handler: func(d time.Duration) string {
			return "in the future"
		},
	},
	{
		secondsthreshold: 60,
		handler: func(d time.Duration) string {
			return "just now"
		},
	},
	{
		secondsthreshold: 3600,
		handler: func(d time.Duration) string {
			return prettyfyduration(
				"minute",
				math.Round(d.Minutes()),
			)
		},
	},
	{
		secondsthreshold: 86400,
		handler: func(d time.Duration) string {
			return prettyfyduration(
				"hour",
				math.Round(d.Hours()),
			)
		},
	},
	{
		secondsthreshold: 172800,
		handler: func(d time.Duration) string {
			return "yesterday"
		},
	},
	{
		secondsthreshold: 604800,
		handler: func(d time.Duration) string {
			return prettyfyduration(
				"day",
				math.Round(d.Hours()/86400),
			)
		},
	},
}

func prettyfyduration(unit string, qty float64) string {
	if qty > 1 {
		unit += "s"
	}
	return fmt.Sprintf(
		"%.f %v ago",
		qty,
		unit,
	)
}

func prettyTime(t time.Time) string {
	timenow := time.Now()
	diff := timenow.Sub(t)
	duration := diff.Seconds()

	for _, pd := range prettydurations {
		if duration <= pd.secondsthreshold {
			return pd.handler(diff)
		}
	}

	return "long, long ago"
}
