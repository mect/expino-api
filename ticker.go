package main

// TickerItem is the configuration of an item to be shown in the Ticker
type TickerItem struct {
	// ${TICKCER_HOST}/diff/${setup}/${metric}/${interval}/${back}
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Setup    string `json:"setup"`
	Metric   string `json:"metric"`
	Interval string `json:"interval"`
	Back     string `json:"back"`
}
