package model

import (
	"time"
)

type URL struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	ShortCode string    `gorm:"uniqueIndex;type:varchar(10)" json:"short_code"`
	LongURL   string    `gorm:"type:text" json:"long_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type ShortenRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}