package main

// GraphItem contains the info for a Grafana graph to be embedded
type GraphItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
