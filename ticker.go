package main

import (
	"time"
)

// TickerItem is the configuration of an item to be shown in the Ticker
type TickerItem struct {
	// ${TICKCER_HOST}/diff/${setup}/${metric}/${interval}/${back}
	ID       int           `json:"id"`
	Setup    string        `json:"setup"`
	Metric   string        `json:"metric"`
	Interval time.Duration `json:"interval"`
	Back     time.Duration `json:"back"`
}
