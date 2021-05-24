package db

import (
	"gorm.io/gorm"
)

type Display struct {
	gorm.Model
	Name       string `json:"name"`
	Token      string `json:"token"`
	TickerRSS  string `json:"tickerRSS"`
	TickerText string `json:"tickerText"`
	URL        string `json:"url"`
}
