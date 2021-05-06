package db

import (
	"time"

	"gorm.io/gorm"
)

type NewsItem struct {
	gorm.Model
	Name          string             `json:"name"`
	From          *time.Time         `json:"from"`
	To            *time.Time         `json:"to"`
	LanguageItems []NewsLanguageItem `json:"languageItems"`
	SlideTime     int                `json:"slideTime"`
	Order         int                `json:"order"`
	Display       Display            `json:"display"`
	DisplayID     int                `json:"displayID"`
	Hidden        bool               `json:"hidden"`
}

type NewsLanguageItem struct {
	gorm.Model
	NewsItemID int      `json:"newsItemID"`
	NewsItem   NewsItem `json:"newsItem"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Language   string   `json:"language"`
}
