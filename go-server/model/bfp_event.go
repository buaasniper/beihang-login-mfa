package model

import "time"

type BfpEvent struct {
	ID             uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username       string   `gorm:"size:128;index" json:"username"`
	URL            string   `gorm:"column:url;size:512" json:"url"`
	DeltaTime      *float64 `gorm:"column:delta_time" json:"delta_time"`
	ClickTime      string   `gorm:"column:click_time;size:64" json:"click_time"`
	CookieHash     string   `gorm:"column:cookie_hash;size:64;index" json:"cookie_hash"`
	CanvasHash     string   `gorm:"column:canvas_hash;size:64;index" json:"canvas_hash"`
	WebglHash      string   `gorm:"column:webgl_hash;size:64;index" json:"webgl_hash"`
	FontsHash      string   `gorm:"column:fonts_hash;size:64;index" json:"fonts_hash"`
	UserAgent      string   `gorm:"column:user_agent;type:text" json:"user_agent"`
	RestJSON       string   `gorm:"column:rest_json;type:text" json:"rest_json"`
	CreatedAt      time.Time
}

func (BfpEvent) TableName() string { return "bfp_event" }
